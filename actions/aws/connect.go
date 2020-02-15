package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"strconv"
)

var (
	ErrConfigurationNotValid = errors.New("AWS configuration is not valid")
)

const (
	EnvBucket     = "AWS_S3_BUCKET"
	EnvRegion     = "AWS_REGION"
	EnvAccelerate = "AWS_S3_USE_ACCELERATE"
)

type S3Config struct {
	Region          *string /* AWS region to use for session */
	Bucket          *string /* S3 bucket where are files are stored */
	S3UseAccelerate *bool
}

func GetS3Configuration() S3Config {
	var configuration S3Config

	region := os.Getenv(EnvRegion)
	if region == "" {
		configuration.Region = nil
	} else {
		configuration.Region = aws.String(region)
	}

	bucket := os.Getenv(EnvBucket)
	if bucket == "" {
		configuration.Bucket = nil
	} else {
		configuration.Bucket = aws.String(bucket)
	}

	if ac, err := strconv.ParseBool(os.Getenv(EnvAccelerate)); err == nil {
		configuration.S3UseAccelerate = aws.Bool(ac)
	}
	return configuration
}

// NewS3SessionFromConfig returns the AWS session with configuration.
func NewS3SessionFromConfig(config S3Config) (*session.Session, error) {
	if config.Region == nil || *config.Region == "" {
		return nil, fmt.Errorf("AWS Region is empty: %w", ErrConfigurationNotValid)
	}
	if config.Bucket == nil || *config.Bucket == "" {
		return nil, fmt.Errorf("AWS S3 Bucket Name is empty: %w", ErrConfigurationNotValid)
	}
	return session.NewSession(
		&aws.Config{
			Region:          config.Region,
			Credentials:     credentials.NewEnvCredentials(),
			S3UseAccelerate: config.S3UseAccelerate,
		})
}
