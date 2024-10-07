package db

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Mock database URL for testing
const testDatabaseURL = "host=localhost user=freepdm password=PsqlPassword123! dbname=freepdm port=5432 sslmode=disable"

func TestCreateDefaultTables(t *testing.T) {
	db, err := gorm.Open(postgres.Open(testDatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	if err := CreateDefaultTables(db); err != nil {
		t.Fatalf("error creating tables: %v", err)
	}

	// Additional checks can be added here to verify table creation
}
