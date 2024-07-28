package database

import (
	"testing"
)

// Successfully creates a new database connection with valid driver and DSN
//func TestNewConnection_Success(t *testing.T) {
//	driver := "pgx"
//	dsn := "user=postgres password=secret dbname=testdb sslmode=disable"
//
//	db, err := NewConnection(driver, dsn)
//
//	if err != nil {
//		t.Fatalf("expected no error, got %v", err)
//	}
//
//	if db == nil {
//		t.Fatalf("expected a valid database connection, got nil")
//	}
//
//	if err = db.Ping(); err != nil {
//		t.Fatalf("expected successful ping, got %v", err)
//	}
//
//	db.Close()
//}

// Fails to create a database connection with an invalid driver
func TestNewConnection_InvalidDriver(t *testing.T) {
	driver := "invalid_driver"
	dsn := "user=postgres password=secret dbname=testdb sslmode=disable"

	db, err := NewConnection(driver, dsn)

	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	if db != nil {
		t.Fatalf("expected no database connection, got %v", db)
	}
}
