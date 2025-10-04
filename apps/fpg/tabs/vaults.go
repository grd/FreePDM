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
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/grd/FreePDM/apps/fpg/cfg"
)

func loadConfig() (*cfg.Cfg, error) {
	config, err := cfg.Load()
	if err != nil {
		return nil, fmt.Errorf("invalid config TOML: %w", err)
	}
	return config, nil
}

// ----------------------------------------------------------------------------
// Vaults model & filesystem helpers
// ----------------------------------------------------------------------------

// Vault represents a single vault folder detected under LocalVaultDir.
type Vault struct {
	Name string // folder name
	Path string // absolute path
}

// listVaults returns immediate subdirectories in a given root as vaults,
// skipping hidden directories that start with a dot.
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
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		// Skip dot-prefixed directories (hidden)
		if strings.HasPrefix(name, ".") {
			continue
		}
		p := filepath.Join(root, name)
		out = append(out, Vault{Name: name, Path: p})
	}

	// Sort alphabetically by name for stable UX
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
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
	lastSelIdx    int

	// widgets
	rootEntry   *widget.Entry
	vaultList   *widget.List
	detailName  *widget.Label
	detailPath  *widget.Entry
	openBtn     *widget.Button
	refreshBtn  *widget.Button
	statusLabel *widget.Label
}

// NewVaultsTab creates the Vaults tab.
// - win is used for dialogs.
// - onOpen is an optional callback invoked when the user presses "Open Vault".
func NewVaultsTab(win fyne.Window, onOpen func(v Vault)) *VaultsTab {
	t := &VaultsTab{OnOpen: onOpen, selectedIndex: -1, lastSelIdx: -1}

	// Top: Root path controls
	t.rootEntry = widget.NewEntry()
	t.rootEntry.SetPlaceHolder("Local Vault Root: /home/user/My CAD Vaults")

	t.rootEntry.OnSubmitted = func(s string) {
		t.RootPath = strings.TrimSpace(s)
		t.reload(win)
	}

	t.refreshBtn = widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() { t.reload(win) })

	// Entry vult nu het hele midden; alleen de refresh-knop staat rechts
	topBar := container.NewBorder(
		nil, nil, // top, bottom
		nil, t.refreshBtn, // left, right
		container.NewStack(t.rootEntry), // center (vult breedte)
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
		v := t.Vaults[id]

		openHere := widget.NewButtonWithIcon("Open in FreePDM", theme.FolderOpenIcon(), func() {
			if t.OnOpen != nil {
				t.OnOpen(t.Vaults[id]) // <- roept NewVaultTab via je callback
			}
		})

		content := container.NewVBox(
			widget.NewLabelWithStyle("Vault", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(v.Name),
			widget.NewSeparator(),
			widget.NewLabelWithStyle("Path", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(v.Path),
			openHere,
		)

		dialog.NewCustomConfirm(
			"Open Vault",
			"Open", "Close",
			content,
			func(open bool) {
				if open && t.OnOpen != nil {

					t.OnOpen(v)
				}
			},
			win,
		).Show()

		// Unselect zodat volgende klik weer werkt als nieuwe selectie
		t.vaultList.UnselectAll()
		t.detailName.SetText("")
		t.detailPath.SetText("")
	}

	t.vaultList.OnUnselected = func(widget.ListItemID) {
		// no-op
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
	t.openBtn.Hide()

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
	cfg, err := loadConfig()
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to read config: %w", err), win)
	} else if strings.TrimSpace(cfg.LocalVaultDir) != "" {
		t.RootPath = cfg.LocalVaultDir
		t.rootEntry.SetText(cfg.LocalVaultDir) // prefill path field from config
	}

	// Try initial reload (will show a friendly message if root is empty/invalid)
	t.reload(win)

	// Enter (Return) opens the currently "selected" vault (tracked in selectedIndex)
	win.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyReturn,
		Modifier: 0,
	}, func(fyne.Shortcut) {
		id := t.selectedIndex
		if id < 0 || id >= len(t.Vaults) {
			return
		}
		if t.OnOpen != nil {
			t.OnOpen(t.Vaults[id])
		}
		t.selectedIndex = -1
		t.vaultList.UnselectAll()
		t.detailName.SetText("")
		t.detailPath.SetText("")
	})

	// Some systems report Enter as KeyEnter; bind both.
	win.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyEnter,
		Modifier: 0,
	}, func(fyne.Shortcut) {
		id := t.selectedIndex
		if id < 0 || id >= len(t.Vaults) {
			return
		}
		if t.OnOpen != nil {
			t.OnOpen(t.Vaults[id])
		}
		t.selectedIndex = -1
		t.vaultList.UnselectAll()
		t.detailName.SetText("")
		t.detailPath.SetText("")
	})

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
