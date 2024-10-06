// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
)

// https://www.osohq.com/post/sqlalchemy-role-rbac-basics

// struct for genrating Roles
type TempRole struct {
}

// Create new Role
func (t TempRole) AddRole() {
	fmt.Println("new role created")
}

// Delete existing role
func (t TempRole) RemoveRole() {
	fmt.Println("existing role deleted")
}

// Struct for generating users
// Users are Aliases for roles in SQL see: https://www.postgresql.org/docs/14/sql-createuser.html
type TempUser struct {
}

func (t TempUser) AddUserToSql(username string) {
	fmt.Println("This is basically the interface")
}

// Delete existing user
func (t TempUser) RemoveUserFromSql(user_id int, username string) {
	fmt.Println("existing user deleted")
}

func (t TempUser) AddUserToLdap(username string) {
	fmt.Println("This is basically the interface")
}

// Delete existing user
func (t TempUser) RemoveUserFromLdap(user_id int, username string) {
	fmt.Println("existing user deleted")
}

type Permission int

const (
	General Permission = iota
	CadUser
	SuperUser
	ProjectLeader
	Admin
)

// first implement only:
// - General activities
// - CADUser
// all other options are added later!

// role -> permissions mapping
// See: https://www.osohq.com/academy/role-based-access-control-rbac
func RolePermissions(role Permission) []string {
	switch role {
	case General:
		return []string{"Read Documents", "Read Items", "Read Models"}
	case CadUser:
		return []string{"Check-In", "Check-Out", "Create Document", "Create Item", "Create Model"}
	case SuperUser:
		return []string{"Delete Document", "Delete Item", "Delete Model"} // User activities plus
	case ProjectLeader:
		return []string{"Create Project", "Add User to Project", "Remove User from Project"}
	case Admin:
		return []string{"Create User", "Delete User", "Create Database"}
	}
	return nil
}
