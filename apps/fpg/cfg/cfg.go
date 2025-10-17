// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cfg

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Cfg struct {
	// GUI-specific settings (user-side)
	RsyncTarget     string `toml:"rsync_target"`    // e.g. user@host:/freepdm/vaults
	LocalVaultsRoot string `toml:"local_vault_dir"` // e.g. /home/user/My CAD Vaults
	VaultGroupUID   int    `toml:"vault_group_uid"` // e.g. 125

	// Optional: SSH/rsync tuning
	SSHKeyPath string `toml:"ssh_key_path,omitempty"`
	ExtraArgs  string `toml:"extra_args,omitempty"` // free rsync flags
}

func Default() *Cfg { return &Cfg{} }

// ~/.config/fpg/config.toml (Linux); OS-specific pad
func path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	cfgDir := filepath.Join(dir, "fpg")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(cfgDir, "config.toml"), nil
}

func Load() (*Cfg, error) {
	p, err := path()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return Default(), nil
	}
	if err != nil {
		return nil, err
	}
	var c Cfg
	if err := toml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func Save(c *Cfg) error {
	p, err := path()
	if err != nil {
		return err
	}
	b, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o644)
}
