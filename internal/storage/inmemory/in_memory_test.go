package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
)

// Returns a non-nil instance of InMemory
func TestNewInMemoryReturnsNonNil(t *testing.T) {
	instance := NewInMemory()
	assert.NotNil(t, instance)
}

// Handles low memory conditions gracefully
func TestNewInMemoryHandlesLowMemory(t *testing.T) {
	// Simulate low memory condition by creating a large number of instances
	var instances []*InMemory
	for i := 0; i < 1000000; i++ {
		instance := NewInMemory()
		instances = append(instances, instance)
	}
	// Check if the last instance is still non-nil
	assert.NotNil(t, instances[len(instances)-1])
}

// Successfully stores a new shortening when the short URL does not exist
func TestPutStoresNewShortening(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemory()
	shortening := dto.Shortening{
		CorrelationID: "123",
		ShortURL:      "short123",
		OriginalURL:   "http://example.com",
		UserID:        "user1",
		IsDeleted:     false,
	}

	result, err := repo.Put(ctx, shortening)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ShortURL != shortening.ShortURL {
		t.Errorf("expected short URL %s, got %s", shortening.ShortURL, result.ShortURL)
	}
}

// Returns an error when the short URL already exists
func TestPutReturnsErrorWhenShortURLExists(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemory()
	shortening := dto.Shortening{
		CorrelationID: "123",
		ShortURL:      "short123",
		OriginalURL:   "http://example.com",
		UserID:        "user1",
		IsDeleted:     false,
	}

	_, err := repo.Put(ctx, shortening)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = repo.Put(ctx, shortening)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	expectedErr := "identifier already exists"
	if err.Error() != expectedErr {
		t.Errorf("expected error %s, got %s", expectedErr, err.Error())
	}
}

// Successfully stores a list of shortening objects in memory
func TestPutListStoresShortenings(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemory()
	list := dto.ShorteningList{
		{ShortURL: "short1", OriginalURL: "http://example.com/1", UserID: "user1"},
		{ShortURL: "short2", OriginalURL: "http://example.com/2", UserID: "user2"},
	}

	err := repo.PutList(ctx, list)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, shortening := range list {
		stored, err := repo.Get(ctx, shortening.ShortURL)
		if err != nil {
			t.Fatalf("expected to find %v, got error %v", shortening.ShortURL, err)
		}
		if stored.OriginalURL != shortening.OriginalURL {
			t.Errorf("expected original URL %v, got %v", shortening.OriginalURL, stored.OriginalURL)
		}
	}
}

// Returns an error if any shortening object in the list causes an error in the Put method
func TestPutListReturnsErrorOnDuplicate(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemory()
	list := dto.ShorteningList{
		{ShortURL: "short1", OriginalURL: "http://example.com/1", UserID: "user1"},
		{ShortURL: "short1", OriginalURL: "http://example.com/2", UserID: "user2"},
	}

	err := repo.PutList(ctx, list)
	if err == nil {
		t.Fatalf("expected an error due to duplicate short URL, got nil")
	}
}

// Retrieve existing shortening by identifier
func TestRetrieveExistingShorteningByIdentifier(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemory()
	shortening := dto.Shortening{
		ShortURL:    "short123",
		OriginalURL: "http://example.com",
		UserID:      "user1",
	}
	_, err := repo.Put(ctx, shortening)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result, err := repo.Get(ctx, "short123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ShortURL != "short123" {
		t.Errorf("expected short123, got %s", result.ShortURL)
	}
	if result.OriginalURL != "http://example.com" {
		t.Errorf("expected http://example.com, got %s", result.OriginalURL)
	}
}

// Identifier does not exist in storage
func TestIdentifierDoesNotExistInStorage(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemory()

	result, err := repo.Get(ctx, "nonexistent")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

// Retrieves the correct dto.Shortening when the originalURL exists in the map
func TestGetByURLReturnsCorrectShortening(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemory()
	shortening := dto.Shortening{
		CorrelationID: "123",
		ShortURL:      "short1",
		OriginalURL:   "http://example.com",
		UserID:        "user1",
		IsDeleted:     false,
	}
	repo.m.Store(shortening.ShortURL, shortening)

	result, err := repo.GetByURL(ctx, "http://example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.OriginalURL != "http://example.com" {
		t.Errorf("expected original URL to be 'http://example.com', got %v", result.OriginalURL)
	}
}

// Returns an empty dto.Shortening object when the originalURL does not exist in the map
func TestGetByURLReturnsEmptyShorteningWhenNotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemory()

	result, err := repo.GetByURL(ctx, "http://nonexistent.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.OriginalURL != "" {
		t.Errorf("expected original URL to be empty, got %v", result.OriginalURL)
	}
}

// Retrieves a list of URLs for a given user ID when URLs exist
//func TestGetListByUserWithExistingURLs(t *testing.T) {
//	ctx := context.Background()
//	repo := NewInMemory()
//
//	userID := "user123"
//	shortening1 := dto.Shortening{
//		CorrelationID: "1",
//		ShortURL:      "short1",
//		OriginalURL:   "http://example.com/1",
//		UserID:        userID,
//		IsDeleted:     false,
//	}
//	shortening2 := dto.Shortening{
//		CorrelationID: "2",
//		ShortURL:      "short2",
//		OriginalURL:   "http://example.com/2",
//		UserID:        userID,
//		IsDeleted:     false,
//	}
//
//	_, err := repo.Put(ctx, shortening1)
//	if err != nil {
//		t.Fatalf("expected no error, got %v", err)
//	}
//
//	_, err = repo.Put(ctx, shortening2)
//	if err != nil {
//		t.Fatalf("expected no error, got %v", err)
//	}
//
//	result, err := repo.GetListByUser(ctx, userID)
//
//	if err != nil {
//		t.Fatalf("expected no error, got %v", err)
//	}
//
//	if len(*result) != 2 {
//		t.Fatalf("expected 2 URLs, got %d", len(*result))
//	}
//
//	if (*result)[0].ShortURL != "short1" || (*result)[1].ShortURL != "short2" {
//		t.Fatalf("unexpected URLs in result")
//	}
//}

// Handles the case where the user ID is an empty string
func TestGetListByUserWithEmptyUserID(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemory()

	userID := ""

	result, err := repo.GetListByUser(ctx, userID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(*result) != 0 {
		t.Fatalf("expected 0 URLs, got %d", len(*result))
	}
}

// Deletes URLs marked for a specific user
func TestDeleteUserURLDeletesURLsForSpecificUser(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemory()

	shortening1 := dto.Shortening{
		ShortURL:    "short1",
		OriginalURL: "http://example.com/1",
		UserID:      "user1",
		IsDeleted:   false,
	}
	shortening2 := dto.Shortening{
		ShortURL:    "short2",
		OriginalURL: "http://example.com/2",
		UserID:      "user1",
		IsDeleted:   false,
	}

	_, err := repo.Put(ctx, shortening1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	_, err = repo.Put(ctx, shortening2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	repo.DeleteUserURL(ctx, []string{"short1", "short2"}, 1)

	result1, _ := repo.Get(ctx, "short1")
	result2, _ := repo.Get(ctx, "short2")

	if !result1.IsDeleted || !result2.IsDeleted {
		t.Errorf("URLs were not marked as deleted")
	}
}

// No URLs provided in the shortURL slice
func TestDeleteUserURLNoURLsProvided(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemory()

	shortening := dto.Shortening{
		ShortURL:    "short1",
		OriginalURL: "http://example.com/1",
		UserID:      "user1",
		IsDeleted:   false,
	}

	_, err := repo.Put(ctx, shortening)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	repo.DeleteUserURL(ctx, []string{}, 1)

	result, _ := repo.Get(ctx, "short1")

	if result.IsDeleted {
		t.Errorf("URL should not be marked as deleted when no URLs are provided")
	}
}
