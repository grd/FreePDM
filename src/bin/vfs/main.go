// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"log"
	"os"

	vfs "github.com/grd/FreePDM/src/bin/vfs/vcsfs"
	"github.com/jacobsa/fuse"
	"github.com/jacobsa/timeutil"
)

var fMountPoint = flag.String("mount_point", "", "Path to mount point.")
var fReadyFile = flag.Uint64("ready_file", 0, "FD to signal when ready.")

var fDebug = flag.Bool("debug", false, "Enable debug logging.")

func main() {
	flag.Parse()

	// Create an appropriate file system.
	server, err := vfs.NewVFS(timeutil.RealClock())
	if err != nil {
		log.Fatalf("makeFS: %v", err)
	}

	// Mount the file system.
	if *fMountPoint == "" {
		log.Fatalf("You must set --mount_point.")
	}

	cfg := &fuse.MountConfig{
		ReadOnly: false,
	}

	if *fDebug {
		cfg.DebugLogger = log.New(os.Stderr, "fuse: ", 0)
	}

	mfs, err := fuse.Mount(*fMountPoint, server, cfg)
	if err != nil {
		log.Fatalf("Mount: %v", err)
	}

	// Wait for it to be unmounted.
	if err = mfs.Join(context.Background()); err != nil {
		log.Fatalf("Join: %v", err)
	}
}
