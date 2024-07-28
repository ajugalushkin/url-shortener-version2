package file

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
)

// Path does not exist
func TestNewStorageWithNonExistentPath(t *testing.T) {
	path := "testdata/non_existent_path.json"
	storage := NewStorage(path)

	if storage == nil {
		t.Fatalf("Expected storage to be initialized, got nil")
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("Expected file to not exist at path %s, but it does", path)
	}
}

// Successfully store a new shortening when the identifier does not exist
func TestPutNewShortening(t *testing.T) {
	ctx := context.Background()
	config.GetConfig().FileStoragePath = "/tmp/test_storage.json"

	if _, err := os.Stat(config.GetConfig().FileStoragePath); !errors.Is(err, os.ErrNotExist) {
		file, err := os.ReadFile(config.GetConfig().FileStoragePath)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		} else {
			if len(file) > 0 {
				os.Remove(config.GetConfig().FileStoragePath)
			}
		}
	}

	storage := NewStorage(config.GetConfig().FileStoragePath)
	shortening := dto.Shortening{
		ShortURL:    "short1",
		OriginalURL: "http://example.com",
		UserID:      "user1",
	}

	result, err := storage.Put(ctx, shortening)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ShortURL != shortening.ShortURL || result.OriginalURL != shortening.OriginalURL {
		t.Fatalf("expected %v, got %v", shortening, result)
	}
}

// Attempt to store a shortening with an identifier that already exists
func TestPutExistingShortening(t *testing.T) {
	ctx := context.Background()
	storage := NewStorage("/tmp/test_storage.json")
	shortening := dto.Shortening{
		ShortURL:    "short1",
		OriginalURL: "http://example.com",
		UserID:      "user1",
	}

	_, _ = storage.Put(ctx, shortening)
	_, err := storage.Put(ctx, shortening)

	if err == nil || err.Error() != "identifier already exists" {
		t.Fatalf("expected error 'identifier already exists', got %v", err)
	}
}

// Successfully saves a list of shortening objects
func TestPutList_Success(t *testing.T) {
	ctx := context.Background()
	config.GetConfig().FileStoragePath = "/tmp/test_storage.json"

	file, err := os.ReadFile(config.GetConfig().FileStoragePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	} else {
		if len(file) > 0 {
			os.Remove(config.GetConfig().FileStoragePath)
		}
	}

	storage := NewStorage(config.GetConfig().FileStoragePath)
	list := dto.ShorteningList{
		{ShortURL: "short1", OriginalURL: "http://example.com/1"},
		{ShortURL: "short2", OriginalURL: "http://example.com/2"},
	}

	err = storage.PutList(ctx, list)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, shortening := range list {
		_, err := storage.Get(ctx, shortening.ShortURL)
		if err != nil {
			t.Fatalf("expected to find %v, got error %v", shortening.ShortURL, err)
		}
	}
}

// Returns an error if any shortening object in the list causes the Put method to fail
func TestPutList_ErrorOnDuplicate(t *testing.T) {
	ctx := context.Background()
	config.GetConfig().FileStoragePath = "/tmp/test_storage.json"
	storage := NewStorage(config.GetConfig().FileStoragePath)
	list := dto.ShorteningList{
		{ShortURL: "short1", OriginalURL: "http://example.com/1"},
		{ShortURL: "short1", OriginalURL: "http://example.com/2"},
	}

	err := storage.PutList(ctx, list)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "identifier already exists" {
		t.Fatalf("expected 'identifier already exists' error, got %v", err)
	}
}

// Retrieve an existing shortening by its identifier
func TestRetrieveExistingShortening(t *testing.T) {
	ctx := context.Background()
	config.GetConfig().FileStoragePath = "testdata/testfile.json"
	storage := NewStorage(config.GetConfig().FileStoragePath)

	if _, err := os.Stat(config.GetConfig().FileStoragePath); !errors.Is(err, os.ErrNotExist) {
		file, err := os.ReadFile(config.GetConfig().FileStoragePath)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		} else {
			if len(file) > 0 {
				os.Remove(config.GetConfig().FileStoragePath)
			}
		}
	}

	shortening := dto.Shortening{
		ShortURL:    "short123",
		OriginalURL: "http://example.com",
		UserID:      "user1",
	}

	_, err := storage.Put(ctx, shortening)
	if err != nil {
		t.Fatalf("not expected error, got %v", err)
	}

	result, err := storage.Get(ctx, "short123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ShortURL != "short123" || result.OriginalURL != "http://example.com" {
		t.Fatalf("expected %v, got %v", shortening, result)
	}
}

// Return an error when the identifier does not exist
func TestRetrieveNonExistingShortening(t *testing.T) {
	ctx := context.Background()
	storage := NewStorage("testdata/testfile.json")

	_, err := storage.Get(ctx, "nonexistent")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "not found" {
		t.Fatalf("expected 'not found' error, got %v", err)
	}
}

// Retrieves an empty list when no URLs are associated with the user
func TestGetListByUserReturnsEmptyList(t *testing.T) {
	ctx := context.Background()
	userID := "nonexistent_user"
	storage := NewStorage("testdata/testfile.json")

	result, err := storage.GetListByUser(ctx, userID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(*result) != 0 {
		t.Fatalf("expected empty list, got %v", result)
	}
} // Handles non-existent user IDs without crashing
func TestGetListByUserHandlesNonExistentUserID(t *testing.T) {
	ctx := context.Background()
	userID := "nonexistent_user"
	storage := NewStorage("testdata/testfile.json")

	result, err := storage.GetListByUser(ctx, userID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatalf("expected non-nil result, got nil")
	}
}

// Successfully deletes a list of URLs for a given user
//func TestDeleteUserURLSuccessfully(t *testing.T) {
//	ctx := context.Background()
//	config.GetConfig().FileStoragePath = "testdata/testfile.json"
//
//	file, err := os.ReadFile(config.GetConfig().FileStoragePath)
//	if err != nil {
//		t.Fatalf("expected no error, got %v", err)
//	} else {
//		if len(file) > 0 {
//			os.Remove(config.GetConfig().FileStoragePath)
//		}
//	}
//
//	storage := NewStorage(config.GetConfig().FileStoragePath)
//	userID := 1
//	shortURLs := []string{"short1", "short2"}
//
//	// Prepopulate storage with URLs
//	storage.Put(ctx, dto.Shortening{ShortURL: "short1", OriginalURL: "http://example.com/1"})
//	storage.Put(ctx, dto.Shortening{ShortURL: "short2", OriginalURL: "http://example.com/2"})
//
//	// Call DeleteUserURL
//	storage.DeleteUserURL(ctx, shortURLs, userID)
//
//	// Verify URLs are deleted
//	for _, url := range shortURLs {
//		_, err := storage.Get(ctx, url)
//		if err == nil {
//			t.Errorf("expected URL %s to be deleted", url)
//		}
//	}
//}

// User ID does not exist
func TestDeleteUserURLUserIDNotExist(t *testing.T) {
	ctx := context.Background()
	config.GetConfig().FileStoragePath = "testdata/testfile.json"

	if _, err := os.Stat(config.GetConfig().FileStoragePath); !errors.Is(err, os.ErrNotExist) {
		file, err := os.ReadFile(config.GetConfig().FileStoragePath)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		} else {
			if len(file) > 0 {
				os.Remove(config.GetConfig().FileStoragePath)
			}
		}
	}

	storage := NewStorage(config.GetConfig().FileStoragePath)
	userID := 999 // Non-existent user ID
	shortURLs := []string{"short1", "short2"}

	// Prepopulate storage with URLs
	_, err := storage.Put(ctx, dto.Shortening{ShortURL: "short1", OriginalURL: "http://example.com/1"})
	if err != nil {
		t.Fatalf(" not expected error, got nil")
	}
	_, err = storage.Put(ctx, dto.Shortening{ShortURL: "short2", OriginalURL: "http://example.com/2"})
	if err != nil {
		t.Fatalf("not expected error, got %v", err)
	}

	// Call DeleteUserURL
	storage.DeleteUserURL(ctx, shortURLs, userID)

	// Verify URLs are not deleted
	for _, url := range shortURLs {
		_, err := storage.Get(ctx, url)
		if err != nil {
			t.Errorf("expected URL %s to still exist", url)
		}
	}
}

// Successfully saves a map of URLs to a file
func TestSaveSuccessfullySavesMapToFile(t *testing.T) {
	fileName := "test_output/successful_save.txt"
	urls := &sync.Map{}
	urls.Store("1", dto.Shortening{
		ShortURL:    "http://short.url/1",
		OriginalURL: "http://original.url/1",
	})
	urls.Store("2", dto.Shortening{
		ShortURL:    "http://short.url/2",
		OriginalURL: "http://original.url/2",
	})

	err := save(fileName, urls)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fatalf("Expected file to be created, but it does not exist")
	}

	os.RemoveAll("test_output")
}

// Handles empty sync.Map without errors
func TestSaveHandlesEmptySyncMap(t *testing.T) {
	fileName := "test_output/empty_map.txt"
	urls := &sync.Map{}

	err := save(fileName, urls)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fatalf("Expected file to be created, but it does not exist")
	}

	os.RemoveAll("test_output")
}

//// Successfully reads a file and loads its contents into the sync.Map
//func TestLoadSuccessfullyReadsFile(t *testing.T) {
//	var files sync.Map
//	fileName := "testdata/testfile.txt"
//	data := `{"correlation_id":"123","short_url":"short1","original_url":"http://example.com/1","user_id":"user1","is_deleted":false}
//             {"correlation_id":"124","short_url":"short2","original_url":"http://example.com/2","user_id":"user2","is_deleted":false}`
//	err := os.WriteFile(fileName, []byte(data), 0644)
//	if err != nil {
//		t.Fatalf("Failed to write test file: %v", err)
//	}
//	defer os.Remove(fileName)
//
//	err = load(&files, fileName)
//	if err != nil {
//		t.Fatalf("Expected no error, got %v", err)
//	}
//
//	_, ok := files.Load("short1")
//	if !ok {
//		t.Errorf("Expected short1 to be loaded into sync.Map")
//	}
//
//	_, ok = files.Load("short2")
//	if !ok {
//		t.Errorf("Expected short2 to be loaded into sync.Map")
//	}
//}
//
//// Handles an empty file without errors
//func TestLoadHandlesEmptyFile(t *testing.T) {
//	var files sync.Map
//	fileName := "testdata/emptyfile.txt"
//	err := os.WriteFile(fileName, []byte(""), 0644)
//	if err != nil {
//		t.Fatalf("Failed to write test file: %v", err)
//	}
//	defer os.Remove(fileName)
//
//	err = load(&files, fileName)
//	if err != nil {
//		t.Fatalf("Expected no error, got %v", err)
//	}
//
//	files.Range(func(key, value interface{}) bool {
//		t.Errorf("Expected no entries in sync.Map, found key: %v", key)
//		return false
//	})
//}
