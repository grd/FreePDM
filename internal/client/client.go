// Copyright 2024 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grd/FreePDM/internal/shared"
)

func sendCommand(command string, params map[string]string) (*shared.CommandResponse, error) {
	req := shared.CommandRequest{
		User:    user,
		Vault:   currentVault,
		Command: command,
		Params:  params,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	fmt.Println(BrightBlue + string(data) + Reset)

	resp, err := http.Post("http://localhost:8080/command", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check StatusCode
	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("invalid response from server: %w", err)
		}
		return nil, fmt.Errorf("server error: %s", errResp["error"])
	}

	// Decode JSON reply
	var cmdResp shared.CommandResponse
	if err := json.NewDecoder(resp.Body).Decode(&cmdResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &cmdResp, nil
}
