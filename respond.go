package platform

import (
	"encoding/json"
	"net/http"

	"github.com/titpetric/platform/pkg/telemetry"
)

// JSON writes any payload as JSON. If the payload is nil, the write is omitted.
// If an error occurs in encoding, a telemetry error is logged.
func JSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		telemetry.CaptureError(r.Context(), json.NewEncoder(w).Encode(data))
	}
}

// ErrorResponse is our JSON choice of an error response.
type ErrorResponse struct {
	Error ErrorResponseBody `json:"error"`
}

// ErrorResponseBody is an inner type for ErrorResponse.Error.
type ErrorResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error writes an error payload as JSON.
func Error(w http.ResponseWriter, r *http.Request, status int, data error) {
	var response ErrorResponse
	response.Error.Code = status
	if data != nil {
		response.Error.Message = data.Error()
	}

	JSON(w, r, status, response)
}
