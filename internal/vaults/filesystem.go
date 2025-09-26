// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package vaults

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/grd/FreePDM/internal/config"
	"github.com/grd/FreePDM/internal/util"
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
	index          FileIndex
	vaultDir       string
	dataDir        string
	vaultUid       int
	user           string
	userUid        int
	lockedFilesCvs string
	lockedIndex    []LockedIndex
}

// Open implements fs.FS.
func (fs FileSystem) Open(name string) (fs.File, error) {
	panic("unimplemented")
}

const (
	LockedFileCsv = "LockedFiles.csv"
	EmptyFile     = ".empty_file"
)

var (
	vaults, vaultsData string // vaults is the root directory, vaultsdata is the administration part
)

func init() {
	// check for config file
	vaults = config.VaultsDir()

	// or check of enivronment variable
	if vaults == "" {
		var ok bool
		vaults, ok = os.LookupEnv("FREEPDM_VAULTS_DIR")
		if !ok {
			fmt.Println("Error: No FREEPDM_VAULTS_DIR environment set. See installation manual.")
			vaults = "/samba/vaults"
		}
	}
	vaultsData = path.Join(vaults, "/.data")

	// check for critical directories
	util.CriticalDirExist(vaults)
	util.CriticalDirExist(vaultsData)
}

// Constructor
func NewFileSystem(vaultDir, userName string) (fs *FileSystem, err error) {
	fs = new(FileSystem)

	// Check whether vaults directory contains slaches
	parts := strings.Split(vaultDir, "/")
	if len(parts) != 1 {
		log.Fatalf("multi-level vaultDir parameter: %s", vaultDir)
	}

	// Some settings
	fs.vaultDir = path.Join(vaults, vaultDir)
	fs.dataDir = path.Join(vaultsData, vaultDir)
	fs.user = userName

	// check whether the critical directories exist.
	util.CriticalDirExist(fs.vaultDir)
	util.CriticalDirExist(fs.dataDir)

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

	fs.lockedFilesCvs = path.Join(fs.dataDir, LockedFileCsv)

	// retrieve the values
	if err = fs.ReadLockedIndex(); err != nil {
		return nil, err
	}

	err = os.Chdir(fs.vaultDir)
	util.CheckErr(err)

	log.Printf("Vault dir: %s", fs.VaultDir())

	return fs, nil
}

// Returns the Vault directory
func (fs *FileSystem) VaultDir() string {
	return fs.vaultDir
}

// Updates the locked index by reading from the lockedTxt file.
func (fs *FileSystem) ReadLockedIndex() error {

	buf, err := os.ReadFile(fs.lockedFilesCvs)
	if err != nil {
		return fmt.Errorf("error reading %s. The error is %w", fs.lockedFilesCvs, err)
	}

	r := csv.NewReader(bytes.NewBuffer(buf))
	r.Comma = ':'

	records, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("error processing csv file %s", fs.lockedFilesCvs)
	}

	if len(records) == 0 {
		return fmt.Errorf("too less lines in file %s", fs.lockedFilesCvs)
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

// Writes a file in a read-only directory structure
func (fs *FileSystem) DataWriteFile(name string, data []byte) error {
	// Mutex to ensure only one operation modifies permissions at a time
	var permMutex sync.Mutex

	// Acquire exclusive access
	permMutex.Lock()
	defer permMutex.Unlock()

	if err := os.Chmod(vaultsData, 0777); err != nil {
		return fmt.Errorf("error setting directory permissions %s", fs.vaultDir)
	}
	if err := os.Chmod(fs.dataDir, 0777); err != nil {
		return fmt.Errorf("error setting directory permissions %s", fs.dataDir)
	}

	if err := os.WriteFile(name, data, 0644); err != nil {
		return fmt.Errorf("error writing %s", name)
	}
	if err := os.Chown(fs.lockedFilesCvs, fs.userUid, fs.vaultUid); err != nil {
		return fmt.Errorf("error changing ownership %s", name)
	}

	if err := os.Chmod(fs.dataDir, 0555); err != nil {
		return fmt.Errorf("error setting directory permissions %s", fs.vaultDir)
	}

	if err := os.Chmod(vaultsData, 0555); err != nil {
		return fmt.Errorf("error setting directory permissions %s", fs.dataDir)
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

	if err = fs.DataWriteFile(fs.lockedFilesCvs, buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

// import a file inside the PDM. The file name needs to be complete.
// When you import a file the attributes also gets imported,
// which means uploaded to the server.
// When you import a file or files you are placing the new file in the directory dstDir.
// The new file inside the PDM gets a revision number automatically.
// The function returns the FileList structure of the imported file or an error.
func (fs *FileSystem) ImportFile(dstDir, fileName string) (*FileList, error) {

	// check whether the file exist
	if !util.FileExists(fileName) {
		return nil, fmt.Errorf("file %s could not be found", fileName)
	}

	complDstDir := path.Join(fs.vaultDir, dstDir)

	// check whether the directory exist
	if !util.DirExists(complDstDir) {
		return nil, fmt.Errorf("directory %s could not be found", complDstDir)
	}

	fl, err := fs.index.AddItem(dstDir, fileName)
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

	// Checking out the new file so no one else can see it.

	if err = fs.CheckOut(*fl, FileVersion{0, "0", util.Now()}); err != nil {
		return nil, err
	}

	// Verfication
	if err := fs.Verify(*fl); err != nil {
		return nil, err
	}

	log.Printf("imported %s into %s with version %d", fileName, fl.Name, 0)

	return fl, nil
}

// import a url inside the PDM. For the rest it works the same as ImportFile()
func (fs *FileSystem) ImportUrl(dstDir, url string) (*FileList, error) {
	// Reuse the same HTTP client with an extended timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Check if the URL is reachable with a HEAD request
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach URL: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("URL returned status code %d", resp.StatusCode)
	}

	// Calculate the destination directory
	complDstDir := path.Join(fs.vaultDir, dstDir)

	// Check if the destination directory exists
	if !util.DirExists(complDstDir) {
		return nil, fmt.Errorf("directory %s could not be found", complDstDir)
	}

	// Download the file contents
	resp, err = client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download URL: %v", err)
	}
	defer resp.Body.Close()

	// Extract the file name from the URL
	_, file := path.Split(url)

	// Save the URL contents to a local file
	localFilePath := path.Join(complDstDir, file)
	out, err := os.Create(localFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create local file: %v", err)
	}
	defer out.Close()
	defer os.Remove(localFilePath)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to save URL contents to local file: %v", err)
	}

	// Add the file to the index
	fl, err := fs.index.AddItem(dstDir, file)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to index: %v", err)
	}

	fd := NewFileDirectory(fs, *fl)

	// Create the directory structure
	if err = fd.CreateDirectory(); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	if err = fd.ImportNewFile(localFilePath); err != nil {
		return nil, err
	}

	// Check out the file to make it inaccessible to others
	if err = fs.CheckOut(*fl, FileVersion{0, "0", util.Now()}); err != nil {
		return nil, fmt.Errorf("failed to check out file: %v", err)
	}

	// Verfication
	if err := fs.Verify(*fl); err != nil {
		return nil, err
	}

	log.Printf("imported %s into %s with version %d", url, fl.Name, 0)

	return fl, nil
}

// Allocates an empty file container. This can be handy if you want to setup a structure.
// The new container inside the PDM gets a revision number automatically and is checked out.
// The function returns the FileList structure of the imported file or an error.
func (fs *FileSystem) Allocate(dstDir string) (*FileList, error) {
	file, err := os.CreateTemp("", EmptyFile)
	if err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	i := strings.Index(file.Name(), EmptyFile)
	newfile := file.Name()[0 : i+len(EmptyFile)]
	if err := os.Rename(file.Name(), newfile); err != nil {
		return nil, err
	}

	fl, err := fs.ImportFile(dstDir, newfile)
	if err != nil {
		return nil, err
	}

	return fl, nil
}

// Assigns the file name to the container.
// The new container inside the PDM gets a revision number automatically.
// The function returns the FileList structure of the imported file or an error.
func (fs *FileSystem) Assign(containerNumber, fileName string) error {

	fl, err := fs.index.ContainerNumberToFileList(containerNumber)
	if err != nil {
		return err
	}

	complDstDir := path.Join(fs.vaultDir, fl.Path, fl.ContainerNumber)
	complDstFile := path.Join(complDstDir, "0", fileName)
	ef := path.Join(complDstDir, "0", EmptyFile)

	// check whether the directory exist
	if !util.DirExists(complDstDir) {
		return fmt.Errorf("directory %s could not be found", complDstDir)
	}

	// check whether the file exist
	if !util.FileExists(complDstFile) {
		return fmt.Errorf("file %s could not be found", fileName)
	}

	// check whether the EmptyFile exist
	if !util.FileExists(ef) {
		return fmt.Errorf("file %s could not be found", fileName)
	}

	// rename in the index
	if err = fs.index.renameItem(fl, fileName); err != nil {
		return err
	}

	// remove the empty file
	if err = os.Remove(ef); err != nil {
		return err
	}

	log.Printf("assigned %s into %s", fileName, fl.ContainerNumber)

	return nil
}

// Generates a new version of a file. Returns the FileVersion and an error.
func (fs *FileSystem) NewVersion(fl FileList) (*FileVersion, error) {

	// Check whether src is locked or not

	if name := fs.IsLockedItem(fl.ContainerNumber); name != "" {
		return &FileVersion{}, fmt.Errorf("NewVersion error: File %s is checked out by %s", fl, name)
	}

	fd := NewFileDirectory(fs, fl)

	newVersion, err := fd.NewVersion()
	if err != nil {
		return nil, err
	}

	// Checking out the new file so no one else can see it.

	log.Printf("Created version %d of file %s\n", newVersion.Number, fl.Name)

	if err = fs.CheckOut(fl, *newVersion); err != nil {
		return nil, err
	}

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

	if err := os.Mkdir(dir, 0777); err != nil {
		return err
	}

	if err := os.Chown(dir, fs.userUid, fs.vaultUid); err != nil {
		return err
	}

	log.Printf("Created directory: %s\n", dir)

	return nil
}

// Chdir creates a directory or an error
func (fs *FileSystem) Chdir(dir string) error {

	if !util.DirExists(dir) {
		return fmt.Errorf("dir %s does not exist", dir)
	}

	err := os.Chdir(dir)
	if err != nil {
		return err
	}

	log.Printf("Changed directory: %s\n", dir)

	return nil
}

// ListWD lists the sorted directories and files of the current working directory.
func (fs FileSystem) ListWD() ([]FileInfo, error) {
	return fs.ListDir(fs.Getwd())
}

// func (fs.FileSystem) Open() {

// }

// ListDir lists the sorted directories and files within the specified directory name,
// as long as the directory is inside the vault.
func (fs FileSystem) ListDir(dirName string) ([]FileInfo, error) {
	joinDir := path.Join(fs.vaultDir, dirName)

	dirList, err := os.ReadDir(joinDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirName, err)
	}

	list := make([]FileInfo, 0, len(dirList))

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

	if len(dirList) == 0 { // empty directory. Nothing to add but no error message
		return nil, nil
	}

	list := make([]FileInfo, 0, len(dirList))

	for _, elem := range dirList {
		list = append(list, elem)

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

	if usr := fs.IsLocked(fl.ContainerNumber, version); usr != "" {

		return fmt.Errorf("file %s-%d is locked by user %v", fl.ContainerNumber, version.Number, usr)

	} else {

		fs.lockedIndex = append(fs.lockedIndex, LockedIndex{fl.ContainerNumber, version.Number, fs.user})

		if err := fs.WriteLockedIndex(); err != nil {
			return err
		}

		log.Printf("Checked out version %d of file %s\n", version.Number, fl.Name)

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

	usr := fs.IsLocked(fl.ContainerNumber, version)

	if usr != fs.user {

		return fmt.Errorf("file %s-%d is locked by user %s", fl.ContainerNumber, version.Number, usr)

	} else {

		// Remove item from index

		var nr int
		for i, y := range fs.lockedIndex {
			if y.containerNumber == fl.ContainerNumber && y.version == version.Number {
				nr = i
			}
		}

		fs.lockedIndex = slices.Delete(fs.lockedIndex, nr, nr+1)

		if err := fs.WriteLockedIndex(); err != nil {
			return err
		}

		log.Printf("Checked in version %d of file %s", version.Number, fl.Name)

		return nil
	}
}

// Rename a file, for instance when the user wants to use a file with
// a specified numbering system
func (fs *FileSystem) FileRename(src, dst string) error {
	// Check whether src is empty
	if src == "" {
		return errors.New("empty source")
	}

	// Splitting src
	s, srcFile := path.Split(src)

	srcAbs, err := filepath.Abs(s)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", src, err)
	}

	srcDir, err := fs.AbsNormal(srcAbs)
	if err != nil {
		return err
	}

	srcFl, err := fs.index.FileNameToFileList(srcDir, srcFile)
	if err != nil {
		return fmt.Errorf("failed to get container number for %s: %w", src, err)
	}

	// Check whether dst is empty
	if dst == "" {
		return errors.New("empty destination")
	}

	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", dst, err)
	}

	// Check wether dst ends with a file or directory
	info, err := os.Stat(dstAbs)
	if err == nil { // Note, it is not "!="" but "==" !!!
		if info.IsDir() {
			dstAbs = path.Join(dstAbs, srcFile)
		}
	}

	// Splitting dst
	d, dstFile := path.Split(dstAbs)

	dstDir, err := fs.AbsNormal(d)
	if err != nil {
		return err
	}

	dstFl := FileList{Path: dstDir, Name: dstFile, ContainerNumber: srcFl.ContainerNumber}

	// Check whether dst exists
	if item, err := fs.index.FileNameToFileList(dstDir, dstFile); err == nil {
		return fmt.Errorf("file %s already exists and is stored in %s", dst, item.ContainerNumber)
	}

	// Check for src file is locked
	if name := fs.IsLockedItem(srcFl.ContainerNumber); name != "" {
		return fmt.Errorf("file %s is checked out by %s", src, name)
	}

	// Check whether a src and dst are not the same.
	// If that is the case then a fileRename is required.
	if srcFl.Name != dstFl.Name {

		fd := NewFileDirectory(fs, srcFl)

		if err = fd.fileRename(srcFl.Name, dstFl.Name); err != nil {
			return fmt.Errorf("failed to rename file from %s to %s: %w", src, dst, err)
		}

		// Rename the file in the index
		if err = fs.index.renameItem(srcFl, dstFl.Name); err != nil {
			return fmt.Errorf("failed to rename item in index: %w", err)
		}
	}

	// Check whether dst ends with '/' and dstDir contains a directory.
	// In both cases it is a file move operation.
	if dst[len(dst)-1] == '/' || srcDir != dstDir {
		if err := fs.fileMove(srcFl, dstFl); err != nil {
			return err
		}
	}

	// Verification
	item, err := fs.index.FileNameToFileList(dstDir, dstFile)
	if err != nil {
		return err
	}
	if err = fs.Verify(item); err != nil {
		return err
	}

	// Logging
	log.Printf("File %s renamed to %s\n", path.Join(srcFl.Path, srcFl.Name), path.Join(dstDir, dstFile))

	return nil
}

// Moves a file to a different directory.
func (fs *FileSystem) fileMove(src, dst FileList) error {
	source := path.Join(src.Path, src.ContainerNumber)
	dest := path.Join(dst.Path, src.ContainerNumber)

	pwd := fs.Getwd()
	if err := os.Chdir(fs.vaultDir); err != nil {
		return err
	}
	if err := os.Rename(source, dest); err != nil {
		return fmt.Errorf("failed to move file from %s to %s: %w", source, dest, err)
	}
	if err := os.Chown(dest, fs.userUid, fs.vaultUid); err != nil {
		return fmt.Errorf("failed to change ownership of %s: %w", dest, err)
	}
	if err := os.Chdir(pwd); err != nil {
		return err
	}

	// Move file in FileIndex
	if err := fs.index.MoveItem(src, dst.Path); err != nil {
		return fmt.Errorf("failed to move item in index: %w", err)
	}

	return nil
}

// Copy the latest file from src to dst and returns an error.
func (fs *FileSystem) FileCopy(src, dst string) error {
	// Check whether src is empty
	if src == "" {
		return errors.New("empty source file")
	}

	// Splitting src
	s, srcFile := path.Split(src)

	srcAbs, err := filepath.Abs(s)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", src, err)
	}

	srcDir, err := fs.AbsNormal(srcAbs)
	if err != nil {
		return err
	}

	// Check whether dst is empty
	if dst == "" {
		return errors.New("empty destination file")
	}

	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", dst, err)
	}

	// Check wether dst ends with a file or directory
	info, err := os.Stat(dstAbs)
	if err == nil { // Note, it is not "!="" but "==" !!!
		if info.IsDir() {
			dstAbs = path.Join(dstAbs, srcFile)
		}
	}

	// Splitting dst
	d, dstFile := path.Split(dstAbs)

	dstDir, err := fs.AbsNormal(d)
	if err != nil {
		return err
	}

	// // setting up a fake dstFl, without a container number inside
	// dstFl := FileList{dir: dstDir, fileName: dstFile}

	// Check whether dst exists
	if item, err := fs.index.FileNameToFileList(dstDir, dstFile); err == nil {
		return fmt.Errorf("file %s already exists and is stored in %s", dst, item.Path)
	}

	// Check for src file is locked
	srcFl, err := fs.index.FileNameToFileList(srcDir, srcFile)
	if err != nil {
		return fmt.Errorf("failed to get index for %s: %w", src, err)
	}

	if name := fs.IsLockedItem(srcFl.ContainerNumber); name != "" {
		return fmt.Errorf("file %s is checked out by %s", src, name)
	}

	// Check whether dst is a file or directory
	if len(dst) > 1 {
		if len(dst)-1 == '/' {
			dst = dst + src
		}
	}

	// Check for a subdirectory inside dst
	var dstPath string
	if strings.Contains(dst, "/") {
		num := strings.LastIndexByte(dst, '/')
		dstDir = path.Join(fs.Getwd(), dst[:num])
		dstPath = path.Join(fs.vaultDir, dstDir)
		if !util.DirExists(dstPath) {
			return fmt.Errorf("directory %s doesn't exist", dstPath)
		}
	} else {
		dstDir = fs.Getwd()
	}

	srcFd := NewFileDirectory(fs, srcFl)

	// Check wether dst ends with '/'
	if dst[len(dst)-1] == '/' {
		dstFile = src
	} else {
		_, dstFile = path.Split(dst)
	}

	// Check whether dst exists
	if _, err := fs.index.FileNameToFileList(dstDir, dstFile); err == nil {
		return fmt.Errorf("file %s already exists and is stored in %s", dstFile, dstDir)
	}

	dstFl, err := fs.index.AddItem(dstDir, dstFile)
	if err != nil {
		return err
	}

	dstFd := NewFileDirectory(fs, *dstFl)
	if err := dstFd.CreateDirectory(); err != nil {
		return err
	}

	// Copy the file from src to dest
	version := srcFd.LatestVersion()
	newFile := path.Join(srcFd.dir, version.Pretty, src)

	if err = dstFd.ImportNewFile(newFile); err != nil {
		return err
	}

	// Rename file, but only when dst doesn't end with '/' (which means a directory)
	if dst[len(dst)-1] != '/' {
		dstVer := dstFd.LatestVersion()
		dstStr := path.Join(dstFd.dir, dstVer.Pretty)
		if err = os.Rename(path.Join(dstStr, src), path.Join(dstStr, dstFile)); err != nil {
			return fmt.Errorf("failed to rename file from %s to %s: %w", src, dstFile, err)
		}
	}

	// Verification
	item, err := fs.index.FileNameToFileList(dstDir, dstFile)
	if err != nil {
		return err
	}
	if err = fs.Verify(item); err != nil {
		return err
	}

	// Logging
	log.Printf("File %s copied to %s\n", path.Join(srcFl.Path, srcFl.Name), path.Join(dstDir, dstFile))
	// log.Printf("File %s copied to %s\n", src, dst)

	return nil
}

// Copy a directory.
func (fs *FileSystem) DirectoryCopy(src, dst string) error {
	// Check whether src is empty
	if src == "" {
		return errors.New("empty source directory")
	}

	// Check whether dst is empty
	if dst == "" {
		return errors.New("empty destination directory")
	}

	// Check whether dst is a number
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

	// append the src directory to srcFiles
	srcZero := make([]FileInfo, 1)
	srcZero[0] = FileInfo{dir: src, isDir: true}
	srcFiles = append(srcZero, srcFiles...)

	//
	// Copy the directory
	//

	dstFiles := make([]FileInfo, len(srcFiles))

	// Populating dstFiles with data from srcFiles
	l := len(srcFiles[0].dir)
	for k, v := range srcFiles {
		dstFiles[k].isDir = srcFiles[k].isDir
		dstFiles[k].dir = dst + v.dir[l:]
		dstFiles[k].name = v.name
	}

	cwd := fs.Getwd()

	for k, elem := range dstFiles {
		dstJoinPath := path.Join(dstFiles[k].Path(), dstFiles[k].Name())
		if elem.IsDir() {
			if err := fs.Mkdir(path.Join(fs.vaultDir, elem.dir, elem.name)); err != nil {
				return err
			}
		} else {
			fs.Chdir(srcFiles[k].Path())

			if err := fs.FileCopy(srcFiles[k].Name(),
				path.Join("..", dstFiles[k].dir, dstFiles[k].name)); err != nil {
				return err
			}
		}

		// Verification
		dstPathContainer := path.Join(fs.vaultDir, dstFiles[k].dir, dstFiles[k].containerNumber)
		_, err := os.Stat(dstPathContainer)
		if err != nil {
			return fmt.Errorf("unable to locate file %s, error %s", dstJoinPath, err)
		}
	}

	fs.Chdir(cwd)

	// Logging
	log.Printf("Directory %s copied to %s\n", src, dst)

	return nil
}

// Move a directory.
func (fs FileSystem) DirectoryRename(src, dst string) error {
	// Check whether src is empty
	if src == "" {
		return errors.New("empty source directory")
	}

	// Check whether dst is empty
	if dst == "" {
		return errors.New("empty destination directory")
	}

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
		dstFiles[k] = v
		l := strings.Index(v.Path(), src)
		if len(v.Path())-(l+len(src)) != 0 {
			rest := v.Path()[l+len(src)+1:]
			dstFiles[k].dir = path.Join(dst, rest)
		} else {
			dstFiles[k].dir = dst
		}
	}

	// Adding the src directory itself to the list of files
	srcFirst := FileInfo{name: src, isDir: true}
	srcFiles = slices.Insert(srcFiles, 0, srcFirst)

	// Adding the dest directory itself to the list of files
	dstFirst := FileInfo{name: dst, isDir: true}
	dstFiles = slices.Insert(dstFiles, 0, dstFirst)

	// Moving the files from src to dst
	for k, v := range srcFiles {
		srcJoinPath := path.Join(srcFiles[k].Path(), srcFiles[k].Name())
		dstJoinPath := path.Join(dstFiles[k].Path(), dstFiles[k].Name())
		if v.IsDir() {
			if err := os.Mkdir(dstJoinPath, 0775); err != nil {
				return err
			}
		} else {
			if err := fs.FileRename(srcJoinPath, dstJoinPath); err != nil {
				return err
			}
		}

		// Verification
		dstPathContainer := path.Join(fs.vaultDir, dstFiles[k].dir, dstFiles[k].containerNumber)
		_, err := os.Stat(dstPathContainer)
		if err != nil {
			return fmt.Errorf("unable to locate file %s, error %s", dstJoinPath, err)
		}
	}

	// Removing source dir(s)
	if err := os.RemoveAll(src); err != nil {
		return err
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

// Root returns the root of the vaults directory
func Root() string {
	return vaults
}

// RootData returns the root of the vaultsdata directory
func RootData() string {
	return vaultsData
}

// VaultName returns the name of the current vault
func (fs *FileSystem) VaultName() string {
	i := len(Root())
	return fs.vaultDir[i+1:]
}

// ListVaults returns the vault names
func ListVaults() ([]string, error) {
	dir, err := os.Open(RootData())
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	strList := make([]string, len(files))
	for i, file := range files {
		if file.IsDir() {
			strList[i] = file.Name()
		}
	}
	return strList, nil
}

// Returns true or false
func (fs FileSystem) DirExists(dir string) bool {
	str := path.Join(fs.vaultDir, dir)
	return util.DirExists(str)
}

// Getwd returns the working directory
func (fs FileSystem) Getwd() string {
	str, _ := os.Getwd()
	if len(str) > len(fs.vaultDir) {
		return str[len(fs.vaultDir)+1:]
	} else {
		return fs.vaultDir
	}
}

// AbsNormal returns a "normalized" path, with the offset from the vault directory
func (fs FileSystem) AbsNormal(absolutePath string) (string, error) {
	if len(absolutePath) < 1 {
		return "", errors.New("string is too short")
	}

	if absolutePath[len(absolutePath)-1] == '/' {
		absolutePath = absolutePath[:len(absolutePath)-1]
	}

	var dir string
	if len(absolutePath) == len(fs.vaultDir) { // case when file moved to the root
		dir = "."
	} else if strings.HasPrefix(absolutePath, fs.vaultDir) {
		dir = absolutePath[len(fs.vaultDir)+1:]
	} else { // This only happens when the path is not part of the filesystem...
		return "", fmt.Errorf("wrong path %s", absolutePath)
	}

	return dir, nil
}

// Verifies wether the operation succeeded well. It checks
// the location and also the entry inside the FileList.
func (fs FileSystem) Verify(fl FileList) error {
	pathContainer := path.Join(fs.vaultDir, fl.Path, fl.ContainerNumber)
	joinPath := path.Join(fl.Path, fl.Name)

	_, err := os.Stat(pathContainer)
	if err != nil {
		return fmt.Errorf("unable to locate file %s, error %s", joinPath, err)
	}

	for _, v := range fs.index.fileList {
		if fl.ContainerNumber == v.ContainerNumber {
			if strings.Compare(v.Name, fl.Name) != 0 {
				return fmt.Errorf("fileName comparison mismatch %s != %s", v.Path, fl.Path)
			}
			if strings.Compare(v.PreviousName, fl.PreviousName) != 0 {
				return fmt.Errorf("previousName comparison mismatch %s != %s", v.Path, fl.Path)
			}
			if strings.Compare(v.Path, fl.Path) != 0 {
				return fmt.Errorf("dir comparison mismatch %s != %s", v.Path, fl.Path)
			}
			if strings.Compare(v.PreviousPath, fl.PreviousPath) != 0 {
				return fmt.Errorf("previousDir comparison mismatch %s != %s", v.Path, fl.Path)
			}

			// all tests completed. Returning nil
			return nil
		}
	}

	// Nothing found...
	return fmt.Errorf("container %s not found", fl.ContainerNumber)
}

// Removes the container from disk. Returns nil or an error.
func (fs *FileSystem) FileRemove(containerNumber string) error {
	// Check empty string
	if containerNumber == "" {
		return errors.New("empty container number")
	}

	// Check whether dest is a number
	if !util.IsNumber(containerNumber) {
		return fmt.Errorf("container numbers are a number, not %s", containerNumber)
	}

	// Check whether the item exists in the list
	fl, err := fs.index.ContainerNumberToFileList(containerNumber)
	if err != nil {
		return err
	}

	cmpl := path.Join(fs.vaultDir, fl.Path, fl.ContainerNumber)

	// Check if source directory exists
	if !util.DirExists(cmpl) {
		return fmt.Errorf("container directory %s does not exist", containerNumber)
	}

	// Check if file versions are Checked-Out
	if user := fs.IsLockedItem(containerNumber); user != "" {
		return fmt.Errorf("container %s is checked out by %s", containerNumber, user)
	}

	// Remove from the index
	if err = fs.index.ContainerNumberRemove(containerNumber); err != nil {
		return err
	}

	// Remove the container from disk
	fd := NewFileDirectory(fs, fl)
	if err = fd.DeleteAll(); err != nil {
		return err
	}

	// Verification
	// TODO

	// Log the successful move operation
	log.Printf("Successfully removed container %s", containerNumber)

	return nil
}

// Removes the version of a container from disk. Returns nil or an error.
func (fs *FileSystem) FileRemoveVersion(containerNumber string, version int16) error {
	// check, check, check
	// do, do, do
	// verify

	// Log the successful move operation
	log.Printf("Successfully removed item %d from container %s", version, containerNumber)

	return nil
}

// Removes all files and directories from the vault.
// Be careful. Files can't be recovered.
func RemoveAll(vault string) error {
	// Check wheter the vault has '/'
	if strings.ContainsRune(vault, '/') {
		return errors.New("invalid vault")
	}
	vd := path.Join(Root(), vault)      // vault directory
	vdd := path.Join(RootData(), vault) // vault data directory

	// Check wether there is a vault
	if !util.DirExists(vd) {
		return fmt.Errorf("vault %s does not exist", vault)
	}

	// Setting the Root directoriy to 0777
	if err := os.Chmod(Root(), 0777); err != nil {
		return err
	}

	// Setting the vault directoriy to 0777
	if err := os.Chmod(vd, 0777); err != nil {
		return err
	}

	// Setting all directories of the root data dir to 0777
	if err := os.Chmod(RootData(), 0777); err != nil {
		return err
	}

	// Setting all directories of the vault data dir to 0777
	if err := os.Chmod(vdd, 0777); err != nil {
		return err
	}

	// Setting all directories of the vault to 0777
	err := filepath.WalkDir(vd, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %v", path, err)
		}
		if d.IsDir() {
			if err := os.Chmod(path, 0777); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Remove the vault
	if err := os.RemoveAll(vd); err != nil {
		return err
	}

	// Remove the vault data dir
	if err := os.RemoveAll(vdd); err != nil {
		return err
	}

	// Setting all directories of the root data dir to 0555
	if err := os.Chmod(RootData(), 0555); err != nil {
		return err
	}

	// Setting all directories of the root dir to 0555
	if err := os.Chmod(Root(), 0555); err != nil {
		return err
	}

	return nil
}
