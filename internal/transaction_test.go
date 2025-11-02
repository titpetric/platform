package internal

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform/internal/require"
)

func TestTransaction(t *testing.T) {
	provider := NewDatabaseProvider(sqlx.Open)
	provider.Register("test", "sqlite://:memory:")

	db, err := provider.Connect(t.Context(), "test")

	require.NotNil(t, db)
	require.NoError(t, err)

	err = Transaction(t.Context(), db, func(context.Context, *sqlx.Tx) error {
		return sql.ErrNoRows
	})
	require.ErrorIs(t, err, sql.ErrNoRows)
}
