// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	ex "github.com/grd/FreePDM/src/extras"
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

	f, err := os.Open(self.fileListCsv)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	// skip first line
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

// Writes the values to "FileList.csv"
func (self *FileIndex) Write() {

	file, err := os.OpenFile(self.fileListCsv, os.O_WRONLY|os.O_CREATE, 0644)
	ex.CheckErr(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	firstRecord := []string{
		"Index", "FileName", "PreviousFile", "Dir", "PreviousDir",
	}

	if err := writer.Write(firstRecord); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	record := make([]string, 5)

	for _, item := range self.fileList {

		record[0] = ex.I64toa(item.index)
		record[1] = item.file
		record[2] = item.previousFile
		record[3] = item.dir
		record[4] = item.previousDir

		if err := writer.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

	err = os.Chown(self.fileListCsv, self.userUid, self.vaultUid)
	ex.CheckErr(err)
}

// Write a record to a file
func addRecord(fname string, column []string) {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	ex.CheckErr(err)

	writer := csv.NewWriter(file)
	writer.Write(column)

	writer.Flush()
	file.Close()
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
	return "", fmt.Errorf("File %s not found.", fileName)
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
	return "", fmt.Errorf("File %s not found.", fileName)
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
	return "", fmt.Errorf("File %d not found.", fileName)
}

// Returns the file name of a file number.
func (self *FileIndex) FileName(fileName int64) string {
	for _, v := range self.fileList {
		if fileName == v.index {
			return v.file
		}
	}
	return fmt.Sprintf("File %d not found.", fileName)
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

	// getting rid of the path
	_, split_file := path.Split(filename)
	fname := split_file

	self.Read() // refreshing the index

	index := self.increase_index_number()
	fl := FileList{index: index, file: fname, dir: dirname}

	self.fileList = append(self.fileList, fl)

	// Appending to file

	var record = make([]string, 5)

	record[0] = ex.I64toa(fl.index)
	record[1] = fl.file
	record[2] = fl.previousFile
	record[3] = fl.dir
	record[4] = fl.previousDir

	addRecord(self.fileListCsv, record)

	err := os.Chown(self.fileListCsv, self.userUid, self.vaultUid)
	ex.CheckErr(err)

	return index
}

// Returns the index number of the file name,
// or -1 when the file is not found.
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
	for _, v := range self.fileList {
		if fileNname == v.file {
			v.previousDir = v.dir
			v.dir = directory
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

	for _, v := range self.fileList {
		if v.file == src {
			v.previousFile = v.file
			v.file = dest
			break
		}
	}

	// save

	self.Write()

	return nil
}
