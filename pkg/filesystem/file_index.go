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

	"github.com/grd/FreePDM/pkg/util"
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

func NewFileIndex(fs *FileSystem) (ret FileIndex, err error) {

	ret.fs = fs

	ret.fileListCsv = path.Join(ret.fs.dataDir, "FileList.csv")
	ret.indexNumberTxt = path.Join(ret.fs.dataDir, "IndexNumber.txt")

	// check wether the critical files exists.

	util.CriticalFileExist(ret.fileListCsv)
	util.CriticalFileExist(ret.indexNumberTxt)

	ret.getIndexNumber()

	ret.fileList = make([]FileList, 0, ret.indexNumber)

	if err := ret.Read(); err != nil { // read the indexes
		return ret, err
	}

	return ret, nil
}

// Reads the values from "FileList.txt"
func (fix *FileIndex) Read() error {
	// Clear the existing file list
	fix.fileList = nil

	// Read and parse the CSV records
	records, err := fix.readFileListCsv()
	if err != nil {
		return fmt.Errorf("failed to read file list: %w", err)
	}

	// Process each record
	for _, record := range records {
		// Check for correct number of fields
		if len(record) < 5 {
			return fmt.Errorf("invalid record format: %v", record)
		}

		// Parse the record into FileList struct
		index, err := util.Atoi64(record[0])
		if err != nil {
			return fmt.Errorf("invalid index format in record %v: %w", record, err)
		}

		fix.fileList = append(fix.fileList, FileList{
			index:        index,
			file:         record[1],
			previousFile: record[2],
			dir:          record[3],
			previousDir:  record[4],
		})
	}

	return nil
}

// readFileListCsv reads the CSV contents of the "FileList.txt" file
func (fix FileIndex) readFileListCsv() ([][]string, error) {
	// Read the file into a buffer
	buf, err := os.ReadFile(fix.fileListCsv)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fix.fileListCsv, err)
	}

	// Parse the CSV using a reader with colon as delimiter
	r := csv.NewReader(bytes.NewBuffer(buf))
	r.Comma = ':' // Using ':' as a delimiter

	// Read all records from the CSV
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing CSV in file %s: %w", fix.fileListCsv, err)
	}

	// Ensure there are records and skip the header
	if len(records) <= 1 {
		return nil, fmt.Errorf("file %s contains no data or is missing a header", fix.fileListCsv)
	}

	// Skip the header and return the data
	return records[1:], nil
}

// Writes the values to "FileList.csv"
func (fix *FileIndex) Write() error {
	// Initialize the CSV header
	records := [][]string{
		{"Index", "FileName", "PreviousFile", "Dir", "PreviousDir"},
	}

	// Add records from fileList
	for _, item := range fix.fileList {
		records = append(records, []string{
			util.I64toa(item.index),
			item.file,
			item.previousFile,
			item.dir,
			item.previousDir,
		})
	}

	// Create a buffer for writing CSV data
	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)
	writer.Comma = ':'

	// Write the CSV records
	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("failed to write CSV data: %w", err)
	}
	writer.Flush()

	// Check for errors during CSV writing
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %w", err)
	}

	// Write the buffer content to the file
	if err := os.WriteFile(fix.fileListCsv, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fix.fileListCsv, err)
	}

	// Set the correct ownership of the file
	if err := os.Chown(fix.fileListCsv, fix.fs.userUid, fix.fs.vaultUid); err != nil {
		return fmt.Errorf("failed to change ownership of %s: %w", fix.fileListCsv, err)
	}

	return nil
}

func NewFileList(index int64, file, previousFile, dir, previousDir string) FileList {
	return FileList{
		index:        index,
		file:         file,
		previousFile: previousFile,
		dir:          dir,
		previousDir:  previousDir}
}

func (fl FileList) Index() string {
	return fmt.Sprintf("%d", fl.index)
}

func (fl FileList) PathAndIndex() string {
	return path.Join(fl.Path(), fl.Index())
}

func (fl FileList) Name() string {
	return fl.file
}

func (fl FileList) PreviousName() string {
	return fl.previousFile
}

func (fl FileList) Path() string {
	return fl.dir
}

func (fl FileList) PreviousPath() string {
	return fl.previousDir
}

// Reads the index number and stores it.
func (fix *FileIndex) getIndexNumber() {

	buf, err := os.ReadFile(fix.indexNumberTxt)
	util.CheckErr(err)

	_, err = fmt.Sscanf(string(buf), "%d", &fix.indexNumber)
	util.CheckErr(err)
}

// Returns the complete directory name of a file number, or an error when not found.
func (fix *FileIndex) DirIndex(index int64) (string, error) {
	for _, item := range fix.fileList {
		if index == item.index {
			num := fmt.Sprintf("%d", item.index)
			str := path.Join(item.dir, num)
			return str, nil
		}
	}
	return "", fmt.Errorf("index number %d not found", index)
}

// Returns the index number of the file name,
// or an error when the file is not found.
func (fx *FileIndex) FileNameToIndex(fileName string) (int64, error) {
	for _, item := range fx.fileList {
		if fileName == item.file {
			return item.index, nil
		}
	}

	return 0, fmt.Errorf("file %s is not found in the index", fileName)
}

// Input parameter is the file name.
// Returns the path and name of a file, or an error when not found.
func (fix *FileIndex) FileNameToFileList(fileName string) (FileList, error) {
	for _, item := range fix.fileList {
		if fileName == item.file {
			return item, nil
		}
	}
	return FileList{}, fmt.Errorf("file %s not found", fileName)
}

// Returns the file name of a file number.
func (fix *FileIndex) ContainerName(index string) (FileList, error) {
	num, err := util.Atoi64(index)
	if err != nil {
		return FileList{}, err
	}
	for _, item := range fix.fileList {
		if num == item.index {
			return item, nil
		}
	}
	return FileList{}, fmt.Errorf("file %s not found in the index", index)
}

// Returns the file name of a file number.
func (fix *FileIndex) IndexToFileName(index int64) (string, error) {
	for _, item := range fix.fileList {
		if index == item.index {
			return item.file, nil
		}
	}
	return "", fmt.Errorf("file %d not found", index)
}

// Increases the index number, that is stored in the file 'IndexNumber.txt'
// in the root directory of the vault.
func (fix *FileIndex) increase_index_number() int64 {

	buf, err := os.ReadFile(fix.indexNumberTxt)
	util.CheckErr(err)
	_, err = fmt.Sscanf(string(buf), "%d", &fix.indexNumber)
	util.CheckErr(err)

	// Increase index number

	fix.indexNumber++

	str := fmt.Sprintf("%d", fix.indexNumber)

	err = os.WriteFile(fix.indexNumberTxt, []byte(str), 0644)
	util.CheckErr(err)

	err = os.Chown(fix.indexNumberTxt, fix.fs.userUid, fix.fs.vaultUid)
	util.CheckErr(err)

	return fix.indexNumber
}

// Adds the item from both self.fileIndex and self.fileLocationList
// and store the information on disk. Returns the index number.
// It does not add a file on disk.
func (fix *FileIndex) AddItem(filename, dirname string) int64 {

	if err := fix.Read(); err != nil { // refreshing the index
		log.Fatalf("error reading FileIndex.csv, %v", err)
	}
	// getting rid of the path
	_, split_file := path.Split(filename)
	fname := split_file

	index := fix.increase_index_number()

	fl := FileList{index: index, file: fname, dir: dirname}

	fix.fileList = append(fix.fileList, fl)

	if err := fix.Write(); err != nil {
		log.Fatalf("error writing FileIndex.csv, %v", err)
	}

	return index
}

// Moves the filename to an other directory,
// but only in the FileList, not on disk.
func (fix *FileIndex) moveItem(fileNname, directory string) error {

	var movefile *FileList

	for index, v := range fix.fileList {
		if fileNname == v.file {
			movefile = &fix.fileList[index]
			movefile.previousDir = v.dir
			movefile.dir = directory
			break
		}
	}

	if err := fix.Write(); err != nil {
		return err
	}

	return nil
}

// Renames the filename from src to dest,
// but only in the FileList, not on disk.
func (fix *FileIndex) renameItem(src, dest string) error {

	if err := fix.Read(); err != nil { // refreshing index
		return err
	}

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

	if err := fix.Write(); err != nil {
		return err
	}

	return nil
}
