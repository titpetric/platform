package pidfile

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

// TestPidfile covers writing and removing, including the cases where the
// file belongs to another process.
func TestPidfile(t *testing.T) {
	t.Run("an empty path does nothing", func(t *testing.T) {
		dir := t.TempDir()

		var f Pidfile
		assert.NoError(t, f.Write())
		assert.NoError(t, f.Remove())

		entries, err := os.ReadDir(dir)
		assert.NoError(t, err)
		assert.Empty(t, entries)
	})

	t.Run("writes the pid and a newline", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "run.pid")

		f := New(path)
		assert.NoError(t, f.Write())

		data, err := os.ReadFile(path)
		assert.NoError(t, err)
		assert.Equal(t, strconv.Itoa(os.Getpid())+"\n", string(data))
	})

	t.Run("is readable and not executable", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "run.pid")

		f := New(path)
		assert.NoError(t, f.Write())

		info, err := os.Stat(path)
		assert.NoError(t, err)
		// The umask owns the rest of the mode, so asserting 0644 exactly
		// fails under a umask that is not 022.
		assert.True(t, info.Mode().Perm()&0o400 != 0)
		assert.True(t, info.Mode().Perm()&0o111 == 0)
	})

	t.Run("overwrites a file left by a dead process", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "run.pid")
		assert.NoError(t, os.WriteFile(path, []byte("999999\n"), 0o644))

		f := New(path)
		assert.NoError(t, f.Write())

		pid, err := Read(path)
		assert.NoError(t, err)
		assert.Equal(t, os.Getpid(), pid)
	})

	t.Run("remove deletes the file it wrote", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "run.pid")

		f := New(path)
		assert.NoError(t, f.Write())
		assert.NoError(t, f.Remove())

		_, err := os.Stat(path)
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("remove leaves a file another process claimed", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "run.pid")

		f := New(path)
		assert.NoError(t, f.Write())

		// An overlapping restart: the next process wrote the file before
		// this one reached its shutdown. Removing it would leave the live
		// process with nothing to signal.
		assert.NoError(t, os.WriteFile(path, []byte("999999\n"), 0o644))
		assert.NoError(t, f.Remove())

		pid, err := Read(path)
		assert.NoError(t, err)
		assert.Equal(t, 999999, pid)
	})

	t.Run("remove without a write does nothing", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "run.pid")
		assert.NoError(t, os.WriteFile(path, []byte("999999\n"), 0o644))

		// A start that failed before the write has a path and no pid, and
		// must not delete a file it never made.
		f := New(path)
		assert.NoError(t, f.Remove())

		_, err := os.Stat(path)
		assert.NoError(t, err)
	})

	t.Run("a missing directory is an error and is not created", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "absent")

		f := New(filepath.Join(dir, "run.pid"))
		assert.ErrorIs(t, f.Write(), fs.ErrNotExist)

		_, err := os.Stat(dir)
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})
}

// TestReadPidFile covers the format a command line reads back to signal the
// process.
func TestRead(t *testing.T) {
	dir := t.TempDir()

	write := func(t *testing.T, name, contents string) string {
		t.Helper()
		path := filepath.Join(dir, name)
		assert.NoError(t, os.WriteFile(path, []byte(contents), 0o644))
		return path
	}

	t.Run("reads what write recorded", func(t *testing.T) {
		path := filepath.Join(dir, "round-trip.pid")
		f := New(path)
		assert.NoError(t, f.Write())

		pid, err := Read(path)
		assert.NoError(t, err)
		assert.Equal(t, os.Getpid(), pid)
	})

	t.Run("tolerates surrounding whitespace", func(t *testing.T) {
		pid, err := Read(write(t, "spaced.pid", "  4321\n\n"))
		assert.NoError(t, err)
		assert.Equal(t, 4321, pid)
	})

	t.Run("rejects contents that are not a number", func(t *testing.T) {
		_, err := Read(write(t, "garbage.pid", "not a pid\n"))
		assert.Error(t, err)
	})

	t.Run("rejects a pid of zero or less", func(t *testing.T) {
		_, err := Read(write(t, "zero.pid", "0\n"))
		assert.Error(t, err)

		_, err = Read(write(t, "negative.pid", "-1\n"))
		assert.Error(t, err)
	})

	t.Run("reports a missing file", func(t *testing.T) {
		_, err := Read(filepath.Join(dir, "absent.pid"))
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})
}
