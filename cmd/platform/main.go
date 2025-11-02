package main

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/telemetry"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	// Add platform modules.
	_ "github.com/titpetric/platform/module/autoload"
)

func main() {
	ctx := context.Background()
	err := start(ctx)
	telemetry.CaptureError(ctx, err)
}

var options = platform.NewOptions()

func start(ctx context.Context) error {
	// Register common middleware.
	platform.Use(middleware.Logger)
	platform.Use(telemetry.Middleware("platform"))

	p, err := platform.Start(ctx, options)
	if err != nil {
		return fmt.Errorf("exit error: %v", err)
	}
	p.Wait()
	return nil
}
