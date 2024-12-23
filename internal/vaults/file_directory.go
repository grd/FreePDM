// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package vaults

// The FileDirectory struct deals with file versions.
// Each file that is stored inside the vault has its own version.

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

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
// TODO This functionality hasn't been implemented yet...
// For the time being it just reports the string presentation of 'number'
//
//	field: 'date' means the time of a new version with the format "YYYY-MM-DD H:M:S"
type FileVersion struct {
	Number int16
	Pretty string
	Date   string
}

// Initializes the FileDirectory struct. Parameters:
// fsm is necessary because of this struct
// containerNumber means the directory in where to put the file structure
// fileNumber means the file number, which is an int64
func NewFileDirectory(fsm *FileSystem, fl FileList) FileDirectory {
	return FileDirectory{fs: fsm, fl: fl,
		dir: path.Join(fsm.vaultDir, fl.dir, fl.containerNumber)}
}

// Returns the FileList because this field is unexported
func (fd FileDirectory) FileList() FileList {
	return fd.fl
}

// Creates a new directory inside the current working directory.
func (fd *FileDirectory) CreateDirectory() error {

	err := os.Mkdir(fd.dir, 0755)
	util.CheckErr(err)
	err = os.Chown(fd.dir, fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)

	// Write version file
	fd.writeInitialVersionFile()

	log.Printf("Created file structure %s\n",
		path.Join(fd.fl.containerNumber, fd.fl.fileName))

	return nil
}

// Imports a file from an external source.
func (fd FileDirectory) ImportNewFile(fname string) error {

	// create a new version string

	new_version := fd.LatestVersion().Number + 1

	version := fmt.Sprintf("%d", new_version)

	versionDir := path.Join(fd.dir, version)

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

	copiedFile = path.Join(versionDir, copiedFile)

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
func (fd FileDirectory) NewVersion() FileVersion {

	// create a new version string

	oldVersion := fd.LatestVersion()
	newVersion := FileVersion{Number: oldVersion.Number + 1, Date: util.Now()}
	newVersion.Pretty = util.I16toa(newVersion.Number)
	versionDir := path.Join(fd.dir, newVersion.Pretty)

	fd.increaseVersionNumber(newVersion.Pretty)

	// generate the new file name
	fname := path.Join(fd.dir, oldVersion.Pretty, fd.fl.fileName)

	// create a new version dir

	err := os.Mkdir(versionDir, 0755)
	util.CheckErr(err)
	err = os.Chown(versionDir, fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)

	// create a new file reference text inside FileName.txt

	_, copiedFile := path.Split(fname)

	copiedFile = path.Join(versionDir, copiedFile)

	// copy the file inside the new version

	util.CopyFile(fname, copiedFile)

	err = os.Chown(copiedFile, fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)

	return newVersion
}

// Stores the description and long description text.
func (fd FileDirectory) StoreData(version FileVersion, descr, longDescr string) {

	// create a version directory

	versionDir := path.Join(fd.dir, version.Pretty)

	if !util.DirExists(versionDir) {
		log.Fatalf("Directory %s doesn't exist.", versionDir)
	}

	// create a new description text

	if len(descr) > 0 {
		descriptionFile := path.Join(versionDir, Description)
		dsc := []byte(descr)

		err := os.WriteFile(descriptionFile, dsc, 0444)
		util.CheckErr(err)

		err = os.Chown(descriptionFile, fd.fs.userUid, fd.fs.vaultUid)
		util.CheckErr(err)
	}

	// create a new long description text

	if len(longDescr) > 0 {
		longDescriptionFile := path.Join(versionDir, LongDescription)
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
	buf, err := os.ReadFile(path.Join(version.Pretty, Properties))
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
	err := os.WriteFile(path.Join(version.Pretty, Properties), buf, 0644)
	util.CheckErr(err)
	err = os.Chown(path.Join(version.Pretty, Properties), fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)
}

// Returns the latest version.
func (fd *FileDirectory) LatestVersion() FileVersion {

	versions, err := fd.AllFileVersions()
	if err != nil {
		log.Fatalf("Error reading file %s, version %v", fd.fl.fileName, err)
	}

	if len(versions) == 1 {
		return versions[0]
	}

	return versions[len(versions)-1]
}

// Returns all file versions name from file or an error.
func (fd *FileDirectory) AllFileVersions() ([]FileVersion, error) {

	version := path.Join(fd.dir, Ver)

	file, err := os.Open(version)
	util.CheckErr(err)
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ':'

	records, err := r.ReadAll()
	util.CheckErr(err)

	if len(records) == 0 {
		return nil, fmt.Errorf("file %s is empty",
			path.Join(fd.dir, Ver))
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

// Delete one version, not the file on its own.
// The stored file should be put in the archive, but not deleted AND
// also not visible.
func (fd *FileDirectory) DeleteVersion(item int) {
	// TODO: How? By giving the directory file mode 0700.
	// And to set a field inside the database.
	// User Admin should be able to undo this action.
}

// Restore one version.
// Also restore the item from the filesystem indexes.
func (fd *FileDirectory) Restoreversion(item int) {
	// TODO: Implement this. Undo the function DeleteVersion().
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

	dirVersion := path.Join(fd.dir, version.Pretty)

	err := os.Chown(dirVersion, fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)

	// Filemode 0700 means that only that guy can edit the file.
	err = os.Chmod(dirVersion, 0700)
	util.CheckErr(err)

	// And that guy has filemode 0644 for the file itself.
	file := path.Join(dirVersion, fd.fl.fileName)
	err = os.Chmod(file, 0644)
	util.CheckErr(err)
}

// Closes item number for editing.
func (fd *FileDirectory) CloseItemVersion(version FileVersion) {

	dirVersion := path.Join(fd.dir, version.Pretty)

	// Filemode 0755 means that the directory is open for anyone.
	err := os.Chown(dirVersion, fd.fs.userUid, fd.fs.vaultUid)
	util.CheckErr(err)
	err = os.Chmod(dirVersion, 0755)
	util.CheckErr(err)

	// And the file can't be edited anymore with filemode 0444.
	file := path.Join(dirVersion, fd.fl.fileName)
	err = os.Chmod(file, 0444)
	util.CheckErr(err)
}

func (fd *FileDirectory) writeInitialVersionFile() {

	ver := path.Join(fd.dir, Ver)

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

	ver := path.Join(fd.dir, Ver)

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
	err = os.Chmod(ver, 0444)
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
		if err = os.Rename(path.Join(fd.dir, version.Pretty, src), path.Join(fd.dir, version.Pretty, dst)); err != nil {
			return err
		}
	}

	return nil
}