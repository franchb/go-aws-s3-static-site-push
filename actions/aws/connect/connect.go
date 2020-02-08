package connect

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"strconv"
)

const (
	EnvBucket    = "AWS_S3_BUCKET"
	EnvRegion    = "AWS_REGION"
	EnvAccelerate = "AWS_S3_USE_ACCELERATE"
)

type S3Config struct {
	Region         *string /* AWS region to use for session */
	Bucket         *string /* S3 bucket where are files are stored */
	S3UseAccelerate *bool
}

func GetS3Configuration() S3Config {
	region :=  os.Getenv(EnvRegion)
	bucket := os.Getenv(EnvBucket)
	configuration := S3Config{
		Region:         &region,
		Bucket:         &bucket}

	if ac, err := strconv.ParseBool(os.Getenv(EnvAccelerate)); err == nil {
		configuration.S3UseAccelerate = &ac
	}
	return configuration
}

func NewS3SessionFromConfig(config *S3Config) (*session.Session, error) {
	return session.NewSession(
		&aws.Config{
			Region:      config.Region,
			Credentials: credentials.NewEnvCredentials(),
			S3UseAccelerate: config.S3UseAccelerate,
		})
}
