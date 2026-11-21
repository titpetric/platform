package internal

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

func TestSQLiteDatabaseOption(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		maxOpen int
		maxIdle int
	}{
		{name: "file database", dsn: "app.db", maxOpen: 10, maxIdle: 2},
		{name: "memory database", dsn: ":memory:", maxOpen: 1, maxIdle: 1},
		{name: "memory URI", dsn: "file::memory:?cache=shared", maxOpen: 1, maxIdle: 1},
		{name: "named shared memory database", dsn: "file:app?mode=memory&cache=shared", maxOpen: 1, maxIdle: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := databaseOption("sqlite", tt.dsn)
			assert.Equal(t, tt.maxOpen, option.MaxOpenConns)
			assert.Equal(t, tt.maxIdle, option.MaxIdleConns)
		})
	}
}

// TestDatabaseOptionUnknownDriver covers the pool a driver with no entry of its
// own gets. Apply calls the setters unconditionally, and the zero value would
// cap the pool at unlimited open and 0 idle connections.
func TestDatabaseOptionUnknownDriver(t *testing.T) {
	for _, driver := range []string{"pgx", "postgres", "sqlite3", "clickhouse", ""} {
		t.Run(driver, func(t *testing.T) {
			option := databaseOption(driver, "")

			assert.Greater(t, option.MaxOpenConns, 0, "an unknown driver must get an open limit")
			assert.Greater(t, option.MaxIdleConns, 0, "an unknown driver must retain idle connections")
		})
	}
}

// TestDatabaseOptionKnownDrivers pins the entries that do exist, and the
// fallback above cannot quietly start applying to them.
func TestDatabaseOptionKnownDrivers(t *testing.T) {
	mysql := databaseOption("mysql", "user@tcp(h)/db")
	assert.Equal(t, 10, mysql.MaxOpenConns)
	assert.Equal(t, 10, mysql.MaxIdleConns)

	file := databaseOption("sqlite", "/tmp/app.db")
	assert.Equal(t, 10, file.MaxOpenConns)
	assert.Equal(t, 2, file.MaxIdleConns)

	// One connection, so all queries reach the same in-memory database.
	memory := databaseOption("sqlite", ":memory:")
	assert.Equal(t, 1, memory.MaxOpenConns)
	assert.Equal(t, 1, memory.MaxIdleConns)
}

// TestDatabaseOptionApply covers the pool settings a DatabaseOption writes and
// the nil receiver Apply tolerates.
func TestDatabaseOptionApply(t *testing.T) {
	option := &DatabaseOption{MaxOpenConns: 7, MaxIdleConns: 3}

	db := sql.OpenDB(nopConnector{})
	t.Cleanup(func() { _ = db.Close() })

	option.Apply(db)
	assert.Equal(t, 7, db.Stats().MaxOpenConnections)

	// An unknown driver path can hand over a nil option.
	var nilOption *DatabaseOption
	nilOption.Apply(db)
	assert.Equal(t, 7, db.Stats().MaxOpenConnections)
}

type nopConnector struct{}

func (nopConnector) Connect(context.Context) (driver.Conn, error) { return nopConn{}, nil }
func (nopConnector) Driver() driver.Driver                        { return nopDriver{} }

type nopDriver struct{}

func (nopDriver) Open(string) (driver.Conn, error) { return nopConn{}, nil }

type nopConn struct{}

func (nopConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("not implemented") }
func (nopConn) Close() error                        { return nil }
func (nopConn) Begin() (driver.Tx, error)           { return nil, errors.New("not implemented") }
