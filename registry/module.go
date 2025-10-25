package registry

import "github.com/go-chi/chi/v5"

// Module is the implementation contract for modules.
//
// The interface should only be used to enforce the API contract as
// shown below. It's also used to provide AddModule().
type Module interface {
	Name() string
	Mount(chi.Router)
	Close()
}

// Assert *Plugin implements the Module interface.
var _ Module = (*Plugin)(nil)
