// Copyright 2024 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shared

type CommandRequest struct {
	User    string            `json:"user"`
	Vault   string            `json:"vault,omitempty"`
	Command string            `json:"command"`
	Params  map[string]string `json:"params,omitempty"`
}

type CommandResponse struct {
	Status  string      `json:"status"` // e.g., "success" or "error"
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
