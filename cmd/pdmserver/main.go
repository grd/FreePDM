// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/logs"
	"github.com/grd/FreePDM/internal/shared"
)

func main() {
	// start logging
	logs.StartLogging()

	// database
	db, err := db.InitializeDB("")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := db.UserRepo{DB: db}

	// Running the server
	RunServer()

	// Periodically search for new files
	shared.ImportSharedFiles()
}
