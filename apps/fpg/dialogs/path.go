// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dialogs

import (
	"context"
	"errors"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/grd/FreePDM/apps/fpg/cfg"
	"github.com/grd/FreePDM/apps/fpg/state"
	"github.com/grd/FreePDM/internal/sync"
)

// BuildPathSettingsView creates a reusable Path Settings UI.
// You can use the returned CanvasObject inside a tab, a dialog, or a separate window.
func BuildPathSettingsView(parent fyne.Window, st *state.AppState) fyne.CanvasObject {
	// Entries with current values (if any)
	rsyncEntry := widget.NewEntry()
	localEntry := widget.NewEntry()
	if st.Cfg != nil {
		rsyncEntry.SetText(st.Cfg.RsyncTarget)
		localEntry.SetText(st.Cfg.LocalVaultDir)
	}
	rsyncEntry.SetPlaceHolder("user@host:/srv/freepdm/vaults/Main   (or a local path)")
	localEntry.SetPlaceHolder("/home/you/FreePDM/vaults/Main")

	// Folder picker for local vault directory
	browseBtn := widget.NewButton("Choose…", func() {
		fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				localEntry.SetText(uri.Path())
			}
		}, parent)
		fd.Show()
	})

	// Test remote/target (non-blocking)
	testLabel := widget.NewLabel("") // feedback label
	testBtn := widget.NewButton("Test", func() {
		testLabel.SetText("Testing…")
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
			defer cancel()

			// NOTE: For a simple test we use plain rsync (no forced -e ssh).
			// TestTarget signature here is: TestTarget(ctx, source string, extra []string)
			target := strings.TrimSpace(rsyncEntry.Text)

			if err := sync.TestTarget(ctx, target, nil); err != nil {
				fyne.Do(func() { testLabel.SetText("Unreachable or invalid") })
				return
			}
			fyne.Do(func() { testLabel.SetText("OK") })
		}()
	})

	// Save handler (validates and persists to ~/.config/fpg/config.toml)
	save := func() {
		config := &cfg.Cfg{
			RsyncTarget:   strings.TrimSpace(rsyncEntry.Text),
			LocalVaultDir: strings.TrimSpace(localEntry.Text),
		}
		if config.RsyncTarget == "" {
			dialog.ShowError(errors.New("rsync target cannot be empty"), parent)
			return
		}
		// Minimal validation: accept local paths or rsync-style remote (user@host:/path)
		// If you want to require user@host:/path only, add: !strings.Contains(cfg.RsyncTarget, ":")
		if config.LocalVaultDir == "" {
			dialog.ShowError(errors.New("local vault folder cannot be empty"), parent)
			return
		}
		if err := cfg.Save(config); err != nil {
			dialog.ShowError(err, parent)
			return
		}
		st.SetConfig(config)
		dialog.ShowInformation("Settings", "Saved", parent)
	}

	// Build the form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Rsync target", Widget: rsyncEntry},
			{
				Text:   "Local vault folder",
				Widget: container.NewBorder(nil, nil, nil, browseBtn, localEntry),
			},
		},
		SubmitText: "Save",
		CancelText: "Close",
		OnSubmit:   save,
	}

	// Final layout (+ make it comfortably large)
	content := container.NewVBox(
		widget.NewLabelWithStyle("Path Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		form,
		container.NewHBox(testBtn, testLabel),
	)
	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(720, 520))
	return scroll
}

// ShowPathDialog shows the Path Settings inside a modal dialog.
func ShowPathDialog(w fyne.Window, st *state.AppState) {
	content := BuildPathSettingsView(w, st)
	dialog.NewCustom("Path Settings", "Close", content, w).Show()
}

// ShowPathWindow opens the Path Settings in a new, movable window.
func ShowPathWindow(st *state.AppState) {
	win := fyne.CurrentApp().NewWindow("Path Settings")
	win.SetContent(BuildPathSettingsView(win, st))
	win.Resize(fyne.NewSize(900, 600))
	win.Show()
}
