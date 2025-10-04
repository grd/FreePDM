// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package state

import (
	"sync"

	"github.com/grd/FreePDM/apps/fpg/cfg"
	"github.com/grd/FreePDM/internal/client"
)

type AppState struct {
	mu   sync.RWMutex
	API  *client.API
	User string
	Cfg  *cfg.Cfg

	// Example domain state
	vaults []client.Vault
}

func New() *AppState { return &AppState{Cfg: cfg.Default()} }

func (s *AppState) SetAPI(api *client.API) { s.mu.Lock(); s.API = api; s.mu.Unlock() }

func (s *AppState) SetConfig(c *cfg.Cfg) { s.mu.Lock(); s.Cfg = c; s.mu.Unlock() }

func (s *AppState) SetVaults(v []client.Vault) {
	s.mu.Lock()
	s.vaults = append([]client.Vault(nil), v...)
	s.mu.Unlock()
}
func (s *AppState) Vaults() []client.Vault {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]client.Vault(nil), s.vaults...)
}
