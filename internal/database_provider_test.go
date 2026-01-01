package internal

import (
	"errors"
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
