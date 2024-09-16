// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"testing"

	fsm "github.com/grd/FreePDM/src/filesystem"
	ex "github.com/grd/FreePDM/src/utils"

	"github.com/grd/FreePDM/src/config"
)

const testpdm = "testpdm"

var (
	fs fsm.FileSystem

	testvaults     = path.Join(fsm.Vaults, testpdm)
	testvaultsdata = path.Join(fsm.VaultsData, testpdm)

	file1, file2, file3, file4, file5, file6 string

	freePdmDir = path.Join(os.Getenv("HOME"), "FreePDM")
)

func TestInitFileSystem(t *testing.T) {
	userName, err := user.Current()
	ex.CheckErr(err)
	fs = fsm.InitFileSystem(testpdm, userName.Name)
}

func TestMkDir(t *testing.T) {
	err := fs.Mkdir("Standard Parts")
	if err != nil {
		t.Errorf("Error message = %s", err)
	}

	fs.Mkdir("Projects")
	fs.Mkdir("test")

}

func TestImportFile(t *testing.T) {
}

func TestTheRest(t *testing.T) {

	userName, err := user.Current()
	ex.CheckErr(err)

	fmt.Printf("Vault dir: %s\n", testvaults)
	fmt.Printf("Vault data dir: %s\n", testvaultsdata)
	fmt.Printf("User name: %s\n", userName.Name)

	filesDir := path.Join(freePdmDir, "ConceptOfDesign/TestFiles")
	file1 = path.Join(filesDir, "0001.FCStd")
	file2 = path.Join(filesDir, "0002.FCStd")
	file3 = path.Join(filesDir, "0003.FCStd")
	file4 = path.Join(filesDir, "0004.FCStd")
	file5 = path.Join(filesDir, "0005.FCStd")
	file6 = path.Join(filesDir, "0006.FCStd")

	// fileInfo := fs.ListWD()
	// for _, info := range fileInfo {
	// 	fmt.Println(info.FileName)
	// }

	err = fs.Chdir("Projects")
	ex.CheckErr(err)

	f1 := fs.ImportFile(file1)
	fs.CheckIn(f1, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf1-0", "Testf1-0")

	ver := fs.NewVersion(f1)
	fs.CheckIn(f1, ver, "Testf1-1", "Test1-1")

	ver = fs.NewVersion(f1)
	fs.CheckIn(f1, ver, "Testf1-2", "Test1-2")

	ver = fs.NewVersion(f1)
	fs.CheckIn(f1, ver, "Testf1-3", "Test1-3")

	f2 := fs.ImportFile(file2)
	fs.CheckIn(f2, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf2-0", "")

	f3 := fs.ImportFile(file3)
	fs.CheckIn(f3, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf3-0", "")

	ver = fs.NewVersion(f3)
	fs.CheckIn(f3, ver, "Testf3-1", "Test3-1")

	//
	// rename, copy and move requires a bit of love...
	//
	// these functions require that all the files are checked in and will fail when they are not.
	//

	// err = fs.FileRename("0001.FCStd", "0007.FCStd")
	// ex.CheckErr(err)

	// err = fs.FileCopy("0002.FCStd", "0008.FCStd")
	// ex.CheckErr(err)

	// err = fs.FileCopy("0002.FCStd", "../test/0008a.FCStd")
	// ex.CheckErr(err)

	// fileInfo = fs.ListWD()
	// for _, info := range fileInfo {
	// 	fmt.Println(info.FileName)
	// }

	err = fs.Mkdir("temp")
	ex.CheckErr(err)

	err = fs.FileMove("0002.FCStd", "temp")
	ex.CheckErr(err)

	err = fs.FileMove("0003.FCStd", "..")
	ex.CheckErr(err)

	// err = fs.Chdir("..")
	// ex.CheckErr(err)

	// fileInfo = fs.ListWD()
	// for _, info := range fileInfo {
	// 	fmt.Println(info.FileName)
	// }

	err = fs.Chdir("../Standard Parts")
	ex.CheckErr(err)

	f4 := fs.ImportFile(file4)
	fs.CheckIn(f4, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf4-0", "Testf4-0")

	f5 := fs.ImportFile(file5)
	fs.CheckIn(f5, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf5-0", "Testf5-0")

	f6 := fs.ImportFile(file6)
	fs.CheckIn(f6, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf6-0", "Testf6-0")

	// fileInfo = fs.ListWD()
	// for _, info := range fileInfo {
	// 	fmt.Println(info.FileName)
	// }

	err = fs.Mkdir("Projects2")
	ex.CheckErr(err)

	// err = fs.DirectoryCopy("Projects", "Projects2")
	// ex.CheckErr(err)

	// err = fs.Chdir("Projects (copy)")
	// ex.CheckErr(err)

	// fileInfo = fs.ListWD()
	// for _, info := range fileInfo {
	// 	fmt.Println(info.FileName)
	// }
}

func TestMain(m *testing.M) {

	setup()

	m.Run()
}

// A simple "quick and dirty" remove all files from testvault
// but it also initializes the startup files again.

func setup() {

	// cleanup the testvaults

	err := os.RemoveAll(testvaults)
	ex.CheckErr(err)

	err = os.Mkdir(testvaults, 0775)
	ex.CheckErr(err)

	vaultUid := config.GetUid("vault")

	err = os.Chown(testvaults, os.Geteuid(), vaultUid)
	ex.CheckErr(err)

	// cleanup the testvaultsdata

	err = os.RemoveAll(testvaultsdata)
	ex.CheckErr(err)

	err = os.Mkdir(testvaultsdata, 0775)
	ex.CheckErr(err)

	err = os.Chown(testvaultsdata, os.Geteuid(), vaultUid)
	ex.CheckErr(err)

	// writing three files...

	writeIndexFileList()

	writeLockedFiles()

	writeIndexNumber()
}

func writeIndexFileList() {

	fileList := path.Join(testvaultsdata, "FileList.csv")

	f, err := os.Create(fileList)
	ex.CheckErr(err)
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

	err = os.Chown(fileList, os.Geteuid(), vaultUid)
	ex.CheckErr(err)
}

func writeLockedFiles() {

	lockedfile := path.Join(testvaultsdata, "LockedFiles.csv")

	f, err := os.Create(lockedfile)
	ex.CheckErr(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Comma = ':'

	firstRecord := []string{"FileNumber", "Version", "UserName"}

	if err := w.Write(firstRecord); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	vaultUid := config.GetUid("vault")

	err = os.Chown(lockedfile, os.Geteuid(), vaultUid)
	ex.CheckErr(err)
}

func writeIndexNumber() {

	vaultUid := config.GetUid("vault")

	idxfile := path.Join(testvaultsdata, "IndexNumber.txt")

	err := os.WriteFile(idxfile, []byte{'0'}, 0644)
	ex.CheckErr(err)

	err = os.Chown(idxfile, os.Geteuid(), vaultUid)
	ex.CheckErr(err)
}
