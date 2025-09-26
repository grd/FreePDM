// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
)

type TabManager struct {
	win  fyne.Window
	Tabs *container.DocTabs
	curr *container.TabItem // track current tab ourselves
}

func NewTabManager(win fyne.Window) *TabManager {
	tm := &TabManager{
		win:  win,
		Tabs: container.NewDocTabs(),
	}
	tm.Tabs.SetTabLocation(container.TabLocationTop)

	// Keep track of the current tab
	tm.Tabs.OnSelected = func(ti *container.TabItem) {
		tm.curr = ti
	}
	// Optional: intercept close or react on close
	// tm.Tabs.OnClosed = func(ti *container.TabItem) {}

	// Hotkeys
	c := win.Canvas()

	// Ctrl+W — close current
	c.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyW,
		Modifier: fyne.KeyModifierControl,
	}, func(fyne.Shortcut) { tm.CloseCurrent() })

	// Ctrl+Tab — next
	c.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyTab,
		Modifier: fyne.KeyModifierControl,
	}, func(fyne.Shortcut) { tm.NextTab() })

	// Ctrl+Shift+Tab — previous
	c.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyTab,
		Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
	}, func(fyne.Shortcut) { tm.PrevTab() })

	// Ctrl+1..9 — select by index (1-based)
	keys := []fyne.KeyName{
		fyne.Key1, fyne.Key2, fyne.Key3, fyne.Key4, fyne.Key5,
		fyne.Key6, fyne.Key7, fyne.Key8, fyne.Key9,
	}
	for i, k := range keys {
		idx := i // capture
		c.AddShortcut(&desktop.CustomShortcut{
			KeyName:  k,
			Modifier: fyne.KeyModifierControl,
		}, func(fyne.Shortcut) { tm.SelectIndex(idx) })
	}

	return tm
}

func BuildMainMenu(win fyne.Window, tm *TabManager) *fyne.MainMenu {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Open Vaults", func() {
			vt := NewVaultsTab(win, func(v Vault) {
				// TODO: replace with real open logic
				dialog.ShowInformation("Open Vault", "Opening: "+v.Path, win)
			})
			tm.AddTabItem(vt.Tab)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Close Tab", func() { tm.CloseCurrent() }),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() { win.Close() }),
	)

	viewMenu := fyne.NewMenu("View",
		fyne.NewMenuItem("Next Tab", func() { tm.NextTab() }),
		fyne.NewMenuItem("Previous Tab", func() { tm.PrevTab() }),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {
			dialog.ShowInformation("About", "FreePDM GUI (fpg) — alpha prototype", win)
		}),
	)

	return fyne.NewMainMenu(fileMenu, viewMenu, helpMenu)
}

func (tm *TabManager) AddTab(title string, content fyne.CanvasObject) *container.TabItem {
	ti := container.NewTabItem(title, content)
	tm.Tabs.Append(ti)
	tm.Tabs.Select(ti)
	tm.curr = ti
	return ti
}

func (tm *TabManager) AddTabItem(ti *container.TabItem) *container.TabItem {
	tm.Tabs.Append(ti)
	tm.Tabs.Select(ti)
	tm.curr = ti
	return ti
}

func (tm *TabManager) CloseCurrent() {
	if tm.Tabs == nil || len(tm.Tabs.Items) == 0 {
		return
	}
	ti := tm.curr
	if ti == nil {
		// fallback: if nothing tracked, try first tab
		if len(tm.Tabs.Items) == 0 {
			return
		}
		ti = tm.Tabs.Items[0]
	}
	tm.Tabs.Remove(ti)
	// After remove, Fyne will select something; OnSelected will update tm.curr.
}

func (tm *TabManager) NextTab() {
	items := tm.Tabs.Items
	if len(items) <= 1 {
		return
	}
	idx := tm.indexOf(tm.curr)
	if idx < 0 {
		tm.Tabs.SelectIndex(0)
		return
	}
	tm.Tabs.SelectIndex((idx + 1) % len(items))
	// OnSelected callback updates tm.curr
}

func (tm *TabManager) PrevTab() {
	items := tm.Tabs.Items
	if len(items) <= 1 {
		return
	}
	idx := tm.indexOf(tm.curr)
	if idx < 0 {
		tm.Tabs.SelectIndex(0)
		return
	}
	if idx == 0 {
		tm.Tabs.SelectIndex(len(items) - 1)
		return
	}
	tm.Tabs.SelectIndex(idx - 1)
	// OnSelected updates tm.curr
}

func (tm *TabManager) SelectIndex(i int) {
	if i < 0 || i >= len(tm.Tabs.Items) {
		return
	}
	tm.Tabs.SelectIndex(i)
	// OnSelected updates tm.curr
}

func (tm *TabManager) CloseByTitle(title string) bool {
	for _, ti := range tm.Tabs.Items {
		if ti.Text == title {
			tm.Tabs.Remove(ti)
			return true
		}
	}
	return false
}

func (tm *TabManager) indexOf(ti *container.TabItem) int {
	if ti == nil {
		return -1
	}
	for i, it := range tm.Tabs.Items {
		if it == ti {
			return i
		}
	}
	return -1
}
