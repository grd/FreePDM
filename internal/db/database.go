// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseGen is a generic SQL Database class
type DatabaseGen struct {
	DriverName   string
	Username     string
	Password     string
	Host         string
	Port         int
	DatabaseName string
	Engine       *gorm.DB
}

// NewDatabaseGen initializes a new DatabaseGen
func NewDatabaseGen() *DatabaseGen {
	log.Println("Generic DataBase")
	return &DatabaseGen{}
}

// MakeURL creates a new connection URL
func (db *DatabaseGen) MakeURL() string {
	return "host=" + db.Host +
		" port=" + strconv.Itoa(db.Port) +
		" user=" + db.Username +
		" password=" + db.Password +
		" dbname=" + db.DatabaseName +
		" sslmode=disable"
}

// CreateDB not implemented yet
func (db *DatabaseGen) CreateDB() error {
	return fmt.Errorf("func CreateDB is not implemented yet")
}

// CreateTable not implemented yet
func (db *DatabaseGen) CreateTable() error {
	return fmt.Errorf("func CreateTable is not implemented yet")
}

// DeleteTable not implemented yet
func (db *DatabaseGen) DeleteTable() error {
	return fmt.Errorf("func DeleteTable is not implemented yet")
}

// GetTables not implemented yet
func (db *DatabaseGen) GetTables() error {
	return fmt.Errorf("func GetTables is not implemented yet")
}

// DataBasePostgreSQL struct feeds forward of generic SQL functions to PostgreSQL
type DataBasePostgreSQL struct {
	*DatabaseGen
}

// NewDataBasePostgreSQL initializes a new PostgreSQL database
func NewDataBasePostgreSQL() *DataBasePostgreSQL {
	log.Println("PostgreSQL")
	return &DataBasePostgreSQL{DatabaseGen: NewDatabaseGen()}
}

// StartEngine starts the PostgreSQL engine
func (db *DataBasePostgreSQL) StartEngine(echo bool) error {
	dsn := db.MakeURL()
	var err error
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Adjust log level as needed
	}

	db.Engine, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return err
	}
	return nil
}
