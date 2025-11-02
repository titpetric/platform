package storage

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform"
)

// DB will return the sqlx.DB in use for the user module.
// This enables reuse from outside without exposing implementation detail.
func DB(ctx context.Context) (*sqlx.DB, error) {
	return platform.Database.Connect(ctx)
}
