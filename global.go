package platform

import (
	"context"
	"database/sql"

	chi "github.com/go-chi/chi/v5"

	"github.com/titpetric/platform/internal"
)

// global holds the package state New clones into each platform.
var global = struct {
	registry *Registry
	db       *internal.DatabaseProvider
}{
	registry: &Registry{},
	db:       internal.NewDatabaseProvider(internal.Open),
}

// Router is a local shim that aliases the chi router interface.
type Router = chi.Router

// DatabaseProvider opens named connections, "default" when no name is given.
type DatabaseProvider interface {
	// Open takes a context, a different implementation may use it for context awareness.
	// For example, Open may retrieve connection details from the database. The context
	// is used for tracing that operation, and the error returned may be the
	// result of a query issued against the database.
	Open(ctx context.Context, names ...string) (*sql.DB, error)

	// Connect has the semantics of Open + PingContext against the database,
	// verifying that the connection is live. An error is returned if
	// the storage is unreachable.
	//
	// A real database provider options may dictate additional behaviour.
	// For example, if a connection fails, it may retry a number of times,
	// it may back-off, it may continue to retry the connection in a
	// blocking way. It may also grab additional fail-over information from
	// keys based on the input parameters. An example of that would be to
	// define `name` and `name/failover` connections. If the connection to
	// the first fails, the second upstream is used.
	Connect(ctx context.Context, names ...string) (*sql.DB, error)
}

// Database is a holder of the database provider api in package namespace.
var Database DatabaseProvider = global.db

// Register adds a module to the global registry. Deprecated: use RegisterFunc.
func Register(m Module) {
	global.registry.register(m)
}

// RegisterFunc registers a module constructor that runs once per platform.
func RegisterFunc(f func() Module) {
	global.registry.RegisterFunc(f)
}

// Use adds a middleware to the platform router. Call it in main, not a test.
func Use(mw Middleware) {
	global.registry.Use(mw)
}
