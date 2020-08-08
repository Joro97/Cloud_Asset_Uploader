package server

import (
	"net/http"
	"strconv"

	"CloudAssetUploader/data"

	"CloudAssetUploader/config"
	"CloudAssetUploader/constants"
	"CloudAssetUploader/responses"
)

func RequestUploadURL(env *config.Env) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)

		assetName := r.URL.Query().Get("name")

		url, err := env.AssetUploader.GetSignedUploadURL(assetName)
		if err != nil {
			responses.WriteInternalServerErrorResponse(wr, constants.InternalServerErrorMessage)
			return
		}

		id, err := env.Store.AddNewAsset(assetName, url)
		if err != nil {
			responses.WriteInternalServerErrorResponse(wr, constants.InternalServerErrorMessage)
			return
		}

		resp := &responses.UploadUrlResponse{
			Id:  id,
			Url: url,
		}
		responses.WriteOkResponse(wr, resp)
	}
}

func SetUploadStatus(env *config.Env) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)

		assetId := r.URL.Query().Get("id")
		status := r.URL.Query().Get("status")

		asset, err := env.Store.SetAssetStatus(assetId, status)
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
			Id:     asset.Id,
			Status: asset.UploadStatus,
		}
		responses.WriteOkResponse(wr, resp)
	}
}

func GetDownloadURL(env *config.Env) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)

		id := r.URL.Query().Get("id")
		timeout := r.URL.Query().Get("timeout")
		intTimeout := 60
		if val, err := strconv.Atoi(timeout); err == nil {
			if val < 0 || val > 3600 {
				val = 60
			}
			intTimeout = val
		}

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

		url, err := env.AssetUploader.GetSignedDownloadURL(asset.Name, intTimeout)
		if err != nil {
			responses.WriteInternalServerErrorResponse(wr, constants.InternalServerErrorMessage)
			return
		}

		resp := &responses.DownloadUrlResponse{
			Id:          id,
			DownloadUrl: url,
		}
		responses.WriteOkResponse(wr, resp)
	}
}
