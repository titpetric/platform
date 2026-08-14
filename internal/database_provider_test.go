package internal

import (
	"errors"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform/pkg/require"
)

func TestDatabaseProvider_Connect(t *testing.T) {
	provider := NewDatabaseProvider(sqlx.Open)
	provider.Register("test", "sqlite://:memory:")

	db, err := provider.Connect(t.Context(), "test")

	require.NotNil(t, db)
	require.NoError(t, err)

	db2, err := provider.Connect(t.Context(), "test")

	require.NotNil(t, db2)
	require.NoError(t, err)

	require.Equal(t, db, db2)
}

func TestDatabaseProvider_Open(t *testing.T) {
	provider := NewDatabaseProvider(func(string, string) (*sqlx.DB, error) {
		return nil, errors.New("test")
	})
	provider.Register("test", "sqlite://:memory:")

	db, err := provider.Open(t.Context(), "test")
	require.Error(t, err)
	require.Nil(t, db)
}

func TestDatabaseProviderFileSQLiteDefaults(t *testing.T) {
	provider := NewDatabaseProvider(sqlx.Open)
	provider.Register("test", "sqlite://"+filepath.Join(t.TempDir(), "test.db"))

	db, err := provider.Connect(t.Context(), "test")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.Equal(t, 10, db.Stats().MaxOpenConnections)

	first, err := db.Connx(t.Context())
	require.NoError(t, err)
	defer first.Close()
	second, err := db.Connx(t.Context())
	require.NoError(t, err)
	defer second.Close()

	for _, connection := range []*sqlx.Conn{first, second} {
		var journalMode string
		require.NoError(t, connection.GetContext(t.Context(), &journalMode, "PRAGMA journal_mode"))
		require.Equal(t, "wal", journalMode)

		var busyTimeout int
		require.NoError(t, connection.GetContext(t.Context(), &busyTimeout, "PRAGMA busy_timeout"))
		require.Equal(t, 5000, busyTimeout)
	}
}
