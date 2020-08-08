package test

import (
	"CloudAssetUploader/constants"
	"CloudAssetUploader/data"
)

const (
	MockID = "c0703c92-9161-4c6a-947a-77519bedaceb"
	MockNonExistentId = "583195a1-10ee-4608-9cc8-00fb0a32feb0"
	MockURL = "aws.signed.url.should.be.here"
	MockAssetName = "Theseus"
	MockInvalidAssetName = "ThisIsTooLongNameForAnS3AssetAndShouldThrowAnError"
)

//
type MockDb struct {
	Err error
}

//
type MockUploader struct {
	Err error
}

func (m *MockDb) AddNewAsset(assetName, url string) (id string, err error) {
	return MockID, m.Err
}

func (m *MockDb) SetAssetStatus(assetId, status string) (*data.AssetInfo, error) {
	if m.Err == nil {
		return &data.AssetInfo{
			Id:           assetId,
			Name:         MockAssetName,
			Url:          MockURL,
			UploadStatus: status,
		}, nil
	}
	return nil, m.Err
}

func (m *MockDb) GetAsset(assetId string) (*data.AssetInfo, error) {
	return &data.AssetInfo{
		Id:           assetId,
		Name:         MockAssetName,
		Url:          MockURL,
		UploadStatus: constants.AssetStatusUploaded,
	}, m.Err
}

func (m *MockUploader) GetSignedUploadURL(assetName string) (url string, err error) {
	return MockURL, m.Err
}

func (m *MockUploader) GetSignedDownloadURL(assetName string, timeout int) (url string, er error) {
	return MockURL, m.Err
}
