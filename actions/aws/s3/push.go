package s3

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	awsActions "github.com/franchb/go-aws-s3-static-site-push/actions/aws"
	"github.com/franchb/go-aws-s3-static-site-push/pkg/filechecksum"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

const (
	EnvSourceDir = "SOURCE_DIR"
	// EnvPublicRead makes uploaded files publicly readable,
	// bucket settings should be set to public.
	EnvPublicRead = "ACL_PUBLIC_READ"
	// EnvSlackUserName  fixes some symbolic link problems
	EnvSlackUserName = "FOLLOW_SYMLINKS"
	// EnvDeleteOld - permanently deletes files in the S3 bucket
	// that are not present in the latest version of folder provided.
	EnvDeleteOld = "DELETE_OLD"
	// EnvDestDir - the directory inside of the S3 bucket to sync/upload to,
	// defaults to the root of the bucket.
	EnvDestDir = "DEST_DIR"
)

type S3Push struct {
	session *session.Session
	s3      *s3.S3
	config  config
}

func NewS3PushAction() *S3Push {
	return &S3Push{}
}

type config struct {
	awsActions.S3Config
	sourceDir string
}

func (a *S3Push) Name() string {
	return "S3 Push"
}

func (a *S3Push) Help() string {
	return "implement me"
}

func (a *S3Push) Open() error {
	if err := a.getEnvironment(); err != nil {
		return err
	}
	if err := a.check(); err != nil {
		return err
	}
	a.s3 = s3.New(a.session)
	return nil
}

func (a *S3Push) Do(ctx aws.Context) error {
	checksums, err := getS3ObjectsChecksumsMap(a.s3, *a.config.Bucket)
	if err != nil {
		return err
	}
	changedCh := filechecksum.GetListOfChangedFilesChan("", checksums)

	changedCount, updatedCount := 0, 0
	for c := range changedCh {
		changedCount++
		filename, checksum := c[0], c[1]
		err := a.uploadFileToS3(ctx, filename, checksum)
		if err != nil {
			// TODO
		}
		updatedCount++
	}

	if changedCount == 0 {
		log.Info().Msg("all files are up to date")
		return nil
	}
	return nil
}

func (a *S3Push) Close() error {
	return errors.New("not implemented")
}

func (a *S3Push) getEnvironment() error {
	a.config.S3Config = awsActions.GetS3Configuration()
	return nil
}

func (a *S3Push) check() error {
	s, err := awsActions.NewS3SessionFromConfig(a.config.S3Config)
	if err != nil {
		return err
	}
	a.session = s
	return nil
}

// uploadFileToS3 uploads file to S3
func (a *S3Push) uploadFileToS3(ctx aws.Context, fileName string, md5 string) error {

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close file to upload")
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to obtain file stats: %w", err)
	}
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", fileName, err)
	}

	inp := s3.PutObjectInput{
		ACL:                     aws.String("public-read"),
		Body:                    bytes.NewReader(buffer),
		Bucket:                  a.config.Bucket,
		ContentDisposition:      aws.String("attachment"),
		ContentEncoding:         nil,
		ContentLanguage:         nil,
		ContentLength:           aws.Int64(size),
		ContentMD5:              aws.String(md5),
		ContentType:             aws.String(http.DetectContentType(buffer)),
		Key:                     aws.String(fileName),
		Metadata:                nil,
		StorageClass:            aws.String("INTELLIGENT_TIERING"),
		Tagging:                 nil,
		WebsiteRedirectLocation: nil,
	}
	_, err = a.s3.PutObjectWithContext(ctx, &inp)
	if err != nil {
		return fmt.Errorf("failed to put file %s to S3: %w", fileName, err)
	}

	return nil
}
