package platform

import "os"

// Options is a configuration struct for platform behaviour.
type Options struct {
	// ServerAddr is the address the server listens to.
	ServerAddr string

	// Quiet turns down the verbosity in the Platform logging code, set to true in tests.
	Quiet bool
}

// NewOptions provides default options for the platform.
func NewOptions() *Options {
	opt := &Options{}
	opt.ServerAddr = opt.env("PLATFORM_SERVER_ADDR", ":8080")
	return opt
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
