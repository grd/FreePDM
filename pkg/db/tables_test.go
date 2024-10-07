package db

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("host=localhost user=your_user password=your_password dbname=your_test_db sslmode=disable"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&PdmUser{}, &PdmRole{}, &PdmProject{})
	return db, nil
}

func TestPdmUserCreation(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	user := PdmUser{
		UserName:         "test_user",
		UserFirstName:    "Test",
		UserLastName:     "User",
		UserEmailAddress: "test@example.com",
	}

	result := db.Create(&user)
	if result.Error != nil {
		t.Fatalf("failed to create userYour SQLAlchemy ORM code looks well")

	}
}
