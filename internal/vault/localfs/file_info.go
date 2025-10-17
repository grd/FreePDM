// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package localfs

import "path/filepath"

// This shows only the latest version of a file and is handy when being used with
// the web server. If the file is a directory then only the filename is visible.

// The FileInfo struct
type FileInfo struct {
	containerNumber string
	isDir           bool // Is it a directory or a file?
	name            string
	dir             string
	fileDescription string
	fileSecondDescr string
	fileVersion     string
	fileLocked      bool   // Is the file locked out?
	fileLockedOutBy string // Who checked out the file?
	// fileIcon        []byte
	fileProperties []FileProperties
}

// The File Properties
type FileProperties struct {
	Key, Value string
}

func (fi FileInfo) String() string {
	return filepath.Join(fi.Path(), fi.Name())
}
func (fi FileInfo) IsDir() bool {
	return fi.isDir
}

func (fi FileInfo) DirSort() int {
	if fi.IsDir() {
		return 0
	} else {
		return 1
	}
}

// Returns the directory or file name
func (fi FileInfo) Name() string {
	return fi.name
}

// Returns the path of the file or dir
func (fi FileInfo) Path() string {
	return fi.dir
}

func (fi FileInfo) Description() string {
	return fi.fileDescription
}

func (fi FileInfo) SecondDescription() string {
	return fi.fileSecondDescr
}

func (fi FileInfo) Version() string {
	return fi.fileVersion
}

func (fi FileInfo) IsLocked() bool {
	return fi.fileLocked
}

func (fi FileInfo) LockedOutBy() string {
	return fi.fileLockedOutBy
}

func (fi FileInfo) Properties() []FileProperties {
	return fi.fileProperties
}
