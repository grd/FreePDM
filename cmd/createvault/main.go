// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/grd/FreePDM/pkg/util"
)

// Script for creating a new PDM vault.

const (
	fileListCsv   = "FileList.csv"
	indexNumber   = "IndexNumber.txt"
	lockedFileCsv = "LockedFiles.csv"
)

func setup() {
	fmt.Println("")
	fmt.Println("Welcome to the create new vault program.")
	fmt.Println("")
	fmt.Println("Now a small questionaire is coming for the creation of the new vault.")
	fmt.Println("For more info look at https://github.com/grd/FreePDM/tests/README.md")
	fmt.Println("")
	fmt.Println("The necessary information is the login user, user uid, the vault uid,")
	fmt.Println("and the name of the new vault.")
	fmt.Println("")
	fmt.Println("Are you ready? [Y/n]")
	var mounting_point string
	fmt.Scan(&mounting_point)
	if strings.ToLower(mounting_point) == "n" {
		fmt.Println("Program terminated.")
		os.Exit(-1)
	}

	user_uid := os.Geteuid()

	fmt.Println("Enter the vault uid")
	var vault_uid int
	fmt.Scan(&vault_uid)

	os.Chdir("/vault")

	dir, err := os.ReadDir(".")
	util.CheckErr(err)
	if len(dir) > 0 {
		fmt.Println("")
		fmt.Println("This is a list of exisisting vaults:")
		for _, file := range dir {
			fmt.Println(file)
		}
	}

	fmt.Println("\nInput the directory of your new vault...")
	var vault_dir string
	fmt.Scan(&vault_dir)

	if util.DirExists(vault_dir) == true {
		log.Fatalf("The vault \"%v\" already exist.", vault_dir)
	}
	err = os.Mkdir(vault_dir, 0775)
	util.CheckErr(err)
	err = os.Chown(vault_dir, user_uid, vault_uid)
	util.CheckErr(err)

	err = os.Chdir(vault_dir)
	util.CheckErr(err)

	err = os.WriteFile(indexNumber, []byte{'0'}, 0644)
	util.CheckErr(err)
	os.Chown(indexNumber, user_uid, vault_uid)
	util.CheckErr(err)

	WriteIndexNumber()

	WriteLockedFiles()

	err = os.Mkdir("PDM", 0775)
	util.CheckErr(err)
	err = os.Chown("PDM", user_uid, vault_uid)
	util.CheckErr(err)

	fmt.Printf(`Three files have been created: %s,\n
		%s and %s, and the directory PDM.\n\n`, fileListCsv, indexNumber, lockedFileCsv)
	fileList, err := os.ReadDir(".")
	util.CheckErr(err)
	for _, list := range fileList {
		fmt.Println(list.Name())
	}

	fmt.Println("")
	fmt.Println("If that is correct, then the vault has been created.")

	fmt.Println("Program ended successfully.")
}

// A simple "quick and dirty" remove all files from testvault script
// but it also initializes the startup files again.

func WriteIndexNumber() {

	f, err := os.Create(fileListCsv)
	util.CheckErr(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	firstRecord := []string{
		"Index", "FileName", "PreviousFile", "Directory", "PreviousDir",
	}

	if err := w.Write(firstRecord); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	os.Chown(fileListCsv, os.Geteuid(), 1002)
}

func WriteLockedFiles() {

	f, err := os.Create(lockedFileCsv)
	util.CheckErr(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	firstRecord := []string{"FileNumber", "Version", "UserName"}

	if err := w.Write(firstRecord); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	os.Chown(lockedFileCsv, os.Geteuid(), 1002)
}

func main() {
	setup()
}
