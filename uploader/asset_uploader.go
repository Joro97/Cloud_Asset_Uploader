package uploader

import (
	"time"

	"CloudAssetUploader/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
)

//
type Uploader interface {
	GetSignedUploadURL(assetName string) (url string, err error)
	GetSignedDownloadURL(assetName string, timeout int) (url string, er error)
}

//
type AwsAssetUploader struct {
	AWSClient *session.Session
	S3Manager *s3.S3
}

//
func (upld *AwsAssetUploader) GetSignedUploadURL(assetName string) (string, error) {
	resp, _ := upld.S3Manager.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(constants.DEFAULT_BUCKET_NAME),
		Key:    aws.String(assetName),
	})

	url, err := resp.Presign(3 * time.Minute)
	if err != nil {
		return "", err
	}

	return url, nil
}

//
func (upld *AwsAssetUploader) GetSignedDownloadURL(assetName string, timeout int) (string, error) {
	resp, _ := upld.S3Manager.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(constants.DEFAULT_BUCKET_NAME),
		Key:    aws.String(assetName),
	})

	url, err := resp.Presign(time.Duration(timeout) * time.Second)
	if err != nil {
		log.Error().Msgf("Could not create a download URL for given object. Err: %s", err)
		return "", err
	}
	return url, nil
}
