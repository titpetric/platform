package internal

import (
	"database/sql"
	"fmt"
)

// Open creates a *sql.DB from a driver name and a DSN, uninstrumented.
func Open(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return db, nil
}
