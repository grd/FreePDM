// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// NewHomeTab creates a basic Home tab.
// onOpenVaults is optional; if provided, it will be called when the button is pressed.
func NewHomeTab(onOpenVaults func()) *container.TabItem {
	title := widget.NewLabelWithStyle("FreePDM", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	sub := widget.NewLabel("Welcome — this is the Home page. Use the button below to open your vaults.")

	openBtn := widget.NewButtonWithIcon("Open Vaults", theme.FolderOpenIcon(), func() {
		if onOpenVaults != nil {
			onOpenVaults()
		}
	})

	content := container.NewVBox(
		title,
		sub,
		widget.NewSeparator(),
		openBtn,
	)

	return container.NewTabItemWithIcon("Home", theme.HomeIcon(), container.NewPadded(content))
}
