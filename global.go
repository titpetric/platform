package platform

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/telemetry"
)

// global is a value to prevent pollution of the global package namespace.
// for testing purposes the global state should be empty. The state is
// cloned as necessary into the *Platform object, on server creation.
var global = struct {
	registry *Registry
	db       *internal.DatabaseProvider
}{
	registry: &Registry{},
	db:       internal.NewDatabaseProvider(telemetry.Open),
}

// DatabaseProvider is the implementation interface for working with named connections.
// If no connection name is passed, the "default" connection will be used.
type DatabaseProvider interface {
	Open(...string) (*sqlx.DB, error)
	Connect(context.Context, ...string) (*sqlx.DB, error)
}

// Database is a holder of the database provider api in package namespace.
var Database DatabaseProvider = global.db

// Register will register a module in the platform global registry.
// It should not be relied upon in tests, keeping global state empty.
// This enables registering modules using blank imports.
func Register(m Module) {
	global.registry.Register(m)
}

// Use will add a middleware to the platform router.
// It should not be relied upon in tests, keeping global state empty.
// This should be used from main() to define any global middleware.
func Use(mw Middleware) {
	global.registry.Use(mw)
}
