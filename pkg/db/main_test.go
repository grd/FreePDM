package db

import (
	"fmt"
	"log"
	"strings"
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
func main() {
	// Example usage
	config := DatabaseConfig{
		User:     "freepdm",
		Password: "PsqlPassword123!",
		Host:     "localhost",
		Port:     5432,
		DBName:   "freepdm",
	}

	db, err := StartYourEngine("postgresql", config)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	log.Println("Database connection established:", db)
}

// TestStartYourEngine tests the StartYourEngine function.
func TestStartYourEngine(t *testing.T) {
	config := DatabaseConfig{
		User:     "freepdm",
		Password: "PsqlPassword123!",
		Host:     "localhost",
		Port:     5432,
		DBName:   "freepdm",
	}

	db, err := StartYourEngine("postgresql", config)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if db == nil {
		t.Fatal("expected a database connection, got nil")
	}
}

// TestInvalidDBType tests the function with an invalid database type.
func TestInvalidDBType(t *testing.T) {
	config := DatabaseConfig{
		User:     "freepdm",
		Password: "PsqlPassword123!",
		Host:     "localhost",
		Port:     5432,
		DBName:   "freepdm",
	}

	_, err := StartYourEngine("mysql", config)
	if err == nil {
		t.Fatal("expected an error, got none")
	}
	if !strings.Contains(err.Error(), "mysql is not a valid database type") {
		t.Fatalf("expected error about invalid db type, got %v", err)
	}
}

// TestDBConnection tests the database connection with invalid credentials.
func TestDBConnectionWithInvalidCredentials(t *testing.T) {
	config := DatabaseConfig{
		User:     "invalid_user",
		Password: "invalid_password",
		Host:     "localhost",
		Port:     5432,
		DBName:   "freepdm",
	}

	_, err := StartYourEngine("postgresql", config)
	if err == nil {
		t.Fatal("expected an error, got none")
	}
}
