package autoload

import (
	"log"

	"github.com/titpetric/platform/module"
)

func init() {
	err := module.Load()
	if err != nil {
		log.Fatalf("init error loading modules: %v", err)
	}
}
