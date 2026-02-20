package platform

import (
	"io/fs"
	"os"
	"strings"
)

// Options is a configuration struct for platform behaviour.
type Options struct {
	// ServerAddr is the address the server listens to.
	ServerAddr string

	// Quiet turns down the verbosity in the Platform logging code, set to true in tests.
	Quiet bool

	// Modules controls which modules get loaded. If the list
	// is empty (unconfigured, zero value), all modules load.
	Modules []string

	// ThemeFS provides access to application-level theme files (theme.yml, data/*.yml).
	// It can be used for configuration purposes by modules. It's optional and may be nil.
	ThemeFS fs.FS
}

// NewOptions provides default options for the platform.
func NewOptions() *Options {
	opt := &Options{}
	opt.ServerAddr = opt.env("PLATFORM_SERVER_ADDR", ":8080")
	opt.Modules = opt.envCSV("PLATFORM_MODULES")
	return opt
}

func (*Options) envCSV(name string) []string {
	if v := os.Getenv(name); v != "" {
		return strings.Split(v, ",")
	}
	return nil
}

func (*Options) env(name string, def string) string {
	result := def
	if v := os.Getenv(name); v != "" {
		result = v
	}
	return result
}

// NewTestOptions produces default options for tests.
func NewTestOptions() *Options {
	return &Options{
		ServerAddr: "127.0.0.1:0",
		Quiet:      true,
	}
}
