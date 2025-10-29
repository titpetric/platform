package module

import (
	"errors"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/module/user"
)

// Load will load the default modules and add them to the platform.
func Load() error {
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

	addModule(user.NewHandler())

	return errors.Join(errs...)
}
