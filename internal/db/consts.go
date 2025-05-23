// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

type (
	Role             string
	RBAC             string
	ProjectState     string
	RevisionState    string
	TracebilityState string
	DensityUnit      string
	VolumeUnit       string
	WeightUnit       string
	AreaUnit         string
)

const (
	Admin          Role = "admin"
	Designer       Role = "designer"
	SeniorDesigner Role = "senior"
	Viewer         Role = "viewer"
	Editor         Role = "editor"
	Approver       Role = "approver"
	ProjectLead    Role = "projectlead"
	Qa             Role = "qa"
	Guest          Role = "guest"
)

const (
	CheckIn               RBAC = "Check-In"
	CheckOut              RBAC = "Check-Out"
	CreateDocument        RBAC = "Create Document"
	CreateItem            RBAC = "Create Item"
	CreateModel           RBAC = "Create Model"
	DeleteDocument        RBAC = "Delete Document"
	DeleteItem            RBAC = "Delete Item"
	DeleteModel           RBAC = "Delete Model"
	CreateProject         RBAC = "Create Project"
	AddUserToProject      RBAC = "Add User to Project"
	RemoveUserFromProject RBAC = "Remove User from Project"
	CreateUser            RBAC = "Create User"
	DeleteUser            RBAC = "Delete User"
	CreateDatabase        RBAC = "Create Database"
	ReadDocuments         RBAC = "Read Documents"
	ReadItems             RBAC = "Read Items"
	ReadModels            RBAC = "Read Models"
)

const (
	// ProjectState

	// https://www.indeed.com/career-advice/career-development/project-statuses
	// https://support.ptc.com/help/wnc/r11.2.0.0/en/index.html#page/Windchill_Help_Center/ProjMgmtPhaseState.html

	New        ProjectState = "New"
	Upcoming   ProjectState = "Upcoming"
	Pending    ProjectState = "Pending"
	NotStarted ProjectState = "Not Started"
	Draft      ProjectState = "Draft"
	Active     ProjectState = "Active"
	Priority   ProjectState = "Priority"
	Canceld    ProjectState = "Canceled"
	Onhold     ProjectState = "On-Hold"
	Archived   ProjectState = "Archived"

	// RevisionState

	Concept     RevisionState = "Concept"
	Underreview RevisionState = "Under Review"
	Released    RevisionState = "Released"
	Inwork      RevisionState = "In-Work"
	Depreciated RevisionState = "Depreciated"

	// TracebilityState

	Lot       TracebilityState = "Lot"
	LotSerial TracebilityState = "Lot and Serial"
	Serial    TracebilityState = "Serial"
	NotTaced  TracebilityState = "Not-Traced"

	// DensityUnit

	D_gmm3        DensityUnit = "Gram / mm^3"
	D_gcm3        DensityUnit = "Gram / cm^3"
	D_kgm3        DensityUnit = "Kilo Gram / m^3"
	D_tonnem3     DensityUnit = "Tonne / m^3"
	D_metrictonm3 DensityUnit = "Metric Ton / m^3" // equal to tonne
	D_poundft3    DensityUnit = "Pound / Ft^3"
	D_poundinch3  DensityUnit = "Pound / inch^3"

	// VolumeUnit

	V_mm3   VolumeUnit = "Cubic mm [mm^3]"
	V_cm3   VolumeUnit = "Cubic cm [cm^3]"
	V_m3    VolumeUnit = "Cubic m [m^3]"
	V_litre VolumeUnit = "Litre"
	V_ft3   VolumeUnit = "Cubic ft [Ft^3]"
	V_inch3 VolumeUnit = "Cubic inch [Inch^3]"

	// WeightUnit

	W_g         WeightUnit = "Gram [g]"
	W_kg        WeightUnit = "Kilo Gram [kg]"
	W_tonne     WeightUnit = "Tonne [t]"
	W_metricton WeightUnit = "Metric Ton" // equal to tonne
	W_pound     WeightUnit = "Pound [p]"
	W_slug      WeightUnit = "Slug"

	// AreaUnit

	A_mm2   AreaUnit = "Square mm [mm^2]"
	A_cm2   AreaUnit = "Square cm [cm^2]"
	A_m2    AreaUnit = "Square m [m^2]"
	A_ft2   AreaUnit = "Square ft [Ft^2]"
	A_inch2 AreaUnit = "Square inch [Inch^2]"
)
