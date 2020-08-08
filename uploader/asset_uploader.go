package uploader

import (
	"fmt"
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

type ErrorInvalidAssetName struct {
	Name string
}

func (err *ErrorInvalidAssetName) Error() string {
	return fmt.Sprintf("The asset name should be between %d and %d characters long. Current name: %s",
		constants.AssetMinNameLength, constants.AssetMaxNameLength, err.Name)
}

//
func (upld *AwsAssetUploader) GetSignedUploadURL(assetName string) (string, error) {
	assetName, err := validateAssetName(assetName)
	if err != nil {
		return "", err
	}

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

func validateAssetName(name string) (string, error) {
	if len(name) < constants.AssetMinNameLength || len(name) > constants.AssetMaxNameLength {
		return name, &ErrorInvalidAssetName{Name: name}
	}
	return name, nil
}
