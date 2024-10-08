// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database struct encapsulates the GORM DB instance.
type Database struct {
	Engine *gorm.DB
}

// UserAccount represents the user_accounts table.
type UserAccount struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:100"`
}

// UserRole represents the user_roles table.
type UserRole struct {
	ID       uint   `gorm:"primaryKey"`
	RoleName string `gorm:"size:100"`
}

// CreateDefaultTables creates default tables in the database.
func (d *Database) CreateDefaultTables() error {
	err := d.Engine.AutoMigrate(&UserAccount{}, &UserRole{}, &Project{}, &Item{})
	if err != nil {
		return err
	}
	return nil
}

// StartYourEngine initializes the database connection.
func StartYourEngine(url string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Database{Engine: db}, nil
}
