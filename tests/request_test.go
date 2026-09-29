package platform_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/internal/assert"
)

func TestParam(t *testing.T) {
	svc := platform.New(platform.NewTestOptions())

	svc.Register(&platform.UnimplementedModule{
		NameFn: func() string { return "TestParam" },
		MountFn: func(_ context.Context, mux platform.Router) error {
			mux.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
				id := platform.Param(r, "id")
				foo := platform.Param(r, "foo")
				_, _ = w.Write([]byte("user: " + id + " foo: " + foo))
			})
			return nil
		},
	})

	assert.NoError(t, svc.Start(t.Context()))

	resp, err := http.Get(svc.URL() + "/user/test-id?foo=bar")
	assert.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, resp.Body.Close()) })

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	want := "user: test-id foo: bar"

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, want, string(body))

	svc.Stop()
}
