// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import "time"

// Lock describes a pessimistic lock held by a user.
type Lock struct {
	Vault     string
	RelPath   string
	Holder    string
	Since     time.Time
	ExpiresAt time.Time // zero if unlimited/no lease
	Status    string    // "locked" | "stale" | "released"
	Note      string
}
