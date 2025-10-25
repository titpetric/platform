package module

import (
	"errors"

	"github.com/titpetric/platform/module/theme"
	"github.com/titpetric/platform/module/user"
	"github.com/titpetric/platform/registry"
)

// Assert implementation contracts.
var (
	_ registry.Module = (*user.Handler)(nil)
)

// LoadModules will load the default modules and add them to the registry.
func LoadModules() error {
	var (
		errs []error

		// addModule is a readability closure, deduplicating error checks for modules.
		addModule = func(m registry.Module, err error) {
			if err != nil {
				errs = append(errs, err)
				return
			}
			registry.AddModule(m)
		}
	)

	addModule(user.NewHandler(theme.TemplateFS))

	return errors.Join(errs...)
}
