// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sync

import (
	"time"

	"github.com/grd/FreePDM/internal/domain/models"
)

// SSHOptions configures SSH transport for rsync (or similar).
type SSHOptions struct {
	Port         int
	IdentityFile string
	KnownHosts   string
	ExtraArgs    []string
}

// Report summarizes a sync operation.
type Report struct {
	StartedAt time.Time
	EndedAt   time.Time
	Changed   int
	Bytes     int64
	ExitCode  int
	Log       []string
}

// Transport moves blobs; authorization/locking is enforced elsewhere.
type Transport interface {
	SetTarget(target string) // e.g., "user@host:/srv/freepdm/vaults"
	SetSSHOptions(opts SSHOptions)

	DryRunPush(vault models.VaultInfo, srcRel string) (Report, error)
	Push(vault models.VaultInfo, srcRel string) (Report, error)
	Pull(vault models.VaultInfo, rel string) (Report, error)
}
