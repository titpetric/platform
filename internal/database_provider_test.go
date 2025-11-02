package internal

import (
	"testing"

	_ "modernc.org/sqlite"

	"github.com/titpetric/platform/internal/require"
)

func TestDatabaseProvider(t *testing.T) {
	provider := NewDatabaseProvider()
	provider.Register("test", "sqlite://:memory:")

	db, err := provider.Connect(t.Context(), "test")

	require.NotNil(t, db)
	require.NoError(t, err)

	db2, err := provider.Connect(t.Context(), "test")

	require.NotNil(t, db2)
	require.NoError(t, err)

	require.Equal(t, db, db2)
}
