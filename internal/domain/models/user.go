// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

// UserIdentity identifies the current user and capabilities.
type UserIdentity struct {
	UserID      string
	DisplayName string
	Groups      []string
	Authz       map[string]bool // coarse roles/claims
}
