package constants

// HTTP related.
const (
	RequestMethodGet  = "GET"
	RequestMethodPost = "POST"
	RequestMethodPut  = "PUT"
	HeaderContentType = "Content-Type"
	ApplicationJSON   = "application/json"
)

// Common error API responses.
const (
	InternalServerErrorMessage = "An internal error has occurred. Please retry your request later."
)

// Common error internal messages.
const (
	AssetNotFoundMessage = "No asset found with the given id"
	InvalidStatusMessage = "Status can can be one of"
	UnsetStatusMessage = "Asset cannot be downloaded if its status is not set to UPLOADED."
)

// Asset upload status related constants.
const (
	AssetStatusCreated  = "CREATED"
	AssetStatusUploaded = "UPLOADED"
)
