// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package localfs

// The FileDirectory struct deals with file versions.
// Each file that is stored inside the vault has its own version.

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/grd/FreePDM/internal/util"
)

// Some handy file names
const (
	Properties      = "Properties.txt"
	Description     = "Description.txt"
	LongDescription = "LongDescription.txt"
	Ver             = "VER.txt"
)

// File Directory related struct.
type FileDirectory struct {
	fs  *FileSystem
	fl  FileList
	dir string
}

// File versions struct
//
//	field: 'number', an increment
//	field: 'pretty' means a version presentation, such as A.1, 2.5.0, 4.0, 3 or A.
//
// See https://github.com/grd/FreePDM/discussions/93 for a proposal
//
// TODO This functionality hasn't been implemented yet...
//
// For the time being it just reports the string presentation of 'number'
//
//	field: 'date' means the time of a new version with the format "YYYY-MM-DD H:M:S"
type FileVersion struct {
	Number int16
	Pretty string
	Date   string
}

// Initializes the FileDirectory struct. Parameters:
// fs is necessary because of this struct
// containerNumber means the directory in where to put the file structure
// fileNumber means the file number, which is an int64
func NewFileDirectory(fs *FileSystem, fl FileList) FileDirectory {
	return FileDirectory{fs: fs, fl: fl,
		dir: filepath.Join(fs.vaultDir, fl.Path, fl.ContainerNumber)}
}

// Returns the FileList because this field is unexported
func (fd FileDirectory) FileList() FileList {
	return fd.fl
}

// Creates a new directory inside the container directory.
func (fd *FileDirectory) CreateDirectory() error {

	err := os.Mkdir(fd.dir, 0755)
	util.CheckErr(err)
	err = os.Chown(fd.dir, fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)

	// Write version file
	fd.writeInitialVersionFile()

	log.Printf("Created file structure %s\n",
		filepath.Join(fd.fl.ContainerNumber, fd.fl.Name))

	return nil
}

// Imports a file from an external source.
func (fd FileDirectory) ImportNewFile(fname string) error {

	// create a new version string
	new_version := fd.LatestVersion().Number + 1
	version := fmt.Sprint(new_version)
	versionDir := filepath.Join(fd.dir, version)

	fd.increaseVersionNumber(version)

	// create a new version dir
	if err := os.Mkdir(versionDir, 0777); err != nil {
		return err
	}
	if err := os.Chown(versionDir, fd.fs.userUid, fd.fs.vaultUid); err != nil {
		return err
	}

	// create a new file reference text inside FileName.txt
	_, copiedFile := path.Split(fname)
	copiedFile = filepath.Join(versionDir, copiedFile)

	// copy the file inside the new version
	if err := util.CopyFile(fname, copiedFile); err != nil {
		return err
	}
	if err := os.Chown(copiedFile, fd.fs.userUid, fd.fs.vaultUid); err != nil {
		return err
	}

	return nil
}

// Creates a new version from copying the previous version.
func (fd FileDirectory) NewVersion() (*FileVersion, error) {

	// create a new version string
	oldVersion := fd.LatestVersion()
	newVersion := FileVersion{Number: oldVersion.Number + 1, Date: util.Now()}
	newVersion.Pretty = util.I16toa(newVersion.Number)
	versionDir := filepath.Join(fd.dir, newVersion.Pretty)

	fd.increaseVersionNumber(newVersion.Pretty)

	// generate the new file name
	fname := filepath.Join(fd.dir, oldVersion.Pretty, fd.fl.Name)

	// create a new version dir
	if err := os.Mkdir(versionDir, 0755); err != nil {
		return nil, err
	}
	if err := os.Chown(versionDir, fd.fs.userUid, fd.fs.vaultUid); err != nil {
		return nil, err
	}

	// create a new file reference text inside FileName.txt
	_, copiedFile := path.Split(fname)
	copiedFile = filepath.Join(versionDir, copiedFile)

	// copy the file inside the new version
	if err := util.CopyFile(fname, copiedFile); err != nil {
		return nil, err
	}
	if err := os.Chown(copiedFile, fd.fs.userUid, fd.fs.vaultUid); err != nil {
		return nil, err
	}

	return &newVersion, nil
}

// Stores the description and long description text.
func (fd FileDirectory) StoreData(version FileVersion, descr, longDescr string) {

	// create a version directory
	versionDir := filepath.Join(fd.dir, version.Pretty)

	if !util.DirExists(versionDir) {
		log.Fatalf("Directory %s doesn't exist.", versionDir)
	}

	// create a new description text
	if len(descr) > 0 {
		descriptionFile := filepath.Join(versionDir, Description)
		dsc := []byte(descr)

		err := os.WriteFile(descriptionFile, dsc, 0444)
		util.CheckErr(err)

		err = os.Chown(descriptionFile, fd.fs.userUid, fd.fs.vaultUid)
		util.CheckErr(err)
	}

	// create a new long description text
	if len(longDescr) > 0 {
		longDescriptionFile := filepath.Join(versionDir, LongDescription)
		buf2 := []byte(longDescr)

		err := os.WriteFile(longDescriptionFile, buf2, 0444)
		util.CheckErr(err)

		err = os.Chown(longDescriptionFile, fd.fs.userUid, fd.fs.vaultUid)
		util.CheckErr(err)
	}
}

// Returns the file properties of the latest version
func (fd FileDirectory) LatestProperties() []FileProperties {
	release := fd.LatestVersion()
	return fd.Properties(release)
}

// Returns the file properties of the specific version
func (fd FileDirectory) Properties(version FileVersion) []FileProperties {
	buf, err := os.ReadFile(filepath.Join(version.Pretty, Properties))
	util.CheckErr(err)
	str := string(buf)
	// check for latest '\n'
	if str[len(str)-1] == '\n' {
		str = str[:len(str)-2]
	}
	splitStr := strings.Split(str, "\n")
	props := make([]FileProperties, len(splitStr))
	for _, v := range splitStr {
		kv := strings.Split(v, " = ")
		props = append(props, FileProperties{Key: kv[0], Value: kv[1]})
	}
	return props
}

// Sets the file properties of the latest version
func (fd FileDirectory) SetLatestProperties(props []FileProperties) {
	release := fd.LatestVersion()
	fd.SetProperties(release, props)
}

// Sets the file properties of the specific version
func (fd FileDirectory) SetProperties(version FileVersion, props []FileProperties) {
	buf := make([]byte, len(props)*20)
	for _, v := range props {
		str := []byte(fmt.Sprintf("%s = %s\n", v.Key, v.Value))
		buf = append(buf, str...)
	}
	err := os.WriteFile(filepath.Join(version.Pretty, Properties), buf, 0644)
	util.CheckErr(err)
	err = os.Chown(filepath.Join(version.Pretty, Properties), fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)
}

// Returns the latest version.
func (fd *FileDirectory) LatestVersion() FileVersion {

	versions, err := fd.AllFileVersions()
	if err != nil {
		log.Fatalf("Error reading file %s, version %v", fd.fl.Name, err)
	}

	if len(versions) == 1 {
		return versions[0]
	}

	return versions[len(versions)-1]
}

// Returns all file versions name from file or an error.
func (fd *FileDirectory) AllFileVersions() ([]FileVersion, error) {

	version := filepath.Join(fd.dir, Ver)

	file, err := os.Open(version)
	util.CheckErr(err)
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ':'

	records, err := r.ReadAll()
	util.CheckErr(err)

	if len(records) == 0 {
		return nil, fmt.Errorf("file %s is empty",
			filepath.Join(fd.dir, Ver))
	}

	if len(records) == 1 {
		fv_slice := append([]FileVersion{}, FileVersion{Number: -1, Pretty: "-1", Date: util.Now()})
		return fv_slice, nil
	}

	records = records[1:]

	ret := make([]FileVersion, len(records))

	for i, record := range records {
		fmt.Sscanf(record[0], "%d", &ret[i].Number)
		ret[i].Pretty = record[1]
		ret[i].Date = record[2]
	}

	return ret, nil
}

// Delete all versions and files of the container from disk.
// Returns nil or an error.
func (fd *FileDirectory) DeleteAll() error {
	// Get all file versions
	versions, err := fd.AllFileVersions()
	if err != nil {
		return err
	}

	// Put 0777 for the main directory
	if err := os.Chmod(fd.dir, 0777); err != nil {
		return err
	}

	// Put 0777 for all file version directories
	for _, item := range versions {
		s := filepath.Join(fd.dir, fmt.Sprint(item.Number))
		if err := os.Chmod(s, 0777); err != nil {
			return err
		}
	}

	// Remove All
	if err := os.RemoveAll(fd.dir); err != nil {
		return err
	}

	return nil
}

// TODO: How?
// Delete one version of the container.
// Returns nil or an error.
func (fd *FileDirectory) DeleteVersion(item int16) error {
	// Check
	// Do
	return nil
}

// Opens the latest item for editing the SMB mount.
// This "Checks Out" the item.
func (fd *FileDirectory) OpenLatestsVersion() {
	ver := fd.LatestVersion()
	fd.OpenItemVersion(ver)
}

// Closes the latest version for editing.
func (fd *FileDirectory) CloseLatestsVersion() {
	ver := fd.LatestVersion()
	fd.CloseItemVersion(ver)
}

// Opens the item number for editing.
// This "Checks Out" the item.
func (fd *FileDirectory) OpenItemVersion(version FileVersion) {

	dirVersion := filepath.Join(fd.dir, version.Pretty)

	err := os.Chown(dirVersion, fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)

	// Directory mode 0700 means that only that guy can edit the file.
	err = os.Chmod(dirVersion, 0700)
	util.CheckErr(err)

	// // And that guy has filemode 0644 for the file itself.
	// file := filepath.Join(dirVersion, fd.fl.Name)
	// err = os.Chmod(file, 0644)
	// util.CheckErr(err)
}

// Closes item number for editing.
func (fd *FileDirectory) CloseItemVersion(version FileVersion) {

	dirVersion := filepath.Join(fd.dir, version.Pretty)

	// Filemode 0555 means that the directory is read only for anyone.
	err := os.Chown(dirVersion, fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)
	err = os.Chmod(dirVersion, 0555)
	util.CheckErr(err)

	// // And the file can't be edited anymore with filemode 0444.
	// file := filepath.Join(dirVersion, fd.fl.Name)
	// err = os.Chmod(file, 0444)
	// util.CheckErr(err)
}

func (fd *FileDirectory) writeInitialVersionFile() {

	ver := filepath.Join(fd.dir, Ver)

	records := [][]string{{"Version", "Pretty", "Date"}}

	file, err := os.OpenFile(ver, os.O_WRONLY|os.O_CREATE, 0644)
	util.CheckErr(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ':'

	err = writer.WriteAll(records) // calls Flush internally
	util.CheckErr(err)

	os.Chown(ver, fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)
}

// Increase the version number
func (fd *FileDirectory) increaseVersionNumber(version string) {

	date := util.Now()

	ver := filepath.Join(fd.dir, Ver)

	record := []string{version, version, date}

	err := os.Chmod(ver, 0644)
	util.CheckErr(err)

	file, err := os.OpenFile(ver, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	util.CheckErr(err)

	writer := csv.NewWriter(file)
	writer.Comma = ':'

	writer.Write(record)

	writer.Flush()
	file.Close()

	err = os.Chown(ver, fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)
}

// Renames the filename. Returns an error when unsuccessful.
func (fd *FileDirectory) fileRename(src, dst string) error {
	versions, err := fd.AllFileVersions()
	if err != nil {
		return err
	}

	// Rename all versioned files
	for _, version := range versions {
		if err = fd.WithTempPermissions(version.Pretty, func(subDir string) error {
			return os.Rename(filepath.Join(subDir, src), filepath.Join(subDir, dst))
		}); err != nil {
			return err
		}
	}

	return nil
}

// Function that changes permissions, performs an operation, and restores permissions
func (fd FileDirectory) WithTempPermissions(version string, operation func(subDir string) error) error {
	// Mutex to ensure only one operation modifies permissions at a time
	var permMutex sync.Mutex

	// Acquire exclusive access
	permMutex.Lock()
	defer permMutex.Unlock()

	// joining fd.dir with the version number
	versionDir := filepath.Join(fd.dir, version)

	// Step 1: Set permissions to 0777 (allow full access)
	if err := os.Chmod(fd.dir, 0777); err != nil {
		return fmt.Errorf("error setting mode 0777 on %s: %v", fd.dir, err)
	}
	if err := os.Chmod(versionDir, 0777); err != nil {
		return fmt.Errorf("error setting mode 0777 on %s: %v", versionDir, err)
	}

	// Step 2: Perform the requested operation
	if err := operation(versionDir); err != nil {
		return fmt.Errorf("error executing operation: %v", err)
	}

	// Step 3: Restore permissions to 0555 (read-only access)
	if err := os.Chmod(versionDir, 0555); err != nil {
		return fmt.Errorf("error setting mode 0555 on %s: %v", versionDir, err)
	}
	if err := os.Chmod(fd.dir, 0555); err != nil {
		return fmt.Errorf("error setting mode 0555 on %s: %v", fd.dir, err)
	}

	return nil
}

// func main() {
// 	rootDir := "/path/to/directory"
// 	subDir := rootDir + "/subdir"

// 	// Call the function with an operation (e.g., creating a file)
// 	err := withTempPermissions(rootDir, subDir, func(subDir string) error {
// 		filePath := subDir + "/new_file.txt"
// 		file, err := os.Create(filePath)
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()
// 		fmt.Println("File created:", filePath)
// 		return nil
// 	})

// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	} else {
// 		fmt.Println("Operation completed successfully!")
// 	}
// }
