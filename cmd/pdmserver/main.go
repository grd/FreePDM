// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"

	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/logs"
	"github.com/grd/FreePDM/internal/server"
	"github.com/grd/FreePDM/internal/shared"
)

func main() {
	logs.StartLogging()

	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	userRepo := db.NewUserRepo(dbConn)
	srv := server.NewServer(userRepo)
	srv.Routes(mux)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

	// Periodically search for new files
	shared.ImportSharedFiles()
}
