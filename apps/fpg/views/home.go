// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MakeHomeView() fyne.CanvasObject {
	return container.NewCenter(widget.NewLabel("Welcome to FreePDM"))
}
