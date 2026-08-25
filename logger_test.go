package platform_test

import (
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/require"
)

// recordLogger implements platform.Logger and keeps the messages it
// received. The platform logs from the SIGTERM goroutine as well, so
// access is guarded.
type recordLogger struct {
	mu       sync.Mutex
	messages []string
}

func (l *recordLogger) Info(msg string, _ ...any)  { l.record(msg) }
func (l *recordLogger) Error(msg string, _ ...any) { l.record(msg) }

func (l *recordLogger) record(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.messages = append(l.messages, msg)
}

func (l *recordLogger) has(msg string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	return slices.Contains(l.messages, msg)
}

// verboseTestOptions produces test options that log. NewTestOptions is
// quiet, which is what the rest of the suite wants, but a test that asserts
// on output needs the platform to produce some.
func verboseTestOptions() *platform.Options {
	options := platform.NewTestOptions()
	options.Quiet = false
	return options
}

func TestLogger(t *testing.T) {
	t.Run("New sets a logger", func(t *testing.T) {
		require.NotNil(t, platform.New(platform.NewTestOptions()).Logger)
		require.NotNil(t, platform.New(verboseTestOptions()).Logger)
	})

	t.Run("injected logger receives platform output", func(t *testing.T) {
		log := &recordLogger{}

		svc := platform.New(verboseTestOptions())
		t.Cleanup(svc.Stop)

		svc.Logger = log
		svc.Register(platform.NewUnimplementedModule("TestLogger"))

		require.NoError(t, svc.Start(t.Context()))

		require.True(t, log.has("modules started"), "registry should log through the platform logger")
		require.True(t, log.has("server listening"), "listener should log through the platform logger")
		require.True(t, log.has("routes registered"), "route printing should log through the platform logger")
	})

	t.Run("nil logger discards output", func(t *testing.T) {
		svc := platform.New(verboseTestOptions())
		t.Cleanup(svc.Stop)

		svc.Logger = nil

		require.NoError(t, svc.Start(t.Context()))
	})

	// Stop releases the signal handler, and that cancels the same context
	// a delivered signal does. Only a signal is worth a line.
	t.Run("stop is not a caught signal", func(t *testing.T) {
		log := &recordLogger{}

		svc := platform.New(verboseTestOptions())
		svc.Logger = log

		require.NoError(t, svc.Start(t.Context()))

		svc.Stop()
		svc.Wait()

		// The signal goroutine wakes on its own schedule.
		time.Sleep(50 * time.Millisecond)

		require.False(t, log.has("caught sigterm, stopping server"))
	})
}
