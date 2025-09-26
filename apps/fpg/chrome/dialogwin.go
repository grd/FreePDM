// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chrome

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DialogOptions struct {
	Width, Height                 float32
	PanelMinWidth, PanelMinHeight float32
	DimBackground                 bool
	Footer                        fyne.CanvasObject
	HeaderLeft, HeaderRight       fyne.CanvasObject
}

// NewDialogWindow creates a movable, resizable window that looks like a dialog.
func NewDialogWindow(app fyne.App, title string, body fyne.CanvasObject, opts *DialogOptions) fyne.Window {
	if opts == nil {
		opts = &DialogOptions{
			Width: 720, Height: 480,
			PanelMinWidth: 420, PanelMinHeight: 240,
			DimBackground: true,
		}
	}
	win := app.NewWindow(title)

	panel := buildDialogPanel(title, body, opts)
	scroll := container.NewScroll(panel)
	scroll.SetMinSize(fyne.NewSize(opts.PanelMinWidth, opts.PanelMinHeight))

	var content fyne.CanvasObject
	if opts.DimBackground {
		dimmer := canvas.NewRectangle(color.NRGBA{A: 82}) // ~32% black
		content = container.NewStack(
			dimmer,
			container.NewCenter(scroll),
		)
	} else {
		content = container.NewCenter(scroll)
	}

	win.SetContent(content)
	win.Resize(fyne.NewSize(opts.Width, opts.Height))
	return win
}

// BuildDialogSurface returns the dialog-like surface as a CanvasObject (for tabs, etc.).
func BuildDialogSurface(_ fyne.App, title string, body fyne.CanvasObject, opts *DialogOptions) fyne.CanvasObject {
	if opts == nil {
		opts = &DialogOptions{
			PanelMinWidth: 420, PanelMinHeight: 240,
			DimBackground: false,
		}
	}
	panel := buildDialogPanel(title, body, opts)
	scroll := container.NewScroll(panel)
	scroll.SetMinSize(fyne.NewSize(opts.PanelMinWidth, opts.PanelMinHeight))

	if opts.DimBackground {
		dimmer := canvas.NewRectangle(color.NRGBA{A: 82})
		return container.NewStack(dimmer, container.NewCenter(scroll))
	}
	return container.NewCenter(scroll)
}

// ---- internals ----

func buildDialogPanel(title string, body fyne.CanvasObject, opts *DialogOptions) fyne.CanvasObject {
	// Theme-aware colors (modern API: single-arg theme.Color)
	bgCol := theme.Color(theme.ColorNameBackground)

	titleLbl := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleLbl.Wrapping = fyne.TextWrapWord

	header := container.NewBorder(nil, nil, opts.HeaderLeft, opts.HeaderRight, titleLbl)

	inner := container.NewVBox(
		header,
		widget.NewSeparator(),
		container.NewPadded(body),
	)
	if opts.Footer != nil {
		inner = container.NewVBox(
			header,
			widget.NewSeparator(),
			container.NewPadded(body),
			widget.NewSeparator(),
			container.NewPadded(opts.Footer),
		)
	}

	// Background rectangle using current theme background color
	bg := canvas.NewRectangle(bgCol)

	// Padded panel for a subtle card-like look
	panel := container.NewStack(
		bg,
		container.NewPadded(inner),
	)

	// Outer padding so the panel never touches window edges
	return container.NewPadded(panel)
}
