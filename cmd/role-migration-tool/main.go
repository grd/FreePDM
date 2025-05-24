// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/grd/FreePDM/internal/db"
)

func main() {
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	repo := db.NewUserRepo(dbConn)

	fmt.Println("Starting role normalization...")

	changedCount, totalCount, err := repo.NormalizeAllUsersVerbose()
	if err != nil {
		log.Fatalf("Failed to normalize all users: %v", err)
	}

	fmt.Printf("✅ Role normalization completed: %d users processed, %d users modified.\n", totalCount, changedCount)
}
