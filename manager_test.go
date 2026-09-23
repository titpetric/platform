package platform_test

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/internal/pidfile"
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

func TestManagerPidFile(t *testing.T) {
	newManager := func(tb testing.TB, path string) *platform.Manager {
		tb.Helper()

		m := platform.NewManager(testOptions(path))
		m.Setup = func(p *platform.Platform) error {
			p.Use(platform.TestMiddleware())
			return nil
		}
		tb.Cleanup(m.Stop)
		return m
	}

	t.Run("written on start and removed on stop", func(t *testing.T) {
		path := pidPath(t)

		m := newManager(t, path)
		require.NoError(t, m.Start(t.Context()))

		pid, err := pidfile.Read(path)
		require.NoError(t, err)
		require.Equal(t, os.Getpid(), pid)

		m.Stop()
		m.Wait()

		_, err = os.Stat(path)
		require.ErrorIs(t, err, fs.ErrNotExist)
	})

	// A reload replaces the platform, not the process, so there is no
	// window in which a command reading the file finds it missing.
	t.Run("a reload leaves the file alone", func(t *testing.T) {
		path := pidPath(t)

		m := newManager(t, path)
		require.NoError(t, m.Start(t.Context()))

		before := m.Platform()
		require.NotNil(t, before)

		require.NoError(t, m.Reload(t.Context()))

		require.NotEqual(t, before, m.Platform())

		pid, err := pidfile.Read(path)
		require.NoError(t, err)
		require.Equal(t, os.Getpid(), pid)
	})

	t.Run("a path that cannot be written fails the start", func(t *testing.T) {
		m := newManager(t, filepath.Join(t.TempDir(), "absent", "run.pid"))

		require.Error(t, m.Start(t.Context()))
		require.Nil(t, m.Platform())
	})
}

// TestManagerCheck covers the hook that decides whether a reload happens,
// separating a refusal that keeps serving from a failure that does not.
func TestManagerCheck(t *testing.T) {
	// Check is read from the handler's goroutine, so it goes on before Start.
	newManager := func(tb testing.TB, mod platform.Module, check func() error) *platform.Manager {
		tb.Helper()

		m := platform.NewManager(platform.NewTestOptions())
		m.Setup = func(p *platform.Platform) error {
			p.Register(mod)
			p.Use(platform.TestMiddleware())
			return nil
		}
		m.Check = check
		tb.Cleanup(m.Stop)
		require.NoError(tb, m.Start(tb.Context()))
		return m
	}

	t.Run("a refused reload leaves the generation serving", func(t *testing.T) {
		mod := &countingModule{}
		refused := errors.New("the configuration does not load")
		m := newManager(t, mod, func() error { return refused })

		serving := m.Platform()
		require.NotNil(t, serving)

		require.ErrorIs(t, m.Reload(t.Context()), refused)

		// Nothing was torn down to find out the configuration was bad.
		require.Equal(t, serving, m.Platform())
		require.Equal(t, int64(1), mod.starts.Load())
		require.Equal(t, int64(0), mod.stops.Load())

		status, body := get(t, m.URL()+"/generation")
		require.Equal(t, 200, status)
		require.Equal(t, "1", string(body))
	})

	t.Run("a passing check reloads", func(t *testing.T) {
		mod := &countingModule{}

		var called atomic.Int64
		m := newManager(t, mod, func() error {
			called.Add(1)
			return nil
		})

		require.NoError(t, m.Reload(t.Context()))
		require.Equal(t, int64(1), called.Load())
		require.Equal(t, int64(2), mod.starts.Load())

		_, body := get(t, m.URL()+"/generation")
		require.Equal(t, "2", string(body))
	})

	t.Run("the error names the refusal", func(t *testing.T) {
		m := newManager(t, &countingModule{}, func() error {
			return errors.New("mail.host is not a server name")
		})

		err := m.Reload(t.Context())
		require.ErrorContains(t, err, "reload refused")
		require.ErrorContains(t, err, "mail.host is not a server name")
	})

	// The handler decides from what is serving rather than from the error,
	// so a refused SIGHUP must not take the process down.
	t.Run("a refused reload survives sighup", func(t *testing.T) {
		mod := &countingModule{}
		m := newManager(t, mod, func() error {
			return errors.New("the configuration does not load")
		})

		serving := m.Platform()

		require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGHUP))

		// The handler runs on a goroutine; give it time, then assert
		// nothing moved.
		time.Sleep(500 * time.Millisecond)

		require.Equal(t, serving, m.Platform())
		require.Equal(t, int64(0), mod.stops.Load())

		select {
		case <-m.Context().Done():
			t.Fatal("the manager stopped on a refused reload")
		default:
		}

		status, _ := get(t, m.URL()+"/generation")
		require.Equal(t, 200, status)
	})
}
