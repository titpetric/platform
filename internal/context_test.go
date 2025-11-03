package internal_test

import (
	"net/http/httptest"
	"testing"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/pkg/assert"
)

func TestContextValue_GetSet(t *testing.T) {
	type TestContext struct {
		Message string
	}

	a := assert.New(t)

	// black-box: key type is unexported from internal package; define local key type
	type testContextKey struct{}
	key := testContextKey{}

	// create manager for *TestContext values (pointer type)
	manager := internal.NewContextValue[*TestContext](key)

	// create request
	req := httptest.NewRequest("GET", "/", nil)

	// GET before Set should return nil (zero value for pointer)
	got := manager.Get(req)
	a.Nil(got, "expected nil when value not set in request context")

	// Set a pointer value
	want := &TestContext{Message: "hello"}
	returnedReq := manager.Set(req, want)
	a.NotNil(returnedReq, "Set should return a request")

	// GET after Set should return the same pointer value
	got2 := manager.Get(req) // manager.Set mutates the original request as well
	a.Equal(want, got2, "expected to get the pointer value that was set")

	// also verify using the returned request (defensive)
	got3 := manager.Get(returnedReq)
	a.Equal(want, got3, "expected to get the pointer value from the returned request")
}
