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
	containerNumber string // The number of the file
	version         int16  // The number of the version
	userName        string // Who checked this file out
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

	// check whether the critical directories exist.

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

	// retrieve the values
	if err = fs.ReadLockedIndex(); err != nil {
		return nil, err
	}

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
func (fs *FileSystem) ReadLockedIndex() error {

	buf, err := os.ReadFile(fs.lockedCvs)
	if err != nil {
		return fmt.Errorf("error reading %s. The error is %w", fs.lockedCvs, err)
	}

	r := csv.NewReader(bytes.NewBuffer(buf))
	r.Comma = ':'

	records, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("error processing csv file %s", fs.lockedCvs)
	}

	if len(records) == 0 {
		return fmt.Errorf("too less lines in file %s", fs.lockedCvs)
	}

	records = records[1:]

	fs.lockedIndex = nil

	var list = LockedIndex{}

	for _, record := range records {
		list.containerNumber = record[0]
		list.version, _ = util.Atoi16(record[1])
		list.userName = record[2]

		fs.lockedIndex = append(fs.lockedIndex, list)
	}

	return nil
}

func (fs *FileSystem) WriteLockedIndex() error {

	records := [][]string{
		{"ContainerNumber", "Version", "UserName"},
	}

	for _, list := range fs.lockedIndex {
		records = append(records, []string{
			list.containerNumber,
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

	if err = os.WriteFile(fs.lockedCvs, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing %s", fs.lockedCvs)
	}

	err = os.Chown(fs.lockedCvs, fs.userUid, fs.vaultUid)
	util.CheckErr(err)

	return nil
}

// import a file inside the PDM. When you import a file the attributes also gets imported,
// which means uploaded to the server.
// When you import a file or files you are placing the new file in the current directory.
// The new file inside the PDM gets a revision number automatically.
// The function returns the number of the imported file or an error.
func (fs *FileSystem) ImportFile(fileName string) (*FileList, error) {

	// check whether a file exist

	if !util.FileExists(fileName) {
		return nil, fmt.Errorf("file %s could not be found", fileName)
	}

	fl, err := fs.index.AddItem(fs.currentWorkingDir, fileName)
	if err != nil {
		return nil, err
	}

	fd := NewFileDirectory(fs, *fl)

	if err = fd.CreateDirectory(); err != nil {
		return nil, err
	}

	if err = fd.ImportNewFile(fileName); err != nil {
		return nil, err
	}

	log.Printf("imported %s into %s with version %d", fileName, fl.fileName, 0)

	// Checking out the new file so no one else can see it.

	if err = fs.CheckOut(*fl, FileVersion{0, "0", util.Now()}); err != nil {
		return nil, err
	}

	return fl, nil
}

// Generates a new version of a file. Returns the FileVersion and an error.
func (fs *FileSystem) NewVersion(fl FileList) (FileVersion, error) {

	// Check whether src is locked or not

	if name := fs.IsLockedItem(fl.containerNumber); name != "" {
		return FileVersion{}, fmt.Errorf("NewVersion error: File %s is checked out by %s", fl, name)
	}

	fd := NewFileDirectory(fs, fl)

	newVersion := fd.NewVersion()

	// Checking out the new file so no one else can see it.

	log.Printf("Created version %d of file %s\n", newVersion.Number, fl.fileName)

	err := fs.CheckOut(fl, newVersion)
	util.CheckErr(err)

	return newVersion, nil
}

// Returns the number of an item
func (fs FileSystem) GetItem(dir, file string) (FileList, error) {
	return fs.index.FileNameToFileList(dir, file)
}

// Creates a new directory inside the current directory, with the correct uid and gid.
func (fs FileSystem) Mkdir(dir string) error {

	// Check whether dname is an int. We don't want that, because the number could interfere with the fileindex.
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
		list = append(list, FileInfo{isDir: true, name: ".."})
	}

	for _, subDir := range dirList {
		if _, err := util.Atoi64(subDir.Name()); err == nil {
			// Handle file entries
			fileName, err := fs.index.ContainerNumberToFileName(subDir.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to convert index to file name: %w", err)
			}

			idx, err := fs.index.FileNameToContainerNumber(dirName, fileName)
			if err != nil {
				return nil, fmt.Errorf("failed to convert file name to index: %w", err)
			}

			lockedUser := fs.IsLockedItem(idx)
			list = append(list, FileInfo{
				isDir:           false,
				name:            fileName,
				dir:             dirName,
				fileLocked:      lockedUser != "",
				fileLockedOutBy: lockedUser,
			})
		} else {
			// Handle directory entries
			list = append(list, FileInfo{isDir: true, name: subDir.Name(), dir: dirName})
		}
	}

	// Sort the list by directory and then by file name
	sort.Slice(list, func(i, j int) bool {
		if list[i].isDir != list[j].isDir {
			return list[i].isDir // Directories first
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

// Check whether the conatiner number is locked.
// Returns the name of the user who locked it or empty when not locked.
func (fs FileSystem) IsLocked(containerNumber string, version FileVersion) string {
	for _, item := range fs.lockedIndex {
		if item.containerNumber == containerNumber && item.version == version.Number {
			return item.userName
		}
	}
	return "" // Nothing found
}

// Check whether the container number is locked.
// Returns the name of the user who locked it or empty when not locked.
func (fs FileSystem) IsLockedItem(containerNumber string) string {
	for _, item := range fs.lockedIndex {
		if item.containerNumber == containerNumber {
			return item.userName
		}
	}
	return "" // Nothing found
}

// Checkout means locking a conainer number so that only you can use it.
func (fs *FileSystem) CheckOut(fl FileList, version FileVersion) error {
	// update the index
	if err := fs.ReadLockedIndex(); err != nil {
		return err
	}

	// Set file mode 0700

	fd := NewFileDirectory(fs, fl)
	fd.OpenItemVersion(version)

	// check whether the itemnr is locked

	if usr := fs.IsLocked(fl.containerNumber, version); usr != "" {

		return fmt.Errorf("file %s-%d is locked by user %v", fl.containerNumber, version.Number, usr)

	} else {

		fs.lockedIndex = append(fs.lockedIndex, LockedIndex{fl.containerNumber, version.Number, fs.user})

		if err := fs.WriteLockedIndex(); err != nil {
			return err
		}

		log.Printf("Checked out version %d of file %s\n", version.Number, fl.fileName)

		return nil
	}
}

// Checkin means unlocking a container number.
// The description and long description are meant for storage.
func (fs *FileSystem) CheckIn(fl FileList, version FileVersion, descr, longdescr string) error {

	// Set file mode 0755

	fd := NewFileDirectory(fs, fl)

	fd.StoreData(version, descr, longdescr)

	fd.CloseItemVersion(version)

	// check whether the itemnr is locked

	usr := fs.IsLocked(fl.containerNumber, version)

	if usr != fs.user {

		return fmt.Errorf("file %s-%d is locked by user %s", fl.containerNumber, version.Number, usr)

	} else {

		// Remove item from index

		var nr int
		for i, y := range fs.lockedIndex {
			if y.containerNumber == fl.containerNumber && y.version == version.Number {
				nr = i
			}
		}

		fs.lockedIndex = slices.Delete(fs.lockedIndex, nr, nr+1)

		if err := fs.WriteLockedIndex(); err != nil {
			return err
		}

		log.Printf("Checked in version %d of file %s", version.Number, fl.fileName)

		return nil
	}
}

// Rename a file, for instance when the user wants to use a file with
// a specified numbering system
func (fs *FileSystem) FileRename(src, dst string) error {
	// Check whether src is locked or not
	cn, err := fs.index.FileNameToContainerNumber(fs.currentWorkingDir, src)
	if err != nil {
		return fmt.Errorf("failed to get container number for %s: %w", src, err)
	}

	// Check for file locked
	if name := fs.IsLockedItem(cn); name != "" {
		return fmt.Errorf("file %s is checked out by %s", src, name)
	}

	// Check whether dst ends with '/'. In that case it is a file move.
	if dst[len(dst)-1] == '/' {
		return fs.fileMove(src, dst)
	}

	// Splitting dst for later use
	dstDir, dstFile := path.Split(dst)

	// Check whether dst exists
	if item, err := fs.index.FileNameToFileList(path.Join(fs.currentWorkingDir, dstDir), dst); err == nil {
		return fmt.Errorf("file %s already exists and is stored in %s", dst, item.ContainerNumber())
	}

	// Check whether dst contains a directory. In that case it is a file move,
	// and after that a file rename.
	if dstDir != "" {
		if err = fs.fileMove(src, dstDir); err != nil {
			return err
		}
		if dstFile != src { // file rename
			wd, err := os.Getwd()
			wcd := fs.currentWorkingDir
			if err != nil {
				return err
			}
			if err = fs.Chdir(dstDir); err != nil {
				return err
			}
			if err = fs.FileRename(src, dstFile); err != nil {
				return err
			}
			if err = os.Chdir(wd); err != nil {
				return err
			}
			fs.currentWorkingDir = wcd
			return nil
		}
	}

	// Rename the file from src to dest in the FileList
	fl, err := fs.index.FileNameToFileList(fs.currentWorkingDir, src)
	if err != nil {
		return fmt.Errorf("failed to get file list for %s: %w", src, err)
	}

	fd := NewFileDirectory(fs, fl)

	if err = fd.fileRename(src, dst); err != nil {
		return fmt.Errorf("failed to rename file from %s to %s: %w", src, dst, err)
	}

	// Rename the file in the index
	if err = fs.index.renameItem(src, dst); err != nil {
		return fmt.Errorf("failed to rename item in index: %w", err)
	}

	// Logging
	log.Printf("File %s renamed to %s\n", src, dst)

	return nil
}

// Copy the latest file from src to dst and returns an error.
func (fs *FileSystem) FileCopy(src, dst string) error {
	// Check whether src is locked
	srcFl, err := fs.index.FileNameToFileList(fs.currentWorkingDir, src)
	if err != nil {
		return fmt.Errorf("failed to get index for %s: %w", src, err)
	}

	if name := fs.IsLockedItem(srcFl.containerNumber); name != "" {
		return fmt.Errorf("FileCopy error: File %s is checked out by %s", src, name)
	}
	fmt.Println("hello")

	// Check whether dest is a file or directory
	if len(dst) > 1 {
		if len(dst)-1 == '/' {
			dst = dst + src
			fmt.Printf("dest = %s\n", dst)
		}
	}

	var dstPath, dstDir string
	if strings.Contains(dst, "/") {
		num := strings.LastIndexByte(dst, '/')
		dstDir = path.Join(fs.currentWorkingDir, dst[:num])
		dstPath = path.Join(fs.vaultDir, dstDir)
		if !util.DirExists(dstPath) {
			return fmt.Errorf("directory %s doesn't exist", dstPath)
		}
	} else {
		dstDir = fs.currentWorkingDir
	}

	srcFd := NewFileDirectory(fs, srcFl)

	_, dstFile := path.Split(dst)

	destFl, err := fs.index.AddItem(dstDir, dstFile)
	if err != nil {
		return err
	}

	dstFd := NewFileDirectory(fs, *destFl)
	if err := dstFd.CreateDirectory(); err != nil {
		return err
	}

	// Copy the file from src to dest
	version := srcFd.LatestVersion()
	srcFile := path.Join(srcFd.dir, version.Pretty, src)

	if err = dstFd.ImportNewFile(srcFile); err != nil {
		return err
	}

	// Rename file
	dstVer := dstFd.LatestVersion()
	dstStr := path.Join(dstFd.dir, dstVer.Pretty)
	if err = os.Rename(path.Join(dstStr, src), path.Join(dstStr, dstFile)); err != nil {
		return fmt.Errorf("failed to rename file from %s to %s: %w", src, dstFile, err)
	}

	// Logging
	log.Printf("File %s copied to %s\n", src, dst)

	return nil
}

// Moves a file to a different directory.
func (fs *FileSystem) fileMove(fileName, destDir string) error {
	// Check whether src is locked or not
	fileNr, err := fs.index.FileNameToContainerNumber(fs.currentWorkingDir, fileName)
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
	fname, err := fs.index.FileNameToFileList(fs.currentWorkingDir, fileName)
	if err != nil {
		return fmt.Errorf("failed to get file list for %s: %w", fileName, err)
	}

	dest := path.Join(destDir, fname.ContainerNumber())
	if err := os.Rename(fname.ContainerNumber(), dest); err != nil {
		return fmt.Errorf("failed to move file from %s to %s: %w", fname.ContainerNumber(), dest, err)
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
func (fs *FileSystem) DirectoryCopy(src, dst string) error {

	// Check whether dest is a number
	if util.IsNumber(dst) {
		return fmt.Errorf("directory %s is a number", dst)
	}

	// Check if source directory exists
	if !util.DirExists(src) {
		return fmt.Errorf("source directory %s does not exist", src)
	}

	// Check whether dest directory exists
	if util.DirExists(dst) {
		return fmt.Errorf("directory %s exists", dst)
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
	// Copy the directory
	//

	dstFiles := make([]FileInfo, len(srcFiles))

	// Populating dstFiles with data from srcFiles
	for k, v := range srcFiles {
		l := strings.Index(v.Path(), src)
		dstFiles[k].isDir = srcFiles[k].isDir
		if len(v.Path())-(l+len(src)) != 0 {
			rest := v.Path()[l+len(src)+1:]
			if !dstFiles[k].IsDir() {
				dstFiles[k].dir = path.Join(dst, rest)
				idx := strings.Index(srcFiles[k].Name(), ".")
				str := srcFiles[k].name[0:idx] + " (copy)" + srcFiles[k].name[idx:]
				dstFiles[k].name = str
			} else {
				dstFiles[k].dir = path.Join(dst, rest)
			}
		} else {
			if !dstFiles[k].IsDir() {
				dstFiles[k].dir = dst
				idx := strings.Index(srcFiles[k].Name(), ".")
				str := srcFiles[k].name[0:idx] + " (copy)" + srcFiles[k].name[idx:]
				dstFiles[k].name = str
			} else {
				dstFiles[k].dir = dst
			}
		}
	}

	fmt.Printf("%v\n", srcFiles)
	fmt.Printf("%v\n", dstFiles)

	for k, elem := range srcFiles {
		if elem.IsDir() {
			if err := fs.Mkdir(dstFiles[k].dir); err != nil {
				return err
			}
		} else {
			if err := fs.FileCopy(path.Join(srcFiles[k].Path(), srcFiles[k].Name()),
				path.Join(dstFiles[k].dir, srcFiles[k].name)); err != nil {
				return err
			}
		}
	}

	// Check whether files and directories are copied to the right place
	for _, elem := range dstFiles {
		_, err := os.Stat(path.Join(fs.vaultDir, elem.name))
		if err != nil {
			return fmt.Errorf("unable to locate file %s, error %s", elem.name, err)
		}
	}

	// Logging
	log.Printf("Directory %s copied to %s\n", src, dst)

	return nil
}

// Move a directory.
func (fs FileSystem) DirectoryRename(src, dst string) error {

	// Check whether dest is a number
	if util.IsNumber(dst) {
		return fmt.Errorf("destination directory %s cannot be a number", dst)
	}

	// Check if source directory exists
	if !util.DirExists(src) {
		return fmt.Errorf("source directory %s does not exist", src)
	}

	// Check whether dest directory exists
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

	dstFiles := make([]FileInfo, len(srcFiles))

	// Populating dstFiles with data from srcFiles
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

	// Check whether files and directories are moved and in the right place
	for _, elem := range dstFiles {
		_, err := os.Stat(path.Join(fs.vaultDir, elem.name))
		if err != nil {
			return fmt.Errorf("unable to locate file %s, error %s", elem.name, err)
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
