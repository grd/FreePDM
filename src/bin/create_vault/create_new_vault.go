// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	ex "github.com/grd/FreePDM/src/extras"
)

// Script for creating a new PDM vault.

func setup() {
	fmt.Println("")
	fmt.Println("Welcome to the create new vault program.")
	fmt.Println("")
	fmt.Println("Now a small questionaire is coming for the creation of the new vault.")
	fmt.Println("For more info look at https://github.com/grd/FreePDM/tests/README.md")
	fmt.Println("")
	fmt.Println("The necessary information is the user uid, the vault uid,")
	fmt.Println("and the name of the new vault.")
	fmt.Println("")
	fmt.Println("Are you ready? [Y/n]")
	var mounting_point string
	fmt.Scan(&mounting_point)
	if strings.ToLower(mounting_point) == "n" {
		fmt.Println("Program terminated.")
		os.Exit(-1)
	}

	fmt.Println("\nEnter the user name")
	var user_name string
	fmt.Scan(&user_name)

	fmt.Println("\nEnter the user uid")
	var user_uid int
	fmt.Scan(&user_uid)

	fmt.Println("Enter the vault uid")
	var vault_uid int
	fmt.Scan(&vault_uid)

	os.Chdir("/vault")

	dir, err := os.ReadDir(".")
	ex.CheckErr(err)
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

	if ex.DirExists(vault_dir) == true {
		log.Fatalf("The vault \"%v\" already exist.", vault_dir)
	}
	err = os.Mkdir(vault_dir, 0777)
	ex.CheckErr(err)
	err = os.Chown(vault_dir, user_uid, vault_uid)
	ex.CheckErr(err)

	err = os.Chdir(vault_dir)
	ex.CheckErr(err)

	// creating three files: "FileList.txt",
	// "IndexNumber.txt" and "Locked.txt"

	fileListTxt := "FileList.txt"
	index_number := "IndexNumber.txt"
	locked_file := "Locked.txt"

	f1, err := os.Create(fileListTxt)
	f1.Close()
	ex.CheckErr(err)
	os.Chown(fileListTxt, user_uid, vault_uid)

	_, err = os.Create(locked_file)
	ex.CheckErr(err)
	os.Chown(locked_file, user_uid, vault_uid)

	d1 := []byte("0")
	err = os.WriteFile(index_number, d1, 0644)
	ex.CheckErr(err)
	os.Chown(index_number, user_uid, vault_uid)

	err = os.Mkdir("PDM", 0777)
	ex.CheckErr(err)
	err = os.Chown("PDM", user_uid, vault_uid)
	ex.CheckErr(err)

	fmt.Println("Three files have been created: FileList.txt,")
	fmt.Println("IndexNumber.txt and Locked.txt, and the directory PDM")
	fmt.Println("")
	fileList, err := os.ReadDir(".")
	ex.CheckErr(err)
	for _, list := range fileList {
		fmt.Println(list.Name())
	}

	fmt.Println("")
	fmt.Println("If that is correct, then the vault has been created.")

	fmt.Println("Program ended successfully.")
}

func main() {
	setup()
}
