//go:build integration
// +build integration

package storage

import (
	"database/sql"
	"testing"

	_ "github.com/titpetric/platform/drivers"

	"github.com/titpetric/platform/internal/require"
	"github.com/titpetric/platform/module/user/model"
)

func TestNewUserStorage_integration(t *testing.T) {
	db, err := DB(t.Context())
	require.NoError(t, err)

	s := NewUserStorage(db)
	require.NotNil(t, s)

	{
		user, err := s.Authenticate(t.Context(), model.UserAuth{})
		require.Nil(t, user)
		require.ErrorContains(t, err, "missing authentication info:")
	}

	{
		user, err := s.Authenticate(t.Context(), model.UserAuth{
			Email:    "me@titpetric.com",
			Password: "correct horse battery staple",
		})
		require.Nil(t, user)
		require.ErrorIs(t, err, sql.ErrNoRows)
	}
}
