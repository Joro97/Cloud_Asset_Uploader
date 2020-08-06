package server

import (
	"net/http"

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
			Url: url,
			Id: id,
		}
		responses.WriteOkResponse(wr, resp)
	}
}

func GetDownloadURL(env *config.Env) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)


	}
}

func SetUploadStatus(env *config.Env) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)

		/*assetId := r.URL.Query().Get("id")
		status := r.URL.Query().Get("status")*/

	}
}
