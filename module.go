package platform

import "context"

// Module is the implementation contract for modules.
//
// The interface should only be used to enforce the API contract as
// shown below. It's also used to provide `platform.Register()`.
type Module interface {
	// Name should return a meaningful name for your module.
	Name() string

	// Start is used to create any goroutines or otherwise
	// set up the module by starting a server. It allows
	// to implement a lifecycle of the service.
	Start(context.Context) error

	// Stop should clean up any goroutines, clean up leaks.
	Stop(context.Context) error

	// Mount runs before the server starts, and allows you to
	// register new routes to your module.
	Mount(context.Context, Router) error
}
