// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	ex "github.com/grd/FreePDM/src/extras"
	fs "github.com/grd/FreePDM/src/filesystem"
)

func main() {
	vaultDirectory := flag.String("v", "", "The existing vault directory, stored in /vault")
	userName := flag.String("u", "", "The user- or login name")
	flag.Parse()
	fmt.Printf("Vault dir: %s\n", *vaultDirectory)
	fmt.Printf("User name: %s\n", *userName)

	err := os.Chdir("/media/nas/FreePDM")
	ex.CheckErr(err)
	fd_dir, err := os.Getwd()
	ex.CheckErr(err)
	test_dir := path.Join(fd_dir, "ConceptOfDesign/TestFiles")
	file1 := path.Join(test_dir, "0001.FCStd")
	file2 := path.Join(test_dir, "0002.FCStd")
	file3 := path.Join(test_dir, "0003.FCStd")
	file4 := path.Join(test_dir, "0004.FCStd")
	file5 := path.Join(test_dir, "0005.FCStd")
	file6 := path.Join(test_dir, "0006.FCStd")

	fmt.Println("This program populates a vault with existing data.")

	vaultDir := path.Join("/vault", *vaultDirectory) // storage of temporary data
	err = os.Chdir(vaultDir)
	ex.CheckErr(err)

	fs := fs.InitFileSystem(vaultDir, *userName)
	fmt.Printf("Root directory of the vault: %s\n\n", vaultDir)
	fs.Mkdir("Standard Parts")
	fs.Mkdir("Projects")
	fs.Chdir("Projects")

	f1 := fs.ImportFile(file1)
	fs.CheckIn(f1, 0, "Testf1-0", "Testf1-0")

	ver := fs.NewVersion(f1)
	fs.CheckIn(f1, ver, "Testf1-1", "Test1-1")

	ver = fs.NewVersion(f1)
	fs.CheckIn(f1, ver, "Testf1-2", "Test1-2")

	ver = fs.NewVersion(f1)
	fs.CheckIn(f1, ver, "Testf1-3", "Test1-3")

	f2 := fs.ImportFile(file2)
	fs.CheckIn(f2, 0, "Testf2-0", "")

	f3 := fs.ImportFile(file3)
	fs.CheckIn(f3, 0, "Testf3-0", "")

	ver = fs.NewVersion(f3)
	fs.CheckIn(f3, ver, "Testf3-1", "Test3-1")

	//
	// rename, copy and move requires a bit of love...
	//
	// these functions require that all the files are checked in and will fail when they are not.
	//

	err = fs.FileRename("0001.FCStd", "0007.FCStd")
	ex.CheckErr(err)

	err = fs.FileCopy("0002.FCStd", "0008.FCstd")
	ex.CheckErr(err)

	fs.Mkdir("temp")

	// err = fs.FileMove("0003.FCStd", "temp")
	// ex.CheckErr(err)

	fs.Chdir("..")
	fs.Chdir("Standard Parts")

	f4 := fs.ImportFile(file4)
	fs.CheckIn(f4, 0, "Testf4-0", "Testf4-0")

	f5 := fs.ImportFile(file5)
	fs.CheckIn(f5, 0, "Testf5-0", "Testf5-0")

	f6 := fs.ImportFile(file6)
	fs.CheckIn(f6, 0, "Testf6-0", "Testf6-0")

	fmt.Println("This script successfully finished.")
}
