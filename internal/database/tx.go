package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// WithTxFunc описание функции.
type WithTxFunc func(ctx context.Context, tx *sqlx.Tx) error

// WithTx функция реализует логику транзакции
func WithTx(ctx context.Context, db *sqlx.DB, fn WithTxFunc) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "db.BeginTxx()")
	}

	if err = fn(ctx, tx); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errors.Wrap(err, "Tx.Rollback")
		}

		return errors.Wrap(err, "Tx.WithTxFunc")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "Tx.Commit")
	}

	return nil
}
