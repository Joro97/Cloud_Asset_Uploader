package main

import (
	"net/http"

	"CloudAssetUploader/config"
	"CloudAssetUploader/constants"
	"CloudAssetUploader/server"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.REGION),
	})
	env := config.NewEnv(sess)

	r := chi.NewRouter()

	r.Post("/assets", server.RequestUploadURL(env))
	r.Put("/status", server.SetUploadStatus(env))
	r.Get("/assets", server.GetDownloadURL(env))

	log.Info().Msgf("Starting server on port %s", ":8090")
	err = http.ListenAndServe(":8090", r)
	if err != nil {
		log.Error().Msgf("Could not start server: %s", err)
	}
}