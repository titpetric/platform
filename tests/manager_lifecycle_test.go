package platform_test

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/internal/assert"
)

// TestManagerReloadAfterStop covers Reload building a generation nothing could
// retire. Stop's once is spent by then, so the modules stayed started, the
// watch goroutine never returned, and Platform() was non-nil after Stop.
func TestManagerReloadAfterStop(t *testing.T) {
	mod := &countingModule{}

	m := platform.NewManager(platform.NewTestOptions())
	m.Setup = func(p *platform.Platform) error {
		p.Register(mod)
		return nil
	}

	assert.NoError(t, m.Start(t.Context()))
	assert.NotNil(t, m.Platform())

	m.Stop()
	assert.Nil(t, m.Platform(), "a stopped manager serves nothing")

	err := m.Reload(t.Context())
	assert.Error(t, err, "a stopped manager does not start again")
	assert.Contains(t, err.Error(), "manager stopped")

	assert.Nil(t, m.Platform(), "the refused reload must not leave a generation")
	assert.Equal(t, mod.starts.Load(), mod.stops.Load(), "every module that started was stopped")
}

// TestManagerStopContextLive covers the teardown order. cancel used to run
// before the registry closed, so every Module.Stop was handed a context that
// was already done and Wait returned while a Stop was still running.
func TestManagerStopContextLive(t *testing.T) {
	var (
		seen    atomic.Pointer[error]
		flushed atomic.Bool
	)

	p := platform.New(platform.NewTestOptions())
	p.Register(&platform.UnimplementedModule{
		NameFn: func() string { return "flusher" },
		StopFn: func(ctx context.Context) error {
			err := ctx.Err()
			seen.Store(&err)

			// Long enough that a Wait released early returns first.
			time.Sleep(200 * time.Millisecond)
			flushed.Store(true)
			return nil
		},
	})

	assert.NoError(t, p.Start(t.Context()))

	waited := make(chan struct{})
	go func() {
		p.Wait()
		close(waited)
	}()

	p.Stop()

	<-waited
	assert.True(t, flushed.Load(), "Wait returned before the module finished flushing")

	err := seen.Load()
	assert.NotNil(t, err, "Stop was never called")
	assert.NoError(t, *err, "Module.Stop must get a live context")
}

// TestManagerPidFileFailureReleasesSocket covers the socket a failed pidfile
// write used to leak. setup binds the listener before the pidfile is written,
// and Serve runs after, so the listener was one http.Server.Shutdown knew
// nothing about and p.served had no writer: Stop waited out its whole timeout
// and the port stayed bound for the life of the process.
func TestManagerPidFileFailureReleasesSocket(t *testing.T) {
	options := platform.NewTestOptions()
	options.PidFile = "/proc/definitely/not/writable/run.pid"

	p := platform.New(options)

	assert.Error(t, p.Start(t.Context()), "an unwritable pidfile fails the start")

	// The old code took the full five second timeout here.
	start := time.Now()
	p.Stop()
	assert.True(t, time.Since(start) < 2*time.Second, "Stop waited on a channel with no writer")
}

// The environment the parent hands the child half of the reload signal test.
const (
	reloadChildEnv = "PLATFORM_TEST_RELOAD_CHILD"

	reloadReady   = "reload-ready"
	reloadStopped = "reload-stopped"
)

// TestManagerReloadSignalChild runs only as the child
// TestManagerSigtermDuringReload starts. Its module takes long enough to start
// that the parent can land a SIGTERM inside the reload window.
func TestManagerReloadSignalChild(t *testing.T) {
	if os.Getenv(reloadChildEnv) == "" {
		t.Skip("runs only as the child of TestManagerSigtermDuringReload")
	}

	m := platform.NewManager(verboseTestOptions())
	m.Setup = func(p *platform.Platform) error {
		p.Register(&platform.UnimplementedModule{
			NameFn: func() string { return "slow" },
			StartFn: func(context.Context) error {
				// Only the generations a reload builds are slow, so the
				// first start does not stall the handshake.
				if os.Getenv(reloadChildEnv) == "reloading" {
					time.Sleep(3 * time.Second)
				}
				return nil
			},
		})
		return nil
	}

	if err := m.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	os.Stdout.WriteString(reloadReady + "\n")
	m.Wait()
	os.Stdout.WriteString(reloadStopped + "\n")
}

// TestManagerSigtermDuringReload covers the window a reload used to leave with
// no SIGTERM handler. The generation armed the process signals, and retire
// released them, so a SIGTERM between one generation stopping and the next
// arming its own took the default disposition and killed the process.
//
// It signals a child rather than raising in process: signal.Notify is
// process-wide, and an unarmed moment kills the test binary.
func TestManagerSigtermDuringReload(t *testing.T) {
	if testing.Short() {
		t.Skip("re-executes the test binary")
	}

	executable, err := os.Executable()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, executable, "-test.run=^TestManagerReloadSignalChild$", "-test.v")
	cmd.Env = append(os.Environ(), reloadChildEnv+"=reloading")

	var stderr strings.Builder
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	assert.NoError(t, err)

	assert.NoError(t, cmd.Start())

	lines := bufio.NewScanner(stdout)
	assert.True(t, waitFor(lines, reloadReady), "the child never reported ready")

	// SIGHUP retires the generation and starts one whose module sleeps, so
	// the SIGTERM below lands while the reload is in flight.
	assert.NoError(t, cmd.Process.Signal(syscall.SIGHUP))
	time.Sleep(500 * time.Millisecond)
	assert.NoError(t, cmd.Process.Signal(syscall.SIGTERM))

	assert.True(t, waitFor(lines, reloadStopped), "the child did not stop gracefully")

	// An error here is the child killed by the signal rather than catching
	// it, which is the defect this test exists for.
	assert.NoError(t, cmd.Wait())
	assert.Contains(t, stderr.String(), "caught sigterm, stopping server")
}
