package platform_test

import (
	"io"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/require"
)

func NewTestPlatform(tb testing.TB) *platform.Platform {
	svc, err := platform.Start(tb.Context(), platform.NewTestOptions())

	require.NoError(tb, err)
	require.NotNil(tb, svc)

	tb.Cleanup(svc.Stop)
	return svc
}

func TestPlatform(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		svc := platform.New(platform.NewTestOptions())
		defer svc.Stop()

		svc.Register(&platform.UnimplementedModule{
			NameFn: func() string {
				return "TestPlatform"
			},
			MountFn: func(r platform.Router) error {
				r.Get("/404", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("You found a valid route"))
				})
				return nil
			},
		})
		svc.Use(platform.TestMiddleware())

		t.Run("find", func(t *testing.T) {
			var mod *platform.UnimplementedModule

			require.True(t, svc.Find(&mod))
			require.Equal(t, "TestPlatform", mod.Name())
		})

		plugins, mws := svc.Stats()
		require.Equal(t, 1, plugins)
		require.Equal(t, 1, mws)

		require.NoError(t, svc.Start(t.Context()))

		resp, err := http.Get(svc.URL() + "/404")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, string(body), "You found a valid route")
	})

	t.Run("multi", func(t *testing.T) {
		NewTestPlatform(t)
		NewTestPlatform(t)
		NewTestPlatform(t)
		NewTestPlatform(t)
	})
}

// This test case is an eyeball test. It starts and stops platforms in a loop and prints
// how many goroutines are alive. It doesn't make any assertion on the goroutine count,
// as tests are run in parallel. The eyeball test confirms stable goroutine levels.
func TestPlatform_goroutine_leaks(t *testing.T) {
	if !testing.Verbose() {
		t.Skip()
		return
	}

	t.Run("stress", func(t *testing.T) {
		t.Logf("start: %d", runtime.NumGoroutine())
		for i := 0; i < 30; i++ {
			svc, err := platform.Start(t.Context(), platform.NewTestOptions())

			require.NoError(t, err)
			require.NotNil(t, svc)

			svc.Stop()

			t.Logf("run[%d]: %d", i, runtime.NumGoroutine())
		}

		time.Sleep(time.Second)
		runtime.GC()

		t.Logf("final: %d", runtime.NumGoroutine())
		// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	})
}

func BenchmarkPlatform(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			svc, err := platform.Start(b.Context(), platform.NewTestOptions())

			require.NoError(b, err)
			require.NotNil(b, svc)

			svc.Stop()
		}
	})
}
