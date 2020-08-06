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
	GetSignedUploadURL(assetName string) (url string, id string, err error)
}

//
type AwsAssetUploader struct {
	AWSClient *session.Session
	S3Manager *s3.S3
}

//
func (upld *AwsAssetUploader) GetSignedUploadURL(assetName string) (string, string, error) {
	resp, _ := upld.S3Manager.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(constants.DEFAULT_BUCKET_NAME),
		Key: aws.String(assetName),
	})

	//fmt.Printf("The output is: %+v\n", output)
	//fmt.Printf("The response is: %+v\n", resp)
	url, err := resp.Presign(3 * time.Minute)
	if err != nil {
		return "", "", err
	}
	return url, "", nil
}
