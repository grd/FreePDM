// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"

	"github.com/grd/FreePDM/internal/logs"
	"github.com/grd/FreePDM/internal/server"
	"github.com/grd/FreePDM/internal/shared"
)

func main() {
	// start logging
	logs.StartLogging()

	http.HandleFunc("/command", server.CommandHandler)

	log.Println("Server running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// Periodically search for new files
	shared.ImportSharedFiles()
}
