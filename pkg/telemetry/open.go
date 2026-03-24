package telemetry

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// sqlOpen is the low-level sql.Open function, replaced by initOpenTelemetry
// with an instrumented variant when OTEL is enabled.
var sqlOpen = sql.Open

// Open creates a *sqlx.DB using the active sql open function.
// When OpenTelemetry is enabled, the connection is instrumented automatically.
func Open(driver, dsn string) (*sqlx.DB, error) {
	db, err := sqlOpen(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return sqlx.NewDb(db, driver), nil
}
