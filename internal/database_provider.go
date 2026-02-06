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
		return nil, err
	}
	return db, err
}

// Open is the same as sql.Open. It creates a client from a named connection.
func (r *DatabaseProvider) Open(_ context.Context, names ...string) (*sqlx.DB, error) {
	db, err := r.cached(r.open, names...)
	return db, err
}

// cached will return a singleton *db.DB from a named connection.
func (r *DatabaseProvider) cached(connector func(string, string) (*sqlx.DB, error), names ...string) (*sqlx.DB, error) {
	if len(names) == 0 {
		names = []string{"default"}
	}

	for _, name := range names {
		r.mu.Lock()
		db, ok := r.cache[name]
		r.mu.Unlock()
		if ok {
			return db, nil
		}
	}

	db, err := r.with(connector, names...)
	if err == nil {
		r.mu.Lock()
		r.cache[names[0]] = db
		r.mu.Unlock()
	}
	return db, err
}

// with will create a *db.DB given the connector (sqlx.Connect/Open).
func (r *DatabaseProvider) with(connector func(string, string) (*sqlx.DB, error), names ...string) (*sqlx.DB, error) {
	if len(names) == 0 {
		names = []string{"default"}
	}

	for _, name := range names {
		if value, ok := r.credentials[name]; ok {
			driver, dsn := r.parseCredential(value)
			client, err := connector(driver, dsn)
			if err != nil {
				return nil, err
			}

			opt, _ := databaseOptions[driver]
			opt.Apply(client)
			return client, nil
		}
	}
	return nil, fmt.Errorf("No configuration found for database: %v", names)
}

func (r *DatabaseProvider) parseCredential(credential string) (driver string, dsn string) {
	driver, dsn = "mysql", credential

	// allow specifying the driver with url notation,
	// in the follwing form: <driver>://<dsn>.
	if sepIndex := strings.Index(dsn, "://"); sepIndex != -1 {
		driver = dsn[:sepIndex]
		dsn = dsn[sepIndex+3:]
	}

	return driver, cleanDSN(dsn)
}
