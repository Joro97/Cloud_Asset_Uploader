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

//
type AssetInfo struct {
	Name         string `bson:"awsName,omitempty"`
	Url          string `bson:"url,omitempty"`
	UploadStatus string `bson:"uploadStatus,omitempty"`
}

//
type ErrorNoAssetFound struct {
	Id string
}

//
func (err *ErrorNoAssetFound) Error() string {
	return fmt.Sprintf("%s:%s", constants.AssetNotFoundMessage, err.Id)
}

type ErrorInvalidStatus struct {
	Status string
}

func (err *ErrorInvalidStatus) Error() string {
	return fmt.Sprintf("Bad status:%s %s: %s, %s",
		err.Status, constants.InvalidStatusMessage, constants.AssetStatusCreated, constants.AssetStatusUploaded)
}

type ErrorDownloadForNotUploadedAsset struct {
}

func (err *ErrorDownloadForNotUploadedAsset) Error() string {
	return constants.UnsetStatusMessage
}

//
func (db *DB) AddNewAsset(assetName, url string) (string, error) {
	assetInfo := &AssetInfo{
		Name:         assetName,
		Url:          url,
		UploadStatus: constants.AssetStatusCreated,
	}

	assetInfoCollection := db.Client.Database(constants.AssetUploaderDatabaseName).Collection(constants.AssetUploaderCollectionName)
	_, err := assetInfoCollection.InsertOne(context.Background(), assetInfo)
	if err != nil {
		return "", err
	}

	return assetName, nil
}

//
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
		return nil, &ErrorNoAssetFound{Id: awsName}
	}
	if err != nil {
		return nil, err
	}

	return asset, err
}

//
func (db *DB) GetAsset(assetId string) (*AssetInfo, error) {
	assetInfoCollection := db.Client.Database(constants.AssetUploaderDatabaseName).Collection(constants.AssetUploaderCollectionName)

	asset := &AssetInfo{}
	err := assetInfoCollection.FindOne(
		context.Background(),
		bson.M{
			"awsName": assetId,
		}).Decode(asset)

	if err == mongo.ErrNoDocuments {
		return nil, &ErrorNoAssetFound{Id: assetId}
	}
	if err != nil {
		log.Error().Msgf("Could not fetch asset with id: %s. err: %s", assetId, err)
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
