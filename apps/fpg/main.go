// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/grd/FreePDM/apps/fpg/hotkeys"
	"github.com/grd/FreePDM/apps/fpg/state"
	"github.com/grd/FreePDM/apps/fpg/ui"
	"github.com/grd/FreePDM/apps/fpg/views"
	"github.com/grd/FreePDM/internal/fpg/config"
)

func main() {
	a := app.NewWithID("org.freepdm.gui")
	w := a.NewWindow("FreePDM")
	w.SetMaster()

	// App state
	st := state.New()
	if cfg, err := config.Load(); err == nil {
		st.SetConfig(cfg)
	}

	// MDI with tabs
	tm := ui.NewTabManager()
	w.SetContent(tm.Tabs)

	hotkeys.Bind(w, hotkeys.Handlers{
		CloseTab: tm.CloseSelected,
		NextTab:  tm.SelectNext,
		PrevTab:  tm.SelectPrev,
	})
	// Menubar (needs tm so menu actions can open tabs)
	w.SetMainMenu(ui.NewMainMenu(w, st, tm))

	// Optional: start with a simple Home tab
	tm.AddTab("Home", views.MakeHomeView())

	w.Resize(fyne.NewSize(1000, 700))
	w.ShowAndRun()
}
