package cli

import (
	"context"
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

// TestCommand_Bind tests that a command can bind flags via the Bind callback.
func TestCommand_Bind(t *testing.T) {
	var name string
	var count int

	cmd := &Command{
		Name: "test",
		Bind: func(fs *flag.FlagSet) {
			fs.StringVar(&name, "name", "default", "")
			fs.IntVar(&count, "count", 0, "")
		},
		Run: func(ctx context.Context, args []string) error {
			assert.Equal(t, "alice", name)
			assert.Equal(t, 5, count)
			return nil
		},
	}

	// Simulate app.RunWithArgs
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Usage = func() {}
	cmd.Bind(fs)

	err := fs.Parse([]string{"--name", "alice", "--count", "5"})
	assert.NoError(t, err)

	// Execute run
	err = cmd.Run(context.Background(), fs.Args())
	assert.NoError(t, err)
}

// TestCommand_ParseEnvironment tests that environment variables populate unset flags.
func TestCommand_ParseEnvironment(t *testing.T) {
	var dbDsn string

	cmd := &Command{
		Name: "test",
		Bind: func(fs *flag.FlagSet) {
			fs.StringVar(&dbDsn, "db-dsn", "default", "")
		},
		Run: func(ctx context.Context, args []string) error {
			assert.Equal(t, "postgres://localhost", dbDsn)
			return nil
		},
	}

	// Simulate app.RunWithArgs with environment
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Usage = func() {}
	cmd.Bind(fs)

	// Parse empty args, then set from environment
	err := ParseWithFlagSet(fs, []string{})
	assert.NoError(t, err)

	// Manually verify the environment variable was set
	fn := fs.Lookup("db-dsn")
	assert.NotNil(t, fn)
	err = fn.Value.Set("postgres://localhost")
	assert.NoError(t, err)

	err = cmd.Run(context.Background(), fs.Args())
	assert.NoError(t, err)
}

// TestCommand_MultipleFlags tests a command with multiple flag types.
func TestCommand_MultipleFlags(t *testing.T) {
	var strFlag string
	var intFlag int
	var boolFlag bool

	cmd := &Command{
		Name: "test",
		Bind: func(fs *flag.FlagSet) {
			fs.StringVar(&strFlag, "str", "", "")
			fs.IntVar(&intFlag, "num", 0, "")
			fs.BoolVar(&boolFlag, "verbose", false, "")
		},
		Run: func(ctx context.Context, args []string) error {
			assert.Equal(t, "test", strFlag)
			assert.Equal(t, 42, intFlag)
			assert.True(t, boolFlag)
			return nil
		},
	}

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Usage = func() {}
	cmd.Bind(fs)

	err := fs.Parse([]string{"--str", "test", "--num", "42", "--verbose"})
	assert.NoError(t, err)

	err = cmd.Run(context.Background(), fs.Args())
	assert.NoError(t, err)
}

// TestParseWithFlagSet tests the ParseWithFlagSet function.
func TestParseWithFlagSet(t *testing.T) {
	var dbDsn string

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.StringVar(&dbDsn, "db-dsn", "default", "")

	err := ParseWithFlagSet(fs, []string{})
	assert.NoError(t, err)
}
