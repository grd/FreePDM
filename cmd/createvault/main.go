// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/grd/FreePDM/internal/server/config"
	"github.com/grd/FreePDM/internal/util"
)

// Script for creating a new PDM vault.

const (
	fileListCsv     = "FileList.csv"
	containerNumber = "ContainerNumber.txt"
	lockedFileCsv   = "LockedFiles.csv"
)

var (
	vaults     string
	vaultsData string
	newVault   string
	userName   string
	userUid    int
	vaultUid   int
)

func setup() {
	fmt.Println("")
	fmt.Println("Welcome to the new vault creation program.")
	fmt.Println("")
	fmt.Println("Some information will be necessary prior to creating a new vault.")
	fmt.Println("For more info please refer to:")
	fmt.Println("https://github.com/grd/FreePDM/blob/main/ConceptOfDesign/FreePDM_Install/Install.md")
	fmt.Println("")
	fmt.Println("The necessary information:")
	fmt.Println("- the new vault directory (for store your information),")
	fmt.Println("")
	fmt.Println("Continue? [Y/n]")
	var mounting_point string
	fmt.Scan(&mounting_point)
	if strings.ToLower(mounting_point) == "n" {
		fmt.Println("Program terminated.")
		os.Exit(-1)
	}

	// Setting vaults and vaultsData
	vaults = config.Conf.VaultsDirectory
	vaultsData = path.Join(vaults, "/.data")

	if !util.DirExists(vaults) {
		log.Fatal("Houston, we have a problem.")
	}

	userName = os.Getenv("USER")
	userUid = config.Conf.Users[userName]
	vaultUid = config.Conf.Users["vault"]

	if err := os.Chdir(vaults); err != nil {
		log.Fatal(err)
	}

	dir, err := os.ReadDir(".")
	util.CheckErr(err)
	if len(dir) > 1 {
		dir = dir[1:] // skip ".data"
		fmt.Println("")
		fmt.Println("This is a list of exisisting vaults:")
		for _, file := range dir {
			fmt.Println(file.Name())
		}
	}

	fmt.Println("\nInput the directory of your new vault...")
	fmt.Scan(&newVault)

	//
	// Creating directory structure
	//

	if util.DirExists(newVault) {
		log.Fatalf("The vault \"%v\" already exist.", newVault)
	}
	err = os.Chmod(vaults, 0775)
	util.CheckErr(err)

	err = os.Mkdir(newVault, 0775)
	util.CheckErr(err)

	err = os.Chown(newVault, userUid, vaultUid)
	util.CheckErr(err)

	if !util.DirExists(vaultsData) {
		if err := os.Mkdir(vaultsData, 0775); err != nil {
			log.Fatalf("error creating dir %s", err)
		}
	}

	if err := os.Chmod(vaultsData, 0775); err != nil {
		log.Fatalf("error chmod dir %s", err)
	}

	// Creating data inside vaultsdata dir
	err = os.Chdir(vaultsData)
	util.CheckErr(err)

	err = os.Mkdir(newVault, 0775)
	util.CheckErr(err)

	err = os.Chdir(newVault)
	util.CheckErr(err)

	err = os.WriteFile(containerNumber, []byte{'0'}, 0644)
	util.CheckErr(err)
	os.Chown(containerNumber, userUid, vaultUid)
	util.CheckErr(err)

	WriteIndexNumber()

	WriteLockedFiles()

	err = os.Chdir("..")
	util.CheckErr(err)

	err = os.Chmod(newVault, 0555)
	util.CheckErr(err)

	err = os.Chdir("..")
	util.CheckErr(err)

	err = os.Chmod(vaultsData, 0555)
	util.CheckErr(err)

	err = os.Chdir("..")
	util.CheckErr(err)

	err = os.Chmod(vaults, 0555)
	util.CheckErr(err)

	fmt.Printf("Three files have been created: %s, %s and %s\n", fileListCsv, containerNumber, lockedFileCsv)
	fmt.Println("")
	fmt.Println("The vault has been created.")

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
		"ContainerNumber", "FileName", "PreviousFile", "Directory", "PreviousDir",
	}

	if err := w.Write(firstRecord); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	os.Chown(fileListCsv, userUid, vaultUid)
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

	os.Chown(lockedFileCsv, userUid, vaultUid)
}

func main() {
	setup()
}
