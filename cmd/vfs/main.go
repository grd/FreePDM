// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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

package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/hirochachacha/go-smb2"
)

func main() {
	// Set up SMB connection to your share
	conn, err := net.Dial("tcp", "smb://your-share-name:445")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Get a list of files and directories in the root directory
	filesAndDirs, err := os.ReadDir(conn.Root())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, fileOrDir := range filesAndDirs {
		// Create a new path with 'a' added to the beginning
		newPath := "a" + fileOrDir.Name()

		// Open the directory for reading (in binary mode)
		d, err := conn.OpenFile(fileOrDir.Name(), os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer d.Close()

		// Read the contents of the directory into memory
		buf, err := io.ReadAll(d)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Write the modified buffer back to the original file/directory
		d.Seek(0, 0)
		_, err = d.Write(buf)

		// Recursively traverse subdirectories and files
		for _, entry := range bytes.Split(buf, []byte("\n")) {
			if len(entry) > 1 && entry[0] == 'a' { // skip already processed entries
				continue
			}

			entryPath := fileOrDir.Name() + "/" + string(entry)
			newEntryPath := "a" + entryPath

			// Open the subdirectory or file for reading (in binary mode)
			subdir, err := conn.OpenFile(newEntryPath, os.O_RDONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer subdir.Close()

			// Read the contents of the subdirectory or file into memory
			buf2, err := io.ReadAll(subdir)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Write the modified buffer back to the original subdirectory/file
			subdir.Seek(0, 0)
			_, err = subdir.Write(buf2)

			// Recursively traverse deeper levels if necessary
			if len(entry) > 1 && entry[0] == '/' { // found a directory
				setup(conn, newPath+"/"+newEntryPath) // recursive call!
			}
		}
	}
}

func setup(conn *smb2.Conn, path string) {
	// Get a list of files and directories in the current directory
	filesAndDirs, err := os.ReadDir(conn.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644))
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, fileOrDir := range filesAndDirs {
		// Create a new path with 'a' added to the beginning
		newPath := "a" + fileOrDir.Name()

		// Open the directory for reading (in binary mode)
		d, err := conn.OpenFile(path+"/"+fileOrDir.Name(), os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer d.Close()

		// Read the contents of the directory into memory
		buf, err := io.ReadAll(d)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Write the modified buffer back to the original file/directory
		d.Seek(0, 0)
		_, err = d.Write(buf)

		// Recursively traverse subdirectories and files
		for _, entry := range bytes.Split(buf, []byte("\n")) {
			if len(entry) > 1 && entry[0] == 'a' { // skip already processed entries
				continue
			}

			entryPath := path + "/" + fileOrDir.Name() + "/" + string(entry)
			newEntryPath := "a" + entryPath

			// Open the subdirectory or file for reading (in binary mode)
			subdir, err := conn.OpenFile(newEntryPath, os.O_RDONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer subdir.Close()

			// Read the contents of the subdirectory or file into memory
			buf2, err := io.ReadAll(subdir)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Write the modified buffer back to the original subdirectory/file
			subdir.Seek(0, 0)
			_, err = subdir.Write(buf2)

			// Recursively traverse deeper levels if necessary
			if len(entry) > 1 && entry[0] == '/' { // found a directory
				setup(conn, newPath+"/"+newEntryPath) // recursive call!
			}
		}
	}
}
