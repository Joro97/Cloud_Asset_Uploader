package server

import (
	"net/http"

	"CloudAssetUploader/config"
	"CloudAssetUploader/constants"
)

func RequestUploadURL(env *config.Env) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)
	}
}
