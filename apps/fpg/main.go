// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/grd/FreePDM/apps/fpg/tabs"
)

func main() {
	a := app.New()
	w := a.NewWindow("FreePDM")

	// 1) Eerst de TabManager maken
	tm := tabs.NewTabManager(w)

	// 2) Main menu pas daarna (menu-acties gebruiken tm)
	w.SetMainMenu(tabs.BuildMainMenu(w, tm))

	// 3) Content = AppTabs uit de manager
	w.SetContent(container.NewBorder(nil, nil, nil, nil, tm.Tabs))

	// 4) Optioneel: start met Home tab
	home := tabs.NewHomeTab(func() {
		vt := tabs.NewVaultsTab(w, func(v tabs.Vault) {
			// TODO: echte open-logic hier
		})
		tm.AddTabItem(vt.Tab)
	})
	tm.AddTabItem(home)

	w.Resize(fyne.NewSize(1100, 700))
	w.ShowAndRun()
}
