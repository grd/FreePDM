// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/grd/FreePDM/internal/db"
)

func (s *Server) handlePdm(w http.ResponseWriter, r *http.Request) {
}

// Sample CreateDocumentHandler. It is not defined yet, but only to show the RBAC features
func (s *Server) CreateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequest(r, s.UserRepo)
	if user == nil || !user.HasPermission(db.CreateDocument) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// ... logic to create a document ...
}

func getUserFromRequest(r *http.Request, repo *db.UserRepo) *db.PdmUser {
	cookie, err := r.Cookie("PDM_Session")
	if err != nil {
		return nil
	}

	username := cookie.Value
	user, err := repo.LoadUser(username)
	if err != nil {
		return nil
	}

	return user
}
