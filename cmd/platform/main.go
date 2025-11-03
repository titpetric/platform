package main

import (
	"context"

	platformcmd "github.com/titpetric/platform/cmd"

	// Add platform drivers and modules.
	_ "github.com/titpetric/platform/pkg/drivers"
)

func main() {
	platformcmd.Main(context.Background())
}
