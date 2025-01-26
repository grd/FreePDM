// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// This is a helper module of arbitrary functions.

// Convert a string to a slice of strings
func StringToSlice(input string) []string {
	return []string{input}
}

// today returns the date of today in the format "YYYY-MM-DD".
func Today() string {
	t := time.Now().String()
	return t[0:10]
}

// Now returns the time of now in the format "YYYY-MM-DD HH:MM:SS"
func Now() string {
	currentTime := time.Now()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d",
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second())
}

// Fatal when having an error
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Check whether a critical file exists.
func CriticalFileExist(fileName string) {
	_, err := os.Stat(fileName)
	if err != nil {
		log.Fatalf("%s file doesn't exist\n", fileName)
	}
}

// Check whether a critical directory exists.
func CriticalDirExist(dirName string) {
	_, err := os.Stat(dirName)
	if err != nil {
		log.Fatalf("%s directory does not exist\n", dirName)
	}
}

// Check whether a directory exists. Returns true when the directory exists, otherwise false.
func DirExists(dirName string) bool {
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Check whether a file exists. Returns true when the file exists, otherwise false.
func FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Copies a file from src to dest with file permission 0666.
// Returns an error.
func CopyFile(src, dest string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	return nil
}

// Converts a string into an int16
func Atoi16(str string) (int16, error) {
	var num int16
	if _, err := fmt.Sscanf(str, "%d", &num); err != nil {
		return 0, err
	}
	return num, nil
}

// Converts an int16 into a string
func I16toa(i int16) string {
	return fmt.Sprintf("%d", i)
}

// Converts a string into an int64
func Atoi64(str string) (int64, error) {
	var num int64
	if _, err := fmt.Sscanf(str, "%d", &num); err != nil {
		return 0, err
	}
	return num, nil
}

// Converts an int64 into a string
func I64toa(i int64) string {
	return fmt.Sprintf("%d", i)
}

func SplitFileExtension(file string) (base, ext string) {
	ext = filepath.Ext(file)
	base = file[:len(file)-len(ext)]
	return
}

// Returns true when the string is a number.
func IsNumber(name string) bool {
	if _, err := strconv.Atoi(name); err == nil {
		return true
	} else {
		return false
	}
}
