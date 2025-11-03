package platform

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform/internal"
	_ "github.com/titpetric/platform/pkg/drivers"
	"github.com/titpetric/platform/pkg/require"
)

func TestTransaction(t *testing.T) {
	provider := internal.NewDatabaseProvider(sqlx.Open)
	provider.Register("test", "sqlite://:memory:")

	db, err := provider.Connect(t.Context(), "test")
	require.NoError(t, err)
	require.NotNil(t, db)

	err = Transaction(t.Context(), db, func(context.Context, *sqlx.Tx) error {
		return sql.ErrNoRows
	})
	require.ErrorIs(t, err, sql.ErrNoRows)
}
