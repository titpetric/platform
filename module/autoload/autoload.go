package autoload

import (
	"github.com/titpetric/platform"
	"github.com/titpetric/platform/module/user"
)

func init() {
	platform.Register(user.NewHandler())
}
