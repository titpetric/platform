package platform

// ErrorResponse is the JSON structure for error responses.
type ErrorResponse struct {
	Error ErrorResponseBody `json:"error"`
}

// ErrorResponseBody is an inner type for ErrorResponse.Error.
type ErrorResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
