// Copyright 2024 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func main() {
	url := "http://localhost:3000/command"

	// Voorbeeld: ls uitvoeren in de root directory
	req := Request{Command: "ls", Path: "/"}
	reqBody, _ := json.Marshal(req)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	if result.Success {
		fmt.Println("Files:", result.Data)
	} else {
		fmt.Println("Error:", result.Error)
	}
}
