package cli

import (
	"context"
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

// TestApp_AddCommand tests that commands are registered and can be found.
func TestApp_AddCommand(t *testing.T) {
	app := NewApp("testapp")

	app.AddCommand("hello", "Say hello", func() *Command {
		cmd := &Command{
			Name: "hello",
			Run: func(ctx context.Context, args []string) error {
				return nil
			},
		}
		return cmd
	})

	cmd, err := app.findCommand([]string{"hello"}, "")
	assert.NoError(t, err)
	assert.Equal(t, "hello", cmd.Name)
}

// TestApp_FindCommand_NotFound tests that findCommand returns an error for unknown commands.
func TestApp_FindCommand_NotFound(t *testing.T) {
	app := NewApp("testapp")

	_, err := app.findCommand([]string{"nonexistent"}, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown command")
}

// TestApp_FindCommand_NoCommand tests that findCommand returns an error when no command is specified.
func TestApp_FindCommand_NoCommand(t *testing.T) {
	app := NewApp("testapp")

	_, err := app.findCommand([]string{}, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no command specified")
}

// TestParseCommands tests that parseCommands extracts command names up to the first flag.
func TestParseCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "single command",
			args:     []string{"hello"},
			expected: []string{"hello"},
		},
		{
			name:     "command with flags",
			args:     []string{"hello", "--name", "world"},
			expected: []string{"hello"},
		},
		{
			name:     "multiple commands",
			args:     []string{"sub", "cmd", "--flag"},
			expected: []string{"sub", "cmd"},
		},
		{
			name:     "flag at start",
			args:     []string{"--flag", "hello"},
			expected: []string{},
		},
		{
			name:     "empty args",
			args:     []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCommands(tt.args)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestApp_RunWithArgs_Integration tests the full command execution flow.
func TestApp_RunWithArgs_Integration(t *testing.T) {
	app := NewApp("testapp")
	executed := false

	app.AddCommand("test", "Test command", func() *Command {
		var msg string

		cmd := &Command{
			Name: "test",
			Bind: func(fs *flag.FlagSet) {
				fs.StringVar(&msg, "msg", "default", "")
			},
			Run: func(ctx context.Context, args []string) error {
				executed = true
				assert.Equal(t, "hello", msg)
				return nil
			},
		}

		return cmd
	})

	// Simulate command line with flag
	err := app.RunWithArgs([]string{"test", "--msg", "hello"})
	assert.NoError(t, err)
	assert.True(t, executed)
}
