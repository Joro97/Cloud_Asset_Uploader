package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"CloudAssetUploader/config"
	"CloudAssetUploader/constants"
	"CloudAssetUploader/data"
	"CloudAssetUploader/responses"
	"CloudAssetUploader/test"
	"CloudAssetUploader/uploader"

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

func TestRequestUploadURLWithInvalidNameShouldThrowProperError(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPost,
		fmt.Sprintf("%s?name=%s", constants.AssetsURL, test.MockInvalidAssetName), nil)
	require.NoError(t, err)

	db := &test.MockDb{}
	upd := &test.MockUploader{Err: &uploader.ErrorInvalidAssetName{Name: test.MockInvalidAssetName}}
	env := &config.Env{AssetUploader: upd, Store: db}

	RequestUploadURL(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)

	respMsg := ""
	require.NoError(t, json.Unmarshal(buf, &respMsg))
}

func TestRequestUploadURLWithInternalErrorOnGetSignedUploadURLShouldThrowProperError(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPost,
		fmt.Sprintf("%s?name=%s", constants.AssetsURL, test.MockAssetName), nil)
	require.NoError(t, err)

	db := &test.MockDb{}
	upd := &test.MockUploader{Err: errors.New("this should be an internal server error")}
	env := &config.Env{AssetUploader: upd, Store: db}

	RequestUploadURL(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)

	respMsg := ""
	require.NoError(t, json.Unmarshal(buf, &respMsg))
	assert.Equal(t, constants.InternalServerErrorMessage, respMsg)
}

func TestRequestUploadURLWithInternalErrorOnAddNewAssetShouldThrowProperError(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPost,
		fmt.Sprintf("%s?name=%s", constants.AssetsURL, test.MockAssetName), nil)
	require.NoError(t, err)

	db := &test.MockDb{Err: errors.New("this should be an internal server error")}
	upd := &test.MockUploader{}
	env := &config.Env{AssetUploader: upd, Store: db}

	RequestUploadURL(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)

	respMsg := ""
	require.NoError(t, json.Unmarshal(buf, &respMsg))
	assert.Equal(t, constants.InternalServerErrorMessage, respMsg)
}

func TestRequestUploadURLWithOnAddNewAssetShouldThrowProperError(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPost,
		fmt.Sprintf("%s?name=%s", constants.AssetsURL, test.MockAssetName), nil)
	require.NoError(t, err)

	db := &test.MockDb{Err: errors.New("this should be an internal server error")}
	upd := &test.MockUploader{}
	env := &config.Env{AssetUploader: upd, Store: db}

	RequestUploadURL(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)

	respMsg := ""
	require.NoError(t, json.Unmarshal(buf, &respMsg))
	assert.Equal(t, constants.InternalServerErrorMessage, respMsg)
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

func TestSetUploadStatusWithNonExistentAssetIDShouldReturnProperError(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPut,
		fmt.Sprintf("%s?id=%s", constants.StatusURL, test.MockNonExistentId), nil)
	require.NoError(t, err)

	db := &test.MockDb{Err: &data.ErrorNoAssetFound{}}
	upd := &test.MockUploader{}
	env := &config.Env{AssetUploader: upd, Store: db}

	SetUploadStatus(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)
}

func TestSetUploadStatusWithInvalidStatusIDShouldReturnProperError(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPut,
		fmt.Sprintf("%s?id=%s", constants.StatusURL, test.MockNonExistentId), nil)
	require.NoError(t, err)

	db := &test.MockDb{Err: &data.ErrorInvalidStatus{}}
	upd := &test.MockUploader{}
	env := &config.Env{AssetUploader: upd, Store: db}

	SetUploadStatus(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)
}

func TestSetUploadStatusWithErrorInSetAssetStatusShouldReturnProperError(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPut,
		fmt.Sprintf("%s?id=%s", constants.StatusURL, test.MockNonExistentId), nil)
	require.NoError(t, err)

	db := &test.MockDb{Err: errors.New("this should throw and internal server error for the test")}
	upd := &test.MockUploader{}
	env := &config.Env{AssetUploader: upd, Store: db}

	SetUploadStatus(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)

	respMsg := ""
	require.NoError(t, json.Unmarshal(buf, &respMsg))
	require.Equal(t, constants.InternalServerErrorMessage, respMsg)
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