package platform

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// pidFileMode is the mode a pidfile is created with, before the umask. It is
// what every init script and `kill $(cat ...)` expects to be able to read.
const pidFileMode = 0o644

// pidfile is the file a process records its own process id in.
//
// The zero value is a disabled pidfile, which is what an empty
// Options.PidFile asks for, and every method on it is then a no-op. That is
// the whole of how the option is turned off: a caller that named no file
// takes the same code path as one that named it.
type pidfile struct {
	path string

	// pid is what write recorded, and zero until it succeeded. remove
	// consults it so a start that failed before the write does not delete a
	// file it never made.
	pid int
}

// newPidfile returns the pidfile path names, disabled when path is empty.
func newPidfile(path string) pidfile {
	return pidfile{path: path}
}

// write records this process's id.
//
// An existing file is overwritten rather than treated as a running instance.
// A pidfile is not a lock: deciding whether the pid in it is alive means
// kill(pid, 0), which answers wrongly across a pid namespace, after the pid
// counter wraps, and whenever an unrelated program now holds that number.
// Refusing to start on a file left by a killed process is how a service fails
// to come back after a power cut.
func (f *pidfile) write() error {
	if f.path == "" {
		return nil
	}

	pid := os.Getpid()
	if err := os.WriteFile(f.path, []byte(strconv.Itoa(pid)+"\n"), pidFileMode); err != nil {
		return fmt.Errorf("write pidfile: %w", err)
	}

	f.pid = pid
	return nil
}

// remove deletes the file, and only while it still holds the pid write put
// there.
//
// The check is what keeps an overlapping restart from leaving the new process
// unsignalable: if it wrote the file before this one got to its shutdown, the
// file is the new process's and is left alone. Reading and removing is not
// atomic, so this narrows the window rather than closing it, which is the
// same trade nginx makes.
func (f *pidfile) remove() error {
	if f.path == "" || f.pid == 0 {
		return nil
	}

	if pid, err := ReadPidFile(f.path); err != nil || pid != f.pid {
		return nil
	}
	if err := os.Remove(f.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove pidfile: %w", err)
	}
	return nil
}

// ReadPidFile returns the process id recorded in the file Options.PidFile
// names. It is what a command line sending a signal to a running platform
// reads, so the format stays owned by the package that writes it.
func ReadPidFile(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("pidfile %q does not hold a process id", path)
	}
	if pid <= 0 {
		return 0, fmt.Errorf("pidfile %q holds %d, which is not a process id", path, pid)
	}
	return pid, nil
}
