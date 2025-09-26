// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dialogs

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/grd/FreePDM/apps/fpg/state"
	"github.com/grd/FreePDM/internal/client"
)

// ShowLoginWindow opens a movable, resizable login window.
func ShowLoginWindow(st *state.AppState, onSuccess func()) {
	win := fyne.CurrentApp().NewWindow("Login")

	user := widget.NewEntry()
	pass := widget.NewPasswordEntry()
	status := widget.NewLabel("")

	user.SetPlaceHolder("Username")
	pass.SetPlaceHolder("Password")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: user},
			{Text: "Password", Widget: pass},
		},
		SubmitText: "Login",
		CancelText: "Close",
		OnSubmit: func() {
			if user.Text == "" || pass.Text == "" {
				dialog.ShowError(errors.New("please enter username and password"), win)
				return
			}
			status.SetText("Signing in…")

			// TODO: read base URL from st.Cfg later if you store it there.
			api := client.New("http://localhost:8080")

			go func() {
				if err := api.Login(user.Text, pass.Text); err != nil {
					fyne.Do(func() {
						status.SetText("")
						dialog.ShowError(err, win)
					})
					return
				}
				fyne.Do(func() {
					st.SetAPI(api)
					st.User = user.Text
					status.SetText("Signed in")
					win.Close()
					if onSuccess != nil {
						onSuccess()
					}
				})
			}()
		},
	}

	content := container.NewVBox(
		widget.NewLabelWithStyle("Sign in to FreePDM", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		form,
		status,
	)

	win.SetContent(content)
	win.Resize(fyne.NewSize(520, 360)) // window is movable/resizable
	win.Show()
}
