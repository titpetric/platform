package main

import (
	"log"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/titpetric/platform"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	// Add telemetry with OpenTelemetry.
	"github.com/titpetric/platform/module/telemetry"

	// Add platform modules.
	_ "github.com/titpetric/platform/module/autoload"
)

func main() {
	// Register common middleware.
	platform.Use(middleware.Logger)
	platform.Use(telemetry.Middleware)

	if err := platform.Start(); err != nil {
		log.Fatalf("exit error: %v", err)
	}
}
