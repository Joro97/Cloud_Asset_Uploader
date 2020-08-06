package config

import (
	"CloudAssetUploader/data"
	"CloudAssetUploader/uploader"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// The environment required for the application
type Env struct {
	AssetUploader uploader.Uploader
	Store data.Store
}

//
func NewEnv(client *session.Session) *Env {
	return &Env{
		AssetUploader: &uploader.AwsAssetUploader{
			AWSClient: client,
			S3Manager: s3.New(session.Must(client, nil)),
		},
	}
}
