// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	ex "github.com/grd/FreePDM/src/extras"
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
	fs         *FileSystem
	dir        string
	fileNumber int64
}

func InitFileDirectory(fsm *FileSystem, dir string, fileNumber int64) FileDirectory {
	return FileDirectory{fs: fsm, dir: dir, fileNumber: fileNumber}
}

// Creates a new directory inside the current working directory.
func (self *FileDirectory) NewDirectory() FileDirectory {

	dirName := self.dir

	err := os.Mkdir(dirName, 0755)
	ex.CheckErr(err)
	err = os.Chown(dirName, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)

	// Write version file
	self.writeVersionFile()

	log.Printf("Created file structure %s\n",
		path.Join(self.dir, self.fs.index.FileName(self.fileNumber)))

	return *self
}

// Imports a file from an external source.
func (self FileDirectory) ImportNewFile(fname string) FileDirectory {

	// create a new version string

	new_version := self.LatestVersion() + 1

	version := fmt.Sprintf("VER%03d", new_version)

	versionDir := path.Join(self.dir, version)

	self.addVersion(version, ex.Today())

	// create a new version dir

	err := os.Mkdir(versionDir, 0777)
	ex.CheckErr(err)
	err = os.Chown(versionDir, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)

	// create a new file reference text inside FileName.txt

	_, copiedFile := path.Split(fname)

	copiedFile = path.Join(versionDir, copiedFile)

	// copy the file inside the new version

	ex.CopyFile(fname, copiedFile)

	err = os.Chown(copiedFile, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)

	return self
}

// Creates a new version from copying the previous version.
func (self FileDirectory) NewVersion() int16 {

	// create a new version string

	new_version := self.LatestVersion() + 1
	version := fmt.Sprintf("VER%03d", new_version)
	versionDir := path.Join(self.dir, version)

	self.addVersion(version, ex.Today())

	// generate the new file name

	old_version := new_version - 1
	old_directory := fmt.Sprintf("VER%03d", old_version)

	filename := self.fs.index.FileName(self.fileNumber)

	fname := path.Join(self.dir, old_directory, filename)

	// create a new version dir

	err := os.Mkdir(versionDir, 0755)
	ex.CheckErr(err)
	err = os.Chown(versionDir, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)

	// create a new file reference text inside FileName.txt

	_, copiedFile := path.Split(fname)

	copiedFile = path.Join(versionDir, copiedFile)

	// copy the file inside the new version

	ex.CopyFile(fname, copiedFile)

	err = os.Chown(copiedFile, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)

	return new_version
}

// Stores the description and long description text.
func (self FileDirectory) StoreData(ver int16, descr, longDescr string) {

	// create a version directory

	version := fmt.Sprintf("VER%03d", ver)
	versionDir := path.Join(self.dir, version)

	if ex.DirExists(versionDir) == false {
		log.Fatalf("Directory %s doesn't exist.", versionDir)
	}

	// create a new description text

	if len(descr) > 0 {
		descriptionFile := path.Join(versionDir, Description)
		dsc := []byte(descr)

		err := os.WriteFile(descriptionFile, dsc, 0444)
		ex.CheckErr(err)

		err = os.Chown(descriptionFile, self.fs.userUid, self.fs.vaultUid)
		ex.CheckErr(err)
	}

	// create a new long description text

	if len(longDescr) > 0 {
		longDescriptionFile := path.Join(versionDir, LongDescription)
		buf2 := []byte(longDescr)

		err := os.WriteFile(longDescriptionFile, buf2, 0444)
		ex.CheckErr(err)

		err = os.Chown(longDescriptionFile, self.fs.userUid, self.fs.vaultUid)
		ex.CheckErr(err)
	}
}

// The number of the directory
func (self FileDirectory) FileNumber() int64 {
	var num int64
	fmt.Sscanf(self.dir, "%d", &num)
	return num
}

// Returns the file properties of the latest version
func (self FileDirectory) LatestProperties() []FileProperties {
	release := self.LatestVersion()
	return self.Properties(release)
}

// Returns the file properties of the specific version
func (self FileDirectory) Properties(nr int16) []FileProperties {
	version := fmt.Sprintf("VER%03d", nr)
	buf, err := os.ReadFile(path.Join(version, Properties))
	ex.CheckErr(err)
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
func (self FileDirectory) SetLatestProperties(props []FileProperties) {
	release := self.LatestVersion()
	self.SetProperties(release, props)
}

// Sets the file properties of the specific version
func (self FileDirectory) SetProperties(nr int16, props []FileProperties) {
	version := fmt.Sprintf("VER%03d", nr)
	buf := make([]byte, len(props)*20)
	for _, v := range props {
		str := []byte(fmt.Sprintf("%s = %s\n", v.Key, v.Value))
		buf = append(buf, str...)
	}
	err := os.WriteFile(path.Join(version, Properties), buf, 0644)
	ex.CheckErr(err)
	err = os.Chown(path.Join(version, Properties), self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)
}

// Returns the latest version.
func (self *FileDirectory) LatestVersion() int16 {

	versions, err := self.AllFileVersions()
	if err != nil {
		return -1
	}

	return versions[len(versions)-1]
}

// Returns all file versions name from file or an error.
func (self *FileDirectory) AllFileVersions() ([]int16, error) {

	version := path.Join(self.dir, Ver)

	file, err := os.Open(version)
	ex.CheckErr(err)
	defer file.Close()

	r := csv.NewReader(file)

	records, err := r.ReadAll()

	ex.CheckErr(err)

	if len(records) <= 1 {
		return []int16{-1}, fmt.Errorf("File %s is empty or has only one record.\n",
			path.Join(self.fs.currentWorkingDir, self.dir, Ver))
	}

	records = records[1:]

	ret := make([]int16, len(records))

	for i, record := range records {

		fmt.Sscanf(record[0], "VER%d", &ret[i])
	}

	return ret, nil
}

// Delete one version, not the file on its own.
// The stored file should be put in the archive, but not deleted AND
// also not visible.
func (self *FileDirectory) DeleteVersion(item int) {
	// TODO: How? By giving the directory file mode 0700.
	// And to set a field inside the database.
	// User Admin should be able to undo this action.
}

// Restore one version.
// Also restore the item from the filesystem indexes.
func (self *FileDirectory) Restoreversion(item int) {
	// TODO: Implement this. Undo the function DeleteVersion().
}

// Opens the latest item for editing the SMB mount.
// This "Checkes Out" the item.
func (self *FileDirectory) OpenLatestsVersion() {
	ver := self.LatestVersion()
	self.OpenItemVersion(ver)
}

// Closes the latest version for editing.
func (self *FileDirectory) CloseLatestsVersion() {
	ver := self.LatestVersion()
	self.CloseItemVersion(ver)
}

// Opens the item number for editing the SMB mount.
// This "Checkes Out" the item.
func (self *FileDirectory) OpenItemVersion(ver int16) {

	version := fmt.Sprintf("VER%03d", ver)
	dirVersion := path.Join(self.dir, version)

	err := os.Chown(dirVersion, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)

	// Filemode 0700 means that only that guy can edit the file.

	err = os.Chmod(dirVersion, 0700)
	ex.CheckErr(err)

	// And that guy has filemode 644 for the file itself.

	base := path.Base(self.dir)
	num := ex.Atoi64(base)
	file := path.Join(dirVersion, self.fs.index.FileName(num))

	err = os.Chmod(file, 0644)
	ex.CheckErr(err)
}

func (self *FileDirectory) writeVersionFile() {

	ver := path.Join(self.dir, Ver)

	records := [][]string{{"Version", "Date"}}

	file, err := os.OpenFile(ver, os.O_WRONLY|os.O_CREATE, 0644)
	ex.CheckErr(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(records) // calls Flush internally
	ex.CheckErr(err)

	os.Chown(ver, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)
}

// Add the version number and the date
func (self *FileDirectory) addVersion(version, date string) {

	ver := path.Join(self.dir, Ver)

	record := []string{version, date}

	err := os.Chmod(ver, 0644)
	ex.CheckErr(err)

	file, err := os.OpenFile(ver, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	ex.CheckErr(err)

	writer := csv.NewWriter(file)
	writer.Write(record)

	writer.Flush()
	file.Close()

	err = os.Chown(ver, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)
	err = os.Chmod(ver, 0444)
	ex.CheckErr(err)
}

// Closes item number for editing.
func (self *FileDirectory) CloseItemVersion(ver int16) {

	version := fmt.Sprintf("VER%03d", ver)
	dirVersion := path.Join(self.dir, version)

	// Filemode 0755 means that the directory is open for anyone.

	err := os.Chown(dirVersion, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)
	err = os.Chmod(dirVersion, 0755)
	ex.CheckErr(err)

	// And the file can't be edited anymore with filemode 0444.

	base := path.Base(self.dir)
	num := ex.Atoi64(base)
	file := path.Join(dirVersion, self.fs.index.FileName(num))

	err = os.Chmod(file, 0444)
	ex.CheckErr(err)
}

// Renames the filename. Returns an error when not succeed.
func (self *FileDirectory) fileRename(src, dest string) error {

	versions, err := self.AllFileVersions()
	ex.CheckErr(err)

	for _, version := range versions {
		ver := fmt.Sprintf("VER%03d", version)

		// Rename

		err = os.Rename(path.Join(self.dir, ver, src), path.Join(self.dir, ver, dest))
		ex.CheckErr(err)
	}

	return nil
}
