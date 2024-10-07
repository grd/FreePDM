package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CreateDefaultTables creates the default set of tables in the database.
func CreateDefaultTables(db *gorm.DB) error {
	// Create the default set of tables
	if err := db.AutoMigrate(&PdmUser{}, &PdmRole{}, &PdmProject{}, &PdmItem{},
		&PdmModel{}, &PdmDocument{}, &PdmMaterial{}, &PdmHistory{}, &PdmPurchase{},
		&PdmManufacturer{}, &PdmVendor{}); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	// Log successful creation
	fmt.Println("Tables created successfully")
	return nil
}

// InitializeDB initializes the database connection.
func InitializeDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	return db, nil
}
