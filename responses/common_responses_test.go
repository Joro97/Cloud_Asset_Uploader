package responses

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"CloudAssetUploader/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	MockID  = "c0703c92-9161-4c6a-947a-77519bedaceb"
	MockURL = "aws.signed.url.should.be.here"
)

func TestWriteBadRequestFunctionShouldWriteProperStatusCodeAndMessage(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteBadRequest(rec, constants.InvalidStatusMessage)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)
	respMsg := ""
	require.NoError(t, json.Unmarshal(buf, &respMsg))
	assert.Equal(t, constants.InvalidStatusMessage, respMsg)
}

func TestWriteInternalServerErrorResponseFunctionShouldWriteProperStatusCodeAndMessage(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteInternalServerErrorResponse(rec, constants.InternalServerErrorMessage)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)
	respMsg := ""
	require.NoError(t, json.Unmarshal(buf, &respMsg))
	assert.Equal(t, constants.InternalServerErrorMessage, respMsg)
}

func TestWWriteResourceNotFoundResponseFunctionShouldWriteProperStatusCodeAndMessage(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteResourceNotFoundResponse(rec, constants.InvalidStatusMessage)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)
	respMsg := ""
	require.NoError(t, json.Unmarshal(buf, &respMsg))
	assert.Equal(t, constants.InvalidStatusMessage, respMsg)
}

func TestWriteOkResponseFunctionShouldWriteProperStatusCodeAndMessage(t *testing.T) {
	mockResp := &UploadURLResponse{
		ID:  MockID,
		URL: MockURL,
	}
	rec := httptest.NewRecorder()

	WriteOkResponse(rec, mockResp)
	assert.Equal(t, http.StatusOK, rec.Code)

	buf, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)
	resp := &UploadURLResponse{}
	require.NoError(t, json.Unmarshal(buf, resp))
	assert.Equal(t, mockResp.ID, resp.ID)
	assert.Equal(t, mockResp.URL, resp.URL)
}
