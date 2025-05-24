// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"strings"

	"github.com/grd/FreePDM/internal/db"
)

func IsAdmin(user *db.PdmUser) bool {
	return user.HasRole("admin")
}

func IsProjectLead(user *db.PdmUser) bool {
	return user.HasRole("projectlead")
}

func IsSeniorDesigner(user *db.PdmUser) bool {
	return user.HasRole("seniordesigner")
}

// HasAnyRole checks if the user has at least one of the given roles (case-insensitive)
func HasAnyRole(user *db.PdmUser, roles ...string) bool {
	for _, checkRole := range roles {
		for _, userRole := range user.Roles {
			if strings.EqualFold(userRole, checkRole) {
				return true
			}
		}
	}
	return false
}
