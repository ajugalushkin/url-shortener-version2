package storage

import (
	"testing"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage/file"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage/inmemory"
)

// Returns a file storage when FileStoragePath is provided
func TestGetStorageReturnsFileStorage(t *testing.T) {
	// Arrange
	config.GetConfig().FileStoragePath = "mock_path"

	// Act
	storage := GetStorage()

	// Assert
	if _, ok := storage.(*file.Storage); !ok {
		t.Errorf("expected *file.Storage, got %T", storage)
	}

	config.GetConfig().FileStoragePath = ""
}

// Returns an in-memory storage when neither DataBaseDsn nor FileStoragePath is provided
func TestGetStorageReturnsInMemory(t *testing.T) {
	// Arrange

	// Act
	storage := GetStorage()

	// Assert
	if _, ok := storage.(*inmemory.InMemory); !ok {
		t.Errorf("expected *inmemory.InMemory, got %T", storage)
	}
}
