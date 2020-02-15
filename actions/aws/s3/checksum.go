package s3

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
)

// getS3ObjectsChecksumsMap constructs and returns a map of keys (relative filenames)
// to checksums for the given bucket-configured S3 service.
// It is assumed  that the objects have not been multipart-uploaded,
// which will change the checksum.
// TODO: @franchb - We’re concerned with new or updated files –
//  any deletes are ignored and we can handle separately.
func getS3ObjectsChecksumsMap(s3Service *s3.S3, bucket string) (map[string]string, error) {
	objects, err := s3Service.ListObjects(&s3.ListObjectsInput{Bucket: &bucket})
	if err != nil {
		return nil, err
	}

	var items = make(map[string]string, len(objects.Contents))
	// The "Key" in S3 is the filename, while the ETag contains the
	// MD5 checksum
	for i := range objects.Contents {
		eTag := strings.Trim(*objects.Contents[i].ETag, "\"")
		// TODO: @franchb - check for multipart upload hash (md5-numberOfParts)
		//  https://stackoverflow.com/questions/12186993/what-is-the-algorithm-to-compute-the-amazon-s3-etag-for-a-file-larger-than-5gb
		items[*objects.Contents[i].Key] = eTag
	}
	return items, nil
}

