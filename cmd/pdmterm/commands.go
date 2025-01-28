// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ANSI escape codes as constants
const (
	// Reset
	Reset = "\033[0m"

	// Regular Colors
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright Colors
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"
)

var (
	currentVault = "" // Placeholder for the current vault
	currentDir   = "" // Initial directory
	user         = "" // Login name
)

// handleCommand processes the input command and executes corresponding actions.
func handleCommand(input string, directory string) {
	// Split the command and arguments
	parts := strings.Fields(input)
	if len(parts) == 0 {
		fmt.Println(Red + "No command entered." + Reset)
		return
	}

	// command usage
	command := parts[0]
	args := parts[1:]

	switch command {
	case "help":
		handleHelp()
	case "list":
		handleList()
	case "vault":
		handleVault(args[0])
	case "tree":
		handleTree(args[0])
	case "pwd":
		handlePwd()
	case "ls":
		handleLs(directory)
	case "cd":
		if len(args) < 1 {
			fmt.Println(Cyan + "Usage: cd <directory>" + Reset)
			return
		}
		handleCd(args[0])
	case "allocate": // returns the ID
		handleAllocate()
	// Not yet...
	// case "assign":
	// 	handleAssign()
	// case "rm":
	// 	handleFileRemove()
	case "exit", "quit":
		fmt.Println(Cyan + "Exiting the shell." + Reset)
		os.Exit(0)
	default:
		fmt.Printf(Red+"Unknown command: %s\n"+Reset, command)
	}
}

func handleHelp() {
	message := `
Welcome to the FreePDM help!

Commands available:
- help                       : Show help
- list                       : Show the list of vaults
- vault <name>               : Activate a vault (shows in the prompt)
- pwd                        : Show current directory (shows in the prompt)
- ls                         : List files in the current directory
- tree                       : Shows a tree of files and directories from the pwd
- cd <dir>                   : Change to a different directory
- mkdir <dir>                : Create a directory
- rmdir <dir>                : Remove an empty directory
- import <file>              : Import a file in the current directory of the vault
- allocate                   : Allocates a container to store a file from FreeCAD
- assign <cont nr> <file>    : Assignes the file to the container number
- rm <file>                  : Removes a file (container) from the current directory
- mv <src> <dst>             : Move a file. Move file between vaults is not yet available
- rename <src> <dst>         : Rename a file. Rename file between vaults is not yet available
- copy <src> <dst>           : Copy a file. Copy file between vaults is not yet available
- versions <file>            : Returns the number of versions
- newversion <file>          : Creates a new version of a file and check out
- checkout <file> <version>  : Checks out a file. No-one but you can modify it
- checkin <file> <version>   : Check in a file
- info <file> <version>      : Returns the parameters of a file. If no version show the latest
- exit                       : Quit the program
`
	fmt.Println(message)
}

// handleList the list of vaults.
func handleList() {
	resp, err := sendCommand("list", nil)
	if err != nil {
		fmt.Println(err)
	}
	if resp.Error == "Failed to show the list of vaults" {
		fmt.Println(Red + resp.Error + Reset)
	} else {
		fmt.Printf(Cyan+"Data = %s\n"+Reset, resp.Data)
	}
}

// handleVault changes the current vault.
func handleVault(vault string) {
	resp, err := sendCommand("list", nil)
	if err != nil {
		fmt.Println(err)
	}
	if resp.Error == "Failed to show the list of vaults" {
		fmt.Println(Cyan + resp.Error + Reset)
		return
	}

	for _, elem := range resp.Data {
		if vault == elem {
			currentVault = elem
			currentDir = "."
			break
		}
	}
	fmt.Printf(BrightGreen+"Switched to vault: %s\n"+Reset, vault)
}

// handleTree changes the current vault.
func handleTree(path string) {
	resp, err := sendCommand("tree", map[string]string{
		"path": path,
	})
	if err != nil {
		fmt.Println(err)
	}
	if resp.Error == "Failed to show the tree" {
		fmt.Println(Cyan + resp.Error + Reset)
		return
	}

	for _, elem := range resp.Data {
		fmt.Println(elem)
	}
}

// handleLs lists files and directories in the current directory.
func handleLs(directory string) {
	resp, err := sendCommand("ls", map[string]string{
		"path": directory,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Error = %s\n", resp.Error)
	fmt.Printf("Data = %s\n", resp.Data)
}

// handleCd changes the current working directory.
func handleCd(target string) {
	fmt.Printf("dir = %s\n", target)
	if currentVault == "" {
		fmt.Println(Red + "First set the vault with the command vault" + Reset)
		return
	}
	resp, err := sendCommand("direxists", map[string]string{
		"path": target,
	})
	if err != nil {
		fmt.Println(err)
	}
	if resp.Error != "" {
		fmt.Printf(Red+"directory %s does not exist\n"+Reset, target)
		return
	}
	if target == ".." { // Up one level
		if target != "/" {
			lastSlash := strings.LastIndex(target, "/")
			target = (target)[:lastSlash]
			if target == "" {
				target = "/"
			}
		}
	} else { // Subdirectory
		if target == "/" {
			target += target
		} else {
			target += "/" + target
		}
	}
	fmt.Printf(Cyan+"Changed directory to: %s\n", target)
}

// handlePwd prints the current working directory.
func handlePwd() {
	fmt.Println("Current directory:", currentDir)
}

// allocates a new container inside the current vault and path.
// Returns the container number.
func handleAllocate() {
	if currentVault == "" {
		fmt.Println(Red + "First set the vault with the command vault" + Reset)
		return
	}

	resp, err := sendCommand("allocate", map[string]string{
		"vault": currentVault,
		"path":  currentDir,
	})
	if err != nil {
		fmt.Println(err)
	}
	if resp.Error == "Failed to allocate a container" {
		fmt.Println(Red + resp.Error + Reset)
	}

	if len(resp.Data) > 0 {
		fmt.Printf(Cyan+"Container number = %s\n"+Reset, resp.Data[0])
	} else {
		fmt.Println(Red + "A serious error..." + Reset)
	}
}

func newPrompt() {
	user = os.Getenv("USER")

	fmt.Println("Welcome to the FreePDM CLI!")
	fmt.Println("If you need any help, type help.")
	fmt.Println("")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Show prompt
		fmt.Printf(BrightCyan+"%s:%s>"+Reset, currentVault, currentDir)

		// Read input
		if scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			handleCommand(input, currentDir)
		} else {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf(Red+"Error reading input: %v\n", err)
	}
}
