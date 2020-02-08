package push

import (

	"github.com/aws/aws-sdk-go/aws/session"
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
	config config
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

func (a *S3Push) GetEnvironment() error {
	a.config.S3Config = connect.GetS3Configuration()
	return nil
}


func (a *S3Push) Check() error {
	s, err := connect.NewS3SessionFromConfig(&a.config.S3Config)
	if err != nil {
		return err
	}
	a.session = s
	return nil
}

func (a *S3Push) Connect() error {
	return nil
}

func (a *S3Push) Do() error {
	panic("implement me")
}
