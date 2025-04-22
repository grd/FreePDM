// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// // https://www.osohq.com/post/sqlalchemy-role-rbac-basics

// // struct for generating Roles
// type TempRole struct {
// }

// // Create new Role
// func (t TempRole) AddRole() {
// 	fmt.Println("new role created")
// }

// // Delete existing role
// func (t TempRole) RemoveRole() {
// 	fmt.Println("existing role deleted")
// }

// // Struct for generating users
// // Users are Aliases for roles in SQL see: https://www.postgresql.org/docs/14/sql-createuser.html
// type TempUser struct {
// }

// func (t TempUser) AddUserToSql(username string) {
// 	fmt.Println("This is basically the interface")
// }

// // Delete existing user
// func (t TempUser) RemoveUserFromSql(user_id int, username string) {
// 	fmt.Println("existing user deleted")
// }

// func (t TempUser) AddUserToLdap(username string) {
// 	fmt.Println("This is basically the interface")
// }

// // Delete existing user
// func (t TempUser) RemoveUserFromLdap(user_id int, username string) {
// 	fmt.Println("existing user deleted")
// }

type RoleList []Role

func (r RoleList) Value() (driver.Value, error) {
	strs := make([]string, len(r))
	for i, role := range r {
		strs[i] = string(role)
	}
	return strings.Join(strs, ","), nil
}

func (r *RoleList) Scan(value interface{}) error {
	if value == nil {
		*r = []Role{}
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan Roles: value is not a string")
	}

	parts := strings.Split(str, ",")
	roles := make([]Role, len(parts))
	for i, s := range parts {
		roles[i] = Role(s)
	}
	*r = roles
	return nil
}

// role -> permissions mapping
// See: https://www.osohq.com/academy/role-based-access-control-rbac
func RolePermissions(role []Role) (ret []string) {
	for _, r := range role {
		switch r {
		case Editor:
			ret = append(ret, []string{"Check-In", "Check-Out", "Create Document", "Create Item", "Create Model"}...)
		case Approver:
			ret = append(ret, []string{"Check-In", "Check-Out", "Create Document", "Create Item", "Create Model"}...)
		case Qa:
			ret = append(ret, []string{"Check-In", "Check-Out", "Create Document", "Create Item", "Create Model"}...)
		case Guest:
			ret = append(ret, []string{"Read Documents", "Read Items", "Read Models"}...)
		case Viewer:
			ret = append(ret, []string{"Read Documents", "Read Items", "Read Models"}...)
		case Designer:
			ret = append(ret, []string{"Check-In", "Check-Out", "Create Document", "Create Item", "Create Model"}...)
		case Senior: // User activities plus
			ret = append(ret, []string{"Delete Document", "Delete Item", "Delete Model"}...)
		case ProjectLead:
			ret = append(ret, []string{"Create Project", "Add User to Project", "Remove User from Project"}...)
		case Admin:
			ret = append(ret, []string{"Create User", "Delete User", "Create Database"}...)
		}
	}
	return ret
}

// Check whether user has a role
func (u *PdmUser) HasRole(role string) bool {
	switch Role(role) {
	case Editor:
		return true
	case Approver:
		return true
	case Qa:
		return true
	case Guest:
		return true
	case Viewer:
		return true
	case Designer:
		return true
	case Senior:
		return true
	case ProjectLead:
		return true
	case Admin:
		return true
	default:
		return false
	}
}

// Check whether user has all the roles
func (u *PdmUser) HasAnyRole(roles []string) bool {
	for _, r := range roles {
		if !u.HasRole(r) {
			return false
		}
	}
	return true
}
