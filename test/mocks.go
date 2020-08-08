package test

import (
	"CloudAssetUploader/constants"
	"CloudAssetUploader/data"
)

const (
	MockID = "c0703c92-9161-4c6a-947a-77519bedaceb"
	MockURL = "aws.signed.url.should.be.here"
	MockAssetName = "Theseus"
)

//
type MockDb struct {

}

//
type MockUploader struct {

}

func (m *MockDb) AddNewAsset(assetName, url string) (id string, err error) {
	return MockID, nil
}

func (m *MockDb) SetAssetStatus(assetId, status string) (*data.AssetInfo, error) {
	return &data.AssetInfo{
		Id:           assetId,
		Name:         MockAssetName,
		Url:          MockURL,
		UploadStatus: status,
	}, nil
}

func (m *MockDb) GetAsset(assetId string) (*data.AssetInfo, error) {
	return &data.AssetInfo{
		Id:           assetId,
		Name:         MockAssetName,
		Url:          MockURL,
		UploadStatus: constants.AssetStatusUploaded,
	}, nil
}

func (m *MockUploader) GetSignedUploadURL(assetName string) (url string, err error) {
	return MockURL, nil
}

func (m *MockUploader) GetSignedDownloadURL(assetName string, timeout int) (url string, er error) {
	return MockURL, nil
}
