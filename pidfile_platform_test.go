package platform_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/internal/pidfile"
	"github.com/titpetric/platform/pkg/require"
)

// pidPath returns an uncreated path for a test to point Options.PidFile at.
func pidPath(tb testing.TB) string {
	tb.Helper()
	return filepath.Join(tb.TempDir(), "run.pid")
}

// testOptions is NewTestOptions with a pidfile.
func testOptions(path string) *platform.Options {
	options := platform.NewTestOptions()
	options.PidFile = path
	return options
}

func TestPlatformPidFile(t *testing.T) {
	t.Run("written on start and removed on stop", func(t *testing.T) {
		path := pidPath(t)

		svc := platform.New(testOptions(path))
		require.NoError(t, svc.Start(t.Context()))

		pid, err := pidfile.Read(path)
		require.NoError(t, err)
		require.Equal(t, os.Getpid(), pid)

		svc.Stop()
		svc.Wait()

		_, err = os.Stat(path)
		require.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("an empty path writes nothing", func(t *testing.T) {
		dir := t.TempDir()

		svc := platform.New(platform.NewTestOptions())
		t.Cleanup(svc.Stop)
		require.NoError(t, svc.Start(t.Context()))

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Empty(t, entries)
	})

	t.Run("a path that cannot be written fails the start", func(t *testing.T) {
		options := testOptions(filepath.Join(t.TempDir(), "absent", "run.pid"))

		svc := platform.New(options)
		require.Error(t, svc.Start(t.Context()))

		svc.Stop()
		svc.Wait()
	})

	t.Run("not written when the socket cannot be bound", func(t *testing.T) {
		path := pidPath(t)

		options := testOptions(path)
		options.ServerAddr = "127.0.0.1:not-a-port"

		svc := platform.New(options)
		require.Error(t, svc.Start(t.Context()))

		// A file naming a process that never came up is what a service
		// manager would act on.
		_, err := os.Stat(path)
		require.ErrorIs(t, err, fs.ErrNotExist)

		svc.Stop()
		svc.Wait()
	})

	t.Run("the shorthand releases what a failed start allocated", func(t *testing.T) {
		options := testOptions(filepath.Join(t.TempDir(), "absent", "run.pid"))
		options.ServerAddr = "127.0.0.1:0"

		svc, err := platform.Start(t.Context(), options)
		require.Error(t, err)
		require.Nil(t, svc)
	})
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
