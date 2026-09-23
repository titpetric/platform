// Package pidfile records a process id in a file, for a service manager or a
// command line that signals the process later.
package pidfile

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Mode is the mode a pidfile is created with, before the umask.
const Mode = 0o644

// Pidfile is the file a process records its own id in.
//
// The zero value is disabled, which is what an empty path asks for, and every
// method on it is then a no-op. A pidfile is a record and not a lock: Write
// overwrites whatever was there, and Remove gives the file up only while it
// still holds the id Write put in it.
type Pidfile struct {
	path string
	pid  int
}

// New returns the pidfile at path, disabled when path is empty.
func New(path string) Pidfile {
	return Pidfile{path: path}
}

// Write records this process's id, overwriting a file left by a dead one.
func (f *Pidfile) Write() error {
	if f.path == "" {
		return nil
	}

	pid := os.Getpid()
	if err := os.WriteFile(f.path, []byte(strconv.Itoa(pid)+"\n"), Mode); err != nil {
		return fmt.Errorf("write pidfile: %w", err)
	}

	f.pid = pid
	return nil
}

// Remove deletes the file, unless another process has since claimed it.
func (f *Pidfile) Remove() error {
	if f.path == "" || f.pid == 0 {
		return nil
	}

	if pid, err := Read(f.path); err != nil || pid != f.pid {
		return nil
	}
	if err := os.Remove(f.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove pidfile: %w", err)
	}
	return nil
}

// Read returns the process id recorded in the file at path.
func Read(path string) (int, error) {
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
