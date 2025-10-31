package telemetry

import (
	"fmt"

	"github.com/XSAM/otelsql"
	"github.com/jmoiron/sqlx"
)

// Open creates a new instrumented *sqlx.DB without pinging.
// It wraps otelsql.Open and sqlx.NewDb for full compatibility.
func Open(driver, dsn string) (*sqlx.DB, error) {
	// Open an instrumented *sql.DB
	db, err := otelsql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("open instrumented db: %w", err)
	}

	return sqlx.NewDb(db, driver), nil
}

// Connect is like Open but verifies the connection (calls Ping).
func Connect(driver, dsn string) (*sqlx.DB, error) {
	db, err := Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
