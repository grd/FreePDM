// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Request struct {
	Command string `json:"command"`
	Path    string `json:"path"`
}

type Response struct {
	Success bool     `json:"success"`
	Data    []string `json:"data,omitempty"`
	Error   string   `json:"error,omitempty"`
}

func lsHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Command != "ls" {
		http.Error(w, "Unsupported command", http.StatusBadRequest)
		return
	}

	files, err := os.ReadDir(req.Path)
	if err != nil {
		resp := Response{Success: false, Error: err.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	resp := Response{Success: true, Data: filenames}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/command", lsHandler)
	fmt.Println("Server running on :3000")
	http.ListenAndServe(":3000", nil)
}
