// Copyright 2024 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFile    *os.File   // Global variable to hold the current log file
	logMutex   sync.Mutex // Mutex to ensure safe access to log file
	freePdmDir string
)

// startLogging initializes and manages the log file
func startLogging() {
	// Create the first log file based on the current date
	createLogFile()

	// Start a goroutine to handle daily log file rotation
	go func() {
		for {
			time.Sleep(1 * time.Minute) // Check every minute

			// Rotate log file if the day has changed
			if isNewDay() {
				logMutex.Lock()
				rotateLogFile()
				logMutex.Unlock()
			}
		}
	}()
}

// createLogFile creates a new log file based on the current date
func createLogFile() {
	logMutex.Lock() // Lock to ensure safe file creation
	defer logMutex.Unlock()

	// Close the existing log file, if any
	if logFile != nil {
		logFile.Close()
	}

	// Determine the filename based on the current date
	today := time.Now().Format("2006-01-02") // Format: YYYY-MM-DD
	logFileName := path.Join(freePdmDir, fmt.Sprintf("logs/%s.log", today))

	fmt.Printf("logFileName = %s\n", logFileName)
	// Open or create the log file
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	// Set the log output to the current log file
	log.SetOutput(logFile)

	// Clean up old log files
	cleanupOldLogs("logs", 7) // Retain log files for the last 7 days
}

// isNewDay checks if the date has changed since the last log file creation
func isNewDay() bool {
	currentDate := time.Now().Format("2006-01-02")                // Current date
	logFileDate := logFile.Name()[len("logs/") : len("logs/")+10] // Extract date from the file name
	return currentDate != logFileDate                             // Return true if dates differ
}

// rotateLogFile closes the current log file and creates a new one
func rotateLogFile() {
	log.Println("Rotating log file...")  // Log the rotation
	createLogFile()                      // Create a new log file
	log.Println("New log file created.") // Log the success of the new file
}

// cleanupOldLogs removes log files older than the specified number of days
func cleanupOldLogs(dir string, days int) {
	// Determine the cutoff date for deleting old files
	cutoff := time.Now().AddDate(0, 0, -days)

	// Walk through the directory and remove outdated files
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && info.ModTime().Before(cutoff) {
			os.Remove(path) // Delete the file
		}
		return nil
	})
}

func init() {
	var ok bool
	freePdmDir, ok = os.LookupEnv("FREEPDM_DIR")
	if !ok {
		log.Fatal("The environment FREEPDM_DIR is not set.")
	}

}
