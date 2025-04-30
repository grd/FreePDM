// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"os"
	"path"

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

	// Start HTTPS
	certPath := path.Join("certs", "localhost.pem")
	keyPath := path.Join("certs", "localhost-key.pem")

	if os.Getenv("USE_HTTPS") == "true" {
		log.Println("Server running with HTTPS on https://localhost:8443")
		log.Fatal(http.ListenAndServeTLS(":8443", certPath, keyPath, mux))
	} else {
		log.Println("Server running with HTTP on http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", mux))
	}

	log.Fatal(http.ListenAndServeTLS(":8443", certPath, keyPath, mux))

	// Periodically search for new files
	shared.ImportSharedFiles()
}
