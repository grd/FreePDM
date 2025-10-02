// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package localfs

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"slices"

	"github.com/grd/FreePDM/internal/util"
)

// The FileList is a struct with five fields:
// * The container number.
// * The file name.
// * The previous name of the file.
// * The dir name of the storage. This is an offset of "mainPdmDir".
// * The previous directory name.
type FileList struct {
	ContainerNumber string
	Name            string
	PreviousName    string
	Path            string
	PreviousPath    string
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
	fi.indexNumberTxt = path.Join(fi.fs.dataDir, "ContainerNumber.txt")

	// check whether the critical files exists.

	util.CriticalFileExist(fi.fileListCsv)
	util.CriticalFileExist(fi.indexNumberTxt)

	fi.getContainerNumber()

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

		fi.fileList = append(fi.fileList, FileList{
			ContainerNumber: record[0],
			Name:            record[1],
			PreviousName:    record[2],
			Path:            record[3],
			PreviousPath:    record[4],
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
		{"Container", "FileName", "PreviousFile", "Dir", "PreviousDir"},
	}

	// Add records from fileList
	for _, item := range fi.fileList {
		records = append(records, []string{
			item.ContainerNumber,
			item.Name,
			item.PreviousName,
			item.Path,
			item.PreviousPath,
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
	if err := fi.fs.DataWriteFile(fi.fileListCsv, buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

// Reads the index number and stores it.
func (fi *FileIndex) getContainerNumber() {

	buf, err := os.ReadFile(fi.indexNumberTxt)
	util.CheckErr(err)

	_, err = fmt.Sscanf(string(buf), "%d", &fi.indexNumber)
	util.CheckErr(err)
}

// Returns the FileList struct from an container number, or an error when not found.
func (fi *FileIndex) ContainerNumberToFileList(containerNumber string) (FileList, error) {
	for _, item := range fi.fileList {
		if containerNumber == item.ContainerNumber {
			return item, nil
		}
	}
	return FileList{}, fmt.Errorf("container number %s not found", containerNumber)
}

// Returns the complete directory name of a container number, or an error when not found.
func (fi *FileIndex) ContainerNumberDir(containerNumber string) (string, error) {
	fl, err := fi.ContainerNumberToFileList(containerNumber)
	if err != nil {
		return "", err
	}
	return path.Join(fl.Path, fl.ContainerNumber), nil
}

// Returns the file name of a container number.
func (fi *FileIndex) ContainerNumberToFileName(containerNumber string) (string, error) {
	fl, err := fi.ContainerNumberToFileList(containerNumber)
	if err != nil {
		return "", err
	}
	return fl.Name, nil
}

// Input parameter is the file name.
// Returns the FileList and an error when not found.
func (fi *FileIndex) FileNameToFileList(dir, fileName string) (FileList, error) {
	if dir == "." {
		dir = ""
	}
	for _, item := range fi.fileList {
		if dir == item.Path && fileName == item.Name {
			return item, nil
		}
	}
	return FileList{}, fmt.Errorf("file %s not found in the FileList", fileName)
}

// Returns the index number of the file name in directory,
// or an error when the file is not found.
func (fi *FileIndex) FileNameToContainerNumber(dir, fileName string) (string, error) {
	fl, err := fi.FileNameToFileList(dir, fileName)
	if err != nil {
		return "", err
	}
	return fl.ContainerNumber, nil
}

// Increases the container number, that is stored in the file 'IndexNumber.txt'
// in the root directory of the vault.
func (fi *FileIndex) increaseContainerNumber() (string, error) {

	buf, err := os.ReadFile(fi.indexNumberTxt)
	if err != nil {
		return "", err
	}
	_, err = fmt.Sscanf(string(buf), "%d", &fi.indexNumber)
	if err != nil {
		return "", err
	}

	// Increase index number

	fi.indexNumber++

	str := fmt.Sprintf("%d", fi.indexNumber)

	if err = fi.fs.DataWriteFile(fi.indexNumberTxt, []byte(str)); err != nil {
		return "", err
	}

	return str, nil
}

// Adds the filename to the filelist. Returns the index number, and an error.
// It does not add a file on disk.
func (fi *FileIndex) AddItem(dirName, fileName string) (*FileList, error) {

	if err := fi.Read(); err != nil { // refreshing the index
		log.Fatalf("error reading FileIndex.csv, %v", err)
	}
	// getting rid of the path
	_, split_file := path.Split(fileName)
	fname := split_file

	index, err := fi.increaseContainerNumber()
	if err != nil {
		return nil, err
	}

	fl := FileList{ContainerNumber: index, Name: fname, Path: dirName}

	fi.fileList = append(fi.fileList, fl)

	if err := fi.Write(); err != nil {
		log.Fatalf("error writing FileIndex.csv, %v", err)
	}

	return &fl, nil
}

// Moves the filename to an other directory,
// but only in the FileList, not on disk.
func (fi *FileIndex) MoveItem(src FileList, directory string) error {
	if directory == "." {
		directory = ""
	}

	i := -1

	for index, v := range fi.fileList {
		if src.ContainerNumber == v.ContainerNumber {
			i = index
			break
		}
	}

	if i == -1 {
		return fmt.Errorf("filename %s is not inside the FileIndex", path.Join(src.Path, src.Name))
	}

	// Moving the item
	fi.fileList[i].PreviousPath = fi.fileList[i].Path
	fi.fileList[i].Path = directory

	if err := fi.Write(); err != nil {
		return err
	}

	return nil
}

// Renames the filename from src to dest,
// but only in the FileList, not on disk.
func (fi *FileIndex) renameItem(src FileList, dest string) error {

	if err := fi.Read(); err != nil { // refreshing index
		return err
	}

	// check whether new name already exist

	for _, v := range fi.fileList {
		if dest == v.Name {
			return fmt.Errorf("duplicate file in index: %s", dest)
		}
	}

	// Rename

	var renamefile *FileList

	for index, v := range fi.fileList {
		if v.ContainerNumber == src.ContainerNumber {
			renamefile = &fi.fileList[index]
			renamefile.PreviousName = v.Name
			renamefile.Name = dest
			break
		}
	}

	// save

	if err := fi.Write(); err != nil {
		return err
	}

	return nil
}

// Removes the container number from the list, or an error.
func (fi *FileIndex) ContainerNumberRemove(containerNumber string) error {
	ok := false
	j := -1
	for i, item := range fi.fileList {
		if containerNumber == item.ContainerNumber {
			ok = true
			j = i
		}
	}
	if !ok {
		return fmt.Errorf("container number %s not found in the index", containerNumber)
	}

	fi.fileList = slices.Delete(fi.fileList, j, j+1)

	if err := fi.Write(); err != nil {
		return err
	}

	return nil
}
