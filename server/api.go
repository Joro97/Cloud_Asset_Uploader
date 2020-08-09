package server

import (
	"net/http"

	"CloudAssetUploader/config"
	"CloudAssetUploader/constants"
	"CloudAssetUploader/data"
	"CloudAssetUploader/responses"
)

// RequestUploadURL is a handler for the /assets POST API endpoint. It uses dependency injection to allow for more robust testing.
func RequestUploadURL(env *config.Env) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)

		timeout := r.URL.Query().Get("timeout")

		awsName, url, err := env.AssetUploader.GetSignedUploadURL(timeout)
		if err != nil {
			responses.WriteInternalServerErrorResponse(wr, constants.InternalServerErrorMessage)
			return
		}

		id, err := env.Store.AddNewAsset(awsName, url)
		if err != nil {
			responses.WriteInternalServerErrorResponse(wr, constants.InternalServerErrorMessage)
			return
		}

		resp := &responses.UploadURLResponse{
			ID:  id,
			URL: url,
		}
		responses.WriteOkResponse(wr, resp)
	}
}

// SetUploadStatus is a handler for the /status PUT API endpoint. It uses dependency injection to allow for more robust testing.
func SetUploadStatus(env *config.Env) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)

		assetID := r.URL.Query().Get("id")
		status := r.URL.Query().Get("status")

		asset, err := env.Store.SetAssetStatus(assetID, status)
		if err != nil {
			switch err.(type) {
			case *data.ErrorNoAssetFound:
				responses.WriteResourceNotFoundResponse(wr, err.Error())
			case *data.ErrorInvalidStatus:
				responses.WriteBadRequest(wr, err.Error())
			default:
				responses.WriteInternalServerErrorResponse(wr, constants.InternalServerErrorMessage)
			}
			return
		}

		resp := responses.StatusUpdateResponse{
			ID:     asset.Name,
			Status: asset.UploadStatus,
		}
		responses.WriteOkResponse(wr, resp)
	}
}

// GetDownloadURL is a handler for the /assets GET API endpoint. It uses dependency injection to allow for more robust testing.
func GetDownloadURL(env *config.Env) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)

		id := r.URL.Query().Get("id")
		timeout := r.URL.Query().Get("timeout")

		asset, err := env.Store.GetAsset(id)
		if err != nil {
			switch err.(type) {
			case *data.ErrorNoAssetFound:
				responses.WriteResourceNotFoundResponse(wr, err.Error())
			default:
				responses.WriteInternalServerErrorResponse(wr, constants.InternalServerErrorMessage)
			}
			return
		}
		if asset.UploadStatus == constants.AssetStatusCreated {
			responses.WriteBadRequest(wr, constants.UnsetStatusMessage)
			return
		}

		url, err := env.AssetUploader.GetSignedDownloadURL(asset.Name, timeout)
		if err != nil {
			responses.WriteInternalServerErrorResponse(wr, constants.InternalServerErrorMessage)
			return
		}

		resp := &responses.DownloadURLResponse{
			ID:          id,
			DownloadURL: url,
		}
		responses.WriteOkResponse(wr, resp)
	}
}
