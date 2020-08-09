package uploader

import (
	"strconv"
	"time"

	"CloudAssetUploader/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Uploader is an interface for interacting with underlying cloud where the assets will be stored.
type Uploader interface {
	GetSignedUploadURL(timeout string) (awsName, url string, err error)
	GetSignedDownloadURL(assetName, timeout string) (url string, er error)
}

// AwsAssetUploader is an Uploader for AWS.
type AwsAssetUploader struct {
	AWSClient *session.Session
	S3Manager *s3.S3
}

// GetSignedUploadURL creates a presigned URL for uploading an asset to AWS. The URL expiry is specified by timeout.
func (upld *AwsAssetUploader) GetSignedUploadURL(timeout string) (string, string, error) {
	secondsTimeout := validateTimeout(timeout)
	awsName := uuid.New().String()

	resp, _ := upld.S3Manager.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(constants.DefaultBucketName),
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
		Bucket: aws.String(constants.DefaultBucketName),
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
