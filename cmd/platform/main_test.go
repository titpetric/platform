package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/titpetric/platform/internal/require"

	"github.com/titpetric/platform"
)

func TestStart(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	// global state override to run start()
	options = platform.NewTestOptions()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		start(ctx)
		wg.Done()
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	wg.Wait()

	require.True(t, true)
}
