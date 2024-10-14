// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"path"
)

// This shows only the latest version of a file and is handy when being used with
// the web server. If the file is a directory then only the filename is visible.

// The FileInfo struct
type FileInfo struct {
	Dir             bool // Is it a directory or a file?
	FileName        string
	FilePath        string
	FileDescription string
	FileSecondDescr string
	FileVersion     string
	FileLocked      bool   // Is the file locked out?
	FileLockedOutBy string // Who checked out the file?
	FileIcon        []byte
	FileProperties  []FileProperties
}

// The File Properties
type FileProperties struct {
	Key, Value string
}

func (fi FileInfo) String() string {
	return path.Join(fi.Path(), fi.Name())
}
func (fi FileInfo) IsDir() bool {
	return fi.Dir
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
	return fi.FileName
}

// Returns the path of the file or dir
func (fi FileInfo) Path() string {
	return fi.FilePath
}

func (fi FileInfo) Description() string {
	return fi.FileDescription
}

func (fi FileInfo) SecondDescription() string {
	return fi.FileSecondDescr
}

func (fi FileInfo) Version() string {
	return fi.FileVersion
}

func (fi FileInfo) IsLocked() bool {
	return fi.FileLocked
}

func (fi FileInfo) LockedOutBy() string {
	return fi.FileLockedOutBy
}

func (fi FileInfo) Properties() []FileProperties {
	return fi.FileProperties
}
