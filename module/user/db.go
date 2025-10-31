package user

import (
	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform"
)

// DB will return the sqlx.DB in use for the user module.
// This enables reuse from outside without exposing implementation detail.
func DB() (*sqlx.DB, error) {
	return platform.Database.Connect()
}
