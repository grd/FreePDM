// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	ex "github.com/grd/FreePDM/src/extras"
)

// Some handy file names
const (
	FileName        = "FileName.txt"
	Properties      = "Properties.txt"
	Description     = "Description.txt"
	LongDescription = "LongDescription.txt"
	Ver             = "VER.txt"
)

// File Directory related struct.
type FileDirectory struct {
	fs  *FileSystem
	dir string
}

func InitFileDirectory(fsm *FileSystem, dir string) FileDirectory {
	return FileDirectory{fs: fsm, dir: dir}
}

// Creates a new directory inside the current working directory.
func (self *FileDirectory) NewDirectory() FileDirectory {

	dirName := self.dir

	err := os.Mkdir(dirName, 0755)
	ex.CheckErr(err)
	err = os.Chown(dirName, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)

	ver := path.Join(dirName, Ver)
	verFile, err := os.Create(ver)
	ex.CheckErr(err)
	verFile.Close()
	os.Chown(ver, self.fs.userUid, self.fs.vaultUid)

	log.Printf("Created file structure %s into %s\n", self.fs.index.FileNameOfString(dirName), self.fs.currentWorkingDir)

	return *self
}

// Imports a file from an external source.
func (self FileDirectory) ImportNewFile(fname string) FileDirectory {

	// going into the directory...

	err := os.Chdir(self.dir)
	ex.CheckErr(err)

	// create a new version string

	err = os.Chmod(Ver, 0644)
	ex.CheckErr(err)

	new_version := self.LatestVersion() + 1
	version := fmt.Sprintf("VER%03d", new_version)
	to_day := ex.Today()
	buf := fmt.Sprintf("%s\n%s\n", version, to_day)

	err = os.Chmod(Ver, 0644)
	ex.CheckErr(err)

	f, err := os.OpenFile(Ver, os.O_APPEND|os.O_WRONLY, 0644)
	ex.CheckErr(err)

	_, err = f.WriteString(buf)
	ex.CheckErr(err)
	f.Close()

	err = os.Chown(Ver, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)
	err = os.Chmod(Ver, 0444)
	ex.CheckErr(err)

	// create a new version dir

	err = os.Mkdir(version, 0777)
	ex.CheckErr(err)
	err = os.Chown(version, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)
	err = os.Chdir(version)
	ex.CheckErr(err)

	// create a new file reference text inside FileName.txt

	_, copiedFile := path.Split(fname)

	cf := []byte(copiedFile)
	err = os.WriteFile(FileName, cf, 0444)
	ex.CheckErr(err)

	err = os.Chown(FileName, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)

	// copy the file inside the new version

	ex.CopyFile(fname, copiedFile)

	err = os.Chown(copiedFile, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)

	// Return to original directory

	err = os.Chdir("../..")
	ex.CheckErr(err)

	return self
}

// Creates a new version from copying the previous version.
func (self FileDirectory) NewVersion() int16 {

	// going into the directory...

	err := os.Chdir(self.dir)
	ex.CheckErr(err)

	// create a new version string

	new_version := self.LatestVersion() + 1
	version := fmt.Sprintf("VER%03d", new_version)
	to_day := ex.Today()
	str := fmt.Sprintf("%s\n%s\n", version, to_day)

	err = os.Chmod(Ver, 0644)
	ex.CheckErr(err)

	f, err := os.OpenFile(Ver, os.O_APPEND|os.O_WRONLY, 0644)
	ex.CheckErr(err)

	_, err = f.WriteString(str)
	ex.CheckErr(err)
	f.Close()

	err = os.Chown(Ver, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)
	err = os.Chmod(Ver, 0444)
	ex.CheckErr(err)

	// generate the new file name

	old_version := new_version - 1
	old_directory := fmt.Sprintf("VER%03d", old_version)
	buf, err := os.ReadFile(path.Join(old_directory, FileName))
	fname := path.Join(self.dir, old_directory, string(buf))

	// create a new version dir

	err = os.Mkdir(version, 0755)
	ex.CheckErr(err)
	err = os.Chown(version, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)
	err = os.Chdir(version)
	ex.CheckErr(err)

	// create a new file reference text inside FileName.txt

	_, copiedFile := path.Split(fname)

	cf := []byte(copiedFile)
	err = os.WriteFile(FileName, cf, 0644)
	ex.CheckErr(err)

	err = os.Chown(FileName, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)

	// copy the file inside the new version

	ex.CopyFile(fname, copiedFile)

	err = os.Chown(copiedFile, self.fs.userUid, self.fs.vaultUid)
	ex.CheckErr(err)

	// Return to original directory

	err = os.Chdir("../..")
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

// String returns the file name of the latest version
func (self FileDirectory) String() string {
	return self.LatestFileName()
}

// Returns the file name of the specific version
func (self FileDirectory) FileName(nr int) string {
	version := fmt.Sprintf("VER%03d", nr)
	file, err := os.ReadFile(path.Join(version, FileName))
	ex.CheckErr(err)

	return string(file)
}

// Returns the file name of the latest version
func (self FileDirectory) LatestFileName() string {
	release := self.LatestVersion()
	version := fmt.Sprintf("VER%03d", release)
	file, err := os.ReadFile(path.Join(version, FileName))
	ex.CheckErr(err)

	return string(file)
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
func (self FileDirectory) LatestVersion() (ret int16) {

	versionFile := Ver // path.Join(self.dir, Ver)

	// check wether the current directory ends with "self.dir"

	wd, _ := os.Getwd()
	if strings.HasSuffix(wd, self.dir) == false {
		versionFile = path.Join(self.dir, Ver)
	}

	// check wether file exist

	if ex.FileExists(versionFile) == false {
		wd, _ := os.Getwd()
		log.Fatalf("File %s doesn't exist inside %s", versionFile, wd)
	}

	err := os.Chmod(versionFile, 0644)
	ex.CheckErr(err)

	buf, err := os.ReadFile(versionFile)
	ex.CheckErr(err)

	if len(buf) > 3 {

		// check for latest '\n'
		if buf[len(buf)-1] == '\n' {
			buf = buf[:len(buf)-2]
		}

		var num int16
		lines := bytes.Split(buf, []byte{'\n'})
		nr := len(lines) - 2
		fmt.Sscanf(string(lines[nr]), "VER%d", &num)
		ret = num
	} else {
		ret = -1
	}

	err = os.Chmod(versionFile, 0444)
	ex.CheckErr(err)

	return ret
}

// Returns all file versions name from file or an error.
func (self FileDirectory) AllFileVersions() ([]int16, error) {
	buf, err := os.ReadFile(path.Join(self.dir, Ver))
	ex.CheckErr(err)

	if len(buf) == 0 {
		return nil, fmt.Errorf("%s, \"%s\" is empty.", self.dir, Ver)
	}

	// check for latest '\n'
	if buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-2]
	}

	lines := bytes.Split(buf, []byte{'\n'})
	ret := make([]int16, len(lines))
	for nr, line := range lines {
		if nr%2 == 0 {
			var num int16
			fmt.Sscanf(string(line), "VER%d", &num)
			if num >= 0 { // Negative numbers are archived.
				ret = append(ret, num)
			}
		}
	}

	return ret, nil
}

// Release is a milestone. It also generates data.
func (self *FileDirectory) Release(nr int) {
	// TODO: Implement this. But TBH I don't know how this releases data and also not what data.
	// Specs:
	// - User Name
	// - Date
	// - Release description name of action. Short one and a long one.
	// - ORM: Update all values.
	// - The "run script" I don't know yet... and I don't know whether it is neccessary.
}

// Dealing with the OOPS factor.
func UnRelease(nr int) {
	// TODO: Implement this.
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

// Renames the filename. It returns an error when not succeed.
func (self *FileDirectory) fileRename(src, dest string) error {
	versions, err := self.AllFileVersions()
	ex.CheckErr(err)
	for version := range versions {
		ver := fmt.Sprintf("VER%03d", version)
		buf := []byte(dest)
		fname := path.Join(ver, FileName)
		// FileName with new name
		err := os.WriteFile(fname, buf, 0644)
		ex.CheckErr(err)
		err = os.Chown(fname, self.fs.userUid, self.fs.vaultUid)
		ex.CheckErr(err)
		// Rename
		err = os.Rename(path.Join(ver, src), path.Join(ver, dest))
		ex.CheckErr(err)
		err = os.Chown(ver, self.fs.userUid, self.fs.vaultUid)
		ex.CheckErr(err)
	}
	return nil
}
