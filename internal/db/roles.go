// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

//
// How to check whether a guy has access:
//
// if !user.HasPermission(CheckIn) {
// 	http.Error(w, "Forbidden", http.StatusForbidden)
// 	return
// }
//

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
func RolePermissions(role []Role) (ret []RBAC) {
	for _, r := range role {
		switch r {
		case Editor, Approver, Qa, Designer:
			ret = append(ret,
				CheckIn,
				CheckOut,
				CreateDocument,
				CreateItem,
				CreateModel,
			)
		case Guest, Viewer:
			ret = append(ret,
				ReadDocuments,
				ReadItems,
				ReadModels,
			)
		case SeniorDesigner:
			ret = append(ret,
				DeleteDocument,
				DeleteItem,
				DeleteModel,
			)
		case ProjectLead:
			ret = append(ret,
				CreateProject,
				AddUserToProject,
				RemoveUserFromProject,
			)
		case Admin:
			ret = append(ret,
				CreateUser,
				DeleteUser,
				CreateDatabase,
			)
		}
	}
	return
}

func DefaultRoles() []Role {
	return []Role{Viewer} // or Guest, depends on the baseline
}

// Check whether user has a role
func (u *PdmUser) HasRole(role string) bool {
	for _, r := range u.Roles {
		if string(r) == role {
			return true
		}
	}
	return false
}

// Check whether user has all the roles
func (u *PdmUser) HasAnyRole(roles []string) bool {
	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}

func (u *PdmUser) HasAllRoles(roles []string) bool {
	for _, role := range roles {
		if !u.HasRole(role) {
			return false
		}
	}
	return true
}

// HasPermission checks if the user has the given RBAC permission.
func (u *PdmUser) HasPermission(permission RBAC) bool {
	permissions := RolePermissions(stringsToRoles(u.Roles))
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// HasAnyPermission checks if the user has at least one of the given permissions.
func (u *PdmUser) HasAnyPermission(perms []RBAC) bool {
	userPerms := RolePermissions(stringsToRoles(u.Roles))
	for _, perm := range perms {
		for _, userPerm := range userPerms {
			if perm == userPerm {
				return true
			}
		}
	}
	return false
}

func stringsToRoles(in []string) []Role {
	roles := make([]Role, 0, len(in))
	for _, r := range in {
		roles = append(roles, Role(r))
	}
	return roles
}
