package platform

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

// TestOptionsFromRequest covers the options a request carries.
func TestOptionsFromRequest(t *testing.T) {
	var options *Options = NewTestOptions()

	req := httptest.NewRequest("GET", "/", nil)
	optionsContext.Set(req, options)

	assert.Equal(t, options, OptionsFromRequest(req))
}

// TestOptionsFromContext covers the options a context carries.
func TestOptionsFromContext(t *testing.T) {
	var options *Options = NewTestOptions()

	ctx := optionsContext.SetContext(context.Background(), options)
	assert.Equal(t, options, OptionsFromContext(ctx))
}

// TestNewOptions covers the environment defaults a platform starts from.
func TestNewOptions(t *testing.T) {
	options := NewOptions()
	assert.NotNil(t, options)
	assert.NotEmpty(t, options.ServerAddr)
	assert.False(t, options.Quiet)
}
