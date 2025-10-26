package platform

import (
	"github.com/jmoiron/sqlx"
	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/internal/repository"
)

// global is a value to prevent pollution of the global package namespace.
// for testing purposes the global state should be empty. The state is
// cloned as necessary into the *Platform object, on server creation.
var global = struct {
	registry *Registry
	db       DatabaseProvider
}{
	registry: &Registry{},
	db:       internal.NewDatabaseProvider(),
}

// DatabaseProvider is the implementation interface for working with named connections.
// If no connection name is passed, the "default" connection will be used.
type DatabaseProvider interface {
	Add(string, string)
	Open(...string) (*sqlx.DB, error)
	Connect(...string) (*sqlx.DB, error)
}

// Database is a holder of the database provider api in package namespace.
// It's intended to be used as `platform.Database.Connect/Open(name string)` to
// get a live connection, or an error if one occured.
var Database = global.db

var (
	// AddModule provides a `package.AddModule` api for the global registry.
	AddModule = global.registry.AddModule
	// AddMiddleware provides a `package.AddMiddleware` api for the global registry.
	AddMiddleware = global.registry.AddMiddleware
)

type (
	// Plugin represents a decomposed `Module`.
	Plugin = repository.Plugin
	// Module is an implementation interface for a plugin.
	Module = repository.Module
	// Registry holds state of middlewares and plugins.
	Registry = repository.Registry
	// Middleware is the function signature for middlewares.
	Middleware = repository.Middleware
	// Router is the router interface (chi.Router).
	Router = repository.Router
)
