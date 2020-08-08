package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"CloudAssetUploader/config"
	"CloudAssetUploader/constants"
	"CloudAssetUploader/responses"
	"CloudAssetUploader/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestUploadURLWithOkRequestShouldReturnProperResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPost, fmt.Sprintf("%s?%s", constants.AssetsURL, test.MockAssetName), nil)
	require.NoError(t, err)

	db := &test.MockDb{}
	upd := &test.MockUploader{}
	env := &config.Env{AssetUploader: upd, Store: db}

	RequestUploadURL(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)

	uploadResp := &responses.UploadUrlResponse{}
	require.NoError(t, json.Unmarshal(buf, uploadResp))
	assert.Equal(t, test.MockID, uploadResp.Id)
	assert.Equal(t, test.MockURL, uploadResp.Url)
}

func TestSetUploadStatusWithOkRequestShouldReturnProperResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPut,
		fmt.Sprintf("%s?id=%s&status=%s", constants.StatusURL, test.MockID, constants.AssetStatusUploaded), nil)
	require.NoError(t, err)

	db := &test.MockDb{}
	upd := &test.MockUploader{}
	env := &config.Env{AssetUploader: upd, Store: db}

	SetUploadStatus(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)

	resp := &responses.StatusUpdateResponse{}
	require.NoError(t, json.Unmarshal(buf, resp))
	assert.Equal(t, test.MockID, resp.Id)
	assert.Equal(t, constants.AssetStatusUploaded, resp.Status)
}

func TestGetDownloadURLWithOkRequestShouldReturnProperResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodGet,
		fmt.Sprintf("%s?id=%s", constants.AssetsURL, test.MockID), nil)
	require.NoError(t, err)

	db := &test.MockDb{}
	upd := &test.MockUploader{}
	env := &config.Env{AssetUploader: upd, Store: db}

	GetDownloadURL(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)

	resp := &responses.DownloadUrlResponse{}
	assert.NoError(t, json.Unmarshal(buf, resp))
	assert.Equal(t, test.MockID, resp.Id)
	assert.Equal(t, test.MockURL, resp.DownloadUrl)
}