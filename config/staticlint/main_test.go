package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadExecutablePath(t *testing.T) {
	appfile, err := os.Executable()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if appfile == "" {
		t.Fatalf("Expected a valid executable path, got an empty string")
	}

	expectedDir := filepath.Dir(appfile)
	if expectedDir == "" {
		t.Fatalf("Expected a valid directory path, got an empty string")
	}
}

func TestMissingConfigFile(t *testing.T) {
	appfile, err := os.Executable()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	configPath := filepath.Join(filepath.Dir(appfile), Config)
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatalf("Expected config file to be missing, but it exists")
	}
}
