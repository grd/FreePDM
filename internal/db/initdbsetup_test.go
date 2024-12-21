package db

import (
	"fmt"
	"testing"
)

func TestDatabaseFunctions(t *testing.T) {
	username := "freepdm"
	password := "PsqlPassword123!"
	dbname := "freepdm"
	url := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", username, password, dbname)

	db, err := StartYourEngine(url)
	if err != nil {
		t.Fatalf("Failed to start database engine: %v", err)
	}

	err = db.CreateDefaultTables()
	if err != nil {
		t.Fatalf("Failed to create default tables: %v", err)
	}
}
