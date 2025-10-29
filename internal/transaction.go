package internal

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Transaction wraps a function in a transaction.
// If the function returns an error, the transaction is rolled back.
// If the function returns nil, the transaction is committed.
func Transaction(db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Ensure rollback if function fails or panics; Commit will override if successful
	defer func() {
		_ = tx.Rollback()
	}()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
