package platform

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

// TestMiddlewarePassthrough covers the pass-through middleware the tests share.
func TestMiddlewarePassthrough(t *testing.T) {
	var mw Middleware = TestMiddleware()

	called := false
	handler := mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	assert.True(t, called)
}
