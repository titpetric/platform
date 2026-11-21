package platform

import (
	"net/http"
	"net/http/httptest"
	"testing"

	chi "github.com/go-chi/chi/v5"

	"github.com/titpetric/platform/internal/assert"
)

// TestRequestParams covers the URL and query parameter helpers on a routed request.
func TestRequestParams(t *testing.T) {
	router := chi.NewRouter()
	router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "42", URLParam(r, "id"))
		assert.Equal(t, "search", QueryParam(r, "q"))
		assert.Equal(t, "42", Param(r, "id"))
		assert.Equal(t, "search", Param(r, "q"))
		assert.Equal(t, "", Param(r, "missing"))
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest("GET", "/users/42?q=search", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
}
