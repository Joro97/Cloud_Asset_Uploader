package constants

// Database related constants.
const (
	AssetUploaderDatabaseName   = "assets"
	AssetUploaderCollectionName = "assetsInfo"
)

// Mock constant values used for testing.
const (
	MockID               = "c0703c92-9161-4c6a-947a-77519bedaceb"
	MockNonExistentID    = "583195a1-10ee-4608-9cc8-00fb0a32feb0"
	MockURL              = "aws.signed.url.should.be.here"
	InvalidStatus        = "Theseus"
	MockInvalidAssetName = "ThisIsTooLongNameForAnS3AssetAndShouldThrowAnError"
	UploadImagePath      = "cassandra.jpeg"
)
