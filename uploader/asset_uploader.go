package uploader

import (
	"time"

	"CloudAssetUploader/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//
type Uploader interface {
	GetSignedUploadURL(assetName string) (url string, err error)
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
