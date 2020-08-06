package constants

// HTTP related.
const (
	RequestMethodGet  = "GET"
	RequestMethodPost = "POST"
	RequestMethodPut = "PUT"
	HeaderContentType = "Content-Type"
	ApplicationJSON   = "application/json"
)

// Common errors.
const (
	InternalServerErrorMessage = "An internal error has occurred. Please retry your request later."
)

// Asset upload status related constants.
const (
	AssetStatusCreated = "STATUS_CREATED"
	AssetStatusUploaded = "STATUS_UPLOADED"
)