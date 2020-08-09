package test

import (
	"CloudAssetUploader/constants"
	"CloudAssetUploader/data"
)

//
type MockDb struct {
	Err                   error
	ShouldStatusBeCreated bool
}

//
type MockUploader struct {
	Err error
}

func (m *MockDb) AddNewAsset(assetName, url string) (id string, err error) {
	return constants.MockID, m.Err
}

func (m *MockDb) SetAssetStatus(assetId, status string) (*data.AssetInfo, error) {
	if m.Err == nil {
		return &data.AssetInfo{
			Name:         assetId,
			Url:          constants.MockURL,
			UploadStatus: status,
		}, nil
	}
	return nil, m.Err
}

func (m *MockDb) GetAsset(assetId string) (*data.AssetInfo, error) {
	if !m.ShouldStatusBeCreated {
		return &data.AssetInfo{
			Name:         assetId,
			Url:          constants.MockURL,
			UploadStatus: constants.AssetStatusUploaded,
		}, m.Err
	}

	return &data.AssetInfo{
		Name:         assetId,
		Url:          constants.MockURL,
		UploadStatus: constants.AssetStatusCreated,
	}, m.Err
}

func (m *MockUploader) GetSignedUploadURL() (awsName, url string, err error) {
	return constants.MockID, constants.MockURL, m.Err
}

func (m *MockUploader) GetSignedDownloadURL(assetName string, timeout int) (url string, er error) {
	return constants.MockURL, m.Err
}
