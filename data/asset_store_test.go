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

	connStr, _ := BuildConnectionStringForMongoDB()
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
	_, err := BuildConnectionStringForMongoDB()
	assert.Error(t, err)
	os.Setenv("MONGO_USERNAME", "mongoadmin")

	os.Unsetenv("MONGO_PASSWORD")
	_, err = BuildConnectionStringForMongoDB()
	assert.Error(t, err)
	os.Setenv("MONGO_PASSWORD", "bigSecret")

	os.Unsetenv("MONGO_CONTAINER_NAME")
	_, err = BuildConnectionStringForMongoDB()
	assert.Error(t, err)
	os.Setenv("MONGO_CONTAINER_NAME", "localhost")

	os.Unsetenv("MONGO_PORT")
	_, err = BuildConnectionStringForMongoDB()
	assert.Error(t, err)
	os.Setenv("MONGO_PORT", "12345")
}

func TestIntegrationNewDBReturnsAnErrorForInvalidConnectionString(t *testing.T) {
	_, err := NewDB(invalidConnStr)
	assert.Error(t, err)
}

func TestIntegrationAddNewAssetAddsAssetWithProperNameAndUrl(t *testing.T) {
	_, err := db.AddNewAsset(constants.MockID, constants.MockURL)
	require.NoError(t, err)
}

func TestIntegrationSetAssetStatusShouldReturnProperAssetWithValidIdAndStatus(t *testing.T) {
	_, err := db.AddNewAsset(constants.MockID, constants.MockURL)
	require.NoError(t, err)

	asset, err := db.SetAssetStatus(constants.MockID, constants.AssetStatusUploaded)
	require.NoError(t, err)
	assert.Equal(t, constants.MockID, asset.Name)
	assert.Equal(t, constants.AssetStatusUploaded, asset.UploadStatus)
}

func TestIntegrationSetAssetStatusShouldReturnProperErrorWithInvalidStatus(t *testing.T) {
	_, err := db.SetAssetStatus(constants.MockID, constants.InvalidStatus)
	require.Error(t, err)
}

func TestIntegrationSetAssetStatusShouldReturnProperErrorWithNonExistentAssetId(t *testing.T) {
	_, err := db.SetAssetStatus(constants.MockNonExistentID, constants.AssetStatusUploaded)
	require.Error(t, err)
}

func TestIntegrationGetAssetShouldReturnProperAssetWhenValidId(t *testing.T) {
	_, err := db.AddNewAsset(constants.MockID, constants.MockURL)
	require.NoError(t, err)

	asset, err := db.GetAsset(constants.MockID)
	require.NoError(t, err)
	assert.Equal(t, constants.MockID, asset.Name)
}

func TestIntegrationGetAssetShouldReturnProperErrorWhenNonExistentIdProvided(t *testing.T) {
	_, err := db.AddNewAsset(constants.MockID, constants.MockURL)
	require.NoError(t, err)

	asset, err := db.GetAsset(constants.MockNonExistentID)
	require.Error(t, err)
	assert.Nil(t, asset)
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

	assetInfoCollection := db.Client.Database(constants.AssetUploaderDatabaseName).Collection(constants.AssetUploaderCollectionName)
	err = assetInfoCollection.Drop(context.Background())
	if err != nil {
		log.Error().Msg("Could not drop database after tests.")
	}
}
