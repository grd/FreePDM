package db

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig holds the database connection parameters.
type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

const dsn = "user=freepdm password=PsqlPassword123 dbname=freepdm host=localhost port=5432 sslmode=disable"

// StartYourEngine initializes the database connection.
func StartYourEnginenow(dbType string, dbConfig DatabaseConfig) (*gorm.DB, error) {
	if dbType != "postgresql" {
		return nil, fmt.Errorf("%s is not a valid database type", dbType)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// main function for testing purposes
func TestStartYourEngine(t *testing.T) {
	db, err := StartYourEngine(dsn)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	log.Println("Database connection established:", db)
}

// TestStartYourEngine tests the StartYourEngine function.
func TestStartYourEngine2(t *testing.T) {
	db, err := StartYourEngine(dsn)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if db == nil {
		t.Fatal("expected a database connection, got nil")
	}
}
