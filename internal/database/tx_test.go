package database

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// Error occurs when beginning transaction
func TestErrorOccursWhenBeginningTransaction(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin().WillReturnError(errors.New("begin error"))

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	fn := func(ctx context.Context, tx *sqlx.Tx) error {
		return nil
	}

	err = WithTx(ctx, sqlxDB, fn)
	if err == nil || err.Error() != "db.BeginTxx(): begin error" {
		t.Fatalf("expected begin error, got %v", err)
	}
}
