// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package client

import (
	"fmt"
	"os"
	"path/filepath"
)

// tree recursively prints the directory structure starting from root
func Tree(path string, indent string) error {
	entries, err := os.ReadDir(path) // ReadDir returns a slice of fs.DirEntry
	if err != nil {
		return err
	}

	for i, entry := range entries {
		isLast := i == len(entries)-1
		prefix := "├── "
		nextIndent := indent + "│   "

		if isLast {
			prefix = "└── "
			nextIndent = indent + "    "
		}

		fmt.Println(indent + prefix + entry.Name())

		if entry.IsDir() {
			err := Tree(filepath.Join(path, entry.Name()), nextIndent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
