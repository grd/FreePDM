// Copyright 2024 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shared

import (
	"encoding/json"
)

type CommandRequest struct {
	User    string            `json:"user"`
	Vault   string            `json:"vault,omitempty"`
	Command string            `json:"command"`
	Params  map[string]string `json:"params,omitempty"`
}

type CommandResponse struct {
	Error string   `json:"error,omitempty"`
	Data  []string `json:"data,omitempty"`
}

func (req CommandRequest) String() string {
	str, _ := json.Marshal(req)
	return string(str)
}

func (res CommandResponse) String() string {
	str, _ := json.Marshal(res)
	return string(str)
}
