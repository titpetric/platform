package main

import (
	"log"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/registry"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	_ "github.com/titpetric/platform/module/autoload"
)

func main() {
	// Register common middleware.
	registry.AddMiddleware(middleware.Logger)

	if err := platform.Start(); err != nil {
		log.Fatalf("exit error: %v", err)
	}
}
