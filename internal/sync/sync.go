// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"os/exec"
	"time"
)

// Optional: check if rsync exists
func HasRsync(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "rsync", "--version").Run()
}

// Build args without explicit -e ssh
func buildArgs(source, dest string, extra []string, dryRun bool) []string {
	args := []string{
		"-a",       // archive
		"--delete", // sync deletions
		"--partial",
		"--human-readable",
	}
	if dryRun {
		args = append(args, "-n", "--list-only")
	}
	args = append(args, extra...)
	if dryRun {
		// only source required for list-only
		args = append(args, source)
		return args
	}
	// trailing slash semantics or rsync:
	// source/ => substance; source => directory only
	args = append(args, source, dest)
	return args
}

func TestTarget(ctx context.Context, source string, extra []string) error {
	args := buildArgs(source, "", extra, true)
	cmd := exec.CommandContext(ctx, "rsync", args...)
	return cmd.Run()
}

func Sync(ctx context.Context, source, dest string, extra []string) error {
	args := buildArgs(source, dest, extra, false)
	cmd := exec.CommandContext(ctx, "rsync", args...)
	return cmd.Run()
}
