package platform_test

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/internal/assert"
)

func TestDatabaseProvider_Connect(t *testing.T) {
	provider := internal.NewDatabaseProvider(sql.Open)
	provider.Register("test", "sqlite://:memory:")

	db, err := provider.Connect(t.Context(), "test")

	assert.NotNil(t, db)
	assert.NoError(t, err)

	db2, err := provider.Connect(t.Context(), "test")

	assert.NotNil(t, db2)
	assert.NoError(t, err)

	assert.Equal(t, db, db2)
}

func TestDatabaseProvider_Open(t *testing.T) {
	provider := internal.NewDatabaseProvider(func(string, string) (*sql.DB, error) {
		return nil, errors.New("test")
	})
	provider.Register("test", "sqlite://:memory:")

	db, err := provider.Open(t.Context(), "test")
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestDatabaseProviderFileSQLiteDefaults(t *testing.T) {
	provider := internal.NewDatabaseProvider(sql.Open)
	provider.Register("test", "sqlite://"+filepath.Join(t.TempDir(), "test.db"))

	db, err := provider.Connect(t.Context(), "test")
	assert.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })
	assert.Equal(t, 10, db.Stats().MaxOpenConnections)

	first, err := db.Conn(t.Context())
	assert.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, first.Close()) })
	second, err := db.Conn(t.Context())
	assert.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, second.Close()) })

	for _, connection := range []*sql.Conn{first, second} {
		var journalMode string
		assert.NoError(t, connection.QueryRowContext(t.Context(), "PRAGMA journal_mode").Scan(&journalMode))
		assert.Equal(t, "wal", journalMode)

		var busyTimeout int
		assert.NoError(t, connection.QueryRowContext(t.Context(), "PRAGMA busy_timeout").Scan(&busyTimeout))
		assert.Equal(t, 5000, busyTimeout)
	}
}
