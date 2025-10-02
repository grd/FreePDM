// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package locks

import "github.com/grd/FreePDM/internal/domain/models"

// Service is the authoritative lock service for files.
type Service interface {
	Status(vault, rel string) (models.Lock, error)
	Checkout(vault, rel string, who models.UserIdentity, ttlMinutes int) (models.Lock, error)
	Heartbeat(vault, rel string, who models.UserIdentity) (models.Lock, error)
	Checkin(vault, rel string, who models.UserIdentity) error
	ForceUnlock(vault, rel string, admin models.UserIdentity, reason string) error
}
