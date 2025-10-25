package platform_test

import (
	"testing"

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

		plugins, mws := svc.Registry.Stats()

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
