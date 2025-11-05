package httpcontext_test

import (
	"net/http/httptest"
	"testing"

	"github.com/titpetric/platform/pkg/assert"
	"github.com/titpetric/platform/pkg/httpcontext"
)

func TestContextValue_GetSet(t *testing.T) {
	type TestContext struct {
		Message string
	}

	a := assert.New(t)

	// black-box: key type is unexported from internal package; define local key type
	type testContextKey struct{}

	// create manager for *TestContext values (pointer type)
	manager := httpcontext.NewValue[*TestContext](testContextKey{})

	// create request
	req := httptest.NewRequest("GET", "/", nil)

	// GET before Set should return nil (zero value for pointer)
	got := manager.Get(req)
	a.Empty(got, "expected nil when value not set in request context")

	// Set a pointer value
	want := &TestContext{Message: "hello"}
	manager.Set(req, want)

	// GET after Set should return the same pointer value
	got2 := manager.Get(req)
	a.Equal(want, got2, "expected to get the pointer value that was set")
}
