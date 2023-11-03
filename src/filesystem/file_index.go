// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path"

	ex "github.com/grd/FreePDM/src/utils"
)

// The FileList is a struct with five fields:
// * The index number.
// * The file name.
// * The previous name of the file, including the date when the file was renamed.
// * The dir name of the storage. This is an offset of "mainPdmDir".
// * The previous directory name, including the date when the file was moves.
type FileList struct {
	index        int64
	file         string
	previousFile string
	dir          string
	previousDir  string
}

// File Index Files in the root
type FileIndex struct {
	vaultDir string
	userUid  int
	vaultUid int

	fileListCsv string
	fileList    []FileList

	indexNumberTxt string
	indexNumber    int64
}

func InitFileIndex(vaultDir string, user_uid, vault_uid int) (ret FileIndex) {

	ret.vaultDir = vaultDir
	ret.userUid = user_uid
	ret.vaultUid = vault_uid

	// check wether the critical directory exists.

	ex.CriticalDirExist(vaultDir)

	ret.fileListCsv = path.Join(ret.vaultDir, "FileList.csv")
	ret.indexNumberTxt = path.Join(ret.vaultDir, "IndexNumber.txt")

	// check wether the critical files exists.

	ex.CriticalFileExist(ret.fileListCsv)
	ex.CriticalFileExist(ret.indexNumberTxt)

	ret.getIndexNumber()

	ret.fileList = make([]FileList, 0, ret.indexNumber)

	ret.Read() // read the indexes

	return ret
}

// Reads the values from "FileList.txt"
func (self *FileIndex) Read() {

	self.fileList = nil

	var fl FileList

	records, err := self.readCsv()
	ex.CheckErr(err)

	for _, record := range records {

		fl.index = ex.Atoi64(record[0])
		fl.file = record[1]
		fl.previousFile = record[2]
		fl.dir = record[3]
		fl.previousDir = record[4]

		self.fileList = append(self.fileList, fl)
	}
}

func (self FileIndex) readCsv() ([][]string, error) {

	buf, err := os.ReadFile(self.fileListCsv)
	ex.CheckErr(err)

	r := csv.NewReader(bytes.NewBuffer(buf))
	r.Comma = ':'

	records, err := r.ReadAll()
	if err != nil {
		return [][]string{}, fmt.Errorf("Error reading file %s with error: %v", self.fileListCsv, err)
	}

	if len(records) == 0 {
		return [][]string{}, fmt.Errorf("Error reading file %s with error: No header included.", self.fileListCsv)
	}

	records = records[1:]

	return records, nil
}

// Writes the values to "FileList.csv"
func (self *FileIndex) Write() {

	records := [][]string{
		{"Index", "FileName", "PreviousFile", "Dir", "PreviousDir"},
	}

	for _, item := range self.fileList {

		records = append(records, []string{
			ex.I64toa(item.index),
			item.file,
			item.previousFile,
			item.dir,
			item.previousDir})
	}

	var buf []byte
	buffer := bytes.NewBuffer(buf)

	writer := csv.NewWriter(buffer)
	writer.Comma = ':'

	err := writer.WriteAll(records) // calls Flush internally
	ex.CheckErr(err)

	err = writer.Error()
	ex.CheckErr(err)

	err = os.WriteFile(self.fileListCsv, buffer.Bytes(), 0644)
	ex.CheckErr(err)

	err = os.Chown(self.fileListCsv, self.userUid, self.vaultUid)
	ex.CheckErr(err)
}

// Reads the index number and stores it.
func (self *FileIndex) getIndexNumber() {

	buf, err := os.ReadFile(self.indexNumberTxt)
	ex.CheckErr(err)

	_, err = fmt.Sscanf(string(buf), "%d", &self.indexNumber)
	ex.CheckErr(err)
}

// Returns the complete directory name of a file, or an error when not found.
func (self *FileIndex) Dir(fileName string) (string, error) {
	for _, v := range self.fileList {
		if fileName == v.file {
			s := fmt.Sprintf("%d", v.index)
			str := path.Join(v.dir, s)
			return str, nil
		}
	}
	return "", fmt.Errorf("[2] File %s not found.", fileName)
}

// Returns the directory name of a file placed inside the current directory,
// or an error when not found.
func (self *FileIndex) CurrentDir(fileName string) (string, error) {
	for _, v := range self.fileList {
		if fileName == v.file {
			str := fmt.Sprintf("%d", v.index)
			return str, nil
		}
	}
	return "", fmt.Errorf("[1] File %s not found.", fileName)
}

// Returns the complete directory name of a file number, or an error when not found.
func (self *FileIndex) DirIndex(fileName int64) (string, error) {
	for _, v := range self.fileList {
		if fileName == v.index {
			s := fmt.Sprintf("%d", v.index)
			str := path.Join(v.dir, s)
			return str, nil
		}
	}
	return "", fmt.Errorf("[3] File %d not found.", fileName)
}

// Returns the file name of a file number.
func (self *FileIndex) ContainerName(fileNumber string) (FileList, error) {
	num := ex.Atoi64(fileNumber)
	for _, item := range self.fileList {
		if num == item.index {
			return item, nil
		}
	}
	return FileList{}, fmt.Errorf("File %s not !!! found in the index.", fileNumber)
}

// Returns the file name of a file number.
func (self *FileIndex) FileName(fileName int64) string {
	for _, v := range self.fileList {
		if fileName == v.index {
			return v.file
		}
	}
	return fmt.Sprintf("[4] File %d not found.", fileName)
}

// Returns the file name from a string (instead of an int64) of a file number.
func (self *FileIndex) FileNameOfString(fileName string) string {
	var i int64
	fmt.Sscanf(fileName, "%d", &i)
	return self.FileName(i)
}

// Increases the index number, that is stored in the file 'IndexNumber.txt'
// in the root directory of the vault.
func (self *FileIndex) increase_index_number() int64 {

	buf, err := os.ReadFile(self.indexNumberTxt)
	ex.CheckErr(err)
	_, err = fmt.Sscanf(string(buf), "%d", &self.indexNumber)
	ex.CheckErr(err)

	// Increase index number

	self.indexNumber++

	str := fmt.Sprintf("%d", self.indexNumber)

	err = os.WriteFile(self.indexNumberTxt, []byte(str), 0644)
	ex.CheckErr(err)

	err = os.Chown(self.indexNumberTxt, self.userUid, self.vaultUid)
	ex.CheckErr(err)

	return self.indexNumber
}

// Adds the item from both self.fileIndex and self.fileLocationList
// and store the information on disk. Returns the index number.
// It does not add a file on disk.
func (self *FileIndex) AddItem(filename, dirname string) int64 {

	self.Read() // refreshing the index

	// getting rid of the path
	_, split_file := path.Split(filename)
	fname := split_file

	index := self.increase_index_number()

	fl := FileList{index: index, file: fname, dir: dirname}

	self.fileList = append(self.fileList, fl)

	self.Write()

	return index
}

// Returns the index number of the file name,
// or an error when the file is not found.
func (self *FileIndex) Index(fileName string) (int64, error) {

	for _, v := range self.fileList {
		if fileName == v.file {
			return v.index, nil
		}
	}

	return 0, fmt.Errorf("File %s is not found in the index.", fileName)
}

// Moves the filename to an other directory,
// but only in the FileList, not on disk.
func (self *FileIndex) moveItem(fileNname, directory string) {

	var movefile *FileList

	for index, v := range self.fileList {
		if fileNname == v.file {
			movefile = &self.fileList[index]
			movefile.previousDir = v.dir
			movefile.dir = directory
			break
		}
	}

	self.Write()
}

// Renames the filename from src to dest,
// but only in the FileList, not on disk.
func (self *FileIndex) renameItem(src, dest string) error {

	self.Read() // refreshing index

	// check whether new name already exist

	for _, v := range self.fileList {
		if dest == v.file {
			return errors.New(fmt.Sprintf("Duplicate file in index: %s", dest))
		}
	}

	// Rename

	var renamefile *FileList

	for index, v := range self.fileList {
		if v.file == src {
			renamefile = &self.fileList[index]
			renamefile.previousFile = v.file
			renamefile.file = dest
			break
		}
	}

	// save

	self.Write()

	return nil
}
