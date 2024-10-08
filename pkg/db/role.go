// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
)

// https://www.osohq.com/post/sqlalchemy-role-rbac-basics

// for generating Roles
type Role struct {
}

// Create new role
func (r Role) AddRole() {
	fmt.Println("new role created")
}

// Delete existing role
func (r Role) RemoveRole() {
	fmt.Println("existing role deleted")
}

// for generating users
// Users are Aliases for roles in SQL see: https://www.postgresql.org/docs/14/sql-createuser.html
type User struct {
}

func (r User) AddUserToSql(username string) {
	fmt.Println("This is basically the interface")
}

// Delete existing user
func (r User) RemoveUserFromSql(user_id int, username string) {
	fmt.Println("existing user deleted")
}

func (r User) AddUserToLdap(username string) {
	fmt.Println("This is basically the interface")
}

// Delete existing user
func (r User) RemoveUserFromLdap(user_id int, username string) {
	fmt.Println("existing user deleted")
}
