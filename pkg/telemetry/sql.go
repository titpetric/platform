package telemetry

import (
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	"github.com/jmoiron/sqlx"
)

// sql_open gets overridden in init() if opentelemetry is enabled
var sql_open = sql.Open

// Open creates a new instrumented *sqlx.DB without pinging.
// It wraps otelsql.Open and sqlx.NewDb for full compatibility.
func Open(driver, dsn string) (*sqlx.DB, error) {
	// Open an instrumented *sql.DB
	db, err := otelsql.Open(driver, dsn, otelsql.WithDisableSkipErrMeasurement(true), otelsql.WithSpanOptions(otelsql.SpanOptions{DisableErrSkip: true}))
	if err != nil {
		return nil, fmt.Errorf("open instrumented db: %w", err)
	}

	return sqlx.NewDb(db, driver), nil
}
