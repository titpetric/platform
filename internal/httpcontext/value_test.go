package httpcontext

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

// TestValue covers the read and write of one value on a request and a context.
func TestValue(t *testing.T) {
	type key struct{}

	value := NewValue[string](key{})
	req := httptest.NewRequest("GET", "/", nil)

	assert.Equal(t, "", value.Get(req))

	value.Set(req, "from-request")
	assert.Equal(t, "from-request", value.Get(req))

	ctx := value.SetContext(context.Background(), "from-context")
	assert.Equal(t, "from-context", value.GetContext(ctx))
}
