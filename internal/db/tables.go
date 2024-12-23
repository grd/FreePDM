// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Base struct to embed in other models
type Base struct {
	ID uint `gorm:"primaryKey"`
}

// PdmUser represents the users table
type PdmUser struct {
	Base
	UserName         string `gorm:"type:varchar(30)"`
	UserFirstName    string `gorm:"type:varchar(30)"`
	UserLastName     string `gorm:"type:varchar(30)"`
	UserFullName     string
	UserEmailAddress string `gorm:"not null"`
	UserPhoneNumber  int    // Consider changing to string if phone numbers include characters
	UserDepartment   string // Enum type can be defined later

	Roles     []*PdmRole    `gorm:"many2many:user_role_link"`
	Projects  []*PdmProject `gorm:"many2many:user_project_link"`
	Items     []PdmItem     `gorm:"foreignKey:UserID"`
	Models    []PdmModel    `gorm:"foreignKey:UserID"`
	Documents []PdmDocument `gorm:"foreignKey:UserID"`
}

// PdmRole represents the roles table
type PdmRole struct {
	Base
	RoleName string `gorm:"type:varchar(32)"`

	Users []*PdmUser `gorm:"many2many:user_role_link"`
}

// PdmUserRoleLink is the association table for users and roles
type PdmUserRoleLink struct {
	UserID uint `gorm:"primaryKey;autoIncrement:false"`
	RoleID uint `gorm:"primaryKey;autoIncrement:false"`
}

// PdmProject represents the projects table
type PdmProject struct {
	Base
	ProjectNumber     string `gorm:"type:varchar(16);not null"`
	ProjectName       string `gorm:"type:varchar(32)"`
	ProjectStatus     string // Enum type can be defined later
	ProjectDateStart  *time.Time
	ProjectDateFinish *time.Time
	ProjectPath       string

	Users []*PdmUser `gorm:"many2many:user_project_link"`
}

// PdmUserProjectLink is the association table for users and projects
type PdmUserProjectLink struct {
	UserID    uint `gorm:"primaryKey;autoIncrement:false"`
	ProjectID uint `gorm:"primaryKey;autoIncrement:false"`
}

// PdmItem represents the items table
type PdmItem struct {
	Base
	ItemNumber            string `gorm:"type:varchar(16)"`
	ItemName              string `gorm:"type:varchar(32)"`
	ItemDescription       string `gorm:"type:varchar(32)"`
	ItemFullDescription   string
	ItemNumberLinkedFiles int
	ItemPath              string `gorm:"not null"`
	ItemPreview           []byte // For LargeBinary

	UserID    uint          `gorm:"foreignKey:UserID"`
	User      PdmUser       `gorm:"foreignKey:UserID"`
	ProjectID uint          `gorm:"foreignKey:ProjectID"`
	Models    []PdmModel    `gorm:"foreignKey:ItemID"`
	Documents []PdmDocument `gorm:"foreignKey:ItemID"`
	// Material   PdmMaterial   `gorm:"foreignKey:ItemID"`
	// Purchasing PdmPurchase   `gorm:"foreignKey:ItemID"`
}

// PdmProjectItemLink is the association table for projects and items
type PdmProjectItemLink struct {
	ProjectID uint `gorm:"primaryKey;autoIncrement:false"`
	ItemID    uint `gorm:"primaryKey;autoIncrement:false"`
}

// PdmModel represents the models table
type PdmModel struct {
	Base
	ModelNumber          int    // Or string?
	ModelName            string `gorm:"type:varchar(32)"`
	ModelDescription     string `gorm:"type:varchar(32)"`
	ModelFullDescription string
	ModelFilename        string `gorm:"type:varchar(253);not null"`
	ModelExt             string `gorm:"type:varchar(253);not null"`
	ModelPreview         []byte // For LargeBinary

	UserID   uint        `gorm:"foreignKey:UserID"`
	User     PdmUser     `gorm:"foreignKey:UserID"`
	ItemID   uint        `gorm:"foreignKey:ItemID"`
	Item     PdmItem     `gorm:"foreignKey:ItemID"`
	Material PdmMaterial `gorm:"foreignKey:ModelID"`
}

// PdmDocument represents the documents table
type PdmDocument struct {
	Base
	DocumentNumber          int    // Or string?
	DocumentName            string `gorm:"type:varchar(32)"`
	DocumentDescription     string `gorm:"type:varchar(32)"`
	DocumentFullDescription string
	DocumentFilename        string `gorm:"type:varchar(253);not null"`
	DocumentExt             string `gorm:"type:varchar(253);not null"`

	UserID uint    `gorm:"foreignKey:UserID"`
	User   PdmUser `gorm:"foreignKey:UserID"`
	ItemID uint    `gorm:"foreignKey:ItemID"`
	Item   PdmItem `gorm:"foreignKey:ItemID"`
}

// PdmMaterial represents the materials table
type PdmMaterial struct {
	Base
	MaterialName            string `gorm:"type:varchar(32)"`
	MaterialFinish          string `gorm:"type:varchar(32)"`
	MaterialDensity         float64
	MaterialDensityUnit     string // Enum can be defined later
	MaterialVolume          float64
	MaterialVolumeUnit      string // Enum can be defined later
	MaterialWeight          float64
	MaterialWeightUnit      string // Enum can be defined later
	MaterialSurfaceArea     float64
	MaterialSurfaceAreaUnit string // Enum can be defined later

	// ModelID uint     `gorm:"foreignKey:ModelID"`
	// Model   PdmModel `gorm:"foreignKey:ModelID"`
	ItemID uint    `gorm:"foreignKey:ItemID"`
	Item   PdmItem `gorm:"foreignKey:ItemID"`
}

// PdmHistory represents the history table
type PdmHistory struct {
	Base
	HistoryDateCreated    *time.Time
	HistoryCreatedBy      string
	HistoryDateLastEdit   *time.Time
	HistoryLastEditBy     string
	HistoryCheckedOutBy   string
	HistoryRevisionState  string // Enum can be defined later
	HistoryRevisionNumber int
	HistoryStoredNumber   int
}

// PdmPurchase represents the purchasing table
type PdmPurchase struct {
	Base
	PurchasingSource       bool
	PurchasingTraceability string // Enum can be defined later

	ItemID         uint              `gorm:"foreignKey:ItemID"`
	Item           PdmItem           `gorm:"foreignKey:ItemID"`
	ManufacturerID uint              `gorm:"foreignKey:ManufacturerID"`
	Manufacturers  []PdmManufacturer `gorm:"foreignKey:ManufacturerID"`
	VendorID       uint              `gorm:"foreignKey:VendorID"`
	Vendors        []PdmVendor       `gorm:"foreignKey:VendorID"`
}

// PdmManufacturer represents the manufacturers table
type PdmManufacturer struct {
	Base
	ManufacturerName string      `gorm:"type:varchar(32)"`
	PurchasingID     uint        `gorm:"foreignKey:PurchasingID"`
	Purchasing       PdmPurchase `gorm:"foreignKey:PurchasingID"`
}

// PdmVendor represents the vendors table
type PdmVendor struct {
	Base
	VendorName    string      `gorm:"type:varchar(32)"`
	VPurchasingID uint        `gorm:"foreignKey:VPurchasingID"`
	VPurchasing   PdmPurchase `gorm:"foreignKey:VPurchasingID"`
}

// To initialize the database
func InitDB() {

	var err error

	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Europe/Amsterdam"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	// Auto migrate the schema
	db.AutoMigrate(&PdmUser{}, &PdmRole{}, &PdmUserRoleLink{}, &PdmProject{}, &PdmUserProjectLink{}, &PdmItem{},
		&PdmProjectItemLink{}, &PdmModel{}, &PdmDocument{}, &PdmMaterial{}, &PdmHistory{}, &PdmPurchase{}, &PdmManufacturer{}, &PdmVendor{})
}