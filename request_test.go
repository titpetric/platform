package platform_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/require"
)

func TestParam(t *testing.T) {
	svc := platform.New(platform.NewTestOptions())

	svc.Register(&platform.UnimplementedModule{
		NameFn: func() string { return "TestParam" },
		MountFn: func(mux platform.Router) error {
			mux.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
				id := platform.Param(r, "id")
				foo := platform.Param(r, "foo")
				w.Write([]byte("user: " + id + " foo: " + foo))
			})
			return nil
		},
	})

	require.NoError(t, svc.Start(t.Context()))

	resp, err := http.Get(svc.URL() + "/user/test-id?foo=bar")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	want := "user: test-id foo: bar"

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, want, string(body))

	svc.Stop()
}
