// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"log"
	"os"

	ex "github.com/grd/FreePDM/src/extras"
)

// A simple "quick and dirty" remove all files from testvault script
// but it also initializes the startup files again.

func WriteIndexNumber() {

	f, err := os.Create("FileList.csv")
	ex.CheckErr(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	firstRecord := []string{
		"Index", "FileName", "PreviousFile", "Directory", "PreviousDir",
	}

	if err := w.Write(firstRecord); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	os.Chown("FileList.csv", os.Geteuid(), 1002)
}

func WriteLockedFiles() {

	f, err := os.Create("LockedFiles.csv")
	ex.CheckErr(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	firstRecord := []string{"FileNumber", "Version", "UserName"}

	if err := w.Write(firstRecord); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	os.Chown("LockedFiles.csv", os.Geteuid(), 1002)
}

func main() {
	// Check wether we are inside the right folder...
	if ex.FileExists("IndexNumber.txt") == false {
		log.Fatal("Wrong directory...")
	}

	err := os.WriteFile("IndexNumber.txt", []byte{'0'}, 0644)
	ex.CheckErr(err)

	WriteIndexNumber()

	WriteLockedFiles()

	err = os.RemoveAll("PDM")
	ex.CheckErr(err)

	err = os.Mkdir("PDM", 0777)
	ex.CheckErr(err)

	os.Chown("PDM", os.Geteuid(), 1002)
}
