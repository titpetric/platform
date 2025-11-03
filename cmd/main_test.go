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
	platform.Register(&platform.UnimplementedModule{})
	platform.Use(platform.TestMiddleware())

	ctx, cancel := context.WithCancel(t.Context())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		Main(ctx, platform.NewTestOptions())
		wg.Done()
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	wg.Wait()

	require.True(t, true)
}
