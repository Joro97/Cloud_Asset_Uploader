package main

import (
	"context"
	"net/http"
	"os"

	"CloudAssetUploader/config"
	"CloudAssetUploader/constants"
	"CloudAssetUploader/data"
	"CloudAssetUploader/server"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func main() {
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = constants.DefaultRegion
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		log.Fatal().Msgf("Could not connect to aws: %s", err)
	}

	connStr, err := data.BuildConnectionStringForDB()
	if err != nil {
		log.Fatal().Msgf("Could not build connection for MongoDB: %s", err)
	}

	db, err := data.NewDB(connStr)
	if err != nil {
		log.Fatal().Msgf("Could not connect to MongoDB: %s", err)
	}
	defer db.Client.Disconnect(context.Background())

	env := config.NewEnv(sess, db)

	r := chi.NewRouter()

	r.Post("/assets", server.RequestUploadURL(env))
	r.Put("/status", server.SetUploadStatus(env))
	r.Get("/assets", server.GetDownloadURL(env))

	srvPort := os.Getenv("SERVER_PORT")
	if srvPort == "" {
		srvPort = constants.DefaultServerPort
	}
	log.Info().Msgf("Starting server on port %s", srvPort)
	err = http.ListenAndServe(srvPort, r)
	if err != nil {
		log.Error().Msgf("Could not start server: %s", err)
	}
}
