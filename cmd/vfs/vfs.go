// Copyright 2024 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	iofs "io/fs"

	"github.com/hirochachacha/go-smb2"
)

// Create a multi-level file system with directories and files (that are directories too)
// with a structure that looks like this:
//
//	file1/        file2/        dir1/        file3/         // rw
//	  1/ 2/         1/ 2/ 3/      file4/      1/            // ro
//	    data          data          ^           data        // ro + rw
//                                  |
//                                  |
//                                  |   In this case file4 also has the multi-level file structure
//                                      and is stored inside dir1.
//

// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Tests for the SMB share functionality.
func runner(path string, d iofs.DirEntry, err error) error {
	fmt.Println(path, d, err)

	return nil
}

func RunTheProgram(fs *smb2.Share) error {

	fmt.Println("here")
	matches, err := iofs.Glob(fs.DirFS("."), "*")
	if err != nil {
		panic(err)
	}
	for _, match := range matches {
		fmt.Println(match)
	}

	err = iofs.WalkDir(fs.DirFS("."), ".", runner)
	if err != nil {
		panic(err)
	}

	return nil
}
