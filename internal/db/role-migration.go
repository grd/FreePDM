// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
	"strings"
)

// NormalizeUserRolesVerbose updates and reports if any roles were changed
func (r *UserRepo) NormalizeUserRolesVerbose(loginname string) (bool, error) {
	var user PdmUser
	if err := r.DB.Where("login_name = ?", loginname).First(&user).Error; err != nil {
		return false, err
	}

	originalRoles := make([]string, len(user.Roles))
	copy(originalRoles, user.Roles)

	var normalized []string
	for _, role := range user.Roles {
		lowerRole := strings.ToLower(role)
		switch lowerRole {
		case string(Admin), string(ProjectLead), string(SeniorDesigner), string(Designer), string(Approver), string(Qa), string(Editor), string(Viewer), string(Guest):
			normalized = append(normalized, lowerRole)
		default:
			fmt.Printf("[WARNING] Unknown role '%s' on user '%s' — preserving as-is\n", role, loginname)
			normalized = append(normalized, lowerRole)
		}
	}

	changed := !equalStringSlices(originalRoles, normalized)
	if changed {
		fmt.Printf("[INFO] Updating roles for user '%s': %v → %v\n", loginname, originalRoles, normalized)
		user.Roles = normalized
		if err := r.DB.Save(&user).Error; err != nil {
			return false, err
		}
	}

	return changed, nil
}

// NormalizeAllUsersVerbose runs over all users and reports counts
func (r *UserRepo) NormalizeAllUsersVerbose() (changedCount, totalCount int, err error) {
	var users []PdmUser
	if err := r.DB.Find(&users).Error; err != nil {
		return 0, 0, err
	}

	totalCount = len(users)
	for _, user := range users {
		changed, err := r.NormalizeUserRolesVerbose(user.LoginName)
		if err != nil {
			fmt.Printf("[ERROR] Failed to normalize roles for user '%s': %v\n", user.LoginName, err)
			continue
		}
		if changed {
			changedCount++
		}
	}

	return changedCount, totalCount, nil
}

// Helper: compare two string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// NormalizeUserRoles updates the roles on a user to match the canonical lowercase Role constants
func (r *UserRepo) NormalizeUserRoles(loginname string) error {
	var user PdmUser
	if err := r.DB.Where("login_name = ?", loginname).First(&user).Error; err != nil {
		return err
	}

	var normalized []string
	for _, role := range user.Roles {
		lowerRole := strings.ToLower(role)
		// Check if it matches a known constant
		switch lowerRole {
		case string(Admin), string(ProjectLead), string(SeniorDesigner), string(Designer), string(Approver), string(Qa), string(Editor), string(Viewer), string(Guest):
			normalized = append(normalized, lowerRole)
		default:
			fmt.Printf("[WARNING] Unknown role '%s' on user '%s' — preserving as-is\n", role, loginname)
			normalized = append(normalized, lowerRole)
		}
	}

	user.Roles = normalized
	return r.DB.Save(&user).Error
}

// NormalizeAllUsers walks over all users and normalizes their roles
func (r *UserRepo) NormalizeAllUsers() error {
	var users []PdmUser
	if err := r.DB.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		if err := r.NormalizeUserRoles(user.LoginName); err != nil {
			fmt.Printf("[ERROR] Failed to normalize roles for user '%s': %v\n", user.LoginName, err)
		}
	}
	return nil
}
