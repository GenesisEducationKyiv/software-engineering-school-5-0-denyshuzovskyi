package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
)

func WithTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
