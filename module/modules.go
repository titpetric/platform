package module

import (
	"github.com/titpetric/platform/module/theme"
	"github.com/titpetric/platform/module/user"
	"github.com/titpetric/platform/registry"
)

// Assert implementation contracts.
var (
	_ registry.Module = (*user.Handler)(nil)
)

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

	registry.AddModule(result.user)

	return nil
}
