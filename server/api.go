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

		url, id, err := env.AssetUploader.GetSignedUploadURL(assetName)
		if err != nil {
			responses.WriteInternalServerErrorResponse(wr, constants.InternalServerErrorMessage)
			return
		}

		resp := struct {
			Url string `json:"upload_url,omitempty"`
			Id 	string  `json:"asset_id,omitempty"`
		}{
			Url: url,
			Id: id,
		}
		responses.WriteOkResponse(wr, resp)
	}
}
