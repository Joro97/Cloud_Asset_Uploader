package data

import (
	"context"
	"os"
	"testing"

	"CloudAssetUploader/constants"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	db             *DB
	invalidConnStr = "notMongodb://username:password@host:port"
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	setLocalhostEnvVars()

	connStr, _ := BuildConnectionStringForDB()
	database, _ := NewDB(connStr)
	db = database
}

func setLocalhostEnvVars() {
	err := os.Setenv("MONGO_USERNAME", "mongoadmin")
	err = os.Setenv("MONGO_PASSWORD", "bigSecret")
	err = os.Setenv("MONGO_CONTAINER_NAME", "localhost")
	err = os.Setenv("MONGO_PORT", "12345")

	if err != nil {
		log.Error().Msg("Could not set environment variables for connection to MongoDB container from local machine. Integration tests will fail.")
	}
}

func TestIntegrationBuildConnectionStringForDBReturnsErrorsForMissingEnvVars(t *testing.T) {
	os.Unsetenv("MONGO_USERNAME")
	_, err := BuildConnectionStringForDB()
	assert.Error(t, err)
	os.Setenv("MONGO_USERNAME", "mongoadmin")

	os.Unsetenv("MONGO_PASSWORD")
	_, err = BuildConnectionStringForDB()
	assert.Error(t, err)
	os.Setenv("MONGO_PASSWORD", "bigSecret")

	os.Unsetenv("MONGO_CONTAINER_NAME")
	_, err = BuildConnectionStringForDB()
	assert.Error(t, err)
	os.Setenv("MONGO_CONTAINER_NAME", "localhost")

	os.Unsetenv("MONGO_PORT")
	_, err = BuildConnectionStringForDB()
	assert.Error(t, err)
	os.Setenv("MONGO_PORT", "12345")
}

func TestIntegrationNewDBReturnsAnErrorForInvalidConnectionString(t *testing.T) {
	_, err := NewDB(invalidConnStr)
	assert.Error(t, err)
}

func TestIntegrationAddNewAssetAddsAssetWithProperNameAndUrl(t *testing.T) {
	_, err := db.AddNewAsset(constants.MockAssetName, constants.MockURL)
	require.NoError(t, err)
}

func tearDown() {
	defer db.Client.Disconnect(context.Background())

	err := os.Unsetenv("MONGO_USERNAME")
	err = os.Unsetenv("MONGO_PASSWORD")
	err = os.Unsetenv("MONGO_CONTAINER_NAME")
	err = os.Unsetenv("MONGO_PORT")
	if err != nil {
		log.Error().Msg("Could not unset all MongoDB env vars properly. Please unset manually.")
	}

/*	assetInfoCollection := db.Client.Database(constants.AssetUploaderDatabaseName).Collection(constants.AssetUploaderCollectionName)
	err = assetInfoCollection.Drop(context.Background())
	if err != nil {
		log.Error().Msg("Could not drop database after tests.")
	}*/
}
