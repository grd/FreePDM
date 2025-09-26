// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tabs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/BurntSushi/toml"
)

// ----------------------------------------------------------------------------
// Config loading
// ----------------------------------------------------------------------------

// appConfig represents the minimal config we need from ~/.config/fpg/config.toml.
type appConfig struct {
	RsyncTarget   string `toml:"RsyncTarget"`   // not used here yet
	LocalVaultDir string `toml:"LocalVaultDir"` // absolute path to local vaults root
}

// loadConfig reads ~/.config/fpg/config.toml and returns the parsed config.
// It expects LocalVaultDir to be an absolute path (per your decision).
func loadConfig() (*appConfig, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	cfgPath := filepath.Join(home, ".config", "fpg", "config.toml")

	cfg := &appConfig{}
	if _, err := os.Stat(cfgPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Config missing: return empty with path so caller can handle UX.
			return cfg, cfgPath, nil
		}
		return nil, "", fmt.Errorf("cannot access config: %w", err)
	}
	if _, err := toml.DecodeFile(cfgPath, cfg); err != nil {
		return nil, "", fmt.Errorf("invalid config TOML: %w", err)
	}
	return cfg, cfgPath, nil
}

// ----------------------------------------------------------------------------
// Vaults model & filesystem helpers
// ----------------------------------------------------------------------------

// Vault represents a single vault folder detected under LocalVaultDir.
type Vault struct {
	Name string // folder name
	Path string // absolute path
}

// listVaults returns immediate subdirectories in a given root as vaults.
// You can customize this (e.g., require a marker file like .freepdm) if desired.
func listVaults(root string) ([]Vault, error) {
	if strings.TrimSpace(root) == "" {
		return nil, errors.New("vault root is empty")
	}
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("cannot stat vault root: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("vault root is not a directory: %s", root)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("cannot read vault root: %w", err)
	}

	out := make([]Vault, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			p := filepath.Join(root, e.Name())
			out = append(out, Vault{Name: e.Name(), Path: p})
		}
	}

	// Sort alphabetically by name for stable UX
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out, nil
}

// ensureDir returns nil if dir exists and is a directory. Otherwise it returns an error.
// We don't auto-create here; we prefer explicit user action for a vault root.
func ensureDir(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return &fs.PathError{Op: "stat", Path: dir, Err: errors.New("path exists but is not a directory")}
	}
	return nil
}

// ----------------------------------------------------------------------------
// Vaults tab UI
// ----------------------------------------------------------------------------

// VaultsTab bundles state & widgets for the Vaults UI.
type VaultsTab struct {
	Tab      *container.TabItem
	RootPath string
	Vaults   []Vault
	OnOpen   func(v Vault) // callback when "Open Vault" is pressed

	// internal state
	selectedIndex int

	// widgets
	rootEntry   *widget.Entry
	vaultList   *widget.List
	detailName  *widget.Label
	detailPath  *widget.Entry
	openBtn     *widget.Button
	refreshBtn  *widget.Button
	browseBtn   *widget.Button
	statusLabel *widget.Label
}

// NewVaultsTab creates the Vaults tab.
// - win is used for dialogs.
// - onOpen is an optional callback invoked when the user presses "Open Vault".
func NewVaultsTab(win fyne.Window, onOpen func(v Vault)) *VaultsTab {
	t := &VaultsTab{OnOpen: onOpen, selectedIndex: -1}

	// Top: Root path controls
	t.rootEntry = widget.NewEntry()
	t.rootEntry.SetPlaceHolder("/home/user/My CAD Vault")
	t.rootEntry.OnSubmitted = func(s string) {
		t.RootPath = strings.TrimSpace(s)
		t.reload(win)
	}

	t.browseBtn = widget.NewButtonWithIcon("Browse…", theme.FolderOpenIcon(), func() {
		// Fyne directory open dialog
		d := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if uri == nil {
				return // user canceled
			}
			// Convert URI to file path
			u := uri.Path()
			if u == "" {
				u = uri.String()
			}
			if after, ok := strings.CutPrefix(u, "file://"); ok {
				u = after
			}
			t.rootEntry.SetText(u)
			t.RootPath = u
			t.reload(win)
		}, win)

		// Initial dir: if we already have a path, set it; else use home.
		start := t.RootPath
		if strings.TrimSpace(start) == "" {
			if home, err := os.UserHomeDir(); err == nil {
				start = home
			}
		}
		// SetLocation needs a ListableURI → obtain via ListerForURI.
		u := storage.NewFileURI(start)
		if l, err := storage.ListerForURI(u); err == nil {
			d.SetLocation(l)
		}
		d.Show()
	})

	t.refreshBtn = widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() { t.reload(win) })

	topBar := container.NewBorder(
		nil, nil, nil, t.refreshBtn,
		container.NewHBox(widget.NewLabel("Local Vault Root:"), t.rootEntry, t.browseBtn),
	)

	// Left: Vault list
	t.vaultList = widget.NewList(
		func() int { return len(t.Vaults) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.FolderIcon()),
				widget.NewLabel("vault"),
			)
		},
		func(i widget.ListItemID, co fyne.CanvasObject) {
			row := co.(*fyne.Container)
			lbl := row.Objects[1].(*widget.Label)
			lbl.SetText(t.Vaults[i].Name)
		},
	)
	t.vaultList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(t.Vaults) {
			return
		}
		t.selectedIndex = int(id)
		v := t.Vaults[id]
		t.detailName.SetText(v.Name)
		t.detailPath.SetText(v.Path)
		t.openBtn.Enable()
	}
	t.vaultList.OnUnselected = func(id widget.ListItemID) {
		t.selectedIndex = -1
		t.detailName.SetText("")
		t.detailPath.SetText("")
		t.openBtn.Disable()
	}

	// Right: Details and actions
	t.detailName = widget.NewLabel("")
	t.detailName.Wrapping = fyne.TextWrapWord

	t.detailPath = widget.NewEntry()
	t.detailPath.Disable()

	t.openBtn = widget.NewButtonWithIcon("Open Vault", theme.ConfirmIcon(), func() {
		id := t.selectedIndex
		if id < 0 || id >= len(t.Vaults) {
			return
		}
		if t.OnOpen != nil {
			t.OnOpen(t.Vaults[id])
		} else {
			dialog.ShowInformation("Open Vault", "No Open handler is wired yet.", win)
		}
	})
	t.openBtn.Disable()

	t.statusLabel = widget.NewLabel("")

	right := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Vault details", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Name:"),
		t.detailName,
		widget.NewLabel("Path:"),
		t.detailPath,
		t.openBtn,
		widget.NewSeparator(),
		t.statusLabel,
	)

	// Split view between list and details
	split := container.NewHSplit(t.vaultList, right)
	split.Offset = 0.45

	content := container.NewBorder(topBar, nil, nil, nil, split)
	t.Tab = container.NewTabItemWithIcon("Vaults", theme.FolderIcon(), content)

	// Initial load: from config.toml if present
	cfg, _, err := loadConfig()
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to read config: %w", err), win)
	} else if strings.TrimSpace(cfg.LocalVaultDir) != "" {
		t.RootPath = cfg.LocalVaultDir
		t.rootEntry.SetText(cfg.LocalVaultDir) // prefill path field from config
	}

	// Try initial reload (will show a friendly message if root is empty/invalid)
	t.reload(win)

	return t
}

// reload updates the vault list based on the current RootPath.
func (t *VaultsTab) reload(win fyne.Window) {
	root := strings.TrimSpace(t.RootPath)
	if root == "" {
		t.Vaults = nil
		t.vaultList.Refresh()
		t.statusLabel.SetText("No vault root set. Use Browse… or enter a path, then press Enter.")
		win.Canvas().Focus(t.rootEntry)
		return
	}

	// Validate directory exists
	if err := ensureDir(root); err != nil {
		t.Vaults = nil
		t.vaultList.Refresh()
		t.statusLabel.SetText("")
		dialog.ShowError(fmt.Errorf("vault root invalid: %w", err), win)
		win.Canvas().Focus(t.rootEntry)
		return
	}

	// Load vaults (subdirectories)
	vaults, err := listVaults(root)
	if err != nil {
		t.Vaults = nil
		t.vaultList.Refresh()
		t.statusLabel.SetText("")
		dialog.ShowError(err, win)
		win.Canvas().Focus(t.rootEntry)
		return
	}
	t.Vaults = vaults
	t.vaultList.Refresh()

	if n := len(vaults); n > 0 {
		t.statusLabel.SetText(fmt.Sprintf("%d vault(s) found.", n))
	} else {
		t.statusLabel.SetText("No vaults found.")
	}

	// Clear selection and details after refresh
	t.vaultList.UnselectAll()
	t.selectedIndex = -1
	t.detailName.SetText("")
	t.detailPath.SetText("")
	t.openBtn.Disable()
}
