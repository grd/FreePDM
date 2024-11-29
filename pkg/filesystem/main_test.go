// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filesystem_test

// Tests for the filesystem.

import (
	"encoding/csv"
	"errors"
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
	"github.com/stretchr/testify/assert"
)

const testpdm = "testpdm"

var (
	fs *fsm.FileSystem

	testvaults     = path.Join(fsm.Root(), testpdm)
	testvaultsdata = path.Join(fsm.RootData(), testpdm)

	file1, file2, file3, file4, file5, file6 string

	freePdmDir = path.Join(os.Getenv("HOME"), "FreePDM")
)

func TestInitFileSystem(t *testing.T) {
	userName, err := user.Current()
	if err != nil {
		t.Fatalf("user.Current() error: %s", err)
	}
	fs, err = fsm.NewFileSystem(testpdm, userName.Username)
	if err != nil {
		log.Fatalf("initialization failed, %v", err)
	}
}

func TestMkDir(t *testing.T) {
	if err := fs.Mkdir("Standard Parts"); err != nil {
		t.Fatalf("Mkdir error message = %s", err)
	}

	fs.Mkdir("Projects")
	fs.Mkdir("test")

}

func TestImportFile(t *testing.T) {
	userName, err := user.Current()
	if err != nil {
		t.Fatalf("user.Current() error: %s", err)
	}

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

	if err = fs.Chdir("Projects"); err != nil {
		t.Fatalf("Chdir error message = %s", err)
	}

	f1, err := fs.ImportFile(file1)
	if err != nil {
		t.Fatalf("ImportFile %s error: %s", file1, err)
	}
	compareFileListLine(1, "1:0001.FCStd::Projects:")
	checkOutStatus(1, 0)

	fs.CheckIn(*f1, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf1-0", "Testf1-0")
	checkInStatus(1, 0)

	ver, _ := fs.NewVersion(*f1)
	checkOutStatus(1, 1)

	fs.CheckIn(*f1, ver, "Testf1-1", "Test1-1")
	checkInStatus(1, 1)

	ver, _ = fs.NewVersion(*f1)
	checkOutStatus(1, 2)

	fs.CheckIn(*f1, ver, "Testf1-2", "Test1-2")
	checkInStatus(1, 2)

	ver, _ = fs.NewVersion(*f1)
	checkOutStatus(1, 3)

	fs.CheckIn(*f1, ver, "Testf1-3", "Test1-3")
	checkInStatus(1, 3)

	f2, err := fs.ImportFile(file2)
	if err != nil {
		t.Fatalf("ImportFile %s error: %s", file1, err)
	}
	compareFileListLine(2, "2:0002.FCStd::Projects:")
	checkOutStatus(2, 0)

	fs.CheckIn(*f2, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf2-0", "")
	checkInStatus(2, 0)

	f3, err := fs.ImportFile(file3)
	if err != nil {
		t.Fatalf("ImportFile %s error: %s", file1, err)
	}
	compareFileListLine(3, "3:0003.FCStd::Projects:")
	checkOutStatus(3, 0)

	fs.CheckIn(*f3, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf3-0", "")
	checkInStatus(3, 0)

	ver, _ = fs.NewVersion(*f3)
	checkOutStatus(3, 1)

	fs.CheckIn(*f3, ver, "Testf3-1", "Test3-1")
	checkInStatus(3, 1)

	if err = fs.Mkdir("temp"); err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}

	f4, err := fs.ImportFile(file4)
	if err != nil {
		t.Fatalf("ImportFile %s error: %s", file1, err)
	}
	fs.CheckIn(*f4, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf4-0", "Testf4-0")
	compareFileListLine(4, "4:0004.FCStd::Projects:")

	f5, err := fs.ImportFile(file5)
	if err != nil {
		t.Fatalf("ImportFile %s error: %s", file1, err)
	}
	fs.CheckIn(*f5, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf5-0", "Testf5-0")
	compareFileListLine(5, "5:0005.FCStd::Projects:")

	f6, err := fs.ImportFile(file6)
	if err != nil {
		t.Fatalf("ImportFile %s error: %s", file1, err)
	}
	fs.CheckIn(*f6, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf6-0", "Testf6-0")
	compareFileListLine(6, "6:0006.FCStd::Projects:")

	fs.Chdir("..")
}

func TestFileRename(t *testing.T) {
	if err := fs.Chdir("Projects"); err != nil {
		t.Fatal("chdir failed")
	}

	// Ordinary file rename
	if err := fs.FileRename("0001.FCStd", "0007.FCStd"); err != nil {
		t.Fatalf("FileRename %s error: %s", file1, err)
	}
	compareFileListLine(1, "1:0007.FCStd:0001.FCStd:Projects:")

	// Source file is equal to dest file
	if err := fs.FileRename("0007.FCStd", "0007.FCStd"); err != nil {
		assert.Equal(t, err, errors.New("file 0007.FCStd already exists and is stored in 1"))
	}

	// Dest is empty
	if err := fs.FileRename("0007.FCStd", ""); err != nil {
		assert.Equal(t, err, errors.New("empty destination"))
	}

	// Ordinary file move
	if err := fs.FileRename("0003.FCStd", "temp/"); err != nil {
		t.Fatalf("FileMove failed: %v", err)
	}
	compareFileListLine(3, "3:0003.FCStd::Projects/temp:Projects")

	// More complex file move
	if err := fs.FileRename("0007.FCStd", "../"); err != nil {
		t.Fatalf("FileMove failed: %v", err)
	}
	compareFileListLine(1, "1:0007.FCStd:0001.FCStd::Projects")

	// file move with a sub directory
	if err := fs.FileRename("0006.FCStd", "temp/0006a.FCStd"); err != nil {
		t.Fatalf("FileMove failed: %v", err)
	}
	compareFileListLine(6, "6:0006a.FCStd:0006.FCStd:Projects/temp:Projects")

	if err := fs.Chdir(".."); err != nil {
		t.Fatal("chdir failed")
	}

	// Test case for locking
	item, err := fs.GetItem(".", "0007.FCStd")
	if err != nil {
		t.Fatalf("unable to find: %v", err)
	}
	f7 := fsm.NewFileDirectory(fs, item)
	ver, err := fs.NewVersion(f7.FileList())
	if err != nil {
		t.Fatalf("Expected lock error, got: %v", err)
	}
	checkOutStatus(1, 4)
	err = fs.FileRename("0007.FCStd", "0001.FCStd")
	if err == nil || err.Error() != "file 0007.FCStd is checked out by user" {
		t.Fatalf("Expected lock error, got: %v", err)
	}
	fs.CheckIn(f7.FileList(), ver, "Test", "Test")
	checkInStatus(1, 4)

	// file move with a sub directory and different name
	if err := fs.FileRename("0007.FCStd", "Projects/0001.FCStd"); err != nil {
		t.Fatalf("FileMove failed: %v", err)
	}
	compareFileListLine(1, "1:0001.FCStd:0007.FCStd:Projects:")
}

func TestFileCopy(t *testing.T) {

	fs.Chdir("Projects")

	// ordanary file copy
	err := fs.FileCopy("0002.FCStd", "0011.FCStd")
	if err != nil {
		t.Fatalf("FileCopy failed: %v", err)
	}
	compareFileListLine(7, "7:0011.FCStd::Projects:")

	// dest is empty
	err = fs.FileCopy("0002.FCStd", "")
	if err != nil {
		assert.Equal(t, err, errors.New("empty destination"))
	}

	// copy to different dir and new file name
	if err = fs.FileCopy("0002.FCStd", "../test/0012.FCStd"); err != nil {
		t.Fatalf("FileCopy %s error: %s", file1, err)
	}
	compareFileListLine(8, "8:0012.FCStd::test:")

	// Test for a file copy with only a directory as destination.
	if err = fs.FileCopy("0002.FCStd", "../test/"); err != nil {
		t.Fatalf("FileCopy %s error: %s", file1, err)
	}
	compareFileListLine(9, "9:0002.FCStd::test:")

	// Test locked file
	item, err := fs.GetItem("Projects", "0002.FCStd")
	if err != nil {
		t.Fatalf("unable to find: %v", err)
	}
	f2 := fsm.NewFileDirectory(fs, item)
	ver, err := fs.NewVersion(f2.FileList())
	if err != nil {
		t.Fatalf("Expected lock error, got: %v", err)
	}
	checkOutStatus(2, 1)
	err = fs.FileCopy("0002.FCStd", "0013.FCStd")
	if err == nil || err.Error() != "file 0002.FCStd is checked out by user" {
		t.Fatalf("Expected lock error, got: %v", err)
	}
	fs.CheckIn(f2.FileList(), ver, "Test", "Test")
	checkInStatus(2, 1)

	// Test for file already existing at destination
	err = fs.FileCopy("0004.FCStd", "0011.FCStd")
	if err == nil || err.Error() != "file 0011.FCStd already exists and is stored in Projects" {
		t.Fatalf("Expected existing file error, got: %v", err)
	}

	// back to root dir
	fs.Chdir("..")
}

func TestDirectoryRename(t *testing.T) {
	// Test if dest is empty
	err := fs.DirectoryRename("Projects", "")
	if err != nil {
		assert.Equal(t, err, errors.New("empty destination"))
	}

	// Test with live data, including sub-dirs
	err = fs.DirectoryRename("Projects", "test/Projects3")
	if err != nil {
		t.Fatalf("DirectoryMove failed: %v", err)
	}

	// Test for destination directory being a number
	err = fs.DirectoryRename("Standard Parts", "123")
	if err == nil || err.Error() != "destination directory 123 cannot be a number" {
		t.Fatalf("Expected number error, got: %v", err)
	}

	// Test for source directory not existing
	err = fs.DirectoryRename("nonExistentDir", "destDir")
	if err == nil || err.Error() != "source directory nonExistentDir does not exist" {
		t.Fatalf("Expected source directory error, got: %v", err)
	}

	// Test for creating new directory
	if err := fs.Mkdir("temp"); err != nil {
		t.Fatalf("Expected to create dir temp got %v", err)
	}

	// Test for empty source directory
	err = fs.DirectoryRename("temp", "destDir")
	if err == nil || err.Error() != "directory is empty" {
		t.Fatalf("Expected source directory to be empty, got: %v", err)
	}

	// Test with removal of the test directory
	err = os.Remove("temp")
	if err != nil {
		t.Fatalf("Removal of directory temp failed: %v", err)
	}

	// Test with live data, including sub-dirs, return to Projects3
	err = fs.DirectoryRename("test/Projects3", "Projects3")
	if err != nil {
		t.Fatalf("DirectoryMove failed: %v", err)
	}

	// Test with live data, including sub-dirs, rename Projects3 to Projects again
	err = fs.DirectoryRename("Projects3", "Projects")
	if err != nil {
		t.Fatalf("DirectoryMove failed: %v", err)
	}

	// Test with checked-out file

	file, err := fs.GetItem("Projects", "0002.FCStd")
	if err != nil {
		t.Fatalf("GetItem %s error: %s", file.Name(), err)
	}
	if err = fs.CheckOut(file, fsm.FileVersion{Number: 0, Pretty: "0"}); err != nil {
		t.Fatalf("Checkout %s error: %s", file.Name(), err)
	}
	checkOutStatus(2, 0)

	err = fs.DirectoryRename("Projects", "Projects4")
	if err == nil || err.Error() != "check out errors: [0002.FCStd is checked out by user]" {
		t.Fatalf("Expected one file checked out, got: %v", err)
	}

	fs.CheckIn(file, fsm.FileVersion{Number: 0, Pretty: "0"}, "", "")
	checkInStatus(2, 0)
}

func TestDirectoryCopy(t *testing.T) {
	// Test case with destination is a number
	err := fs.DirectoryCopy("Projects", "123")
	if err == nil || err.Error() != "directory 123 is a number" {
		t.Fatalf("Expected number error, got: %v", err)
	}

	// Test if destination directory already exists
	err = fs.DirectoryCopy("Projects", "test")
	if err == nil || err.Error() != "directory test exists" {
		t.Fatalf("Expected existing directory error, got: %v", err)
	}

	// Test if source exists
	err = fs.DirectoryCopy("blah", "test")
	if err == nil || err.Error() != "source directory blah does not exist" {
		t.Fatalf("Expected an existing directory error, got: %v", err)
	}

	// Test if dest is empty
	err = fs.DirectoryCopy("Projects", "")
	if err != nil {
		assert.Equal(t, err, errors.New("empty destination"))
	}

	// Test if dest is empty
	err = fs.DirectoryCopy("Projects", "Project2")
	if err != nil {
		assert.Equal(t, err, errors.New("subdirectories are not allowed for directory copy"))
	}

	err = fs.DirectoryRename("Projects/temp", "temp/")
	if err != nil {
		t.Fatalf("DirectoryRename failed: %v", err)
	}

	// Test with live data, including sub-dirs
	err = fs.DirectoryCopy("temp", "Project2")
	if err != nil {
		t.Fatalf("DirectoryCopy failed: %v", err)
	}
}

// func TestListTree(t *testing.T) {
// 	lt, err := fs.ListTree("Projects")
// 	if err != nil {
// 		t.Fatalf("listTree failed: %v", err)
// 	}
// 	list := make([]string, len(lt))
// 	for i, item := range lt {
// 		list[i] = path.Join(item.Path(), item.Name())
// 	}
// 	test := []string{
// 		"Projects/temp",
// 		"Projects/0001.FCStd",
// 		"Projects/0002.FCStd",
// 		"Projects/0003.FCStd",
// 	}
// 	for i := range list {
// 		assert.Equal(t, list[i], test[i])
// 	}
// }

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
		"ContainerNumber", "FileName", "PreviousFile", "Directory", "PreviousDir",
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

	firstRecord := []string{"ContainerNumber", "Version", "UserName"}

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
