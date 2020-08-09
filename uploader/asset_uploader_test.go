package uploader

import (
	"os"
	"testing"

	"CloudAssetUploader/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

var (
	upld Uploader
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.Region),
	})
	if err != nil {
		log.Fatal().Msgf("Could not connect to aws: %s", err)
	}

	upld = &AwsAssetUploader{
		AWSClient: sess,
		S3Manager: s3.New(session.Must(sess, nil)),
	}
}

func TestIntegrationGetSignedUploadURLShouldReturnValidUUID(t *testing.T) {
	id, _, err := upld.GetSignedUploadURL("40")
	require.NoError(t, err)

	_, err = uuid.Parse(id)
	require.NoError(t, err)
}

func TestIntegrationGetSignedDownloadURLShouldNotThrowErrorWhenValidId(t *testing.T) {
	_, err := upld.GetSignedDownloadURL(constants.MockID, "")
	require.NoError(t, err)
}
