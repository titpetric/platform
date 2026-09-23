package platform_test

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/require"
)

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

// TestPlatformSignalChild is the other half of TestPlatformSignals. It runs
// only in the child process the parent re-executes, and is skipped in an
// ordinary run of the suite.
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

// TestPlatformSignals covers the two signals a platform stops on, in a
// process of its own.
//
// The signal cannot be sent inside the test binary. signal.Notify is
// process-wide and additive, so a SIGTERM raised here would be delivered to
// every platform the rest of the suite has running, and a moment when no
// handler is armed kills the test binary rather than failing a test. A child
// process bounds the blast radius to itself.
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

			// Blocking on the ready line is what makes the signal safe to
			// send: the child prints it only after Start returned, and
			// Start arms the handler before it does.
			lines := bufio.NewScanner(stdout)
			require.True(t, waitFor(lines, signalReady), "the child never reported ready")

			// The one place the pid is worth asserting: in process,
			// os.Getpid() would be comparing the test to itself.
			pid, err := platform.ReadPidFile(path)
			require.NoError(t, err)
			require.Equal(t, cmd.Process.Pid, pid)

			require.NoError(t, cmd.Process.Signal(test.signal))

			require.True(t, waitFor(lines, signalStopped), "the child did not stop gracefully")

			// A non-nil error here is the child having been killed by the
			// signal instead of catching it, which is the defect this test
			// exists for.
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
