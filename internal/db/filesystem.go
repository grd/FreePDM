// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
)

// File System related Class
type FileSystem struct{}

func (fs *FileSystem) Init() {
	fmt.Println("Generic File System")
}

// Create new folder
func (fs *FileSystem) CreateFolder() {
	// TODO: Implement this function
	fmt.Println("Create new folder") // Replace with actual implementation
}

// Create file copy
func (fs *FileSystem) CreateFileIdx() {
	// TODO: Implement this function
	// create copy of file with index: for example 12345678.FCStd.2
	fmt.Println("Create file copy") // Replace with actual implementation
}
