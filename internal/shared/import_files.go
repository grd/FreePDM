// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shared

// This allows automatic import of files from the file share
// with drag-and-drop.

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/grd/FreePDM/internal/util"
	"github.com/grd/FreePDM/internal/vaults"
)

func ImportSharedFiles() {
	vaultDir := vaults.Root() // Vault-directory

	// Periodically check processing a file
	ticker := time.NewTicker(250 * time.Millisecond) // Each quarter of a second
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				// Start importprocedure
				err := processFiles(vaultDir)
				if err != nil {
					log.Printf("Error during file processing: %v\n", err)
				}
			case <-done:
				return
			}
		}
	}()

	select {} // Keep running until the user ends it
}

// processFiles walks the share and imports files
func processFiles(vaultDir string) error {
	excludeDir := vaults.RootData() // Exclude this directory

	return filepath.WalkDir(vaultDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %v", path, err)
		}

		// Exclude this directory
		if path == excludeDir {
			return fs.SkipDir
		}

		// Excluse Trash directory
		// if trashDir(path) {
		// 	return fs.SkipDir
		// }

		// Exclude the containers that exist
		if skipContainer(d.Name()) {
			return fs.SkipDir
		}

		// Only process files
		if !d.IsDir() {
			if err := importAndDeleteFile(path); err != nil {
				return err
			}
		}
		return nil
	})
}

// // skip trash director[y|ies]
// func trashDir(path string) bool {
// 	if len(path) < 6 {
// 		return false
// 	}
// 	if path[0:5] == ".Trash" {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// skipContainer checks wheter a file is a container
func skipContainer(file string) bool {
	_, f := path.Split(file)
	_, err := util.Atoi64(f)
	if err != nil {
		return false
	} else {
		return true
	}
}

// import and delete a drag-and-dropped file
func importAndDeleteFile(ppath string) error {
	// getting who put the file there
	user, err := getFileOwner(ppath)
	if err != nil {
		return err
	}

	if strings.Contains(ppath, vaults.Root()+"/") {
		p := ppath[len(vaults.Root())+1:] // skipping the vaults directory
		dir, _ := path.Split(p)
		d := strings.Split(dir, "/")
		if len(d) > 0 {
			vault := d[0]

			rest := dir[len(vault)+1:]     // maybe this can crash
			rest = strings.Trim(rest, "/") // remove first and last '/'

			// fmt.Printf("vault = %s, file =%s, rest=%s\n", vault, ppath, rest)
			v, err := vaults.NewFileSystem(vault, user)
			if err != nil {
				return err
			}

			_, err = v.ImportFile(rest, ppath)
			if err != nil {
				return err
			}
		}
	} else {
		log.Fatalf("something is very wrong, %s is not inside the vaults directory", ppath)
	}

	fmt.Println(ppath)

	// Remove the original file
	if err := os.Remove(ppath); err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}

	fmt.Printf("Successfully imported and deleted: %s\n", ppath)

	return nil
}

func getFileOwner(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("could not get file info: %w", err)
	}

	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return "", fmt.Errorf("not a valid stat_t type")
	}

	// Get user name
	owner, err := user.LookupId(fmt.Sprint(stat.Uid))
	if err != nil {
		return "", fmt.Errorf("could not find user for UID %d: %w", stat.Uid, err)
	}

	return owner.Username, nil
}
