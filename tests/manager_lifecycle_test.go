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

// TestManagerReloadAfterStop covers a Reload against a stopped manager, which
// is refused, leaves no generation, and stops every module it started.
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

// TestManagerStopContextLive covers the teardown order: Module.Stop gets a live
// context, and Wait releases only once every Stop has returned.
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

// TestManagerPidFileFailureReleasesSocket covers a Start that fails on the
// pidfile, after the listener bound and before Serve runs. Stop has to return
// without waiting out its timeout, and the port has to be free.
func TestManagerPidFileFailureReleasesSocket(t *testing.T) {
	options := platform.NewTestOptions()
	options.PidFile = "/proc/definitely/not/writable/run.pid"

	p := platform.New(options)

	assert.Error(t, p.Start(t.Context()), "an unwritable pidfile fails the start")

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

// TestManagerSigtermDuringReload covers a SIGTERM landing while a reload is in
// flight, which the manager has to catch and drain rather than take the default
// disposition on.
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
