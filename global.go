package platform

import (
	"context"

	chi "github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/pkg/telemetry"
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

// Router is a local shim that aliases the chi router interface.
type Router = chi.Router

// DatabaseProvider is the implementation interface for working with named connections.
// If no connection name is passed, the "default" connection will be used.
// There's no assumption to how the database connection is provided, and
// the interface supports using external providers, enforces context awareness.
//
// The connection names, singleton behaviour, retries, fallback mechanisms,
// multiple-name logic and everything else to produce a *sql.DB is left to
// the implementation. The first party implementation uses the process environment
// and decodes `PLATFORM_DB_*` and uses the environment variable name for the key
// and the definition of the connection. It doesn't carry other logic like
// reconnecting.
//
// It's also likely that the first party database provider will be a public
// API package in the future. I'd like this to be swappable so people can
// bring in something like AWS secretsmanager or vault, or other.
type DatabaseProvider interface {
	// Open takes a context, a different implementation may use it for context awareness.
	// For example, Open may retrieve connection details from the database. The context
	// is used for tracing that operation, and the error returned may be the
	// result of a query issued against the database.
	Open(ctx context.Context, names ...string) (*sqlx.DB, error)

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
	Connect(ctx context.Context, names ...string) (*sqlx.DB, error)
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

// Start is a shorthand to create a new *Platform instance and
// immediately starts the server listener and handles requests.
func Start(ctx context.Context, options *Options) (*Platform, error) {
	svc := New(options)
	if err := svc.Start(ctx); err != nil {
		return nil, err
	}
	return svc, nil
}
