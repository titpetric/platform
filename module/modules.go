package module

import (
	"github.com/go-chi/chi/v5"

	"github.com/titpetric/platform/module/theme"
	"github.com/titpetric/platform/module/user"
	"github.com/titpetric/platform/registry"
)

// Assert implementation contracts.
var (
	_ Module = (*user.Handler)(nil)
	_ Module = (*registry.Plugin)(nil)
)

// Module is the implementation contract for interfaces.
// The interface should only be used to enforce the API
// contract as shown above.
type Module interface {
	Mount(chi.Router)
	Close()
}

// Modules is a collection holding all your modules.
type Modules struct {
	user *user.Handler
}

// LoadModules will load the default modules and add them to the registry.
func LoadModules() error {
	var err error

	result := &Modules{}

	result.user, err = user.NewHandler(theme.TemplateFS)
	if err != nil {
		return err
	}

	registry.Add("user", result.user.Mount, result.user.Close)

	return nil
}
