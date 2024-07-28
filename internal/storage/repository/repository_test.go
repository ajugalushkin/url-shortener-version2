package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
)

// Repository is successfully created with a valid database connection
func TestNewRepo_SuccessfullyCreatesRepositoryWithValidDBConnection(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxDB)

	if repo == nil {
		t.Error("expected repository to be created, got nil")
	}
}

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
	if err == nil || !strings.Contains(err.Error(), "connection error") {
		t.Errorf("expected connection error, got %v", err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Retrieves a shortening record successfully when the shortURL exists in the database
func TestGet_SuccessfullyRetrievesShorteningRecord(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := &Repo{db: sqlxDB}

	rows := sqlmock.NewRows([]string{"short_url", "correlation_id", "original_url", "user_id", "is_deleted"}).
		AddRow("short123", "corr123", "http://example.com", "user123", false)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT short_url, correlation_id, original_url, user_id, is_deleted FROM shorten_urls WHERE short_url IN ($1)`).
		WithArgs("short123").
		WillReturnRows(rows)
	mock.ExpectCommit()

	shortening, err := repo.Get(ctx, "short123")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if shortening.ShortURL != "short123" {
		t.Errorf("expected shortURL to be 'short123', got %s", shortening.ShortURL)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Handles the case where the shortURL does not exist in the database
func TestGetShortURLNotExists(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := &Repo{db: sqlxDB}

	rows := sqlmock.NewRows([]string{"short_url", "correlation_id", "original_url", "user_id", "is_deleted"})

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT short_url, correlation_id, original_url, user_id, is_deleted FROM shorten_urls WHERE short_url IN ($1)`).
		WithArgs("nonExistentShortURL").
		WillReturnRows(rows)
	mock.ExpectCommit()

	_, err = repo.Get(ctx, "nonExistentShortURL")
	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected error to be sql.ErrNoRows, got %v", err)
	}
}

// Retrieves a shortening record when the original URL exists in the database
func TestGetByURLReturnsShorteningRecord(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := &Repo{db: sqlxDB}

	rows := sqlmock.NewRows([]string{"short_url", "correlation_id", "original_url", "user_id"}).
		AddRow("short123", "corr123", "http://example.com", "user123")

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT short_url, correlation_id, original_url, user_id FROM shorten_urls WHERE original_url IN ($1)`).
		WithArgs("http://example.com").
		WillReturnRows(rows)
	mock.ExpectCommit()

	shortening, err := repo.GetByURL(ctx, "http://example.com")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if shortening.OriginalURL != "http://example.com" {
		t.Errorf("expected original URL to be 'http://example.com', got %s", shortening.OriginalURL)
	}
}

// Returns sql.ErrNoRows when the original URL does not exist in the database
func TestGetByURLReturnsNoRowsError(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := &Repo{db: sqlxDB}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT short_url, correlation_id, original_url, user_id FROM shorten_urls WHERE original_url IN ($1)`).
		WithArgs("http://nonexistent.com").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	shortening, err := repo.GetByURL(ctx, "http://nonexistent.com")
	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows error, got %v", err)
	}

	if shortening != nil {
		t.Errorf("expected shortening to be nil, got %v", shortening)
	}
}

// Retrieves a list of URL shortenings for a given user ID
func TestGetListByUser_Success(t *testing.T) {
	ctx := context.Background()
	userID := "test_user"
	expectedList := &dto.ShorteningList{
		{CorrelationID: "1", ShortURL: "http://short.url/1", OriginalURL: "http://original.url/1", UserID: userID, IsDeleted: false},
		{CorrelationID: "2", ShortURL: "http://short.url/2", OriginalURL: "http://original.url/2", UserID: userID, IsDeleted: false},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := &Repo{db: sqlxDB}

	rows := sqlmock.NewRows([]string{"short_url", "correlation_id", "original_url", "user_id"}).
		AddRow("http://short.url/1", "1", "http://original.url/1", userID).
		AddRow("http://short.url/2", "2", "http://original.url/2", userID)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT short_url, correlation_id, original_url, user_id FROM shorten_urls WHERE user_id IN ($1)").
		WithArgs(userID).WillReturnRows(rows)
	mock.ExpectCommit()

	result, err := repo.GetListByUser(ctx, userID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expectedList) {
		t.Errorf("expected %v, got %v", expectedList, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Handles the case where the user ID does not exist in the database
func TestGetListByUser_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userID := "non_existent_user"

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := &Repo{db: sqlxDB}

	rows := sqlmock.NewRows([]string{"short_url", "correlation_id", "original_url", "user_id"})

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT short_url, correlation_id, original_url, user_id FROM shorten_urls WHERE user_id IN ($1)").
		WithArgs(userID).WillReturnRows(rows)
	mock.ExpectCommit()

	result, err := repo.GetListByUser(ctx, userID)
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected error %v, got %v", sql.ErrNoRows, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Handles database connection failure gracefully
func TestDeleteUserURLHandlesDBConnectionFailure(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := &Repo{db: sqlxDB}

	shortList := []string{"short1", "short2"}
	userID := 1

	mock.ExpectBegin().WillReturnError(fmt.Errorf("db connection error"))

	repo.DeleteUserURL(ctx, shortList, userID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Channel receives all input data when doneCh is not closed
func TestChannelReceivesAllDataWhenDoneChNotClosed(t *testing.T) {
	doneCh := make(chan struct{})
	input := []string{"data1", "data2", "data3"}

	resultCh := prepareList(doneCh, input)

	var result []string
	for data := range resultCh {
		result = append(result, data)
	}

	if !reflect.DeepEqual(result, input) {
		t.Errorf("expected %v, got %v", input, result)
	}
}

// doneCh is closed before any data is processed
func TestDoneChClosedBeforeProcessing(t *testing.T) {
	doneCh := make(chan struct{})
	input := []string{"data1", "data2", "data3"}

	close(doneCh)
	resultCh := prepareList(doneCh, input)

	var result []string
	for data := range resultCh {
		result = append(result, data)
	}

	if len(result) != 0 {
		t.Errorf("expected no data, got %v", result)
	}
}

// Handles empty inputCh gracefully
func TestSearchURLsHandlesEmptyInput(t *testing.T) {
	ctx := context.Background()
	doneCh := make(chan struct{})
	inputCh := make(chan string)
	defer close(doneCh)
	defer close(inputCh)

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := &Repo{db: sqlxDB}

	resultCh := repo.searchURLs(ctx, doneCh, inputCh)

	select {
	case result, ok := <-resultCh:
		if ok {
			t.Errorf("expected no results, but got %v", result)
		}
	default:
		// No results as expected
	}
}

// Correctly splits input channel into 100 worker channels
func TestSplitCorrectlySplitsInto100WorkerChannels(t *testing.T) {
	ctx := context.Background()
	doneCh := make(chan struct{})
	inputCh := make(chan string)
	defer close(doneCh)
	defer close(inputCh)

	db, _, _ := sqlmock.New()
	defer db.Close()
	repo := &Repo{db: sqlx.NewDb(db, "sqlmock")}

	channels := repo.split(ctx, doneCh, inputCh)

	if len(channels) != 100 {
		t.Errorf("expected 100 channels, got %d", len(channels))
	}
}

// Handles empty input channel gracefully
func TestSplitHandlesEmptyInputChannelGracefully(t *testing.T) {
	ctx := context.Background()
	doneCh := make(chan struct{})
	inputCh := make(chan string)
	defer close(doneCh)
	defer close(inputCh)

	db, _, _ := sqlmock.New()
	defer db.Close()
	repo := &Repo{db: sqlx.NewDb(db, "sqlmock")}

	channels := repo.split(ctx, doneCh, inputCh)

	for _, ch := range channels {
		select {
		case res, ok := <-ch:
			if ok {
				t.Errorf("expected channel to be closed, but received value: %v", res)
			}
		default:
		}
	}
}

// Merging multiple channels into a single channel without data loss
func TestMergeMultipleChannels(t *testing.T) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	ch1 := make(chan *dto.Shortening)
	ch2 := make(chan *dto.Shortening)

	go func() {
		ch1 <- &dto.Shortening{ShortURL: "short1", OriginalURL: "original1"}
		close(ch1)
	}()

	go func() {
		ch2 <- &dto.Shortening{ShortURL: "short2", OriginalURL: "original2"}
		close(ch2)
	}()

	finalCh := merge(doneCh, ch1, ch2)

	var results []*dto.Shortening
	for result := range finalCh {
		results = append(results, result)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0].ShortURL != "short1" || results[1].ShortURL != "short2" {
		t.Fatalf("unexpected results: %+v", results)
	}
}

// Handling an empty list of input channels
func TestMergeEmptyChannels(t *testing.T) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	finalCh := merge(doneCh)

	var results []*dto.Shortening
	for result := range finalCh {
		results = append(results, result)
	}

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
