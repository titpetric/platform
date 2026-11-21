package platform

import (
	"net/http/httptest"
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

// TestFromRequest covers the platform bound to a request context.
func TestFromRequest(t *testing.T) {
	var p *Platform = New(NewTestOptions())

	req := httptest.NewRequest("GET", "/", nil)
	platformContext.Set(req, p)

	assert.Equal(t, p, FromRequest(req))
}

// TestPlatformRegistry covers the instance registry a platform keeps for itself.
func TestPlatformRegistry(t *testing.T) {
	var p *Platform = New(NewTestOptions())

	p.Register(NewUnimplementedModule("module"))
	p.Use(TestMiddleware())

	modules, middleware := p.Stats()
	assert.Equal(t, 1, modules)
	assert.Equal(t, 1, middleware)

	var found Module
	assert.True(t, p.Find(&found))
	assert.Equal(t, "module", found.Name())
}
