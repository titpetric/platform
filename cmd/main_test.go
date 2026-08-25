package cmd

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/require"
)

func TestStart(t *testing.T) {
	platform.RegisterFunc(func() platform.Module {
		return platform.NewUnimplementedModule("test")
	})
	platform.Use(platform.TestMiddleware())

	ctx, cancel := context.WithCancel(t.Context())

	var wg sync.WaitGroup
	wg.Go(func() {
		Main(ctx, platform.NewTestOptions())
	})

	time.Sleep(50 * time.Millisecond)
	cancel()

	wg.Wait()

	require.True(t, true)
}
