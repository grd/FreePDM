// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package vaults

import (
	"context"

	"github.com/grd/FreePDM/internal/domain/models"
	"github.com/grd/FreePDM/internal/ports/vaultfs"
)

// Manager defines the contract for managing a set of vaults and
// performing cross-vault operations. It is intentionally small:
// discovery/registry + copy/move between vaults.
//
// Notes
// - The GUI only depends on this interface (and models).
// - Implementations may be in-memory, config-backed, or remote.
// - Cross-vault operations MUST enforce policy (e.g., destination RO).
type Manager interface {
	// List returns metadata for all known vaults in a stable order.
	List() ([]models.VaultInfo, error)

	// Get returns a VaultFS by name, if present.
	// The returned FS must remain valid for the lifetime documented
	// by the implementation (e.g., until CloseAll or process exit).
	Get(name string) (vaultfs.FS, error)

	// Add registers a new vault implementation under its Info().Name.
	// If a vault with the same name exists, it should be replaced only
	// if the implementation documents that behavior (otherwise return an error).
	Add(ctx context.Context, v vaultfs.FS) error

	// Remove unregisters a vault by name.
	// Implementations may also release resources associated with this vault.
	Remove(ctx context.Context, name string) error

	// Copy performs a cross-vault copy operation:
	//   from: source vault
	//   srcRel: path relative to source vault root
	//   to: destination vault
	//   dstRel: destination path relative to destination vault root
	//
	// Requirements (to be enforced by the implementation):
	// - Destination must not be read-only.
	// - Existing destination path behavior (overwrite/merge/fail) must be documented.
	// - Large directory trees should be processed incrementally (streamed) if possible.
	Copy(ctx context.Context, from vaultfs.FS, srcRel string, to vaultfs.FS, dstRel string) error

	// Move performs a cross-vault move operation.
	// Implementations may optimize (rename when same underlying store)
	// or fallback to Copy + Delete semantics across devices or boundaries.
	Move(ctx context.Context, from vaultfs.FS, srcRel string, to vaultfs.FS, dstRel string) error

	// Root returns the root directory of the vaults
	Root() string
}
