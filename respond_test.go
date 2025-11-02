package platform

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got), "should decode JSON response")

	require.Equal(t, want, got)
}
