package platform_test

import (
	"errors"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/require"
)

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
