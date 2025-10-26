package module

import (
	"errors"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/module/theme"
	"github.com/titpetric/platform/module/user"
)

// Assert implementation contracts.
var (
	_ platform.Module = (*user.Handler)(nil)
)

// LoadModules will load the default modules and add them to the platform.
func LoadModules() error {
	var (
		errs []error

		// addModule is a readability closure, deduplicating error checks for modules.
		addModule = func(m platform.Module, err error) {
			if err != nil {
				errs = append(errs, err)
				return
			}
			platform.AddModule(m)
		}
	)

	addModule(user.NewHandler(theme.TemplateFS))

	return errors.Join(errs...)
}
