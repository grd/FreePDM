// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"
)

func (s *Server) handlePdm(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "pdm")
}
