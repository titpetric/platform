//go:build integration
// +build integration

package storage

import (
	"database/sql"
	"testing"

	_ "github.com/titpetric/platform/drivers"

	"github.com/titpetric/platform/internal/require"
)

func TestNewSessionStorage_integration(t *testing.T) {
	db, err := DB(t.Context())
	require.NoError(t, err)

	s := NewSessionStorage(db)
	require.NotNil(t, s)

	{
		err := s.Delete(t.Context(), "non-existant")
		require.NoError(t, err)
	}

	{
		user, err := s.Get(t.Context(), "non-existant")
		require.Nil(t, user)
		require.ErrorIs(t, err, sql.ErrNoRows)
	}
}
