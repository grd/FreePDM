// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dialogs

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// ComposeOptions configures the behavior/layout of the compose dialog.
type ComposeOptions struct {
	// If true, show a "Location" row. If false, the row is omitted and filePath in onSubmit will be "".
	ShowLocation bool

	// If ShowLocation is true and AllowPickLocation is true, a "Browse…" button appears (file-save dialog).
	// If false, no picker button is shown (the user can still type/paste a path if the entry is visible).
	AllowPickLocation bool

	// If ShowLocation is true and EnforceNameMatch is true, the basename of the chosen/typed path
	// must match the provided required file name (read-only field).
	EnforceNameMatch bool

	// Optional initial path value for the Location entry (e.g., "<abs>/0/<filename>").
	InitialPath string

	// Optional dialog title and submit button label; if empty, sensible defaults are used.
	DialogTitle string // default: "Compose descriptions"
	SubmitLabel string // default: "Write files"
}

// ShowComposeDescriptions opens a configurable dialog to capture Title/Details,
// with an optional Location selector, while showing a fixed File name.
// onSubmit receives the selected/typed file path (or "" when ShowLocation=false) and the two descriptions.
func ShowComposeDescriptions(
	win fyne.Window,
	requiredFilename string, // shown read-only; caller provides the exact expected name
	opts ComposeOptions,
	onSubmit func(filePath, shortText, longText string),
) {
	if strings.TrimSpace(requiredFilename) == "" {
		dialog.ShowError(errors.New("missing required file name"), win)
		return
	}
	title := opts.DialogTitle
	if title == "" {
		title = "Compose descriptions"
	}
	okLabel := opts.SubmitLabel
	if okLabel == "" {
		okLabel = "Write files"
	}

	// --- File name (read-only)
	fnEntry := widget.NewEntry()
	fnEntry.SetText(requiredFilename)
	fnEntry.Disable()

	// --- Optional Location row
	var locEntry *widget.Entry
	var locationItem *widget.FormItem

	if opts.ShowLocation {
		locEntry = widget.NewEntry()
		locEntry.SetPlaceHolder("Pick or enter target file path…")
		if opts.InitialPath != "" {
			locEntry.SetText(opts.InitialPath)
		}

		// Picker button (optional) placed to the RIGHT of the entry
		var content fyne.CanvasObject = locEntry
		if opts.AllowPickLocation {
			pickBtn := widget.NewButton("Browse…", func() {
				save := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
					if err != nil {
						dialog.ShowError(err, win)
						return
					}
					if uc == nil {
						return // canceled
					}
					defer uc.Close()

					u := uc.URI()
					if u == nil {
						dialog.ShowError(errors.New("invalid selection (no URI)"), win)
						return
					}
					path := u.Path()

					// Optionally enforce that the basename equals the required filename
					if opts.EnforceNameMatch && filepath.Base(path) != requiredFilename {
						dialog.ShowError(fmt.Errorf(
							"selected file name must be %q (got %q)",
							requiredFilename, filepath.Base(path),
						), win)
						return
					}
					locEntry.SetText(path)
				}, win)

				// Pre-fill filename and extension filter (if any)
				save.SetFileName(requiredFilename)
				if ext := strings.ToLower(filepath.Ext(requiredFilename)); ext != "" {
					save.SetFilter(storage.NewExtensionFileFilter([]string{ext}))
				}
				save.Show()
			})

			// Entry expands; pick button on the right
			content = container.NewBorder(nil, nil, nil, pickBtn, locEntry)
		}

		locationItem = widget.NewFormItem("Location", content)
	}

	// --- Title & Details (empty by default; no prefill)
	short := widget.NewEntry()
	short.SetPlaceHolder("Short description…") // empty prefill
	short.Validator = func(s string) error {
		if strings.TrimSpace(s) == "" {
			return errors.New("short description is required")
		}
		return nil
	}

	long := widget.NewMultiLineEntry()
	long.SetPlaceHolder("Long description…") // empty prefill
	long.Wrapping = fyne.TextWrapWord
	long.SetMinRowsVisible(8)

	// Build the form-items list
	items := []*widget.FormItem{
		widget.NewFormItem("File name", fnEntry),
	}
	if locationItem != nil {
		items = append(items, locationItem)
	}
	items = append(items,
		widget.NewFormItem("Title", short),
		widget.NewFormItem("Details", long),
	)

	// Show dialog
	form := dialog.NewForm(
		title,
		okLabel,
		"Cancel",
		items,
		func(ok bool) {
			if !ok {
				return
			}
			// Validate fields
			if err := short.Validate(); err != nil {
				dialog.ShowError(err, win)
				return
			}
			// Validate Location if present
			filePath := ""
			if opts.ShowLocation {
				path := strings.TrimSpace(locEntry.Text)
				if path == "" {
					dialog.ShowError(errors.New("please specify a file location"), win)
					return
				}
				if opts.EnforceNameMatch && filepath.Base(path) != requiredFilename {
					dialog.ShowError(fmt.Errorf(
						"location's file name must match %q (got %q)",
						requiredFilename, filepath.Base(path),
					), win)
					return
				}
				filePath = path
			}

			onSubmit(filePath, short.Text, long.Text)
		},
		win,
	)
	form.Resize(fyne.NewSize(760, 600))
	form.Show()
}
