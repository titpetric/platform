package main

import (
	"log"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/module"
	"github.com/titpetric/platform/registry"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

func main() {
	if err := start(); err != nil {
		log.Fatalf("exit error: %v", err)
	}
}

func start() error {
	registry.AddMiddleware(middleware.Logger)

	if err := module.LoadModules(); err != nil {
		return err
	}

	return platform.Start()
}
