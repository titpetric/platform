package internal

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
)

// DatabaseProvider holds a list of named sql connection credentials.
type DatabaseProvider struct {
	open func(string, string) (*sqlx.DB, error)

	mu          sync.Mutex
	cache       map[string]*sqlx.DB
	credentials map[string]string
}

// NewDatabaseProvider will allocate a valid `*DatabaseProvider` and return it.
func NewDatabaseProvider(open func(string, string) (*sqlx.DB, error)) *DatabaseProvider {
	return &DatabaseProvider{
		open:        open,
		cache:       make(map[string]*sqlx.DB),
		credentials: make(map[string]string, 1),
	}
}

// List will return the list of credential names.
func (r *DatabaseProvider) List() []string {
	result := make([]string, 0, len(r.credentials))
	for k := range r.credentials {
		result = append(result, k)
	}
	return result
}

// Register will add a new named credential into the provider.
// The function is not concurrency safe, database credentials
// can't be changed during the lifetime of the provider.
func (r *DatabaseProvider) Register(name string, config string) {
	r.credentials[name] = config
}

// Connect issues a PingContext to verify a live connection before returning.
// The context is used to propagate tracing detail so ping is grouped correctly.
func (r *DatabaseProvider) Connect(ctx context.Context, names ...string) (*sqlx.DB, error) {
	db, err := r.Open(ctx, names...)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		r.evict(db)
		return nil, err
	}
	return db, nil
}

// Open is the same as sql.Open. It creates a client from a named connection.
func (r *DatabaseProvider) Open(_ context.Context, names ...string) (*sqlx.DB, error) {
	db, err := r.cached(r.open, names...)
	return db, err
}

// cached will return a singleton *db.DB from a named connection.
//
// The lock is held across the build. Releasing it between the miss and the
// open let every concurrent caller build a connection of its own: the first to
// finish was cached, the rest were handed out and then unreachable, never
// closed, each with a pool of its own. On the default sqlite://:memory: DSN
// they are not even the same database.
//
// The connector is sql.Open, which does not dial, so the critical section is
// short. A custom open function passed to NewDatabaseProvider serialises with
// every other first-time open.
func (r *DatabaseProvider) cached(connector func(string, string) (*sqlx.DB, error), names ...string) (*sqlx.DB, error) {
	if len(names) == 0 {
		names = []string{"default"}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, name := range names {
		if db, ok := r.cache[name]; ok {
			return db, nil
		}
	}

	name, db, err := r.with(connector, names...)
	if err != nil {
		return nil, err
	}

	// Keyed on the name that resolved, not on names[0]. Keying on the first
	// name requested made Open(ctx, "replica", "default") cache the default
	// connection as "replica", so a later Open(ctx, "replica") succeeded for a
	// name that was never registered, and Open(ctx, "default") missed and
	// built a second handle to the same database.
	r.cache[name] = db

	return db, nil
}

// with will create a *db.DB given the connector (sqlx.Connect/Open). It
// returns the name whose credential was used: names is a fallback list, and
// the first entry that has one wins. The caller holds the lock.
func (r *DatabaseProvider) with(connector func(string, string) (*sqlx.DB, error), names ...string) (string, *sqlx.DB, error) {
	if len(names) == 0 {
		names = []string{"default"}
	}

	for _, name := range names {
		if value, ok := r.credentials[name]; ok {
			driver, dsn := r.parseCredential(value)
			client, err := connector(driver, dsn)
			if err != nil {
				return "", nil, err
			}

			opt := databaseOption(driver, dsn)
			opt.Apply(client)
			return name, client, nil
		}
	}
	return "", nil, fmt.Errorf("no configuration found for database: %v", names)
}

// evict drops a handle from the cache and closes it. Without it a Connect
// whose ping failed left the handle cached, and every later Open returned that
// dead connection with a nil error.
func (r *DatabaseProvider) evict(db *sqlx.DB) {
	r.mu.Lock()
	for name, cached := range r.cache {
		if cached == db {
			delete(r.cache, name)
		}
	}
	r.mu.Unlock()

	_ = db.Close()
}

func (r *DatabaseProvider) parseCredential(credential string) (driver string, dsn string) {
	driver, dsn = "mysql", credential

	// allow specifying the driver with url notation,
	// in the follwing form: <driver>://<dsn>.
	if sepIndex := strings.Index(dsn, "://"); sepIndex != -1 {
		driver = dsn[:sepIndex]
		dsn = dsn[sepIndex+3:]
		if driver == "postgres" || driver == "postgresql" {
			driver = "pgx"
			dsn = "postgres://" + dsn
		}
	}

	return driver, cleanDSN(driver, dsn)
}
