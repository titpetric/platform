package internal

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

// DatabaseProvider holds a list of named sql connection credentials.
type DatabaseProvider struct {
	open func(string, string) (*sql.DB, error)

	mu          sync.Mutex
	cache       map[string]*sql.DB
	credentials map[string]string
}

// NewDatabaseProvider will allocate a valid `*DatabaseProvider` and return it.
func NewDatabaseProvider(open func(string, string) (*sql.DB, error)) *DatabaseProvider {
	return &DatabaseProvider{
		open:        open,
		cache:       make(map[string]*sql.DB),
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

// Register adds a named credential. It is not concurrency safe.
func (r *DatabaseProvider) Register(name string, config string) {
	r.credentials[name] = config
}

// Connect pings the connection before returning it. A failed ping is reported,
// and the handle stays cached because *sql.DB reconnects unaided.
func (r *DatabaseProvider) Connect(ctx context.Context, names ...string) (*sql.DB, error) {
	db, err := r.Open(ctx, names...)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

// Open is the same as sql.Open. It creates a client from a named connection.
func (r *DatabaseProvider) Open(_ context.Context, names ...string) (*sql.DB, error) {
	db, err := r.cached(r.open, names...)
	return db, err
}

// cached returns one *sql.DB per name, skipping a name with no credential.
func (r *DatabaseProvider) cached(connector func(string, string) (*sql.DB, error), names ...string) (*sql.DB, error) {
	if len(names) == 0 {
		names = []string{"default"}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, name := range names {
		if _, ok := r.credentials[name]; !ok {
			continue
		}

		if db, ok := r.cache[name]; ok {
			return db, nil
		}

		db, err := r.with(connector, name)
		if err != nil {
			return nil, err
		}

		r.cache[name] = db
		return db, nil
	}

	return nil, fmt.Errorf("no configuration found for database: %v", names)
}

// with creates a *sql.DB through the connector. The caller holds the lock.
func (r *DatabaseProvider) with(connector func(string, string) (*sql.DB, error), names ...string) (*sql.DB, error) {
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

			opt := databaseOption(driver, dsn)
			opt.Apply(client)
			return client, nil
		}
	}
	return nil, fmt.Errorf("no configuration found for database: %v", names)
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
