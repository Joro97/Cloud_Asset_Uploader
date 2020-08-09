package data

import (
	"context"
	"fmt"
	"strings"

	"CloudAssetUploader/constants"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AssetInfo is a struct that represents how an Asset will be stored in MongoDB.
type AssetInfo struct {
	Name         string `bson:"awsName,omitempty"`
	URL          string `bson:"url,omitempty"`
	UploadStatus string `bson:"uploadStatus,omitempty"`
}

// ErrorNoAssetFound is an error that indicates that an asset is missing from the database.
type ErrorNoAssetFound struct {
	ID string
}

// Error returns a description for the missing asset and the ID that was attempted to be fetched.
func (err *ErrorNoAssetFound) Error() string {
	return fmt.Sprintf("%s:%s", constants.AssetNotFoundMessage, err.ID)
}

// ErrorInvalidStatus indicates that an invalid status was attempted to be set.
type ErrorInvalidStatus struct {
	Status string
}

func (err *ErrorInvalidStatus) Error() string {
	return fmt.Sprintf("Bad status:%s %s: %s, %s",
		err.Status, constants.InvalidStatusMessage, constants.AssetStatusCreated, constants.AssetStatusUploaded)
}

// ErrorDownloadForNotUploadedAsset indicates that an asset, which status was not set to UPLOADED was attempted to be downloaded.
type ErrorDownloadForNotUploadedAsset struct {
}

func (err *ErrorDownloadForNotUploadedAsset) Error() string {
	return constants.UnsetStatusMessage
}

// AddNewAsset creates a new AssetInfo entry and stores it in MongoDB.
func (db *DB) AddNewAsset(assetName, url string) (string, error) {
	assetInfo := &AssetInfo{
		Name:         assetName,
		URL:          url,
		UploadStatus: constants.AssetStatusCreated,
	}

	assetInfoCollection := db.Client.Database(constants.AssetUploaderDatabaseName).Collection(constants.AssetUploaderCollectionName)
	_, err := assetInfoCollection.InsertOne(context.Background(), assetInfo)
	if err != nil {
		return "", err
	}

	return assetName, nil
}

// SetAssetStatus sets the status for the given asset or returns an error if an invalid status was passed.
func (db *DB) SetAssetStatus(awsName, status string) (*AssetInfo, error) {
	err := validateStatus(status)
	if err != nil {
		return nil, &ErrorInvalidStatus{Status: status}
	}

	assetInfoCollection := db.Client.Database(constants.AssetUploaderDatabaseName).Collection(constants.AssetUploaderCollectionName)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	asset := &AssetInfo{}
	err = assetInfoCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{
			"awsName": awsName,
		},
		bson.D{
			{"$set", bson.D{{"uploadStatus", status}}},
		}, opts,
	).Decode(asset)

	if err == mongo.ErrNoDocuments {
		return nil, &ErrorNoAssetFound{ID: awsName}
	}
	if err != nil {
		return nil, err
	}

	return asset, err
}

// GetAsset returns the given AssetInfo object specified by an unique id.
func (db *DB) GetAsset(assetID string) (*AssetInfo, error) {
	assetInfoCollection := db.Client.Database(constants.AssetUploaderDatabaseName).Collection(constants.AssetUploaderCollectionName)

	asset := &AssetInfo{}
	err := assetInfoCollection.FindOne(
		context.Background(),
		bson.M{
			"awsName": assetID,
		}).Decode(asset)

	if err == mongo.ErrNoDocuments {
		return nil, &ErrorNoAssetFound{ID: assetID}
	}
	if err != nil {
		log.Error().Msgf("Could not fetch asset with id: %s. err: %s", assetID, err)
		return nil, err
	}

	return asset, nil
}

func validateStatus(status string) error {
	if lower := strings.TrimSpace(strings.ToLower(status)); lower == "uploaded" {
		return nil
	}
	return &ErrorInvalidStatus{Status: status}
}
