package filesystem_test

import (
	"testing"

	fsm "github.com/grd/FreePDM/pkg/filesystem"
)

func TestFileRename(t *testing.T) {
	fs := fsm.FileSystem{
		// Initialize mock FileSystem or dependencies here
	}

	err := fs.FileRename("src.txt", "dest.txt")
	if err != nil {
		t.Errorf("FileRename failed: %v", err)
	}

	// Test cases for locking, existing destination, etc.
	err = fs.FileRename("lockedFile.txt", "newFile.txt")
	if err == nil || err.Error() != "FileRename error: File lockedFile.txt is checked out by user" {
		t.Errorf("Expected lock error, got: %v", err)
	}

	// Further test cases, e.g., file already exists.
}

func TestFileCopy(t *testing.T) {
	fs := fsm.FileSystem{
		// Initialize mock FileSystem or dependencies here
	}

	err := fs.FileCopy("src.txt", "dest.txt")
	if err != nil {
		t.Errorf("FileCopy failed: %v", err)
	}

	// Test locked file
	err = fs.FileCopy("lockedFile.txt", "newFile.txt")
	if err == nil || err.Error() != "FileCopy error: File lockedFile.txt is checked out by user" {
		t.Errorf("Expected lock error, got: %v", err)
	}

	// Test for file already existing at destination
	err = fs.FileCopy("src.txt", "existingFile.txt")
	if err == nil || err.Error() != "file existingFile.txt already exists" {
		t.Errorf("Expected existing file error, got: %v", err)
	}
}

func TestFileMove(t *testing.T) {
	fs := fsm.FileSystem{
		// Initialize mock FileSystem or dependencies here
	}

	err := fs.FileMove("file.txt", "newDir")
	if err != nil {
		t.Errorf("FileMove failed: %v", err)
	}

	// Test moving a locked file
	err = fs.FileMove("lockedFile.txt", "newDir")
	if err == nil || err.Error() != "FileMove error: File lockedFile.txt is checked out by user" {
		t.Errorf("Expected lock error, got: %v", err)
	}

	// Further test cases, e.g., non-existing directory
	err = fs.FileMove("file.txt", "nonExistentDir")
	if err == nil || err.Error() != "directory nonExistentDir doesn't exist" {
		t.Errorf("Expected directory error, got: %v", err)
	}
}

func TestDirectoryCopy(t *testing.T) {
	fs := fsm.FileSystem{
		// Initialize mock FileSystem or dependencies here
	}

	err := fs.DirectoryCopy("srcDir", "destDir")
	if err != nil {
		t.Errorf("DirectoryCopy failed: %v", err)
	}

	// Test cases, e.g., destination is a number
	err = fs.DirectoryCopy("srcDir", "123")
	if err == nil || err.Error() != "directory 123 is a number" {
		t.Errorf("Expected number error, got: %v", err)
	}

	// Test if destination directory already exists
	err = fs.DirectoryCopy("srcDir", "existingDir")
	if err == nil || err.Error() != "directory existingDir exists" {
		t.Errorf("Expected existing directory error, got: %v", err)
	}
}

func TestDirectoryRename(t *testing.T) {
	fs := fsm.FileSystem{
		// Initialize mock FileSystem or dependencies here
	}

	err := fs.DirectoryRename("srcDir", "destDir")
	if err != nil {
		t.Errorf("DirectoryRename failed: %v", err)
	}

	// Test renaming a directory to a number
	err = fs.DirectoryRename("srcDir", "123")
	if err == nil || err.Error() != "directory 123 is a number" {
		t.Errorf("Expected number error, got: %v", err)
	}

	// Test for existing destination directory
	err = fs.DirectoryRename("srcDir", "existingDir")
	if err == nil || err.Error() != "directory existingDir exists" {
		t.Errorf("Expected existing directory error, got: %v", err)
	}
}

func TestDirectoryMove(t *testing.T) {
	fs := fsm.FileSystem{
		// Initialize mock FileSystem or dependencies here
	}

	err := fs.DirectoryMove("srcDir", "destDir")
	if err != nil {
		t.Errorf("DirectoryMove failed: %v", err)
	}

	// Test for destination directory being a number
	err = fs.DirectoryMove("srcDir", "123")
	if err == nil || err.Error() != "destination directory 123 cannot be a number" {
		t.Errorf("Expected number error, got: %v", err)
	}

	// Test for source directory not existing
	err = fs.DirectoryMove("nonExistentDir", "destDir")
	if err == nil || err.Error() != "source directory nonExistentDir does not exist" {
		t.Errorf("Expected source directory error, got: %v", err)
	}
}
