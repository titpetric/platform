package platform

import (
	"context"
	"net/http"
	"testing"

	chi "github.com/go-chi/chi/v5"
	"github.com/titpetric/oida"

	"github.com/titpetric/platform/internal/assert"
)

// TestTelemetryModule covers the recorder and dashboard a default platform registers.
func TestTelemetryModule(t *testing.T) {
	module, err := NewTelemetryModule(oida.Options{})
	assert.NoError(t, err)
	assert.NotNil(t, module)

	var typed *TelemetryModule = module
	assert.NotNil(t, typed.Tracer())
	assert.NotNil(t, typed.Options())

	router := chi.NewRouter()
	assert.NoError(t, typed.Mount(context.Background(), router))

	var middleware Middleware = typed.Middleware
	handler := middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	assert.NotNil(t, handler)
}
