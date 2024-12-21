// Copyright 2024 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	fsm "github.com/grd/FreePDM/internal/vaults"
)

func CommandHandler(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJsonError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println(req)

	// Using the right command
	switch req.Command {
	case "root":
		handleRoot(w)
	case "list":
		handleList(w)
	case "direxists":
		path, ok := req.Params["path"]
		if !ok {
			writeJsonError(w, "Missing parameters", http.StatusBadRequest)
		}
		handleDirexists(w, req.User, req.Vault, path)
	case "ls":
		path, ok := req.Params["path"]
		if !ok {
			writeJsonError(w, "Missing parameters", http.StatusBadRequest)
		}
		handleLs(w, req.User, req.Vault, path)

	// case "rename":
	// 	// Get 'vault', 'src' en 'dst' out of params map
	// 	src := req.Params["src"]
	// 	dst := req.Params["dst"]
	// 	handleRename(w, src, dst)

	default:
		writeJsonError(w, "Unknown command: "+req.Command, http.StatusBadRequest)
	}
}

func writeJsonError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func handleRoot(w http.ResponseWriter) {
	var resp CommandResponse

	root := fsm.Root()
	resp = CommandResponse{
		Data: []string{root},
	}
	json.NewEncoder(w).Encode(resp)
}

// Shows the existing vaults
func handleList(w http.ResponseWriter) {
	var resp CommandResponse

	list, err := fsm.ListVaults()
	if err != nil {
		resp = CommandResponse{
			Error: "Failed to show the list of vaults",
		}
	} else {
		resp = CommandResponse{
			Data: list,
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func handleDirexists(w http.ResponseWriter, user, vault, path string) {
	var resp CommandResponse

	fs, err := fsm.NewFileSystem(vault, user)
	if err != nil {
		log.Fatalf("unable to access the filesystem : %s", err)
	}

	ok := fs.DirExists(path)
	if !ok {
		resp = CommandResponse{
			Error: "Directory " + path + " does not exists",
		}
	} else {
		resp = CommandResponse{
			Error: "Directory " + path + " exists",
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func handleLs(w http.ResponseWriter, user, vault, path string) {
	var resp CommandResponse

	fs, err := fsm.NewFileSystem(vault, user)
	if err != nil {
		log.Fatalf("unable to access the filesystem : %s", err)
	}

	fmt.Printf("path = %s\n", path)

	list, err := fs.ListDir(path)
	if err != nil {
		resp = CommandResponse{
			Error: "Failed to show directory",
		}
	} else {
		files := make([]string, len(list))
		for i, item := range list {
			files[i] = item.Name()
		}
		resp = CommandResponse{
			Data: files,
		}
	}
	json.NewEncoder(w).Encode(resp)
}

// func handleRename(w http.ResponseWriter, src, dst string) {

// 	// Rename the file
// 	if err := os.Rename(src, dst); err != nil {
// 		resp := CommandResponse{
// 			Status:  "error",
// 			Message: "Failed to rename file",
// 			Data:    map[string]interface{}{"error": err.Error()},
// 		}
// 		json.NewEncoder(w).Encode(resp)
// 		return
// 	}

// 	// Send success response
// 	resp := CommandResponse{
// 		Status:  "success",
// 		Message: "File renamed successfully",
// 	}
// 	json.NewEncoder(w).Encode(resp)
// }
