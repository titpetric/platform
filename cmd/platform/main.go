package main

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/telemetry"

	// Add platform drivers and modules.
	_ "github.com/titpetric/platform/drivers"
	_ "github.com/titpetric/platform/module/autoload"
)

func main() {
	start(context.Background())
}

var options = platform.NewOptions()

func start(ctx context.Context) {
	// Register common middleware.
	platform.Use(middleware.Logger)
	platform.Use(telemetry.Middleware("platform"))

	p := platform.New(options)

	if err := p.Start(ctx); err != nil {
		telemetry.Fatal(ctx, fmt.Errorf("exit error: %w", err))
	}

	p.Wait()
}
