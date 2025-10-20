// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package localfs

import "path/filepath"

// This shows only the latest version of a file and is handy when being used with
// the web server. If the file is a directory then only the filename is visible.

// AllocStatus indicates whether a container is allocated and, if so,
// whether a real file candidate is already present alongside the placeholder.
type AllocStatus int

const (
	// AllocNone means: normal container (no placeholder present).
	AllocNone AllocStatus = iota
	// AllocAllocatedEmpty means: placeholder present, no real file candidate detected.
	AllocAllocatedEmpty
	// AllocAllocatedWithCandidate means: placeholder present AND a real file candidate is present.
	AllocAllocatedWithCandidate
)

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

	// New: allocation status for containers (used by GUI for badges/tooltips)
	allocStatus    AllocStatus
	allocCandidate string // non-empty only when allocStatus == AllocAllocatedWithCandidate
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

// Alloc returns the allocation status for this entry.
// GUI can use this to decide which badge/icon to render.
func (fi FileInfo) Alloc() AllocStatus {
	return fi.allocStatus
}

// AllocCandidate returns the suggested filename candidate (if any) that was
// found inside an allocated container alongside the placeholder. The boolean
// indicates whether a candidate is available.
func (fi FileInfo) AllocCandidate() (string, bool) {
	if fi.allocStatus == AllocAllocatedWithCandidate && fi.allocCandidate != "" {
		return fi.allocCandidate, true
	}
	return "", false
}

// IsAllocated is a convenience helper to check if any allocation state applies.
func (fi FileInfo) IsAllocated() bool {
	return fi.allocStatus == AllocAllocatedEmpty || fi.allocStatus == AllocAllocatedWithCandidate
}
