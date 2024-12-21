// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/grd/FreePDM/internal/config"
)

// Script for user management.

type userProp struct {
	userName string
	userUid  int
}

var pages = tview.NewPages()
var form = tview.NewForm()
var flex = tview.NewFlex()

func setup() {

	// if config.Conf.Users == nil {
	// 	config.Conf.Users = make(map[string]int)
	// }

	// var app = tview.NewApplication()

	// var text = tview.NewTextView().
	// 	SetTextColor(tcell.ColorGreen).
	// 	SetText("(q) to quit")

	// app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	if event.Rune() == 113 {
	// 		app.Stop()
	// 	}
	// 	return event
	// })

	// err := app.SetRoot(text, true).EnableMouse(true).Run()
	// ex.CheckErr(err)

	// // if err := app.SetRoot(tview.NewBox(), true).EnableMouse(true).Run(); err != nil {
	// // 	panic(err)
	// // }

	// form.AddInputField("First Name", "", 20, nil, func(firstName string) {
	// 	contact.firstName = firstName
	// })

	// Printing the String method
	fmt.Printf("%#v\n", config.Conf)

	fmt.Printf("%#v\n", config.Conf.Users)

	config.WriteConfig()

}

func main() {
	setup()
}
