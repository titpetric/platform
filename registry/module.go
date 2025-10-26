package registry

// Module is the implementation contract for modules.
//
// The interface should only be used to enforce the API contract as
// shown below. It's also used to provide AddModule().
type Module interface {
	// Name should return a meaningful name for your module.
	Name() string

	// Mount runs before the server starts, and allows you to
	// register new routes to your module.
	Mount(Router)

	// Close will be invoked when the server shuts down.
	// It's mainly important for tests, so that modules don't leak
	// goroutines, and have an opportunity to flush any in-memory
	// data to storage in case of a graceful shutdown.
	Close()
}

// Assert *Plugin implements the Module interface.
var _ Module = (*Plugin)(nil)
