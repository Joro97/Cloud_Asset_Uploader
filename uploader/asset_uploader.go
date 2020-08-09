package uploader

import (
	"os"
	"strconv"
	"time"

	"CloudAssetUploader/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	bucketName string
)

// Uploader is an interface for interacting with underlying cloud where the assets will be stored.
type Uploader interface {
	SetupBucket() error
	GetSignedUploadURL(timeout string) (awsName, url string, err error)
	GetSignedDownloadURL(assetName, timeout string) (url string, er error)
}

// AwsAssetUploader is an Uploader for AWS.
type AwsAssetUploader struct {
	AWSClient *session.Session
	S3Manager *s3.S3
}

// SetupBucket creates an AWS bucket with the specified name from ENV var if it does not exist already.
// If no var is specified, it uses the default name.
func (upld *AwsAssetUploader) SetupBucket() error {
	bucketName = os.Getenv("AWS_BUCKET_NAME")
	if bucketName == "" {
		bucketName = constants.DefaultBucketName
	}

	result, err := upld.S3Manager.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Error().Msgf("Could not list AWS buckets. Err: %s", err)
		return err
	}

	shouldCreateBucket := true
	for _, bucket := range result.Buckets {
		if *bucket.Name == bucketName {
			shouldCreateBucket = false
			break
		}
	}

	if shouldCreateBucket {
		// Create the S3 Bucket
		_, err = upld.S3Manager.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}

		// Wait until bucket is created before finishing
		err = upld.S3Manager.WaitUntilBucketExists(&s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSignedUploadURL creates a presigned URL for uploading an asset to AWS. The URL expiry is specified by timeout.
func (upld *AwsAssetUploader) GetSignedUploadURL(timeout string) (string, string, error) {
	secondsTimeout := validateTimeout(timeout)
	awsName := uuid.New().String()

	resp, _ := upld.S3Manager.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(awsName),
	})

	url, err := resp.Presign(time.Duration(secondsTimeout) * time.Second)
	if err != nil {
		return "", "", err
	}

	return awsName, url, nil
}

// GetSignedDownloadURL creates a presigned URL for downloading a currently existing asset it AWS. URL expiry is specified by the timeout param.
func (upld *AwsAssetUploader) GetSignedDownloadURL(assetName, timeout string) (string, error) {
	secondsTimeout := validateTimeout(timeout)

	resp, _ := upld.S3Manager.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(assetName),
	})

	url, err := resp.Presign(time.Duration(secondsTimeout) * time.Second)
	if err != nil {
		log.Error().Msgf("Could not create a download URL for given object. Err: %s", err)
		return "", err
	}
	return url, nil
}

func validateTimeout(timeout string) int {
	timeoutSeconds := constants.DefaultURLExpiryTimeInSeconds
	if val, err := strconv.Atoi(timeout); err == nil {
		if val < constants.MinimumURLExpiryTimeInSeconds || val > constants.MaximumURLExpiryTimeInSeconds {
			val = constants.DefaultURLExpiryTimeInSeconds
		}
		timeoutSeconds = val
	}
	return timeoutSeconds
}
