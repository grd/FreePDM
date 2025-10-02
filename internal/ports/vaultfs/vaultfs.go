// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package vaultfs

import "github.com/grd/FreePDM/internal/domain/models"

// FS abstracts a single vault's file system operations.
type FS interface {
	Info() models.VaultInfo

	// Browsing
	List(dirRel string) ([]models.Entry, error)
	Stat(rel string) (models.Entry, error)

	// Intra-vault operations
	Rename(srcRel, dstRel string) error
	Move(srcRel, dstRel string) error
	Copy(srcRel, dstRel string) error
	Delete(rel string) error

	// Optional read/write if needed later:
	// ReadFile(rel string) ([]byte, error)
	// WriteFile(rel string, data []byte) error
}
