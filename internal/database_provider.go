package internal

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/platform/module/telemetry"
)

type DatabaseProvider struct {
	cache       map[string]*sqlx.DB
	credentials map[string]string
}

func NewDatabaseProvider() *DatabaseProvider {
	return &DatabaseProvider{
		cache:       make(map[string]*sqlx.DB),
		credentials: make(map[string]string, 1),
	}
}

func (r *DatabaseProvider) Add(name string, config string) {
	r.credentials[name] = config
}

func (r *DatabaseProvider) Connect(names ...string) (*sqlx.DB, error) {
	db, err := r.cached(telemetry.Connect, names...)
	return db, err
}

func (r *DatabaseProvider) Open(names ...string) (*sqlx.DB, error) {
	db, err := r.cached(telemetry.Open, names...)
	return db, err
}

// cached will return a singleton *db.DB from a named connection.
func (r *DatabaseProvider) cached(connector func(string, string) (*sqlx.DB, error), names ...string) (*sqlx.DB, error) {
	if len(names) == 0 {
		names = []string{"default"}
	}
	for _, name := range names {
		db, ok := r.cache[name]
		if ok {
			return db, nil
		}
	}

	db, err := r.with(connector, names...)
	if err == nil {
		r.cache[names[0]] = db
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
			driver, dsn := "mysql", value

			// allow specifying the driver with url notation,
			// in the follwing form: <driver>://<dsn>.
			if sepIndex := strings.Index(dsn, "://"); sepIndex != -1 {
				driver = dsn[:sepIndex]
				dsn = dsn[sepIndex+3:]
			}

			// we should parse/add these options to the DSN as mandatory:
			// ?collation=utf8_general_ci&parseTime=true&loc=Local
			dsn += "?parseTime=true"
			return connector(driver, dsn)
		}
	}
	return nil, fmt.Errorf("No configuration found for database: %v", names)
}
