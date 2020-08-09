package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"CloudAssetUploader/config"
	"CloudAssetUploader/constants"
	"CloudAssetUploader/data"
	"CloudAssetUploader/responses"
	"CloudAssetUploader/server"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

var (
	env        *config.Env
	httpClient *http.Client
)

func setUp() {
	httpClient = &http.Client{}
	setLocalhostEnvVars()

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = constants.DefaultRegion
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		log.Fatal().Msgf("Could not connect to aws: %s", err)
	}

	connStr, err := data.BuildConnectionStringForDB()
	if err != nil {
		log.Fatal().Msgf("Could not build connection for MongoDB: %s", err)
	}

	db, err := data.NewDB(connStr)
	if err != nil {
		log.Fatal().Msgf("Could not connect to MongoDB: %s", err)
	}

	env = config.NewEnv(sess, db)
	err = env.AssetUploader.SetupBucket()
	if err != nil {
		log.Fatal().Msgf("Could not setup AWS bucket env. Err: %s", err)
	}
}

func setLocalhostEnvVars() {
	err := os.Setenv("MONGO_USERNAME", "mongoadmin")
	err = os.Setenv("MONGO_PASSWORD", "bigSecret")
	err = os.Setenv("MONGO_CONTAINER_NAME", "localhost")
	err = os.Setenv("MONGO_PORT", "12345")

	if err != nil {
		log.Error().Msg("Could not set environment variables for connection to MongoDB container from local machine. Integration tests will fail.")
	}
}

func TestAPIFlows(t *testing.T) {
	// First, request an upload URL with valid expiration
	rec := httptest.NewRecorder()
	upldReq, err := http.NewRequest(constants.RequestMethodPost, fmt.Sprintf("%s?timeout=311", constants.AssetsURL), nil)
	require.NoError(t, err)

	server.RequestUploadURL(env).ServeHTTP(rec, upldReq)

	assert.Equal(t, http.StatusOK, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)
	validateContentType(t, rec)

	uploadResp := &responses.UploadURLResponse{}
	require.NoError(t, json.Unmarshal(buf, uploadResp))

	// Now make a PUT request to upload an asset to AWS
	upldFile, err := getUploadBytes()
	require.NoError(t, err)

	awsUploadReq, err := http.NewRequest(constants.RequestMethodPut, uploadResp.URL, bytes.NewBuffer(upldFile))
	require.NoError(t, err)

	awsResp, err := httpClient.Do(awsUploadReq)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, awsResp.StatusCode)

	// Set the status of the asset to uploaded
	statusRec := httptest.NewRecorder()
	statusReq, err := http.NewRequest(constants.RequestMethodPut,
		fmt.Sprintf("%s?id=%s&status=%s", constants.StatusURL, uploadResp.ID, constants.AssetStatusUploaded), nil)
	require.NoError(t, err)

	server.SetUploadStatus(env).ServeHTTP(statusRec, statusReq)

	assert.Equal(t, http.StatusOK, statusRec.Code)
	buf, err = ioutil.ReadAll(statusRec.Body)
	require.NoError(t, err)
	validateContentType(t, statusRec)

	statusResp := &responses.StatusUpdateResponse{}
	require.NoError(t, json.Unmarshal(buf, statusResp))
	assert.Equal(t, uploadResp.ID, statusResp.ID)
	assert.Equal(t, constants.AssetStatusUploaded, statusResp.Status)

	// Make a request for a download URL for the uploaded asset
	downloadRec := httptest.NewRecorder()
	downloadReq, err := http.NewRequest(constants.RequestMethodGet,
		fmt.Sprintf("%s?id=%s", constants.AssetsURL, uploadResp.ID), nil)
	require.NoError(t, err)

	server.GetDownloadURL(env).ServeHTTP(downloadRec, downloadReq)

	assert.Equal(t, http.StatusOK, downloadRec.Code)
	buf, err = ioutil.ReadAll(downloadRec.Body)
	require.NoError(t, err)
	validateContentType(t, downloadRec)

	downloadResp := &responses.DownloadURLResponse{}
	require.NoError(t, json.Unmarshal(buf, downloadResp))
	assert.Equal(t, uploadResp.ID, downloadResp.ID)
	assert.Equal(t, statusResp.ID, downloadResp.ID)

	// Actually download the file from AWS
	awsDownloadResp, err := http.Get(downloadResp.DownloadURL)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, awsDownloadResp.StatusCode)

	expUpldResp := getUploadURLResponse(env, t, "2")
	time.Sleep(2 * time.Second)

	// Now make a PUT request to upload an asset to AWS after the expiry time has passed.
	upldFile, err = getUploadBytes()
	require.NoError(t, err)

	awsExpUploadReq, err := http.NewRequest(constants.RequestMethodPut, expUpldResp.URL, bytes.NewBuffer(upldFile))
	require.NoError(t, err)

	awsExpResp, err := httpClient.Do(awsExpUploadReq)
	require.NoError(t, err)

	assert.Equal(t, http.StatusForbidden, awsExpResp.StatusCode)
}

func getUploadURLResponse(env *config.Env, t *testing.T, timeout string) *responses.UploadURLResponse {
	// Check that upload URL actually expires after the specified timeout
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPost,
		fmt.Sprintf("%s?timeout=%s", constants.AssetsURL, timeout), nil)
	require.NoError(t, err)

	server.RequestUploadURL(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)
	validateContentType(t, rec)

	resp := &responses.UploadURLResponse{}
	require.NoError(t, json.Unmarshal(buf, resp))
	return resp
}

func validateContentType(t *testing.T, rec *httptest.ResponseRecorder) {
	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)
}

func getUploadBytes() ([]byte, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fmt.Sprintf("%s/%s", workDir, constants.UploadImagePath))
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(f)
	return buf, err
}
