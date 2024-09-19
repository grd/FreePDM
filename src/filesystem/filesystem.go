// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/grd/FreePDM/src/config"
	ex "github.com/grd/FreePDM/src/utils"
	"golang.org/x/exp/slices"
)

// LockedIndex is the list of locked files
type LockedIndex struct {
	fileNr   int64  // The number of the file
	version  int16  // The number of the version
	userName string // Who checked this file out
}

// File System related Class
type FileSystem struct {
	index             FileIndex
	vaultDir          string
	dataDir           string
	vaultUid          int
	user              string
	userUid           int
	currentWorkingDir string
	lockedCvs         string
	lockedIndex       []LockedIndex
}

const (
	LockedFileCsv = "LockedFiles.csv"

	Vaults     = "/samba/vaults"
	VaultsData = "/samba/vaultsdata"
)

// Constructor
func InitFileSystem(vaultDir, userName string) (fs FileSystem) {

	parts := strings.Split(vaultDir, "/")
	if len(parts) == 1 {
		fs.vaultDir = path.Join(Vaults, vaultDir)
		fs.dataDir = path.Join(VaultsData, vaultDir)
	} else {
		log.Fatalf("multi-level vaultDir parameter: %s", vaultDir)
	}

	fs.user = userName

	// check wether the critical directories exist.

	ex.CriticalDirExist(fs.vaultDir)
	ex.CriticalDirExist(fs.dataDir)

	fs.currentWorkingDir = ""
	fs.vaultUid = config.GetUid("vault")
	fs.userUid = config.GetUid(userName)

	if fs.userUid == -1 {
		log.Fatal("Username has not been stored into the FreePDM config file. Please follow the setup process.")
	}

	if fs.vaultUid == 0 || fs.vaultUid == -1 {
		log.Fatal("Vault UID has not been stored into the FreePDM config file. Please follow the setup process.")
	}

	fs.index = InitFileIndex(&fs)

	fs.lockedCvs = path.Join(fs.dataDir, LockedFileCsv)

	fs.ReadLockedIndex() // retrieve the values

	err := os.Chdir(fs.vaultDir)
	ex.CheckErr(err)

	log.Printf("Vault dir: %s", fs.currentWorkingDir)

	return fs
}

// Updates the locked index by reading from the lockedTxt file.
func (fs *FileSystem) ReadLockedIndex() {

	buf, err := os.ReadFile(fs.lockedCvs)
	ex.CheckErr(err)

	r := csv.NewReader(bytes.NewBuffer(buf))
	r.Comma = ':'

	records, err := r.ReadAll()
	ex.CheckErr(err)

	if len(records) <= 1 {
		return
	}

	records = records[1:]

	fs.lockedIndex = nil

	var list = LockedIndex{}

	for _, record := range records {

		list.fileNr = ex.Atoi64(record[0])
		list.version = ex.Atoi16(record[1])
		list.userName = record[2]

		fs.lockedIndex = append(fs.lockedIndex, list)
	}
}

func (fs *FileSystem) WriteLockedIndex() {

	records := [][]string{
		{"FileNumber", "Version", "UserName"},
	}

	for _, list := range fs.lockedIndex {

		records = append(records, []string{
			ex.I64toa(list.fileNr),
			ex.I16toa(list.version),
			list.userName})

	}

	var buf []byte
	buffer := bytes.NewBuffer(buf)

	writer := csv.NewWriter(buffer)
	writer.Comma = ':'

	err := writer.WriteAll(records) // calls Flush internally
	ex.CheckErr(err)

	err = writer.Error()
	ex.CheckErr(err)

	err = os.WriteFile(fs.lockedCvs, buffer.Bytes(), 0644)
	ex.CheckErr(err)

	err = os.Chown(fs.lockedCvs, fs.userUid, fs.vaultUid)
	ex.CheckErr(err)
}

// import a file inside the PDM. When you import a file the meta-data also gets imported,
// which means uploaded to the server.
// When you import a file or files you are placing the new file in the current directory.
// The new file inside the PDM gets a revision number automatically.
// The function returns the number of the imported file.
func (fs *FileSystem) ImportFile(fname string) int64 {
	// check wether a file exist
	if !ex.FileExists(fname) {
		log.Fatalf("File %s could not be found.", fname)
	}

	index := fs.index.AddItem(fname, fs.currentWorkingDir)

	dir := fmt.Sprintf("%d", index)

	// fd := InitFileDirectory(self, path.Join(self.mainPdmDir, dir))
	fd := InitFileDirectory(fs, dir, index)

	newDir := fd.NewDirectory()
	newDir.ImportNewFile(fname)

	log.Printf("Imported %s into %s with version %d\n", fname, fs.index.FileNameOfString(fd.dir), 0)

	// Checking out the new file so no one else can see it.

	err := fs.CheckOut(index, FileVersion{0, "0", ex.Now()})
	ex.CheckErr(err)

	return index
}

// Generates a new version of a file and returns the version number.
func (fs *FileSystem) NewVersion(indexNr int64) FileVersion {

	dirIdx, err := fs.index.DirIndex(indexNr)
	ex.CheckErr(err)
	dir := path.Join(fs.vaultDir, dirIdx)

	fd := InitFileDirectory(fs, dir, indexNr)

	ret := fd.NewVersion()

	// Checking out the new file so no one else can see it.

	log.Printf("Created version %d of file %s\n", ret.Number, fs.index.FileName(indexNr))

	err = fs.CheckOut(indexNr, ret)
	ex.CheckErr(err)

	return ret
}

// Creates a new directory inside the current directory, with the correct uid and gid.
func (fs FileSystem) Mkdir(dir string) error {

	// Check wether dname is an int. We don't want that, because the number could interfere with the fileindex.
	if _, err := strconv.Atoi(dir); err == nil {
		return fmt.Errorf("please change %s into a string, now it is a number", dir)
	}

	err := os.Mkdir(dir, 0777)
	ex.CheckErr(err)

	err = os.Chown(dir, fs.userUid, fs.vaultUid)
	ex.CheckErr(err)

	log.Printf("Created directory: %s\n", dir)

	return nil
}

func (fs *FileSystem) Chdir(dir string) error {

	if !ex.DirExists(dir) {
		return fmt.Errorf("dir %s does not exist", dir)
	}

	newPath := path.Join(fs.vaultDir, fs.currentWorkingDir, dir)

	fs.currentWorkingDir = path.Join(fs.currentWorkingDir, dir)

	err := os.Chdir(newPath)
	if err != nil {
		return err
	}

	log.Printf("Changed directory: %s\n", fs.currentWorkingDir)
	log.Printf("Full path: %s\n", newPath)

	return nil
}

// list the sorted directories and files of the current working directory.
func (fs FileSystem) ListWD() []FileInfo {
	return fs.ListDir(path.Join(fs.vaultDir, fs.currentWorkingDir))
}

// list the sorted directories and files, as long as the directory is inside the vault.
func (fs FileSystem) ListDir(dirName string) []FileInfo {
	dir_list, err := os.ReadDir(dirName)
	ex.CheckErr(err)
	var directoryList []FileInfo
	var fileList []FileInfo
	var subDirList []FileInfo

	for _, sub_dir := range dir_list {
		if num, err := strconv.Atoi(sub_dir.Name()); err == nil {
			fileName := fs.index.FileName(int64(num))
			fileList = append(fileList, FileInfo{
				Dir:      false,
				FileName: fileName,
			})
		} else {
			directoryList = append(directoryList, FileInfo{Dir: true, FileName: sub_dir.Name()})
		}
	}

	if path.Clean(dirName) != fs.vaultDir {
		subDirList = append(subDirList, FileInfo{Dir: true, FileName: ".."})
	}

	if directoryList != nil {
		subDirList = append(subDirList, directoryList...)
	}

	if fileList != nil {
		subDirList = append(subDirList, fileList...)
	}

	return subDirList
}

// returns the latest version number of a file in the current
// directory or -1 when the file doesn't exist.
func (fs FileSystem) CheckLatestFileVersion(fname string) int64 {
	file_list, err := os.ReadDir(".")
	ex.CheckErr(err)
	var result int64 = -1
	for _, file := range file_list {
		if ex.DirExists(file.Name()) {
			continue
		}
		file1, ext1 := SplitExt(file.Name())
		if fname == file1 {
			n, err := strconv.ParseInt(ext1[1:], 10, 64)
			ex.CheckErr(err)
			result = n
		}
	}

	return result
}

// Check whether the itemnr is locked.
// Returns the name of the user who locked it or empty when not locked.
func (fs FileSystem) IsLocked(itemNr int64, version FileVersion) string {

	for _, item := range fs.lockedIndex {
		if item.fileNr == itemNr && item.version == version.Number {
			return item.userName
		}
	}

	return "" // Nothing found
}

// Checkout means locking a itemnr so that only you can use it.
func (fs *FileSystem) CheckOut(itemNr int64, version FileVersion) error {

	fs.ReadLockedIndex() // update the index

	// Set file mode 0700

	dir, err := fs.index.DirIndex(itemNr)
	ex.CheckErr(err)

	fd := InitFileDirectory(fs, path.Join(fs.vaultDir, dir), itemNr)
	fd.OpenItemVersion(version)

	// check whether the itemnr is locked

	if usr := fs.IsLocked(itemNr, version); usr != "" {

		return fmt.Errorf("file %d-%d is locked by user %v", itemNr, version.Number, usr)

	} else {

		fs.lockedIndex = append(fs.lockedIndex, LockedIndex{itemNr, version.Number, fs.user})

		fs.WriteLockedIndex()

		log.Printf("Checked out version %d of file %s\n", version.Number, fs.index.FileName(itemNr))

		return nil
	}
}

// Checkin means unlocking an itemnr.
// The description and long description are meant for storage.
func (fs *FileSystem) CheckIn(itemNr int64, version FileVersion, descr, longdescr string) error {

	// Set file mode 0755

	dir, err := fs.index.DirIndex(itemNr)
	ex.CheckErr(err)

	fd := InitFileDirectory(fs, path.Join(fs.vaultDir, dir), itemNr)

	fd.StoreData(version, descr, longdescr)

	fd.CloseItemVersion(version)

	// check whether the itemnr is locked

	usr := fs.IsLocked(itemNr, version)

	if usr != fs.user {

		return fmt.Errorf("file %d-%d is locked by user %s", itemNr, version.Number, usr)

	} else {

		// Remove item from index

		var nr int
		for i, y := range fs.lockedIndex {
			if y.fileNr == itemNr && y.version == version.Number {
				nr = i
			}
		}

		fs.lockedIndex = slices.Delete(fs.lockedIndex, nr, nr+1)

		fs.WriteLockedIndex()

		log.Printf("Checked in version %d of file %s", version.Number, fs.index.FileName(itemNr))

		return nil
	}
}

// Rename a file, for instance when the user wants to use a file
// with a specified numbering system.
// Note that all versions need to be checked in.
func (fs *FileSystem) FileRename(src, dest string) error {

	// Check wether dest exist

	dir, err := fs.index.Dir(dest)
	if err == nil {
		return fmt.Errorf("file %s already exist and is stored in %s", dest, dir)
	}

	// Rename the file from src to dest

	dir, err = fs.index.CurrentDir(src)
	ex.CheckErr(err)

	fileName, err := fs.index.Index(src)
	ex.CheckErr(err)

	fd := InitFileDirectory(fs, dir, fileName)

	err = fd.fileRename(src, dest)
	if err != nil {
		return err
	}

	// Rename the file in the index

	err = fs.index.renameItem(src, dest)
	if err != nil {
		return err
	}

	// Logging

	log.Printf("File %s renamed to %s\n", src, dest)

	return nil
}

// Copy a the latest version of a file.
// Note that all versions need to be checked in.
func (fs *FileSystem) FileCopy(src, dest string) error {

	var dst string

	if strings.Contains(dest, "/") {
		num := strings.LastIndexByte(dest, '/')
		dst = path.Join(fs.currentWorkingDir, dest[:num])
		fmt.Printf("dst = %s\n", dst)
		destPath := path.Join(fs.vaultDir, dst)
		if !ex.DirExists(destPath) {
			return fmt.Errorf("directory %s doesn't exist", destPath)
		}
	}
	_, err := fs.index.Dir(dest)
	if err == nil {
		return fmt.Errorf("file %s already exist", dest)
	}

	file, err := fs.index.ContainerName(src)
	ex.CheckErr(err)

	srcDir, err := fs.index.CurrentDir(src)
	ex.CheckErr(err)

	srcFd := InitFileDirectory(fs, srcDir, file.index)

	_, destFile := path.Split(dest)

	destIndex := fs.index.AddItem(destFile, dst)

	destDir, err := fs.index.CurrentDir(destFile)
	ex.CheckErr(err)

	destFd := InitFileDirectory(fs, destDir, destIndex)

	destFd.NewDirectory()

	// Copy the file from src to dest

	version := srcFd.LatestVersion()

	srcFile := path.Join(srcFd.dir, version.Pretty, src)

	// Copying file from src to dest and also properties

	destFd.ImportNewFile(srcFile)

	err = destFd.fileRename(src, destFile)
	ex.CheckErr(err)

	// Logging

	log.Printf("File %s copied to %s\n", src, dest)

	return nil
}

// Moves a file to a different directory.
// Note that all versions need to be checked in.
func (fs *FileSystem) FileMove(fileName, destDir string) error {

	// "normalize" the destDir

	dir, err := filepath.Abs(destDir)
	ex.CheckErr(err)

	if !ex.DirExists(dir) {
		return fmt.Errorf("directory %s doesn't exist", destDir)
	}
	if len(dir) == len(fs.vaultDir) { // case when file moved to the root
		dir = ""
	}
	if strings.HasPrefix(dir, fs.vaultDir) {
		dir = dir[len(fs.vaultDir)+1:]
	}

	// Move file

	fname, err := fs.index.CurrentDir(fileName)
	ex.CheckErr(err)

	dest := path.Join(destDir, fname)

	err = os.Rename(fname, dest)
	if err != nil {
		return err
	}

	err = os.Chown(destDir, fs.userUid, fs.vaultUid)
	ex.CheckErr(err)

	// Move file in FileIndex

	fs.index.moveItem(fileName, dir)

	// Logging

	log.Printf("File %s moved to to %s\n", fileName, destDir)

	return nil
}

// TODO implement these three, but they are tricky

// TODO also note of the ex.IsNumber() function for copy and rename
//
// Think about recursive ! ! !
// And logging

// Copy a directory.
// Note that all file versions need to be checked in.
func (fs *FileSystem) DirectoryCopy(src, dest string) error {
	if ex.IsNumber(dest) {
		return fmt.Errorf("directory %s is a number", dest)
	}

	dirList, err := os.ReadDir(src)
	ex.CheckErr(err)

	// if ex.DirExists(dest) == false {
	// 	self.Mkdir(dest)
	// }

	for _, item := range dirList {
		if _, err := strconv.Atoi(item.Name()); err != nil {

			// Directory operations

			destDir := path.Join(dest, item.Name())
			err := fs.Mkdir(destDir)
			ex.CheckErr(err)

			err = fs.DirectoryCopy(item.Name(), destDir)
			if err != nil {
				return err
			}

		} else {

			// File operations

			fileNum, err := fs.index.ContainerName(item.Name())
			ex.CheckErr(err)
			base, ext := ex.SplitFileExtension(fileNum.file)
			newFileName := base + " (copy)" + ext
			destFile := path.Join(dest, newFileName)
			err = fs.FileCopy(fileNum.file, destFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Rename a directory.
// Note that all versions need to be checked in.
func (fs *FileSystem) DirectoryRename(src, dest string) error {

	return nil
}

// Move a directory.
// Note that all versions need to be checked in.
func (fs *FileSystem) DirectoryMove(src, dest string) error {
	// Think about recursive ! ! !
	return nil
}

// Splits a file with its base and extension
func SplitExt(path string) (base, ext string) {
	ext = filepath.Ext(path)
	base = path[:len(path)-len(ext)]
	return
}

// Returns the "vault" Uid
func GetVaultUid() int {

	buf, err := os.ReadFile("/etc/group")
	ex.CheckErr(err)

	r := csv.NewReader(bytes.NewBuffer(buf))
	r.Comma = ':'

	records, err := r.ReadAll()
	ex.CheckErr(err)

	for _, record := range records {

		if record[0] == "vault" {
			log.Printf("Vault uid = %s\n", record[2])
			return int(ex.Atoi16(record[2]))
		}
	}

	log.Fatalf("File %s doesn't have an entry \"vault\". See the installation instructions.", "/etc/group")
	return -1
}

// TODO: Implement this inside all functions that need it.

// // Verifies that both the copy and the registry are okay.
// func (self FileSystem) Verify(indexNum int64, fileName, location string, version int16) error {
// 	var index FileList
// 	for _, file := range self.index.fileList {
// 		if indexNum == file.index {
// 			index = file
// 			break
// 		}
// 	}

// 	if index.index == 0 && index.file == "" {
// 		return fmt.Errorf("Critical error: Index not found! The index number %d is not found!\n", indexNum)
// 	}

// 	if fileName != index.file {
// 		return fmt.Errorf("Critical error: Line %d, stored file name %s doesn't match with the file name %s\n", index.index, index.file, fileName)
// 	}

// 	if location != index.dir {
// 		return fmt.Errorf("Critical error: Line %d, stored directory name %s doesn't match with the location name %s\n", index.index, index.dir, location)
// 	}

// 	dirIdx, err := self.index.DirIndex(index.index)
// 	ex.CheckErr(err)
// 	dir := path.Join(self.mainPdmDir, dirIdx)

// 	fd := InitFileDirectory(&self, dir, index.index)

// 	file := path.Join(self.mainPdmDir, index.dir, ex.I64toa(indexNum))

// 	return nil
// }
