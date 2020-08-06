package data

import (
	"CloudAssetUploader/constants"
	"context"
	"github.com/google/uuid"
)

type AssetInfo struct {
	Id           string
	Name         string
	Url          string
	UploadStatus string
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
