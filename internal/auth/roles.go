// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"slices"

	"github.com/grd/FreePDM/internal/db"
)

func IsAdmin(user *db.PdmUser) bool {
	return slices.Contains(user.Roles, "Admin")
}

func IsProjectLead(user *db.PdmUser) bool {
	return slices.Contains(user.Roles, "ProjectLead")
}

func IsSeniorDesigner(user *db.PdmUser) bool {
	return slices.Contains(user.Roles, "SeniorDesigner")
}

// HasAnyRole checks if the user has at least one of the given roles
func HasAnyRole(user *db.PdmUser, roles ...string) bool {
	for _, role := range roles {
		if slices.Contains(user.Roles, role) {
			return true
		}
	}
	return false
}
