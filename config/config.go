package config

import (
	"CloudAssetUploader/uploader"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws/session"
)

// The environment required for the application
type Env struct {
	AssetUploader uploader.Uploader
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
