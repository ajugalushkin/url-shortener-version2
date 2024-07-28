package repository

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
)

// Successfully inserts a new shortening record into the database
func TestPut_SuccessfullyInsertsNewShortening(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := &Repo{db: sqlxDB}

	shorteningInput := dto.Shortening{
		CorrelationID: "123",
		ShortURL:      "short.ly/abc",
		OriginalURL:   "http://example.com",
		UserID:        "user1",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO shorten_urls").
		WithArgs(shorteningInput.ShortURL, shorteningInput.CorrelationID, shorteningInput.OriginalURL, shorteningInput.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := repo.Put(ctx, shorteningInput)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if result.OriginalURL != shorteningInput.OriginalURL {
		t.Errorf("expected %s, got %s", shorteningInput.OriginalURL, result.OriginalURL)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Handles database connection failures gracefully
func TestPut_HandlesDatabaseConnectionFailures(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := &Repo{db: sqlxDB}

	shorteningInput := dto.Shortening{
		CorrelationID: "123",
		ShortURL:      "short.ly/abc",
		OriginalURL:   "http://example.com",
		UserID:        "user1",
	}

	mock.ExpectBegin().WillReturnError(fmt.Errorf("connection error"))

	result, err := repo.Put(ctx, shorteningInput)
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if err == nil || !strings.Contains(err.Error(), "connection error") {
		t.Errorf("expected connection error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
