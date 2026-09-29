package platform_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/internal/assert"
)

func TestTransaction(t *testing.T) {
	provider := internal.NewDatabaseProvider(sqlx.Open)
	provider.Register("test", "sqlite://:memory:")

	db, err := provider.Connect(t.Context(), "test")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = platform.Transaction(t.Context(), db, func(context.Context, *sqlx.Tx) error {
		return sql.ErrNoRows
	})
	assert.ErrorIs(t, err, sql.ErrNoRows)
}
