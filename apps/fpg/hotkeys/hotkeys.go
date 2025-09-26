// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hotkeys

import (
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type Handlers struct {
	CloseTab func()
	NextTab  func()
	PrevTab  func()
}

func shortcutMod() fyne.KeyModifier {
	if runtime.GOOS == "darwin" {
		return fyne.KeyModifierSuper // Cmd on macOS
	}
	return fyne.KeyModifierControl // Ctrl on Win/Linux/BSD
}

func Bind(w fyne.Window, h Handlers) {
	c := w.Canvas()

	add := func(key fyne.KeyName, mod fyne.KeyModifier, fn func()) {
		if fn == nil {
			return
		}
		c.AddShortcut(&desktop.CustomShortcut{
			KeyName:  key,
			Modifier: mod,
		}, func(_ fyne.Shortcut) { fn() })
	}

	mod := shortcutMod()

	// Close tab: Cmd/Ctrl+W
	add(fyne.KeyW, mod, h.CloseTab)

	// Next tab: Cmd/Ctrl+Tab
	add(fyne.KeyTab, mod, h.NextTab)

	// Previous tab: Cmd/Ctrl+Shift+Tab
	add(fyne.KeyTab, mod|fyne.KeyModifierShift, h.PrevTab)

	// Optional extras common on Linux/Windows
	add(fyne.KeyPageDown, fyne.KeyModifierControl, h.NextTab)
	add(fyne.KeyPageUp, fyne.KeyModifierControl, h.PrevTab)
}
