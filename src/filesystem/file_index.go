// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"bytes"
	"encoding/csv"
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
	fs *FileSystem

	fileListCsv string
	fileList    []FileList

	indexNumberTxt string
	indexNumber    int64
}

func InitFileIndex(fs *FileSystem) (ret FileIndex) {

	ret.fs = fs

	ret.fileListCsv = path.Join(ret.fs.dataDir, "FileList.csv")
	ret.indexNumberTxt = path.Join(ret.fs.dataDir, "IndexNumber.txt")

	// check wether the critical files exists.

	ex.CriticalFileExist(ret.fileListCsv)
	ex.CriticalFileExist(ret.indexNumberTxt)

	ret.getIndexNumber()

	ret.fileList = make([]FileList, 0, ret.indexNumber)

	ret.Read() // read the indexes

	return ret
}

// Reads the values from "FileList.txt"
func (fix *FileIndex) Read() {

	fix.fileList = nil

	var fl FileList

	records, err := fix.readCsv()
	ex.CheckErr(err)

	for _, record := range records {

		fl.index = ex.Atoi64(record[0])
		fl.file = record[1]
		fl.previousFile = record[2]
		fl.dir = record[3]
		fl.previousDir = record[4]

		fix.fileList = append(fix.fileList, fl)
	}
}

func (fix FileIndex) readCsv() ([][]string, error) {

	buf, err := os.ReadFile(fix.fileListCsv)
	ex.CheckErr(err)

	r := csv.NewReader(bytes.NewBuffer(buf))
	r.Comma = ':'

	records, err := r.ReadAll()
	if err != nil {
		return [][]string{}, fmt.Errorf("error reading file %s with error: %v", fix.fileListCsv, err)
	}

	if len(records) == 0 {
		return [][]string{}, fmt.Errorf("error reading file %s with error: No header included", fix.fileListCsv)
	}

	records = records[1:]

	return records, nil
}

// Writes the values to "FileList.csv"
func (fix *FileIndex) Write() {

	records := [][]string{
		{"Index", "FileName", "PreviousFile", "Dir", "PreviousDir"},
	}

	for _, item := range fix.fileList {

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

	err = os.WriteFile(fix.fileListCsv, buffer.Bytes(), 0644)
	ex.CheckErr(err)

	err = os.Chown(fix.fileListCsv, fix.fs.userUid, fix.fs.vaultUid)
	ex.CheckErr(err)
}

// Reads the index number and stores it.
func (fix *FileIndex) getIndexNumber() {

	buf, err := os.ReadFile(fix.indexNumberTxt)
	ex.CheckErr(err)

	_, err = fmt.Sscanf(string(buf), "%d", &fix.indexNumber)
	ex.CheckErr(err)
}

// Returns the complete directory name of a file, or an error when not found.
func (fix *FileIndex) Dir(fileName string) (string, error) {
	for _, v := range fix.fileList {
		if fileName == v.file {
			s := fmt.Sprintf("%d", v.index)
			str := path.Join(v.dir, s)
			return str, nil
		}
	}
	return "", fmt.Errorf("[2] file %s not found", fileName)
}

// Returns the directory name of a file placed inside the current directory,
// or an error when not found.
func (fix *FileIndex) CurrentDir(fileName string) (string, error) {
	for _, v := range fix.fileList {
		if fileName == v.file {
			str := fmt.Sprintf("%d", v.index)
			return str, nil
		}
	}
	return "", fmt.Errorf("[1] file %s not found", fileName)
}

// Returns the complete directory name of a file number, or an error when not found.
func (fix *FileIndex) DirIndex(fileName int64) (string, error) {
	for _, v := range fix.fileList {
		if fileName == v.index {
			s := fmt.Sprintf("%d", v.index)
			str := path.Join(v.dir, s)
			return str, nil
		}
	}
	return "", fmt.Errorf("[3] file %d not found", fileName)
}

// Returns the file name of a file number.
func (fix *FileIndex) ContainerName(fileNumber string) (FileList, error) {
	num := ex.Atoi64(fileNumber)
	for _, item := range fix.fileList {
		if num == item.index {
			return item, nil
		}
	}
	return FileList{}, fmt.Errorf("file %s not !!! found in the index", fileNumber)
}

// Returns the file name of a file number.
func (fix *FileIndex) FileName(fileName int64) string {
	for _, v := range fix.fileList {
		if fileName == v.index {
			return v.file
		}
	}
	return fmt.Sprintf("[4] File %d not found.", fileName)
}

// Returns the file name from a string (instead of an int64) of a file number.
func (fix *FileIndex) FileNameOfString(fileName string) string {
	var i int64
	fmt.Sscanf(fileName, "%d", &i)
	return fix.FileName(i)
}

// Increases the index number, that is stored in the file 'IndexNumber.txt'
// in the root directory of the vault.
func (fix *FileIndex) increase_index_number() int64 {

	buf, err := os.ReadFile(fix.indexNumberTxt)
	ex.CheckErr(err)
	_, err = fmt.Sscanf(string(buf), "%d", &fix.indexNumber)
	ex.CheckErr(err)

	// Increase index number

	fix.indexNumber++

	str := fmt.Sprintf("%d", fix.indexNumber)

	err = os.WriteFile(fix.indexNumberTxt, []byte(str), 0644)
	ex.CheckErr(err)

	err = os.Chown(fix.indexNumberTxt, fix.fs.userUid, fix.fs.vaultUid)
	ex.CheckErr(err)

	return fix.indexNumber
}

// Adds the item from both self.fileIndex and self.fileLocationList
// and store the information on disk. Returns the index number.
// It does not add a file on disk.
func (fix *FileIndex) AddItem(filename, dirname string) int64 {

	fix.Read() // refreshing the index

	// getting rid of the path
	_, split_file := path.Split(filename)
	fname := split_file

	index := fix.increase_index_number()

	fl := FileList{index: index, file: fname, dir: dirname}

	fix.fileList = append(fix.fileList, fl)

	fix.Write()

	return index
}

// Returns the index number of the file name,
// or an error when the file is not found.
func (fix *FileIndex) Index(fileName string) (int64, error) {

	for _, v := range fix.fileList {
		if fileName == v.file {
			return v.index, nil
		}
	}

	return 0, fmt.Errorf("file %s is not found in the index", fileName)
}

// Moves the filename to an other directory,
// but only in the FileList, not on disk.
func (fix *FileIndex) moveItem(fileNname, directory string) {

	var movefile *FileList

	for index, v := range fix.fileList {
		if fileNname == v.file {
			movefile = &fix.fileList[index]
			movefile.previousDir = v.dir
			movefile.dir = directory
			break
		}
	}

	fix.Write()
}

// Renames the filename from src to dest,
// but only in the FileList, not on disk.
func (fix *FileIndex) renameItem(src, dest string) error {

	fix.Read() // refreshing index

	// check whether new name already exist

	for _, v := range fix.fileList {
		if dest == v.file {
			return fmt.Errorf("duplicate file in index: %s", dest)
		}
	}

	// Rename

	var renamefile *FileList

	for index, v := range fix.fileList {
		if v.file == src {
			renamefile = &fix.fileList[index]
			renamefile.previousFile = v.file
			renamefile.file = dest
			break
		}
	}

	// save

	fix.Write()

	return nil
}
