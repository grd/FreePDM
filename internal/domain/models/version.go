// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import "time"

// Version links a content hash to an authored revision.
type Version struct {
	ID          string
	Vault       string
	RelPath     string
	ContentHash string
	Size        int64
	Author      string
	CreatedAt   time.Time
	Label       string // e.g. "WIP", "Released"
}
