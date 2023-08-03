// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filesystem

// This shows only the latest version of a file and is handy when being used with
// the web server. If the file is a directory then only the filename is visible.

// The FileInfo struct
type FileInfo struct {
	Dir             bool             `json:"Dir"` // Is it a directory or a file?
	FileName        string           `json:"FileName"`
	FileDescription string           `json:"FileDescription"`
	FileSecondDescr string           `json:"FileSecondDescr"`
	FileVersion     string           `json:"FileVersion"`
	FileLocked      bool             `json:"FileLocked"`      // Is the file locked out?
	FileLockedOutBy string           `json:"FileLockedOutBy"` // Who checked out the file?
	FileIcon        []byte           `json:"FileIcon"`
	FileProperties  []FileProperties `json:"FileProperties"`
}

// The File Properties
type FileProperties struct {
	Key, Value string
}

func (self FileInfo) IsDir() bool {
	return self.Dir
}

// Returns the directory or file name
func (self FileInfo) Name() string {
	return self.FileName
}

func (self FileInfo) Description() string {
	return self.FileDescription
}

func (self FileInfo) SecondDescription() string {
	return self.FileSecondDescr
}

func (self FileInfo) Version() string {
	return self.FileVersion
}

func (self FileInfo) IsLocked() bool {
	return self.FileLocked
}

func (self FileInfo) LockedOutBy() string {
	return self.FileLockedOutBy
}

func (self FileInfo) Properties() []FileProperties {
	return self.FileProperties
}
