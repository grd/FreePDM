// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB loads configuration and initializes the database
func InitDB() (*gorm.DB, error) {
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatalln("no app.env found, use standard settings")
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev" // default
	}

	var db *gorm.DB

	switch env {
	case "prod", "dev":
		// PostgreSQL settings from environment variables
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			getEnv("PG_HOST", "localhost"),
			getEnv("PG_USER", "pdmuser"),
			getEnv("PG_PASSWORD", "pdmsecret"),
			getEnv("PG_NAME", "pdmdb"),
			getEnv("PG_PORT", "5432"),
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "test":
		fallthrough
	default:
		// Use SQLite for testing or fallback
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Auto-create tables if they don't exist
	err = createDefaultTables(db)
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Create admin account (only once)
	err = createAdminAccount(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// getEnv returns a fallback value if key is not set
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

// createDefaultTables creates the default set of tables in the database.
func createDefaultTables(db *gorm.DB) error {
	// Check if a key table already exists

	if db.Migrator().HasTable(&PdmUser{}) {
		// fmt.Println("Tables already exist — skipping creation.")
		return nil
	}

	// Create the default set of tables
	if err := db.AutoMigrate(&PdmUser{}); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}
	// if err := db.AutoMigrate(&PdmUser{}, &PdmProject{}, &PdmItem{},
	// 	&PdmModel{}, &PdmDocument{}, &PdmMaterial{}, &PdmHistory{}, &PdmPurchase{},
	// 	&PdmManufacturer{}, &PdmVendor{}); err != nil {
	// 	return fmt.Errorf("failed to migrate tables: %w", err)
	// }

	// Log successful creation
	fmt.Println("Tables created successfully")
	return nil
}

func createAdminAccount(db *gorm.DB) error {
	var count int64
	db.Model(&PdmUser{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
		admin := PdmUser{
			UserName:           "admin",
			PasswordHash:       string(hashed),
			MustChangePassword: true,
			Roles:              []string{string(Admin)},
		}
		return db.Create(&admin).Error
	}

	return nil
}
