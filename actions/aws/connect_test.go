package aws

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/franchb/go-aws-s3-static-site-push/pkg/testhelpers"
	"os"
	"reflect"
	"testing"
)

func TestGetS3Configuration(t *testing.T) {
	tests := []struct {
		name string
		envs map[string]string
		want S3Config
	}{
		{
			name: "basic1",
			envs: map[string]string{
				EnvBucket:     "SampleBucket",
				EnvRegion:     "eu-central-1",
				EnvAccelerate: "true",
			},
			want: S3Config{
				Region:          aws.String("eu-central-1"),
				Bucket:          aws.String("SampleBucket"),
				S3UseAccelerate: aws.Bool(true),
			},
		},
		{
			name: "basic2",
			envs: map[string]string{
				EnvBucket:     "NoBucket",
				EnvRegion:     "eu-north-1",
				EnvAccelerate: "false",
			},
			want: S3Config{
				Region:          aws.String("eu-north-1"),
				Bucket:          aws.String("NoBucket"),
				S3UseAccelerate: aws.Bool(false),
			},
		},
		{
			name: "empty",
			envs: map[string]string{
				EnvBucket:     "",
				EnvRegion:     "",
				EnvAccelerate: "",
			},
			want: S3Config{
				Region:          nil,
				Bucket:          nil,
				S3UseAccelerate: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer testhelpers.UnsetEnv("AWS")()
			for k, v := range tt.envs {
				if err := os.Setenv(k, v); err != nil {
					t.Errorf("failed to set environment variable %s: %v", k, err)
				}
			}
			got := GetS3Configuration()
			if !reflect.DeepEqual(got, tt.want) {
				g, _ := json.Marshal(got)
				w, _ := json.Marshal(tt.want)
				t.Errorf("NewS3SessionFromConfig() got = %+v, want %+v", string(g), string(w))
			}
		})
	}
}

func TestNewS3SessionFromConfig(t *testing.T) {
	type args struct {
		config S3Config
	}
	type wantArgs struct {
		region *string
		creds credentials.Value
		s3UseAccelerate *bool
	}
	tests := []struct {
		name    string
		args    args
		envs    map[string]string
		want    wantArgs
		wantErr bool
		wantErrCode string
	}{
		{
			name: "basic",
			args: args{
				config: S3Config{
					Region:          aws.String("eu-north-1"),
					Bucket:          aws.String("NoBucket"),
					S3UseAccelerate: aws.Bool(true),
				},
			},
			envs: map[string]string{
				"AWS_ACCESS_KEY_ID": "AKIAIOSFODNN7EXAMPLE",
				"AWS_SECRET_ACCESS_KEY": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			want: wantArgs{
				region: aws.String("eu-north-1"),
				s3UseAccelerate: aws.Bool(true),
				creds: credentials.Value{
					AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
					SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				},
			},
			wantErr: false,
		},
		{
			name: "withEmptyCredentialsEnv",
			args: args{
				config: S3Config{
					Region:          aws.String("eu-north-1"),
					Bucket:          aws.String("NoBucket"),
					S3UseAccelerate: aws.Bool(true),
				},
			},
			envs: map[string]string{
			},
			want: wantArgs{
				region: aws.String("eu-north-1"),
				s3UseAccelerate: aws.Bool(true),
				creds: credentials.Value{
					AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
					SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				},
			},
			wantErr: true,
			wantErrCode: "EnvAccessKeyNotFound",
		},
		{
			name: "withEmptyBucketAndRegionFailed",
			args: args{
				config: S3Config{
					Region:          nil,
					Bucket:          nil,
					S3UseAccelerate: aws.Bool(true),
				},
			},
			envs: map[string]string{
				"AWS_ACCESS_KEY_ID": "AKIAIOSFODNN7EXAMPLE",
				"AWS_SECRET_ACCESS_KEY": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			want: wantArgs{
				region: aws.String("eu-north-1"),
				s3UseAccelerate: aws.Bool(true),
				creds: credentials.Value{
					AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
					SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				},
			},
			wantErr: true,
			wantErrCode: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer testhelpers.UnsetEnv("AWS")()

			for k, v := range tt.envs {
				if err := os.Setenv(k, v); err != nil {
					t.Errorf("failed to set environment variable %s: %v", k, err)
				}
			}
			got, err := NewS3SessionFromConfig(tt.args.config)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewS3SessionFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				} else {
					if !errors.Is(err, ErrConfigurationNotValid) {
						t.Errorf("NewS3SessionFromConfig() error = %v, wantErr %v", err, tt.wantErr)
					}
				}
				return
			}
			if *tt.want.region != *got.Config.Region {
				t.Errorf("NewS3SessionFromConfig() got = %v, want %v", *got.Config.Region, *tt.want.region)
			}
			if *tt.want.region != *got.Config.Region {
				t.Errorf("NewS3SessionFromConfig() got = %v, want %v", *got.Config.Region, *tt.want.region)
			}
			creds, err := got.Config.Credentials.Get()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewS3SessionFromConfig() get creds error = %v, wantErr %v", err, tt.wantErr)
				} else {
					if awsErr, ok := err.(awserr.Error); ok {
						if awsErr.Code() != tt.wantErrCode {
							t.Errorf(
								"NewS3SessionFromConfig() Config.Credentials.Get() error = %v, wantErrCode %s",
								err, tt.wantErrCode)
						}
					}
				}
				return
			}
			if tt.want.creds.AccessKeyID != creds.AccessKeyID {
				t.Errorf("NewS3SessionFromConfig() got = %v, want %v",
					tt.want.creds.AccessKeyID, creds.AccessKeyID)
			}
			if tt.want.creds.SecretAccessKey != creds.SecretAccessKey {
				t.Errorf("NewS3SessionFromConfig() got = %v, want %v",
					tt.want.creds.SecretAccessKey, creds.SecretAccessKey)
			}
		})
	}
}
