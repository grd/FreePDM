// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/shared"
)

func (s *Server) HandleProjectManagement(w http.ResponseWriter, r *http.Request) {
	username, err := shared.GetSessionUsername(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := s.UserRepo.LoadUser(username)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Autorisationcheck
	if !auth.HasAnyRole(user, "Admin", "ProjectLead") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	s.ExecuteTemplate(w, "project-management.html", user)
}
