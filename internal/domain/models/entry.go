// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"io/fs"
	"time"
)

// Entry represents a file or directory inside a vault.
type Entry struct {
	Name        string            // display name
	RelPath     string            // path relative to vault root ("" == root)
	IsDir       bool              // true = directory; false = file (incl. numeric container)
	Container   bool              // true if "file as directory" (numeric-dir semantic)
	Size        int64             // bytes (0 for directories)
	Mode        fs.FileMode       // POSIX mode bits
	ModTime     time.Time         // last modified timestamp
	ContentHash string            // optional: content hash for version pinning
	Meta        map[string]string // arbitrary metadata (index hints, etc.)
}
