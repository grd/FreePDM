// Copyright 2024 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	fsm "github.com/grd/FreePDM/pkg/filesystem"
)

func CommandHandler(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJsonError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println(req.Params)

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
			return
		}
		handleDirexists(w, req.User, req.Vault, path)
	case "ls":
		path, ok := req.Params["path"]
		if !ok {
			writeJsonError(w, "Missing parameters", http.StatusBadRequest)
			return
		}
		// user := req.Params["user"]
		// vault := req.Params["vault"]
		// path := req.Params["path"]

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

	// Logic to list directory contents

	root := fsm.Root()
	resp = CommandResponse{
		Status:  "success",
		Message: "Root directory",
		Data:    map[string]interface{}{"items": root},
	}
	json.NewEncoder(w).Encode(resp)
}

func handleList(w http.ResponseWriter) {
	var resp CommandResponse

	// Logic to list directory contents

	list, err := fsm.ListVaults()
	if err != nil {
		resp = CommandResponse{
			Status:  "error",
			Message: "Failed to show the list of vaults",
			Data:    map[string]interface{}{},
		}
	} else {
		resp = CommandResponse{
			Status:  "success",
			Message: "Directory listed successfully",
			Data:    map[string]interface{}{"items": list},
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
			Status:  "error",
			Message: "Directory " + path + " does not exists",
			Data:    nil,
		}
	} else {
		resp = CommandResponse{
			Status:  "success",
			Message: "Directory " + path + " exists",
			Data:    nil,
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func handleLs(w http.ResponseWriter, user, vault, path string) {
	var resp CommandResponse

	// Logic to list directory contents

	fs, err := fsm.NewFileSystem(vault, user)
	if err != nil {
		log.Fatalf("unable to access the filesystem : %s", err)
	}

	list, err := fs.ListDir(path)
	if err != nil {
		resp = CommandResponse{
			Status:  "error",
			Message: "Failed to show directory",
			Data:    map[string]interface{}{},
		}
	} else {
		files := make([]string, len(list))
		for i, item := range list {
			files[i] = item.Name()
		}
		resp = CommandResponse{
			Status:  "success",
			Message: "Directory listed successfully",
			Data:    map[string]interface{}{"items": files},
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
