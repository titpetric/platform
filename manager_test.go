package platform_test

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/require"
)

// countingModule records how often it was started and stopped, and serves
// the generation it belongs to.
type countingModule struct {
	platform.UnimplementedModule

	starts atomic.Int64
	stops  atomic.Int64
}

func (m *countingModule) Name() string { return "counting" }

func (m *countingModule) Start(context.Context) error {
	m.starts.Add(1)
	return nil
}

func (m *countingModule) Stop(context.Context) error {
	m.stops.Add(1)
	return nil
}

func (m *countingModule) Mount(_ context.Context, r platform.Router) error {
	generation := m.starts.Load()

	r.Get("/generation", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strconv.FormatInt(generation, 10)))
	})
	return nil
}

// newTestManager returns a started manager serving mod, and the modules it
// registers are imperative, so they go in through Setup.
func newTestManager(tb testing.TB, mod platform.Module) *platform.Manager {
	m := platform.NewManager(platform.NewTestOptions())
	m.Setup = func(p *platform.Platform) error {
		p.Register(mod)
		p.Use(platform.TestMiddleware())
		return nil
	}

	tb.Cleanup(m.Stop)
	require.NoError(tb, m.Start(tb.Context()))

	return m
}

func TestManager(t *testing.T) {
	mod := &countingModule{}
	m := newTestManager(t, mod)

	first := m.Platform()
	require.NotNil(t, first)

	status, body := get(t, m.URL()+"/generation")
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, "1", string(body))

	url := m.URL()

	require.NoError(t, m.Reload(t.Context()))

	t.Run("platform is replaced", func(t *testing.T) {
		require.NotNil(t, m.Platform())
		require.True(t, first != m.Platform(), "reload should replace the platform value")

		// The retired generation is stopped, the new one is not.
		require.Error(t, first.Context().Err())
		require.NoError(t, m.Platform().Context().Err())
	})

	t.Run("modules restart", func(t *testing.T) {
		require.Equal(t, int64(2), mod.starts.Load())
		require.Equal(t, int64(1), mod.stops.Load())
	})

	t.Run("address survives the reload", func(t *testing.T) {
		require.Equal(t, url, m.URL())

		status, body := get(t, m.URL()+"/generation")
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, "2", string(body))
	})

	t.Run("stop ends wait", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Go(m.Wait)

		m.Stop()
		wg.Wait()

		require.Nil(t, m.Platform())
		require.Equal(t, int64(2), mod.stops.Load())
	})
}

// TestManager_sighup covers the signal the manager exists for. It has to
// be running throughout: SIGHUP terminates a process that ignores it.
func TestManager_sighup(t *testing.T) {
	mod := &countingModule{}
	m := newTestManager(t, mod)

	first := m.Platform()

	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGHUP))

	// The signal is delivered asynchronously, and the reload it triggers
	// takes as long as the modules take to stop and start. Platform is nil
	// in between, so waiting for a change is not enough.
	deadline := time.Now().Add(5 * time.Second)

	next := m.Platform()
	for (next == nil || next == first) && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
		next = m.Platform()
	}

	require.NotNil(t, next)
	require.True(t, next != first, "sighup should replace the platform value")
	require.Equal(t, int64(2), mod.starts.Load())

	status, body := get(t, m.URL()+"/generation")
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, "2", string(body))
}

// TestManager_reload_drains covers a request in flight while the reload
// retires the generation serving it. The old handler finishes the response.
func TestManager_reload_drains(t *testing.T) {
	release := make(chan struct{})
	serving := make(chan struct{})

	// Only the request that the reload races has to block; the one after
	// it checks that the new generation serves.
	var once sync.Once

	mod := platform.NewUnimplementedModule("draining")
	mod.MountFn = func(_ context.Context, r platform.Router) error {
		r.Get("/slow", func(w http.ResponseWriter, _ *http.Request) {
			once.Do(func() {
				close(serving)
				<-release
			})

			_, _ = w.Write([]byte("drained"))
		})
		return nil
	}

	m := newTestManager(t, mod)

	var (
		wg     sync.WaitGroup
		status int
		body   []byte
	)

	wg.Go(func() {
		status, body = get(t, m.URL()+"/slow")
	})

	<-serving

	reloaded := make(chan error, 1)
	go func() {
		close(release)
		reloaded <- m.Reload(t.Context())
	}()

	wg.Wait()
	require.NoError(t, <-reloaded)

	require.Equal(t, http.StatusOK, status)
	require.Equal(t, "drained", string(body))

	// The new generation serves on the same address.
	status, _ = get(t, m.URL()+"/slow")
	require.Equal(t, http.StatusOK, status)
}

// TestManager_failed_reload covers a module that refuses to start again.
// Nothing serves after that, and the manager is the caller's to stop.
func TestManager_failed_reload(t *testing.T) {
	var starts atomic.Int64

	mod := platform.NewUnimplementedModule("failing")
	mod.StartFn = func(context.Context) error {
		if starts.Add(1) > 1 {
			return context.DeadlineExceeded
		}
		return nil
	}

	m := newTestManager(t, mod)

	require.Error(t, m.Reload(t.Context()))
	require.Nil(t, m.Platform())

	m.Stop()
	m.Wait()
}

// TestManager_context_cancel covers the platform stopping on its own. The
// manager is watching for it, and goes down with the generation it holds.
func TestManager_context_cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	m := platform.NewManager(platform.NewTestOptions())
	t.Cleanup(m.Stop)

	require.NoError(t, m.Start(ctx))

	cancel()
	m.Wait()

	require.Nil(t, m.Platform())
}

// TestManager_start_error covers a manager that never came up. The socket
// is released, and Stop tolerates the state.
func TestManager_start_error(t *testing.T) {
	options := platform.NewTestOptions()
	options.ServerAddr = "127.0.0.1:not-a-port"

	m := platform.NewManager(options)
	require.Error(t, m.Start(t.Context()))
	require.Nil(t, m.Platform())

	m.Stop()
	m.Wait()
}
