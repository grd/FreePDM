// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path"

	ex "github.com/grd/FreePDM/src/extras"
	fs "github.com/grd/FreePDM/src/filesystem"
)

func main() {
	err := os.Chdir("..")
	ex.CheckErr(err)
	fd_dir, err := os.Getwd()
	ex.CheckErr(err)
	test_dir := path.Join(fd_dir, "ConceptOfDesign/TestFiles")
	// file1 := path.Join(test_dir, "0001.FCStd")
	// file2 := path.Join(test_dir, "0002.FCStd")
	// file3 := path.Join(test_dir, "0003.FCStd")
	file4 := path.Join(test_dir, "0004.FCStd")
	file5 := path.Join(test_dir, "0005.FCStd")
	file6 := path.Join(test_dir, "0006.FCStd")

	fmt.Println("This program populates a vault with data.")
	tmpdir := os.TempDir() // storage of temporary data
	err = os.Chdir(tmpdir)
	ex.CheckErr(err)

	// fsm := fs.InitFileSystem("/tmp/testpdm", "user1")
	fsm := fs.InitFileSystem("/samba/vaults/testpdm", "user")

	dir, err := os.ReadDir(".")
	ex.CheckErr(err)
	fmt.Printf("Root directory of the vault: %s", dir)
	fsm.Mkdir("Standard Parts")
	fsm.Mkdir("Projects")
	fmt.Println(fsm.ListDir("."))
	fsm.Chdir("Projects")
	fmt.Println(fsm.ListDir("."))
	// f1 := fsm.ImportNewFile(file1, "ConceptOfDesign/TestFiles", "")
	// f2 := fsm.ImportNewFile(file2, "ConceptOfDesign/TestFiles", "")
	// f3 := fsm.ImportNewFile(file3, "ConceptOfDesign/TestFiles", "")

	// _ = fsm.CheckOut(f1, 0)
	// // _ = fsm.Checkout(f2, 0)
	// // _ = fsm.Checkout(f3, 0)

	// fmt.Println("checkout status of nr:%v = %v", f1, fsm.CheckoutStatus(f1, 0))
	// fsm.Checkin(f1)

	fsm.Mkdir("Temp")
	fsm.Chdir("Temp")

	fsm.ImportFile(file4)
	fsm.ImportFile(file5)
	fsm.ImportFile(file6)

	fsm.Chdir("..")
	fmt.Println(fsm.ListDir("."))

	// rev = fs.check_latest_file_version("0003.FCStd")
	// new_file = "0003.FCStd" + "." + str(rev).zfill(3)
	// fs.revision_file(new_file)
	// fmt.Println("checking file number: " + str(fs.check_latest_file_version("0003.FCStd")))

	// fs.rename("0003.FCStd.003", "0001.FCStd")
}
