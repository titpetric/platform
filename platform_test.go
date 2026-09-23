package platform_test

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/internal/pidfile"
	"github.com/titpetric/platform/pkg/require"
)

func NewTestPlatform(tb testing.TB) *platform.Platform {
	svc, err := platform.Start(tb.Context(), platform.NewTestOptions())

	require.NoError(tb, err)
	require.NotNil(tb, svc)

	tb.Cleanup(svc.Stop)
	return svc
}

func TestPlatform(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		svc := platform.New(platform.NewTestOptions())
		t.Cleanup(svc.Stop)

		svc.Register(&platform.UnimplementedModule{
			NameFn: func() string {
				return "TestPlatform"
			},
			MountFn: func(_ context.Context, r platform.Router) error {
				r.Get("/404", func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte("You found a valid route"))
				})
				return nil
			},
		})
		svc.Use(platform.TestMiddleware())

		t.Run("find", func(t *testing.T) {
			var mod *platform.UnimplementedModule

			require.True(t, svc.Find(&mod))
			require.Equal(t, "TestPlatform", mod.Name())
		})

		plugins, mws := svc.Stats()
		require.Equal(t, 1, plugins)
		require.Equal(t, 1, mws)

		require.NoError(t, svc.Start(t.Context()))

		resp, err := http.Get(svc.URL() + "/404")
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, resp.Body.Close()) })

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, string(body), "You found a valid route")
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
		for i := range 30 {
			svc, err := platform.Start(t.Context(), platform.NewTestOptions())

			require.NoError(t, err)
			require.NotNil(t, svc)

			svc.Stop()

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
			svc, err := platform.Start(b.Context(), platform.NewTestOptions())

			require.NoError(b, err)
			require.NotNil(b, svc)

			svc.Stop()
		}
	})
}

// TestPlatform_stop_after_failed_start covers the lifecycle hole where setup
// fails before the server is built. Stop is still how a caller releases what
// did get allocated, so it has to tolerate that.
func TestPlatform_stop_after_failed_start(t *testing.T) {
	options := platform.NewTestOptions()
	options.ServerAddr = "127.0.0.1:not-a-port"

	svc := platform.New(options)
	require.Error(t, svc.Start(t.Context()))

	svc.Stop()
	svc.Wait()
}

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

// The environment the parent hands the child half of the signal test.
const (
	signalChildEnv   = "PLATFORM_TEST_SIGNAL_CHILD"
	signalPidFileEnv = "PLATFORM_TEST_SIGNAL_PIDFILE"

	// signalReady is printed once Start has returned, which is after the
	// signal handler is armed. The parent waits for it before signalling:
	// a signal delivered before then takes its default disposition and
	// kills the child outright.
	signalReady = "ready"

	// signalStopped is printed once Wait has returned, so the parent can
	// tell a graceful stop from a process that was killed.
	signalStopped = "stopped"
)

// TestPlatformSignalChild runs only as the child TestPlatformSignals starts.
func TestPlatformSignalChild(t *testing.T) {
	if os.Getenv(signalChildEnv) == "" {
		t.Skip("runs only as the child of TestPlatformSignals")
	}

	options := verboseTestOptions()
	options.PidFile = os.Getenv(signalPidFileEnv)

	svc := platform.New(options)
	if err := svc.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	os.Stdout.WriteString(signalReady + "\n")
	svc.Wait()
	os.Stdout.WriteString(signalStopped + "\n")
}

// TestPlatformSignals covers the two signals a platform stops on.
//
// It signals a child rather than raising in process: signal.Notify is
// process-wide, so a SIGTERM raised here would reach every platform the
// suite has running, and an unarmed moment kills the test binary.
func TestPlatformSignals(t *testing.T) {
	if testing.Short() {
		t.Skip("re-executes the test binary")
	}

	executable, err := os.Executable()
	require.NoError(t, err)

	for _, test := range []struct {
		name   string
		signal os.Signal
		logged string
	}{
		{name: "sigterm", signal: syscall.SIGTERM, logged: "caught sigterm, stopping server"},
		{name: "sigint", signal: os.Interrupt, logged: "caught sigterm, stopping server"},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := pidPath(t)

			// A child that never exits is killed with the context rather
			// than left to hang the suite.
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()

			cmd := exec.CommandContext(ctx, executable, "-test.run=^TestPlatformSignalChild$", "-test.v")
			cmd.Env = append(os.Environ(),
				signalChildEnv+"=1",
				signalPidFileEnv+"="+path,
			)

			var stderr strings.Builder
			cmd.Stderr = &stderr

			stdout, err := cmd.StdoutPipe()
			require.NoError(t, err)

			require.NoError(t, cmd.Start())

			// Safe to signal only after this: Start arms the handler
			// before the child prints the line.
			lines := bufio.NewScanner(stdout)
			require.True(t, waitFor(lines, signalReady), "the child never reported ready")

			// The one place the pid means something: in process it would
			// compare the test to itself.
			pid, err := pidfile.Read(path)
			require.NoError(t, err)
			require.Equal(t, cmd.Process.Pid, pid)

			require.NoError(t, cmd.Process.Signal(test.signal))

			require.True(t, waitFor(lines, signalStopped), "the child did not stop gracefully")

			// An error here is the child killed by the signal rather than
			// catching it, which is the defect this test exists for.
			require.NoError(t, cmd.Wait())
			require.Contains(t, stderr.String(), test.logged)

			_, err = os.Stat(path)
			require.Error(t, err, "the pidfile outlived the process")
		})
	}
}

// waitFor reports whether want appears on its own line before the stream
// ends.
func waitFor(lines *bufio.Scanner, want string) bool {
	for lines.Scan() {
		if strings.TrimSpace(lines.Text()) == want {
			return true
		}
	}
	return false
}
