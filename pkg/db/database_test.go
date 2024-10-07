package db

import (
	"testing"
)

func TestMakeURL(t *testing.T) {
	db := NewDataBasePostgreSQL()
	db.Username = "test_user"
	db.Password = "test_password"
	db.Host = "localhost"
	db.Port = 5432
	db.DatabaseName = "test_db"

	expectedURL := "host=localhost port=5432 user=test_user password=test_password dbname=test_db sslmode=disable"
	if db.MakeURL() != expectedURL {
		t.Errorf("MakeURL() = %v; want %v", db.MakeURL(), expectedURL)
	}
}

func TestStartEngine(t *testing.T) {
	db := NewDataBasePostgreSQL()
	db.Username = "your_username"
	db.Password = "your_password"
	db.Host = "localhost"
	db.Port = 5432
	db.DatabaseName = "your_database_name"

	err := db.StartEngine(false) // set echo to false to disable logging
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Optionally, you could test for db.Engine not being nil after successful connection
	if db.Engine == nil {
		t.Error("db.Engine is nil; expected a valid gorm.DB instance")
	}
}
