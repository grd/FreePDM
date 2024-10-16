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
	"sort"
	"strconv"
	"strings"

	"github.com/grd/FreePDM/pkg/config"
	"github.com/grd/FreePDM/pkg/util"
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
func NewFileSystem(vaultDir, userName string) (fs *FileSystem, err error) {

	fs = new(FileSystem)

	parts := strings.Split(vaultDir, "/")
	if len(parts) == 1 {
		fs.vaultDir = path.Join(Vaults, vaultDir)
		fs.dataDir = path.Join(VaultsData, vaultDir)
	} else {
		log.Fatalf("multi-level vaultDir parameter: %s", vaultDir)
	}

	fs.user = userName

	// check wether the critical directories exist.

	util.CriticalDirExist(fs.vaultDir)
	util.CriticalDirExist(fs.dataDir)

	fs.currentWorkingDir = ""
	fs.vaultUid = config.GetUid("vault")
	fs.userUid = config.GetUid(userName)

	if fs.userUid == -1 {
		log.Fatalf("Username %s has not been stored into the FreePDM config file, please follow the setup process", userName)
	}

	if fs.vaultUid == 0 || fs.vaultUid == -1 {
		log.Fatal("Vault UID has not been stored into the FreePDM config file. Please follow the setup process.")
	}

	index, err := NewFileIndex(fs)
	if err != nil {
		return fs, err
	}
	fs.index = index

	fs.lockedCvs = path.Join(fs.dataDir, LockedFileCsv)

	fs.ReadLockedIndex() // retrieve the values

	err = os.Chdir(fs.vaultDir)
	util.CheckErr(err)

	log.Printf("Vault dir: %s", fs.currentWorkingDir)

	return fs, nil
}

// Returns the Vault directory
func (fs *FileSystem) VaultDir() string {
	return fs.vaultDir
}

// Updates the locked index by reading from the lockedTxt file.
func (fs *FileSystem) ReadLockedIndex() {

	buf, err := os.ReadFile(fs.lockedCvs)
	util.CheckErr(err)

	r := csv.NewReader(bytes.NewBuffer(buf))
	r.Comma = ':'

	records, err := r.ReadAll()
	util.CheckErr(err)

	if len(records) <= 1 {
		return
	}

	records = records[1:]

	fs.lockedIndex = nil

	var list = LockedIndex{}

	for _, record := range records {

		list.fileNr, _ = util.Atoi64(record[0])
		list.version, _ = util.Atoi16(record[1])
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
			util.I64toa(list.fileNr),
			util.I16toa(list.version),
			list.userName})

	}

	var buf []byte
	buffer := bytes.NewBuffer(buf)

	writer := csv.NewWriter(buffer)
	writer.Comma = ':'

	err := writer.WriteAll(records) // calls Flush internally
	util.CheckErr(err)

	err = writer.Error()
	util.CheckErr(err)

	err = os.WriteFile(fs.lockedCvs, buffer.Bytes(), 0644)
	util.CheckErr(err)

	err = os.Chown(fs.lockedCvs, fs.userUid, fs.vaultUid)
	util.CheckErr(err)
}

// import a file inside the PDM. When you import a file the attributes also gets imported,
// which means uploaded to the server.
// When you import a file or files you are placing the new file in the current directory.
// The new file inside the PDM gets a revision number automatically.
// The function returns the number of the imported file or an error.
func (fs *FileSystem) ImportFile(fname string) (int64, error) {

	// check wether a file exist

	if !util.FileExists(fname) {
		return -1, fmt.Errorf("file %s could not be found", fname)
	}

	index, err := fs.index.AddItem(fname, fs.currentWorkingDir)
	if err != nil {
		return -1, err
	}

	dir := fmt.Sprintf("%d", index)

	fd := NewFileDirectory(fs, dir, index)

	if err = fd.CreateDirectory(); err != nil {
		return -1, err
	}

	if err = fd.ImportNewFile(fname); err != nil {
		return -1, err
	}

	name, err := fs.index.ContainerName(fd.dir)
	if err != nil {
		return -1, err
	}

	log.Printf("imported %s into %s with version %d", fname, name.file, 0)

	// Checking out the new file so no one else can see it.

	if err = fs.CheckOut(index, FileVersion{0, "0", util.Now()}); err != nil {
		return -1, err
	}

	return index, nil
}

// Generates a new version of a file. Returns the FileVersion and an error.
func (fs *FileSystem) NewVersion(indexNr int64) (FileVersion, error) {

	// Check wether src is locked or not

	dirIdx, err := fs.index.IndexDir(indexNr)
	util.CheckErr(err)

	if name := fs.IsLockedItem(indexNr); name != "" {
		return FileVersion{}, fmt.Errorf("NewVersion error: File %s is checked out by %s", dirIdx, name)
	}

	dir := path.Join(fs.vaultDir, dirIdx)

	fd := NewFileDirectory(fs, dir, indexNr)

	newVersion := fd.NewVersion()

	// Checking out the new file so no one else can see it.

	name, err := fs.index.IndexToFileName(indexNr)
	util.CheckErr(err)

	log.Printf("Created version %d of file %s\n", newVersion.Number, name)

	err = fs.CheckOut(indexNr, newVersion)
	util.CheckErr(err)

	return newVersion, nil
}

// Returns the number of an item
func (fs FileSystem) GetItem(file string) (FileList, error) {
	return fs.index.FileNameToFileList(file)
}

// Creates a new directory inside the current directory, with the correct uid and gid.
func (fs FileSystem) Mkdir(dir string) error {

	// Check wether dname is an int. We don't want that, because the number could interfere with the fileindex.
	if _, err := strconv.Atoi(dir); err == nil {
		return fmt.Errorf("please change %s into a string, now it is a number", dir)
	}

	err := os.Mkdir(dir, 0777)
	util.CheckErr(err)

	err = os.Chown(dir, fs.userUid, fs.vaultUid)
	util.CheckErr(err)

	log.Printf("Created directory: %s\n", dir)

	return nil
}

func (fs *FileSystem) Chdir(dir string) error {

	if !util.DirExists(dir) {
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

// ListWD lists the sorted directories and files of the current working directory.
func (fs FileSystem) ListWD() ([]FileInfo, error) {
	return fs.ListDir(path.Join(fs.vaultDir, fs.currentWorkingDir))
}

// ListDir lists the sorted directories and files within the specified directory name,
// as long as the directory is inside the vault.
func (fs FileSystem) ListDir(dirName string) ([]FileInfo, error) {
	dirList, err := os.ReadDir(dirName)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirName, err)
	}

	var list []FileInfo

	// Add parent directory option if not at the vault root
	if path.Clean(dirName) != fs.VaultDir() {
		list = append(list, FileInfo{Dir: true, FileName: ".."})
	}

	for _, subDir := range dirList {
		if num, err := util.Atoi64(subDir.Name()); err == nil {
			// Handle file entries
			fileName, err := fs.index.IndexToFileName(int64(num))
			if err != nil {
				return nil, fmt.Errorf("failed to convert index to file name: %w", err)
			}

			idx, err := fs.index.FileNameToIndex(fileName)
			if err != nil {
				return nil, fmt.Errorf("failed to convert file name to index: %w", err)
			}

			lockedUser := fs.IsLockedItem(idx)
			list = append(list, FileInfo{
				Dir:             false,
				FileName:        fileName,
				FilePath:        dirName,
				FileLocked:      lockedUser != "",
				FileLockedOutBy: lockedUser,
			})
		} else {
			// Handle directory entries
			list = append(list, FileInfo{Dir: true, FileName: subDir.Name(), FilePath: dirName})
		}
	}

	// Sort the list by directory and then by file name
	sort.Slice(list, func(i, j int) bool {
		if list[i].Dir != list[j].Dir {
			return list[i].Dir // Directories first
		}
		return strings.ToUpper(list[i].Name()) < strings.ToUpper(list[j].Name())
	})

	return list, nil
}

// ListTree lists the sorted directories and files within the specified directory name,
// including subdirectories, as long as the directory is inside the vault.
func (fs FileSystem) ListTree(dirName string) ([]FileInfo, error) {
	return fs.listTree(dirName)
}

// listTree is a helper function to recursively list all directories and files.
func (fs FileSystem) listTree(dirName string) ([]FileInfo, error) {
	dirList, err := fs.ListDir(dirName)
	if err != nil {
		return nil, err
	}

	if len(dirList) >= 1 {
		if dirList[0].Name() == ".." {
			dirList = dirList[1:]
		}
	}

	if len(dirList) == 0 {
		return nil, fmt.Errorf("directory is empty")
	}

	list := make([]FileInfo, len(dirList))

	for i, elem := range dirList {
		list[i] = elem

		if elem.IsDir() {
			// Recursively append subdirectory contents
			subDirContents, err := fs.listTree(path.Join(dirName, elem.Name()))
			if err != nil {
				return nil, err
			}
			list = append(list, subDirContents...)
		}
	}

	return list, nil
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

// Check whether the itemnr is locked.
// Returns the name of the user who locked it or empty when not locked.
func (fs FileSystem) IsLockedItem(itemNr int64) string {

	for _, item := range fs.lockedIndex {
		if item.fileNr == itemNr {
			return item.userName
		}
	}

	return "" // Nothing found
}

// Checkout means locking a itemnr so that only you can use it.
func (fs *FileSystem) CheckOut(itemNr int64, version FileVersion) error {

	fs.ReadLockedIndex() // update the index

	// Set file mode 0700

	dir, err := fs.index.IndexDir(itemNr)
	util.CheckErr(err)

	fd := NewFileDirectory(fs, path.Join(fs.vaultDir, dir), itemNr)
	fd.OpenItemVersion(version)

	// check whether the itemnr is locked

	if usr := fs.IsLocked(itemNr, version); usr != "" {

		return fmt.Errorf("file %d-%d is locked by user %v", itemNr, version.Number, usr)

	} else {

		fs.lockedIndex = append(fs.lockedIndex, LockedIndex{itemNr, version.Number, fs.user})

		fs.WriteLockedIndex()

		name, err := fs.index.IndexToFileName(itemNr)
		util.CheckErr(err)

		log.Printf("Checked out version %d of file %s\n", version.Number, name)

		return nil
	}
}

// Checkin means unlocking an itemnr.
// The description and long description are meant for storage.
func (fs *FileSystem) CheckIn(itemNr int64, version FileVersion, descr, longdescr string) error {

	// Set file mode 0755

	dir, err := fs.index.IndexDir(itemNr)
	util.CheckErr(err)

	fd := NewFileDirectory(fs, path.Join(fs.vaultDir, dir), itemNr)

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

		name, err := fs.index.IndexToFileName(itemNr)
		util.CheckErr(err)

		log.Printf("Checked in version %d of file %s", version.Number, name)

		return nil
	}
}

// Rename a file, for instance when the user wants to use a file
// with a specified numbering system.
func (fs *FileSystem) FileRename(src, dest string) error {
	// Check whether src is locked or not
	fileName, err := fs.index.FileNameToIndex(src)
	if err != nil {
		return fmt.Errorf("failed to get index for %s: %w", src, err)
	}

	if name := fs.IsLockedItem(fileName); name != "" {
		return fmt.Errorf("FileRename error: File %s is checked out by %s", src, name)
	}

	// Check whether dest exists
	if item, err := fs.index.FileNameToFileList(dest); err == nil {
		return fmt.Errorf("file %s already exists and is stored in %s", dest, item.Index())
	}

	// Rename the file from src to dest
	idx, err := fs.index.FileNameToFileList(src)
	if err != nil {
		return fmt.Errorf("failed to get file list for %s: %w", src, err)
	}

	fd := NewFileDirectory(fs, idx.Index(), fileName)

	if err = fd.fileRename(src, dest); err != nil {
		return fmt.Errorf("failed to rename file from %s to %s: %w", src, dest, err)
	}

	// Rename the file in the index
	if err = fs.index.renameItem(src, dest); err != nil {
		return fmt.Errorf("failed to rename item in index: %w", err)
	}

	// Logging
	log.Printf("File %s renamed to %s\n", src, dest)

	return nil
}

// Copy the latest version of a file.
func (fs *FileSystem) FileCopy(src, dest string) error {
	// Check whether src is locked or not
	fileName, err := fs.index.FileNameToIndex(src)
	if err != nil {
		return fmt.Errorf("failed to get index for %s: %w", src, err)
	}

	if name := fs.IsLockedItem(fileName); name != "" {
		return fmt.Errorf("FileCopy error: File %s is checked out by %s", src, name)
	}

	var dst string
	if strings.Contains(dest, "/") {
		num := strings.LastIndexByte(dest, '/')
		dst = path.Join(fs.currentWorkingDir, dest[:num])
		fmt.Printf("dst = %s\n", dst)

		destPath := path.Join(fs.vaultDir, dst)
		if !util.DirExists(destPath) {
			return fmt.Errorf("directory %s doesn't exist", destPath)
		}
	} else {
		dst = fs.currentWorkingDir
	}

	if _, err = fs.index.FileNameToFileList(dest); err == nil {
		return fmt.Errorf("file %s already exists", dest)
	}

	file, err := fs.index.ContainerName(src)
	if err != nil {
		return fmt.Errorf("failed to get container name for %s: %w", src, err)
	}

	srcIndex, err := fs.index.FileNameToFileList(src)
	if err != nil {
		return fmt.Errorf("failed to get file list for %s: %w", src, err)
	}

	srcFd := NewFileDirectory(fs, srcIndex.Index(), file.index)

	_, destFile := path.Split(dest)
	destIndex, err := fs.index.AddItem(destFile, dst)
	if err != nil {
		return err
	}

	destDir, err := fs.index.FileNameToFileList(destFile)
	if err != nil {
		return fmt.Errorf("failed to get file list for %s: %w", destFile, err)
	}

	destFd := NewFileDirectory(fs, destDir.Index(), destIndex)
	if err := destFd.CreateDirectory(); err != nil {
		return err
	}

	// Copy the file from src to dest
	version := srcFd.LatestVersion()
	srcFile := path.Join(srcFd.dir, version.Pretty, src)

	if err = destFd.ImportNewFile(srcFile); err != nil {
		return err
	}

	if err := destFd.fileRename(src, destFile); err != nil {
		return fmt.Errorf("failed to rename file from %s to %s: %w", src, destFile, err)
	}

	// Logging
	log.Printf("File %s copied to %s\n", src, dest)

	return nil
}

// Moves a file to a different directory.
func (fs *FileSystem) FileMove(fileName, destDir string) error {
	// Check whether src is locked or not
	fileNr, err := fs.index.FileNameToIndex(fileName)
	if err != nil {
		return fmt.Errorf("failed to get index for %s: %w", fileName, err)
	}

	if name := fs.IsLockedItem(fileNr); name != "" {
		return fmt.Errorf("FileMove error: File %s is checked out by %s", fileName, name)
	}

	// "normalize" the destDir
	dir, err := filepath.Abs(destDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", destDir, err)
	}

	if !util.DirExists(dir) {
		return fmt.Errorf("directory %s doesn't exist", destDir)
	}

	if len(dir) == len(fs.vaultDir) { // case when file moved to the root
		dir = ""
	} else if strings.HasPrefix(dir, fs.vaultDir) {
		dir = dir[len(fs.vaultDir)+1:]
	}

	// Move file
	fname, err := fs.index.FileNameToFileList(fileName)
	if err != nil {
		return fmt.Errorf("failed to get file list for %s: %w", fileName, err)
	}

	dest := path.Join(destDir, fname.Index())
	if err := os.Rename(fname.Index(), dest); err != nil {
		return fmt.Errorf("failed to move file from %s to %s: %w", fname.Index(), dest, err)
	}

	if err := os.Chown(dest, fs.userUid, fs.vaultUid); err != nil {
		return fmt.Errorf("failed to change ownership of %s: %w", dest, err)
	}

	// Move file in FileIndex
	if err := fs.index.moveItem(fileName, dir); err != nil {
		return fmt.Errorf("failed to move item in index: %w", err)
	}

	// Logging
	log.Printf("File %s moved to %s\n", fileName, destDir)

	return nil
}

// Copy a directory.
func (fs *FileSystem) DirectoryCopy(src, dest string) error {

	// Check wether dest is a number

	if util.IsNumber(dest) {
		return fmt.Errorf("directory %s is a number", dest)
	}

	srcFiles, err := fs.ListTree(src)
	if err != nil {
		return err
	}

	// Check wether dest directory exists

	if util.DirExists(dest) {
		return fmt.Errorf("directory %s exists", dest)
	}

	// Check wether files are Checked-Out

	if err := fs.checkOutFiles(srcFiles); err != nil {
		return err
	}

	//
	// Copy the directory
	//

	// Copy the file in the index

	for _, elem := range srcFiles {
		if elem.IsDir() {
			fs.Mkdir(path.Join(fs.currentWorkingDir, elem.Name()))
		} else {
			if err := fs.FileCopy(elem.Name(), dest); err != nil {
				return err
			}
		}
	}

	// Logging

	log.Printf("Directory %s copied to %s\n", src, dest)

	return nil
}

// Move a directory.
func (fs FileSystem) DirectoryMove(src, dst string) error {

	// Check wether dest is a number
	if util.IsNumber(dst) {
		return fmt.Errorf("destination directory %s cannot be a number", dst)
	}

	// Check if source directory exists
	if !util.DirExists(src) {
		return fmt.Errorf("source directory %s does not exist", src)
	}

	// Check wether dest directory exists
	if util.DirExists(dst) {
		return fmt.Errorf("destination directory %s already exists", dst)
	}

	// List files in the source directory
	srcFiles, err := fs.ListTree(src)
	if err != nil {
		return err
	}

	// Check if files are Checked-Out
	if err := fs.checkOutFiles(srcFiles); err != nil {
		return err
	}

	//
	// Move the directory
	//

	dstFiles := make([]FileList, len(srcFiles))

	for k, v := range srcFiles {
		l := strings.Index(v.Path(), src)
		if len(v.Path())-(l+len(src)) != 0 {
			rest := v.Path()[l+len(src)+1:]
			dstFiles[k].dir = path.Join(dst, rest)

		} else {
			dstFiles[k].dir = dst
		}
	}
	// Move the file(s) in the index
	for k, elem := range srcFiles {
		if !elem.IsDir() {
			if err := fs.index.moveItem(elem.Name(), dstFiles[k].dir); err != nil {
				return err
			}
		}
	}

	// Move the directory and its contents
	if err := os.Rename(src, dst); err != nil {
		return err
	}

	// Check wether files and directories are moved and in the right place
	for _, elem := range dstFiles {
		_, err := os.Stat(path.Join(fs.vaultDir, elem.file))
		if err != nil {
			return fmt.Errorf("unable to locate file %s, error %s", elem.file, err)
		}
	}

	// Log the successful move operation
	log.Printf("Successfully moved directory from %s to %s", src, dst)

	return nil
}

// Check whether a FileInfo list has checkouts.
// Returns each checkout as an error or nil.
func (fs FileSystem) checkOutFiles(list []FileInfo) error {
	var checkOutErrors []error

	for _, item := range list {
		if item.IsLocked() {
			errMsg := fmt.Sprintf("%s is checked out by %s", item.Name(), item.LockedOutBy())
			checkOutErrors = append(checkOutErrors, fmt.Errorf("%v", errMsg))
		}
	}

	if len(checkOutErrors) > 0 {
		return fmt.Errorf("check out errors: %v", checkOutErrors)
	}
	return nil
}

// Splits a file with its base and extension
func SplitExt(path string) (base, ext string) {
	ext = filepath.Ext(path)
	base = path[:len(path)-len(ext)]
	return
}
