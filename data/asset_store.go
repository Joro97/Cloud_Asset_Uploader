package data

import (
	"context"

	"CloudAssetUploader/constants"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AssetInfo struct {
	Id           string `bson:"id,omitempty"`
	Name         string `bson:"name,omitempty"`
	Url          string `bson:"url,omitempty"`
	UploadStatus string `bson:"uploadStatus,omitempty"`
}

func (db *DB) AddNewAsset(assetName, url string) (string, error) {
	id := uuid.New().String()
	assetInfo := &AssetInfo{
		Id:           id,
		Name:         assetName,
		Url:          url,
		UploadStatus: constants.AssetStatusCreated,
	}

	assetInfoCollection := db.Client.Database(constants.AssetUploaderDatabaseName).Collection(constants.AssetUploaderCollectionName)
	_, err := assetInfoCollection.InsertOne(context.Background(), assetInfo)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (db *DB) SetAssetStatus(assetId, status string) (*AssetInfo, error) {
	assetInfoCollection := db.Client.Database(constants.AssetUploaderDatabaseName).Collection(constants.AssetUploaderCollectionName)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	asset := &AssetInfo{}
	err := assetInfoCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{
			"id": assetId,
		},
		bson.D{
			{"$set", bson.D{{"uploadStatus", status}}},
		}, opts,
	).Decode(asset)

	return asset, err
}