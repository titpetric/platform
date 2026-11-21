package platform

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

func TestError(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	want := ErrorResponse{
		Error: ErrorResponseBody{
			Code:    418,
			Message: "something went wrong",
		},
	}

	Error(rec, req, want.Error.Code, errors.New(want.Error.Message))

	var got ErrorResponse
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&got), "should decode JSON response")

	assert.Equal(t, want, got)
}

func TestJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	JSON(rec, req, http.StatusCreated, map[string]string{"id": "42"})
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var got map[string]string
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&got), "should decode JSON response")
	assert.Equal(t, "42", got["id"])
}
