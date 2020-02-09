package push

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/franchb/go-aws-s3-static-site-push/actions/aws/connect"

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
	EnvDestDir    = "DEST_DIR"
)

type S3Push struct {
	session *session.Session
	s3 *s3.S3
	config config
}

func NewS3PushAction() *S3Push {
	return &S3Push{}
}

type config struct {
	connect.S3Config
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

func (a *S3Push) Do() error {
	return errors.New("not implemented")
}

func (a *S3Push) Close() error {
	return errors.New("not implemented")
}


func (a *S3Push) getEnvironment() error {
	a.config.S3Config = connect.GetS3Configuration()
	return nil
}

func (a *S3Push) check() error {
	s, err := connect.NewS3SessionFromConfig(a.config.S3Config)
	if err != nil {
		return err
	}
	a.session = s
	return nil
}