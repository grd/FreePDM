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

func NewFileIndex(fs *FileSystem) (fi FileIndex, err error) {

	fi.fs = fs

	fi.fileListCsv = path.Join(fi.fs.dataDir, "FileList.csv")
	fi.indexNumberTxt = path.Join(fi.fs.dataDir, "IndexNumber.txt")

	// check whether the critical files exists.

	util.CriticalFileExist(fi.fileListCsv)
	util.CriticalFileExist(fi.indexNumberTxt)

	fi.getIndexNumber()

	fi.fileList = make([]FileList, 0, fi.indexNumber)

	if err := fi.Read(); err != nil { // read the indexes
		return fi, err
	}

	return fi, nil
}

// Reads the values from "FileList.txt"
func (fi *FileIndex) Read() error {
	// Clear the existing file list
	fi.fileList = nil

	// Read and parse the CSV records
	records, err := fi.readFileListCsv()
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

		fi.fileList = append(fi.fileList, FileList{
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
func (fi FileIndex) readFileListCsv() ([][]string, error) {
	// Read the file into a buffer
	buf, err := os.ReadFile(fi.fileListCsv)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fi.fileListCsv, err)
	}

	// Parse the CSV using a reader with colon as delimiter
	r := csv.NewReader(bytes.NewBuffer(buf))
	r.Comma = ':' // Using ':' as a delimiter

	// Read all records from the CSV
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing CSV in file %s: %w", fi.fileListCsv, err)
	}

	// Ensure there are records and skip the header
	if len(records) == 0 {
		return nil, fmt.Errorf("file %s contains no data or is missing a header", fi.fileListCsv)
	}

	// Skip the header and return the data
	return records[1:], nil
}

// Writes the values to "FileList.csv"
func (fi *FileIndex) Write() error {
	// Initialize the CSV header
	records := [][]string{
		{"Index", "FileName", "PreviousFile", "Dir", "PreviousDir"},
	}

	// Add records from fileList
	for _, item := range fi.fileList {
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
	if err := os.WriteFile(fi.fileListCsv, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fi.fileListCsv, err)
	}

	// Set the correct ownership of the file
	if err := os.Chown(fi.fileListCsv, fi.fs.userUid, fi.fs.vaultUid); err != nil {
		return fmt.Errorf("failed to change ownership of %s: %w", fi.fileListCsv, err)
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

func (fl FileList) IndexInt64() int64 {
	return fl.index
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
func (fi *FileIndex) getIndexNumber() {

	buf, err := os.ReadFile(fi.indexNumberTxt)
	util.CheckErr(err)

	_, err = fmt.Sscanf(string(buf), "%d", &fi.indexNumber)
	util.CheckErr(err)
}

// Returns the FileList struct from an index number, or an error when not found.
func (fi *FileIndex) IndexToFileList(index int64) (FileList, error) {
	for _, item := range fi.fileList {
		if index == item.index {
			return item, nil
		}
	}
	return FileList{}, fmt.Errorf("index number %d not found", index)
}

// Returns the complete directory name of a file number, or an error when not found.
func (fi *FileIndex) IndexDir(index int64) (string, error) {
	fl, err := fi.IndexToFileList(index)
	if err != nil {
		return "", fmt.Errorf("index number %d not found", index)
	}
	num := fmt.Sprintf("%d", fl.index)
	return path.Join(fl.dir, num), nil
}

// Returns the file name of a file number.
func (fi *FileIndex) IndexToFileName(index int64) (string, error) {
	fl, err := fi.IndexToFileList(index)
	if err != nil {
		return "", fmt.Errorf("file %d not found", index)
	}
	return fl.file, nil
}

// Returns the index number of the file name,
// or an error when the file is not found.
func (fi *FileIndex) FileNameToIndex(fileName string) (int64, error) {
	for _, item := range fi.fileList {
		if fileName == item.file {
			return item.index, nil
		}
	}

	return 0, fmt.Errorf("file %s is not found in the index", fileName)
}

// Input parameter is the file name.
// Returns the path and name of a file, or an error when not found.
func (fi *FileIndex) FileNameToFileList(fileName string) (FileList, error) {
	for _, item := range fi.fileList {
		if fileName == item.file {
			return item, nil
		}
	}
	return FileList{}, fmt.Errorf("file %s not found", fileName)
}

// Returns the file list when a file number is found, else an error.
func (fi *FileIndex) ContainerName(index string) (FileList, error) {
	num, err := util.Atoi64(index)
	if err != nil {
		return FileList{}, err
	}
	for _, item := range fi.fileList {
		if num == item.index {
			return item, nil
		}
	}
	return FileList{}, fmt.Errorf("file %s not found in the index", index)
}

// Increases the index number, that is stored in the file 'IndexNumber.txt'
// in the root directory of the vault.
func (fi *FileIndex) increase_index_number() (int64, error) {

	buf, err := os.ReadFile(fi.indexNumberTxt)
	if err != nil {
		return -1, err
	}
	_, err = fmt.Sscanf(string(buf), "%d", &fi.indexNumber)
	if err != nil {
		return -1, err
	}

	// Increase index number

	fi.indexNumber++

	str := fmt.Sprintf("%d", fi.indexNumber)

	if err = os.WriteFile(fi.indexNumberTxt, []byte(str), 0644); err != nil {
		return -1, err
	}

	if err = os.Chown(fi.indexNumberTxt, fi.fs.userUid, fi.fs.vaultUid); err != nil {
		return -1, err
	}

	return fi.indexNumber, nil
}

// Adds the item from both self.fileIndex and self.fileLocationList
// and store the information on disk. Returns the index number.
// It does not add a file on disk.
func (fi *FileIndex) AddItem(filename, dirname string) (int64, error) {

	if err := fi.Read(); err != nil { // refreshing the index
		log.Fatalf("error reading FileIndex.csv, %v", err)
	}
	// getting rid of the path
	_, split_file := path.Split(filename)
	fname := split_file

	index, err := fi.increase_index_number()
	if err != nil {
		return -1, err
	}

	fl := FileList{index: index, file: fname, dir: dirname}

	fi.fileList = append(fi.fileList, fl)

	if err := fi.Write(); err != nil {
		log.Fatalf("error writing FileIndex.csv, %v", err)
	}

	return index, nil
}

// Moves the filename to an other directory,
// but only in the FileList, not on disk.
func (fi *FileIndex) moveItem(fileNname, directory string) error {

	i := -1

	for index, v := range fi.fileList {
		if fileNname == v.file {
			i = index
			break
		}
	}

	if i == -1 {
		return fmt.Errorf("filename %s is not inside the FileIndex", fileNname)
	}

	// Moving the item
	fi.fileList[i].previousDir = fi.fileList[i].dir
	fi.fileList[i].dir = directory

	if err := fi.Write(); err != nil {
		return err
	}

	return nil
}

// Renames the filename from src to dest,
// but only in the FileList, not on disk.
func (fi *FileIndex) renameItem(src, dest string) error {

	if err := fi.Read(); err != nil { // refreshing index
		return err
	}

	// check whether new name already exist

	for _, v := range fi.fileList {
		if dest == v.file {
			return fmt.Errorf("duplicate file in index: %s", dest)
		}
	}

	// Rename

	var renamefile *FileList

	for index, v := range fi.fileList {
		if v.file == src {
			renamefile = &fi.fileList[index]
			renamefile.previousFile = v.file
			renamefile.file = dest
			break
		}
	}

	// save

	if err := fi.Write(); err != nil {
		return err
	}

	return nil
}
