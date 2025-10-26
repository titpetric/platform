package platform

// Options is a configuration struct for platform behaviour.
type Options struct {
	// ServerAddr is the address the server listens to.
	ServerAddr string

	// Quiet turns down the verbosity in the Platform logging code, set to true in tests.
	Quiet bool
}

// NewOptions provides default options for the platform.
func NewOptions() *Options {
	return &Options{
		ServerAddr: ":8080",
	}
}
