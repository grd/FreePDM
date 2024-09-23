// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package vfs

import (
	"context"
	"io"
	"os"
	"strings"

	fsm "github.com/grd/FreePDM/src/filesystem"
	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"github.com/jacobsa/timeutil"
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

type VcsFS struct {
	fsm.FileSystem

	Clock timeutil.Clock
}

func NewVFS(clock timeutil.Clock) (fuse.Server, error) {
	fs := &VcsFS{
		Clock: clock,
	}

	return fuseutil.NewFileSystemServer(fs), nil
}

// BatchForget implements fuseutil.FileSystem.
func (*VcsFS) BatchForget(context.Context, *fuseops.BatchForgetOp) error {
	panic("unimplemented")
}

// CreateFile implements fuseutil.FileSystem.
func (*VcsFS) CreateFile(context.Context, *fuseops.CreateFileOp) error {
	panic("unimplemented")
}

// CreateLink implements fuseutil.FileSystem.
func (*VcsFS) CreateLink(context.Context, *fuseops.CreateLinkOp) error {
	panic("unimplemented")
}

// CreateSymlink implements fuseutil.FileSystem.
func (*VcsFS) CreateSymlink(context.Context, *fuseops.CreateSymlinkOp) error {
	panic("unimplemented")
}

// Destroy implements fuseutil.FileSystem.
func (*VcsFS) Destroy() {
	panic("unimplemented")
}

// Fallocate implements fuseutil.FileSystem.
func (*VcsFS) Fallocate(context.Context, *fuseops.FallocateOp) error {
	panic("unimplemented")
}

// FlushFile implements fuseutil.FileSystem.
func (*VcsFS) FlushFile(context.Context, *fuseops.FlushFileOp) error {
	panic("unimplemented")
}

// ForgetInode implements fuseutil.FileSystem.
func (*VcsFS) ForgetInode(context.Context, *fuseops.ForgetInodeOp) error {
	panic("unimplemented")
}

// GetXattr implements fuseutil.FileSystem.
func (*VcsFS) GetXattr(context.Context, *fuseops.GetXattrOp) error {
	panic("unimplemented")
}

// ListXattr implements fuseutil.FileSystem.
func (*VcsFS) ListXattr(context.Context, *fuseops.ListXattrOp) error {
	panic("unimplemented")
}

// MkDir implements fuseutil.FileSystem.
func (*VcsFS) MkDir(context.Context, *fuseops.MkDirOp) error {
	return os.Mkdir("bla", 0755)
}

// MkNode implements fuseutil.FileSystem.
func (*VcsFS) MkNode(context.Context, *fuseops.MkNodeOp) error {
	panic("unimplemented")
}

// ReadSymlink implements fuseutil.FileSystem.
func (*VcsFS) ReadSymlink(context.Context, *fuseops.ReadSymlinkOp) error {
	panic("unimplemented")
}

// ReleaseDirHandle implements fuseutil.FileSystem.
func (*VcsFS) ReleaseDirHandle(context.Context, *fuseops.ReleaseDirHandleOp) error {
	panic("unimplemented")
}

// ReleaseFileHandle implements fuseutil.FileSystem.
func (*VcsFS) ReleaseFileHandle(context.Context, *fuseops.ReleaseFileHandleOp) error {
	panic("unimplemented")
}

// RemoveXattr implements fuseutil.FileSystem.
func (*VcsFS) RemoveXattr(context.Context, *fuseops.RemoveXattrOp) error {
	panic("unimplemented")
}

// Rename implements fuseutil.FileSystem.
func (*VcsFS) Rename(context.Context, *fuseops.RenameOp) error {
	panic("unimplemented")
}

// RmDir implements fuseutil.FileSystem.
func (*VcsFS) RmDir(context.Context, *fuseops.RmDirOp) error {
	panic("unimplemented")
}

// SetInodeAttributes implements fuseutil.FileSystem.
func (*VcsFS) SetInodeAttributes(context.Context, *fuseops.SetInodeAttributesOp) error {
	panic("unimplemented")
}

// SetXattr implements fuseutil.FileSystem.
func (*VcsFS) SetXattr(context.Context, *fuseops.SetXattrOp) error {
	panic("unimplemented")
}

// SyncFile implements fuseutil.FileSystem.
func (*VcsFS) SyncFile(context.Context, *fuseops.SyncFileOp) error {
	panic("unimplemented")
}

// Unlink implements fuseutil.FileSystem.
func (*VcsFS) Unlink(context.Context, *fuseops.UnlinkOp) error {
	panic("unimplemented")
}

// WriteFile implements fuseutil.FileSystem.
func (*VcsFS) WriteFile(context.Context, *fuseops.WriteFileOp) error {
	panic("unimplemented")
}

const (
	rootInode fuseops.InodeID = fuseops.RootInodeID + iota
	helloInode
	dirInode
	worldInode
)

type inodeInfo struct {
	attributes fuseops.InodeAttributes

	// File or directory?
	dir bool

	// For directories, children.
	children []fuseutil.Dirent
}

// We have a fixed directory structure.
var gInodeInfo = map[fuseops.InodeID]inodeInfo{
	// root
	rootInode: {
		attributes: fuseops.InodeAttributes{
			Nlink: 1,
			Mode:  0555 | os.ModeDir,
		},
		dir: true,
		children: []fuseutil.Dirent{
			{
				Offset: 1,
				Inode:  helloInode,
				Name:   "hello",
				Type:   fuseutil.DT_File,
			},
			{
				Offset: 2,
				Inode:  dirInode,
				Name:   "dir",
				Type:   fuseutil.DT_Directory,
			},
		},
	},

	// hello
	helloInode: {
		attributes: fuseops.InodeAttributes{
			Nlink: 1,
			Mode:  0444,
			Size:  uint64(len("Hello, world!")),
		},
	},

	// dir
	dirInode: {
		attributes: fuseops.InodeAttributes{
			Nlink: 1,
			Mode:  0555 | os.ModeDir,
		},
		dir: true,
		children: []fuseutil.Dirent{
			{
				Offset: 1,
				Inode:  worldInode,
				Name:   "world",
				Type:   fuseutil.DT_File,
			},
		},
	},
}

func findChildInode(name string, children []fuseutil.Dirent) (fuseops.InodeID, error) {
	for _, e := range children {
		if e.Name == name {
			return e.Inode, nil
		}
	}

	return 0, fuse.ENOENT
}

func (fs *VcsFS) patchAttributes(attr *fuseops.InodeAttributes) {
	now := fs.Clock.Now()
	attr.Atime = now
	attr.Mtime = now
	attr.Crtime = now
}

func (fs *VcsFS) StatFS(ctx context.Context, op *fuseops.StatFSOp) error {
	return nil
}

func (fs *VcsFS) LookUpInode(ctx context.Context, op *fuseops.LookUpInodeOp) error {
	// Find the info for the parent.
	parentInfo, ok := gInodeInfo[op.Parent]
	if !ok {
		return fuse.ENOENT
	}

	// Find the child within the parent.
	childInode, err := findChildInode(op.Name, parentInfo.children)
	if err != nil {
		return err
	}

	// Copy over information.
	op.Entry.Child = childInode
	op.Entry.Attributes = gInodeInfo[childInode].attributes

	// Patch attributes.
	fs.patchAttributes(&op.Entry.Attributes)

	return nil
}

func (fs *VcsFS) GetInodeAttributes(ctx context.Context, op *fuseops.GetInodeAttributesOp) error {
	// Find the info for this inode.
	info, ok := gInodeInfo[op.Inode]
	if !ok {
		return fuse.ENOENT
	}

	// Copy over its attributes.
	op.Attributes = info.attributes

	// Patch attributes.
	fs.patchAttributes(&op.Attributes)

	return nil
}

func (fs *VcsFS) OpenDir(ctx context.Context, op *fuseops.OpenDirOp) error {
	// Allow opening any directory.
	return nil
}

func (fs *VcsFS) ReadDir(ctx context.Context, op *fuseops.ReadDirOp) error {
	// Find the info for this inode.
	info, ok := gInodeInfo[op.Inode]
	if !ok {
		return fuse.ENOENT
	}

	if !info.dir {
		return fuse.EIO
	}

	entries := info.children

	// Grab the range of interest.
	if op.Offset > fuseops.DirOffset(len(entries)) {
		return fuse.EIO
	}

	entries = entries[op.Offset:]

	// Resume at the specified offset into the array.
	for _, e := range entries {
		n := fuseutil.WriteDirent(op.Dst[op.BytesRead:], e)
		if n == 0 {
			break
		}

		op.BytesRead += n
	}

	return nil
}

func (fs *VcsFS) OpenFile(ctx context.Context, op *fuseops.OpenFileOp) error {
	// Allow opening any file.
	return nil
}

func (fs *VcsFS) ReadFile(ctx context.Context, op *fuseops.ReadFileOp) error {
	// Let io.ReaderAt deal with the semantics.
	reader := strings.NewReader("Hello, world!")

	var err error
	op.BytesRead, err = reader.ReadAt(op.Dst, op.Offset)

	// Special case: FUSE doesn't expect us to return io.EOF.
	if err == io.EOF {
		return nil
	}

	return err
}
