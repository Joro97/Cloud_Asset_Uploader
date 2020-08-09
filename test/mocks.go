package test

import (
	"CloudAssetUploader/constants"
	"CloudAssetUploader/data"
)

// MockDb is a mock implementation of the Store interface. Used for tests.
type MockDb struct {
	Err                   error
	ShouldStatusBeCreated bool
}

// MockUploader is a mock implementation of the Uploader interface. Used for tests.
type MockUploader struct {
	Err error
}

// AddNewAsset is a mock implementation that returns a sample ID and the error contained in the struct.
// This error can be used to test different code flows.
func (m *MockDb) AddNewAsset(assetName, url string) (id string, err error) {
	return constants.MockID, m.Err
}

// SetAssetStatus is a mock implementation that returns a sample AssetInfo object and the error contained in the struct.
// This error can be used to test different code flows.
func (m *MockDb) SetAssetStatus(assetID, status string) (*data.AssetInfo, error) {
	if m.Err == nil {
		return &data.AssetInfo{
			Name:         assetID,
			URL:          constants.MockURL,
			UploadStatus: status,
		}, nil
	}
	return nil, m.Err
}

// GetAsset is a mock implementation that returns a sample AssetInfo object and the error contained in the struct.
// This error can be used to test different code flows.
func (m *MockDb) GetAsset(assetID string) (*data.AssetInfo, error) {
	if !m.ShouldStatusBeCreated {
		return &data.AssetInfo{
			Name:         assetID,
			URL:          constants.MockURL,
			UploadStatus: constants.AssetStatusUploaded,
		}, m.Err
	}

	return &data.AssetInfo{
		Name:         assetID,
		URL:          constants.MockURL,
		UploadStatus: constants.AssetStatusCreated,
	}, m.Err
}

func (m *MockUploader) SetupBucket() error {
	return nil
}

// GetSignedUploadURL is a mock implementation that returns a sample ID and an URL.
func (m *MockUploader) GetSignedUploadURL(timeout string) (awsName, url string, err error) {
	return constants.MockID, constants.MockURL, m.Err
}

// GetSignedDownloadURL is a mock implementation that returns a sample URL.
func (m *MockUploader) GetSignedDownloadURL(assetName, timeout string) (url string, er error) {
	return constants.MockURL, m.Err
}
