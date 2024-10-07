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
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/cloudsoda/go-smb2"
)

// Recursively add 'a' to all files and directories
func renameFilesRecursive(share *smb2.Share, path string) error {
	// List all directory entries in the current path
	dirEntries, err := share.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, entry := range dirEntries {
		oldName := filepath.Join(path, entry.Name())
		newName := filepath.Join(path, "a"+entry.Name())

		// Rename the file or directory
		err := share.Rename(oldName, newName)
		if err != nil {
			log.Printf("Failed to rename %s to %s: %v", oldName, newName, err)
			continue
		}

		log.Printf("Renamed %s to %s", oldName, newName)

		// If it's a directory, recurse into it
		if entry.IsDir() {
			err = renameFilesRecursive(share, newName)
			if err != nil {
				log.Printf("Failed to recurse into directory %s: %v", newName, err)
				continue
			}
		}
	}

	return nil
}

func main() {
	// Create an SMB2 session using NTLM authentication
	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     "username", // Replace with your username
			Password: "password", // Replace with your password
			Domain:   "",         // Leave empty if not needed
		},
	}

	// Dial the SMB2 session
	smbSession, err := dialer.Dial(context.Background(), "smb-server-address:445") // Provide the address and port here
	if err != nil {
		log.Fatalf("Failed to start SMB session: %v", err)
	}
	defer smbSession.Logoff()

	// Mount the SMB share
	share, err := smbSession.Mount("share-name") // Replace with your share name
	if err != nil {
		log.Fatalf("Failed to mount SMB share: %v", err)
	}
	defer share.Umount()

	rootDir := "/path/to/remote/directory" // Replace with your remote directory

	// Start renaming files from the root directory
	err = renameFilesRecursive(share, rootDir)
	if err != nil {
		log.Fatalf("Failed to rename files: %v", err)
	}

	log.Println("Renaming completed successfully")
}
