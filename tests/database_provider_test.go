package platform_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/internal/assert"
)

func TestDatabaseProvider_Connect(t *testing.T) {
	provider := internal.NewDatabaseProvider(sqlx.Open)
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
	provider := internal.NewDatabaseProvider(func(string, string) (*sqlx.DB, error) {
		return nil, errors.New("test")
	})
	provider.Register("test", "sqlite://:memory:")

	db, err := provider.Open(t.Context(), "test")
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestDatabaseProviderFileSQLiteDefaults(t *testing.T) {
	provider := internal.NewDatabaseProvider(sqlx.Open)
	provider.Register("test", "sqlite://"+filepath.Join(t.TempDir(), "test.db"))

	db, err := provider.Connect(t.Context(), "test")
	assert.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })
	assert.Equal(t, 10, db.Stats().MaxOpenConnections)

	first, err := db.Connx(t.Context())
	assert.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, first.Close()) })
	second, err := db.Connx(t.Context())
	assert.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, second.Close()) })

	for _, connection := range []*sqlx.Conn{first, second} {
		var journalMode string
		assert.NoError(t, connection.GetContext(t.Context(), &journalMode, "PRAGMA journal_mode"))
		assert.Equal(t, "wal", journalMode)

		var busyTimeout int
		assert.NoError(t, connection.GetContext(t.Context(), &busyTimeout, "PRAGMA busy_timeout"))
		assert.Equal(t, 5000, busyTimeout)
	}
}
