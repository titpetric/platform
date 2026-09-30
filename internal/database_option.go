package internal

import (
	"github.com/jmoiron/sqlx"
)

// DatabaseOption configures database connection pooling settings.
type DatabaseOption struct {
	MaxOpenConns int
	MaxIdleConns int
}

// Apply applies the database option settings to a database connection.
func (o *DatabaseOption) Apply(client *sqlx.DB) {
	if o == nil {
		return
	}
	client.SetMaxOpenConns(o.MaxOpenConns)
	client.SetMaxIdleConns(o.MaxIdleConns)
}

var databaseOptions = map[string]DatabaseOption{
	"sqlite": {
		MaxOpenConns: 10,
		MaxIdleConns: 2,
	},
	"mysql": {
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	},
}

// defaultDatabaseOption is what a driver with no entry of its own gets. The
// zero value is not usable: Apply would call SetMaxIdleConns(0), which retains
// no idle connection at all, so every query on an unnamed driver would dial a
// new one. That is worse than the database/sql default, and pgx has no entry.
var defaultDatabaseOption = DatabaseOption{
	MaxOpenConns: 10,
	MaxIdleConns: 2,
}

func databaseOption(driver, dsn string) DatabaseOption {
	if driver == "sqlite" && isSQLiteMemoryDSN(dsn) {
		return DatabaseOption{
			MaxOpenConns: 1,
			MaxIdleConns: 1,
		}
	}

	if option, ok := databaseOptions[driver]; ok {
		return option
	}

	return defaultDatabaseOption
}
