package platform_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform"
)

var platformTestOptions = &platform.Options{
	ServerAddr: "127.0.0.1:0",
	Quiet:      true,
}

func NewTestPlatform(tb testing.TB) *platform.Platform {
	svc, err := platform.StartPlatform(tb.Context(), platformTestOptions)

	require.NoError(tb, err)
	require.NotNil(tb, svc)

	tb.Cleanup(svc.Close)
	return svc
}

func TestPlatform(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		svc := NewTestPlatform(t)

		plugins, mws := svc.Stats()

		require.Equal(t, 0, plugins)
		require.Equal(t, 0, mws)
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
			svc, err := platform.StartPlatform(t.Context(), platformTestOptions)

			require.NoError(t, err)
			require.NotNil(t, svc)

			svc.Close()

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
			svc, err := platform.StartPlatform(b.Context(), platformTestOptions)

			require.NoError(b, err)
			require.NotNil(b, svc)

			svc.Close()
		}
	})
}
