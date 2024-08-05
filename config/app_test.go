package config

import (
	"testing"
)

// File does not exist or cannot be opened
func TestLoadConfigurationFileNotExist(t *testing.T) {
	// Call the function with a non-existent file path
	config := loadConfiguration("/path/to/nonexistent/file.json")

	// Assert that the returned configuration is empty
	if config.ServerAddress != "" {
		t.Errorf("Expected ServerAddress to be empty, got '%s'", config.ServerAddress)
	}
	if config.BaseURL != "" {
		t.Errorf("Expected BaseURL to be empty, got '%s'", config.BaseURL)
	}
	if config.FlagLogLevel != "" {
		t.Errorf("Expected Log_Level to be empty, got '%s'", config.FlagLogLevel)
	}
	if config.FileStoragePath != "" {
		t.Errorf("Expected File_Storage_PATH to be empty, got '%s'", config.FileStoragePath)
	}
	if config.DataBaseDsn != "" {
		t.Errorf("Expected DataBase_Dsn to be empty, got '%s'", config.DataBaseDsn)
	}
	if config.SecretKey != "" {
		t.Errorf("Expected Secret_Key to be empty, got '%s'", config.SecretKey)
	}
	if config.EnableHTTPS {
		t.Errorf("Expected Enable_HTTPS to be false, got true")
	}
}
