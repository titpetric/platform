package platform

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform/pkg/telemetry"
)

// Transaction wraps a function in a transaction.
// If the function returns an error, the transaction is rolled back.
// If the function returns nil, the transaction is committed.
func Transaction(ctx context.Context, db *sqlx.DB, fn func(context.Context, *sqlx.Tx) error) error {
	ctx, span := telemetry.Start(ctx, "db.Transaction")
	defer span.End()

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Ensure rollback if function fails or panics; Commit will override if successful
	defer func() {
		_ = tx.Rollback()
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
