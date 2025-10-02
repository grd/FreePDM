// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

// VaultInfo describes one logical vault.
type VaultInfo struct {
	Name     string            // logical vault name/id
	RootAbs  string            // optional: absolute local root (for UI only)
	ReadOnly bool              // enforce write protection at vault level
	Labels   map[string]string // arbitrary tags, e.g. {"site":"AMS"}
}
