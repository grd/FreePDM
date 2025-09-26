// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ui

import (
	"fyne.io/fyne/v2"

	"github.com/grd/FreePDM/apps/fpg/dialogs"
	"github.com/grd/FreePDM/apps/fpg/state"
	"github.com/grd/FreePDM/apps/fpg/views"
)

// NewMainMenu builds the main menu and wires actions to dialogs/views.
// Note: we pass tm so menu items can open tabs.
func NewMainMenu(w fyne.Window, st *state.AppState, tm *TabManager) *fyne.MainMenu {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Login…", func() {
			dialogs.ShowLoginWindow(st, func() {
				tm.AddTab("Vaults", views.MakeVaultsView(w, st))
			})
		}),
		fyne.NewMenuItem("Open Vaults", func() {
			tm.AddTab("Vaults", views.MakeVaultsView(w, st))
		}),
		fyne.NewMenuItem("Settings…", func() {
			dialogs.ShowPathDialog(w, st)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() { w.Close() }),
	)

	windowMenu := fyne.NewMenu("Window") // placeholders for future MDI actions (Next/Prev/Close Tab, etc.)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About…", func() {
			dialogs.ShowAboutWindow(fyne.CurrentApp())
		}),
	)

	return fyne.NewMainMenu(fileMenu, windowMenu, helpMenu)
}
