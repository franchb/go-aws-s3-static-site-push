package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	AWS_S3_REGION = ""
	AWS_S3_BUCKET = ""
)

var sess = connectAWS()

func connectAWS() *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(AWS_S3_REGION),
		},
	)
	if err != nil {
		panic(err)
	}
	return sess
}

/* getS3Session returns the AWS session with our configuration.
   For our purposes this could have gone in getS3Service,
   but in case we end up relying on other AWS clients
	     it's handy to have the session available separately.
*/
func getS3Session() *session.Session {
	// Use shared credentials ($HOME/.aws/credentials),
	// but set the region
	return session.Must(session.NewSession(
		&aws.Config{
			Region:      aws.String(getConfiguration().region),
			Credentials: credentials.NewSharedCredentials("", "default"),
		}))
}

/*
    getS3Service returns a service we use to query our bucket
	    and to put new / changed files.
*/
func getS3Service() *s3.S3 {
	return s3.New(getS3Session())
}

type config struct {
	region         string /* AWS region to use for session */
	bucket         string /* S3 bucket where are files are stored */
	localDirectory string /* Local files to sync with S3 bucket */
}

func getConfiguration() *config {
	configuration := config{
		region:         "us-east-1",
		bucket:         "irusin.dev",
		localDirectory: "/home/ubuntu/projects/github.com/franchb/"}
	return &configuration
}

/**
getAwsS3ItemMap constructs and returns a map of keys (relative filenames)
to checksums for the given bucket-configured s3 service.
It is assumed  that the objects have not been multipart-uploaded,
which will change the checksum.
TODO: We’re concerned with new or updated files – any deletes are ignored and we can handle separately.
*/
func getAwsS3ItemMap(s3Service *s3.S3) (map[string]string, error) {

	var loi s3.ListObjectsInput
	loi.SetBucket(getConfiguration().bucket)

	obj, err := s3Service.ListObjects(&loi)

	// Uncomment this to see what AWS returns for ListObjects
	// fmt.Printf("%s\n", obj)

	var items = make(map[string]string)

	// Based on response from aws API, map relative filenames to checksums
	// The "Key" in S3 is the filename, while the ETag contains the
	// checksum
	if err == nil {
		for _, s3obj := range obj.Contents {
			eTag := strings.Trim(*(s3obj.ETag), "\"")
			items[*(s3obj.Key)] = eTag
		}
		return items, nil
	}

	return nil, err
}

func fileWalk(awsItems map[string]string) []string {
	blogDir := getConfiguration().localDirectory
	ch := make(chan string)
	var numfiles int
	filepath.Walk(blogDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			numfiles++
			relativeFile := strings.Replace(path, blogDir, "", 1)
			checksumRemote, _ := awsItems[relativeFile]

			go testFile(path, checksumRemote, ch)
			//fmt.Println(path, info.Size())
		}
		return nil
	})

	var filesToUpload []string

	var f string
	for i := 0; i < numfiles; i++ {
		f = <-ch
		if len(f) > 0 {
			filesToUpload = append(filesToUpload, f)
		}
	}

	return filesToUpload
}

func testFile(filename string, checksumRemote string, ch chan<- string) {

	if checksumRemote == "" {
		ch <- filename
		return
	}

	contents, err := ioutil.ReadFile(filename)
	if err == nil {
		sum := md5.Sum(contents)
		sumString := fmt.Sprintf("%x", sum)
		// checksums don't match, mark for upload
		if sumString != checksumRemote {
			ch <- filename
			return
		} else {
			// Files matched
			ch <- ""
			return
		}
	}
}
