// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package views

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/grd/FreePDM/apps/fpg/state"
)

func MakeVaultsView(w fyne.Window, state *state.AppState) fyne.CanvasObject {
	vaultList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("vault") },
		func(i widget.ListItemID, o fyne.CanvasObject) { /* set text */ },
	)

	refresh := widget.NewButton("Refresh", func() {
		go func() {
			// Non-blocking fetch, then update UI via main thread:
			vaults, err := state.API.ListVaults()
			_ = err
			_ = vaults
			// TODO: bind data to list (use binding if you like)
		}()
	})

	open := widget.NewButton("Open", func() {
		// TODO: navigate to file browser for selected vault
	})

	top := container.NewHBox(refresh, open)
	return container.NewBorder(top, nil, nil, nil, vaultList)
}

// Example: safe cancellation pattern for longer ops
func withCtx() (context.Context, context.CancelFunc) { return context.WithCancel(context.Background()) }
