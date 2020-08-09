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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestUploadURLWithOkRequestShouldReturnProperResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPost, fmt.Sprintf("%s", constants.AssetsURL), nil)
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

	uploadResp := &responses.UploadURLResponse{}
	require.NoError(t, json.Unmarshal(buf, uploadResp))
	assert.Equal(t, constants.MockID, uploadResp.ID)
	assert.Equal(t, constants.MockURL, uploadResp.URL)
}

func TestRequestUploadURLWithInternalErrorOnGetSignedUploadURLShouldThrowProperError(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPost, fmt.Sprintf("%s", constants.AssetsURL), nil)
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
	req, err := http.NewRequest(constants.RequestMethodPost, fmt.Sprintf("%s", constants.AssetsURL), nil)
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
	req, err := http.NewRequest(constants.RequestMethodPost, fmt.Sprintf("%s", constants.AssetsURL), nil)
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
		fmt.Sprintf("%s?id=%s&status=%s", constants.StatusURL, constants.MockID, constants.AssetStatusUploaded), nil)
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
	assert.Equal(t, constants.MockID, resp.ID)
	assert.Equal(t, constants.AssetStatusUploaded, resp.Status)
}

func TestSetUploadStatusWithNonExistentAssetIDShouldReturnProperError(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodPut,
		fmt.Sprintf("%s?id=%s", constants.StatusURL, constants.MockNonExistentID), nil)
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
		fmt.Sprintf("%s?id=%s", constants.StatusURL, constants.MockNonExistentID), nil)
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
		fmt.Sprintf("%s?id=%s", constants.StatusURL, constants.MockNonExistentID), nil)
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
		fmt.Sprintf("%s?id=%s", constants.AssetsURL, constants.MockID), nil)
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

	resp := &responses.DownloadURLResponse{}
	assert.NoError(t, json.Unmarshal(buf, resp))
	assert.Equal(t, constants.MockID, resp.ID)
	assert.Equal(t, constants.MockURL, resp.DownloadURL)
}

func TestGetDownloadURLWithOkTimeoutShouldReturnProperResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodGet,
		fmt.Sprintf("%s?id=%s&timeout=11", constants.AssetsURL, constants.MockID), nil)
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

	resp := &responses.DownloadURLResponse{}
	assert.NoError(t, json.Unmarshal(buf, resp))
	assert.Equal(t, constants.MockID, resp.ID)
	assert.Equal(t, constants.MockURL, resp.DownloadURL)
}

func TestGetDownloadURLWithTooLargeTimeoutAndMissingAssetIDShouldReturnProperResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodGet,
		fmt.Sprintf("%s?id=%s&timeout=111013", constants.AssetsURL, constants.MockID), nil)
	require.NoError(t, err)

	db := &test.MockDb{Err: &data.ErrorNoAssetFound{}}
	upd := &test.MockUploader{}
	env := &config.Env{AssetUploader: upd, Store: db}

	GetDownloadURL(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	_, err = ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)
}

func TestGetDownloadURLWithInternalErrorInGetAssetShouldReturnProperResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodGet,
		fmt.Sprintf("%s?id=%s&timeout=111013", constants.AssetsURL, constants.MockID), nil)
	require.NoError(t, err)

	db := &test.MockDb{Err: errors.New("")}
	upd := &test.MockUploader{}
	env := &config.Env{AssetUploader: upd, Store: db}

	GetDownloadURL(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)

	respMsg := ""
	require.NoError(t, json.Unmarshal(buf, &respMsg))
	assert.Equal(t, constants.InternalServerErrorMessage, respMsg)
}

func TestGetDownloadURLWithAssetThatHasStatusCreatedShouldReturnProperResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodGet,
		fmt.Sprintf("%s?id=%s&timeout=11", constants.AssetsURL, constants.MockID), nil)
	require.NoError(t, err)

	db := &test.MockDb{ShouldStatusBeCreated: true}
	upd := &test.MockUploader{}
	env := &config.Env{AssetUploader: upd, Store: db}

	GetDownloadURL(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)

	respMsg := ""
	require.NoError(t, json.Unmarshal(buf, &respMsg))
	assert.Equal(t, constants.UnsetStatusMessage, respMsg)
}

func TestGetDownloadURLWithInternalErrorInGetSignedDownloadURLShouldReturnProperResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(constants.RequestMethodGet,
		fmt.Sprintf("%s?id=%s&timeout=111013", constants.AssetsURL, constants.MockID), nil)
	require.NoError(t, err)

	db := &test.MockDb{}
	upd := &test.MockUploader{Err: errors.New("")}
	env := &config.Env{AssetUploader: upd, Store: db}

	GetDownloadURL(env).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	contentType := rec.Header().Get(constants.HeaderContentType)
	assert.Equal(t, constants.ApplicationJSON, contentType)

	respMsg := ""
	require.NoError(t, json.Unmarshal(buf, &respMsg))
	assert.Equal(t, constants.InternalServerErrorMessage, respMsg)
}
