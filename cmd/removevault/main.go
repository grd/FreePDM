// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/grd/FreePDM/internal/server/config"
	"github.com/grd/FreePDM/internal/util"
)

// A simple "quick and dirty" remove all files from testvault script
// but it also initializes the startup files again.

func WriteIndexNumber() {

	f, err := os.Create("FileList.csv")
	util.CheckErr(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Comma = ':'

	firstRecord := []string{
		"Index", "FileName", "PreviousFile", "Directory", "PreviousDir",
	}

	if err := w.Write(firstRecord); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	vaultUid := config.GetUid("vault")

	os.Chown("FileList.csv", os.Geteuid(), vaultUid)
}

func WriteLockedFiles() {

	f, err := os.Create("LockedFiles.csv")
	util.CheckErr(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Comma = ':'

	firstRecord := []string{"FileNumber", "Version", "UserName"}

	if err := w.Write(firstRecord); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	vaultUid := config.GetUid("vault")

	os.Chown("LockedFiles.csv", os.Geteuid(), vaultUid)
}

func main() {
	// Check whether we are inside the right folder...
	if !util.FileExists("IndexNumber.txt") {
		log.Fatal("Wrong directory...")
	}

	err := os.WriteFile("IndexNumber.txt", []byte{'0'}, 0644)
	util.CheckErr(err)

	WriteIndexNumber()

	WriteLockedFiles()

	err = os.RemoveAll("PDM")
	util.CheckErr(err)

	err = os.Mkdir("PDM", 0775)
	util.CheckErr(err)

	vaultUid := config.GetUid("vault")

	os.Chown("PDM", os.Geteuid(), vaultUid)
}
