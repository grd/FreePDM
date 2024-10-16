package filesystem_test

import (
	"testing"

	fsm "github.com/grd/FreePDM/pkg/filesystem"
)

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
