// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Tests for the filesystem.

package filesystem_test

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"testing"

	"github.com/grd/FreePDM/pkg/config"
	fsm "github.com/grd/FreePDM/pkg/filesystem"
	"github.com/grd/FreePDM/pkg/util"
)

const testpdm = "testpdm"

var (
	fism *fsm.FileSystem

	testvaults     = path.Join(fsm.Vaults, testpdm)
	testvaultsdata = path.Join(fsm.VaultsData, testpdm)

	file1, file2, file3, file4, file5, file6 string

	freePdmDir = path.Join(os.Getenv("HOME"), "FreePDM")
)

func TestInitFileSystem(t *testing.T) {
	userName, err := user.Current()
	util.CheckErr(err)
	fism, err = fsm.NewFileSystem(testpdm, userName.Username)
	if err != nil {
		log.Fatalf("initialization failed, %v", err)
	}
}

func TestMkDir(t *testing.T) {
	err := fism.Mkdir("Standard Parts")
	if err != nil {
		t.Errorf("Error message = %s", err)
	}

	fism.Mkdir("Projects")
	fism.Mkdir("test")

}

func TestImportFile(t *testing.T) {
}

func TestTheRest(t *testing.T) {

	userName, err := user.Current()
	util.CheckErr(err)

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

	listWd()

	err = fism.Chdir("Projects")
	util.CheckErr(err)

	{
		f1 := fism.ImportFile(file1)
		{
			compareFileListLine(1, "1:0001.FCStd::Projects:")
			checkOutStatus(1, 0)
		}
		fism.CheckIn(f1, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf1-0", "Testf1-0")
		{
			checkInStatus(1, 0)
		}

		ver, _ := fism.NewVersion(f1)

		{
			checkOutStatus(1, 1)
		}
		fism.CheckIn(f1, ver, "Testf1-1", "Test1-1")
		{
			checkInStatus(1, 1)
		}

		ver, _ = fism.NewVersion(f1)
		{
			checkOutStatus(1, 2)
		}
		fism.CheckIn(f1, ver, "Testf1-2", "Test1-2")
		{
			checkInStatus(1, 2)
		}

		ver, _ = fism.NewVersion(f1)
		{
			checkOutStatus(1, 3)
		}
		fism.CheckIn(f1, ver, "Testf1-3", "Test1-3")
		{
			checkInStatus(1, 3)
		}
	}

	{
		f2 := fism.ImportFile(file2)
		{
			compareFileListLine(2, "2:0002.FCStd::Projects:")
			checkOutStatus(2, 0)
		}
		fism.CheckIn(f2, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf2-0", "")
		{
			checkInStatus(2, 0)
		}
	}

	{
		f3 := fism.ImportFile(file3)
		{
			compareFileListLine(3, "3:0003.FCStd::Projects:")
			checkOutStatus(3, 0)
		}
		fism.CheckIn(f3, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf3-0", "")
		{
			checkInStatus(3, 0)
		}

		ver, _ := fism.NewVersion(f3)
		{
			checkOutStatus(3, 1)
		}
		fism.CheckIn(f3, ver, "Testf3-1", "Test3-1")
		{
			checkInStatus(3, 1)
		}
	}

	//
	// rename, copy and move requires a bit of love...
	//
	// these functions require that all the files are checked in and will fail when they are not.
	//

	{
		err = fism.FileRename("0001.FCStd", "0007.FCStd")
		util.CheckErr(err)
		{
			compareFileListLine(1, "1:0007.FCStd:0001.FCStd:Projects:")
		}
	}

	{
		err = fism.FileCopy("0002.FCStd", "0008.FCStd")
		util.CheckErr(err)
		{
			compareFileListLine(4, "4:0008.FCStd::Projects:")
		}
	}

	{
		err = fism.FileCopy("0002.FCStd", "../test/0008a.FCStd")
		util.CheckErr(err)
		{
			compareFileListLine(5, "5:0008a.FCStd::test:")
		}
	}

	listWd()

	err = fism.Mkdir("temp")
	util.CheckErr(err)

	{
		err = fism.FileMove("0002.FCStd", "temp")
		util.CheckErr(err)
		{
			compareFileListLine(2, "2:0002.FCStd::Projects/temp:Projects")
		}
	}

	{
		err = fism.FileMove("0003.FCStd", "..")
		util.CheckErr(err)
		{
			compareFileListLine(3, "3:0003.FCStd:::Projects")
		}
	}

	listWd()

	err = fism.Chdir("../Standard Parts")
	util.CheckErr(err)

	{
		f4 := fism.ImportFile(file4)
		fism.CheckIn(f4, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf4-0", "Testf4-0")
		{
			compareFileListLine(6, "6:0004.FCStd::Standard Parts:")
		}
	}

	{
		f5 := fism.ImportFile(file5)
		fism.CheckIn(f5, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf5-0", "Testf5-0")
		{
			compareFileListLine(7, "7:0005.FCStd::Standard Parts:")
		}
	}

	{
		f6 := fism.ImportFile(file6)
		fism.CheckIn(f6, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf6-0", "Testf6-0")
		{
			compareFileListLine(8, "8:0006.FCStd::Standard Parts:")
		}
	}

	listWd()

	err = fism.Chdir("../Projects")
	util.CheckErr(err)

	listWd()

	fism.Chdir("..")

	list, err := fism.ListWD()
	if err != nil {
		log.Fatal("listwd failed")
	}

	for _, elem := range list {
		fmt.Println(elem)
	}

	// err = fs.Mkdir("Projects2")
	// ex.CheckErr(err)

	// err = fs.DirectoryCopy("Projects", "Projects2")
	// ex.CheckErr(err)

	// err = fs.Chdir("Projects (copy)")
	// ex.CheckErr(err)

	// fmt.Printf("\nstart listtree\n")
	// list = fs.ListTree(fs.VaultDir())
	// for _, elem := range list {
	// 	fmt.Printf("%v\n", elem)
	// }

	// fmt.Printf("\n\n")
	// log.Fatal("oeps")
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
	util.CheckErr(err)

	err = os.Mkdir(testvaults, 0775)
	util.CheckErr(err)

	vaultUid := config.GetUid("vault")

	err = os.Chown(testvaults, os.Geteuid(), vaultUid)
	util.CheckErr(err)

	// cleanup the testvaultsdata

	err = os.RemoveAll(testvaultsdata)
	util.CheckErr(err)

	err = os.Mkdir(testvaultsdata, 0775)
	util.CheckErr(err)

	err = os.Chown(testvaultsdata, os.Geteuid(), vaultUid)
	util.CheckErr(err)

	// writing three files...

	writeIndexFileList()

	writeLockedFiles()

	writeIndexNumber()
}

func writeIndexFileList() {

	fileList := path.Join(testvaultsdata, "FileList.csv")

	f, err := os.Create(fileList)
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

	err = os.Chown(fileList, os.Geteuid(), vaultUid)
	util.CheckErr(err)
}

func writeLockedFiles() {

	lockedfile := path.Join(testvaultsdata, "LockedFiles.csv")

	f, err := os.Create(lockedfile)
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

	err = os.Chown(lockedfile, os.Geteuid(), vaultUid)
	util.CheckErr(err)
}

func writeIndexNumber() {

	vaultUid := config.GetUid("vault")

	idxfile := path.Join(testvaultsdata, "IndexNumber.txt")

	err := os.WriteFile(idxfile, []byte{'0'}, 0644)
	util.CheckErr(err)

	err = os.Chown(idxfile, os.Geteuid(), vaultUid)
	util.CheckErr(err)
}

func listWd() {
	wd, err := os.Getwd()
	util.CheckErr(err)
	fmt.Printf("\n\nList directory of %s\n\n", wd)
	fileInfo, err := fism.ListWD()
	if err != nil {
		log.Fatal("listwd error")
	}
	for _, info := range fileInfo {
		fmt.Println(info.FileName)
	}
	fmt.Println("")

}

func compareFileListLine(num int, line string) {
	content, err := os.ReadFile(path.Join(testvaultsdata, "FileList.csv"))
	if err != nil {
		log.Fatal(err)
	}
	s := string(content[:len(content)-1])
	slice := strings.Split(s, "\n")

	if slice[num] != line {
		log.Fatalf("error comparing lines of FileList.csv\nRead line (line %d): %s\nLine parameter:     %s", num, slice[num], line)
	}
}

func checkStatus(file int64, num int16) bool {
	content, err := os.ReadFile(path.Join(testvaultsdata, "LockedFiles.csv"))
	if err != nil {
		log.Fatal(err)
	}
	s := string(content[:len(content)-1])
	slice := strings.Split(s, "\n")
	slice = slice[1:] // first line
	if len(slice) == 0 {
		return false
	}
	for _, line := range slice {
		segments := strings.Split(line, ":")
		seg0, _ := util.Atoi64(segments[0])
		seg1, _ := util.Atoi16(segments[1])
		if seg0 == file && seg1 == num {
			return true
		}

	}

	return false
}

func checkOutStatus(file int64, num int16) {
	done := checkStatus(file, num)
	if !done {
		log.Fatalf("error file not properly checked out")
	}
}

func checkInStatus(file int64, num int16) {
	done := checkStatus(file, num)
	if done {
		log.Fatalf("error file not properly checked in")
	}
}

// func checkFileComments(file int64, dir string, num int16, descr, longDescr string){

// }
