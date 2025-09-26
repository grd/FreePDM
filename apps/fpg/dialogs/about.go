// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dialogs

import (
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/grd/FreePDM/apps/fpg/chrome"
)

// Optionally set at build time: go build -ldflags "-X 'main.Version=0.1.0'"
var Version = "dev"

func BuildAboutBody(app fyne.App) fyne.CanvasObject {
	// Optional: use app icon if you have one set via app.SetIcon(...)
	icon := app.Icon()
	var iconObj fyne.CanvasObject
	if icon != nil {
		img := canvas.NewImageFromResource(icon)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(64, 64))
		iconObj = img
	} else {
		iconObj = widget.NewIcon(nil)
	}

	title := widget.NewLabelWithStyle("FreePDM", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	ver := widget.NewLabel("Version: " + Version)
	sys := widget.NewLabel("Go: " + runtime.Version() + " • OS/Arch: " + runtime.GOOS + "/" + runtime.GOARCH)

	text := widget.NewRichTextFromMarkdown(
		"FreePDM GUI (fpg)\n\n" +
			"- Cross-platform desktop client built with Fyne\n" +
			"- © 2025 — Gerard van de Schoot\n" +
			"- License: MIT (example)\n",
	)

	meta := container.NewVBox(title, ver, sys)
	header := container.NewHBox(iconObj, container.NewPadded(meta))
	return container.NewVBox(header, widget.NewSeparator(), text)
}

func ShowAboutWindow(app fyne.App) {
	body := BuildAboutBody(app)
	win := chrome.NewDialogWindow(app, "About", body, &chrome.DialogOptions{
		Width: 520, Height: 360,
		PanelMinWidth: 420, PanelMinHeight: 260,
		DimBackground: true,
		// PanelColor: optional custom
	})
	win.Show()
}
