package main

import (
	"log"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/titpetric/platform"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	// Add apmsql wrappers for tracing with build tag `telemetry`.
	_ "github.com/titpetric/platform/internal/tracing"

	// Add platform modules.
	_ "github.com/titpetric/platform/module/autoload"

	// Add instrumentation middleware.
	"go.elastic.co/apm/module/apmchiv5/v2"
)

func main() {
	// Register common middleware.
	platform.Use(middleware.Logger)
	platform.Use(apmchiv5.Middleware())

	if err := platform.Start(); err != nil {
		log.Fatalf("exit error: %v", err)
	}
}
