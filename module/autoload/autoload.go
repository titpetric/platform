package autoload

import (
	"log"

	"github.com/titpetric/platform/module"
)

func init() {
	err := module.LoadModules()
	if err != nil {
		log.Fatalf("init error loading modules: %v", err)
	}
}
