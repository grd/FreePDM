// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"log"
	"net/http"

	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/shared"
)

func (s *Server) HandleAdminUsers(w http.ResponseWriter, r *http.Request) {
	username, err := shared.GetSessionUsername(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := s.UserRepo.LoadUser(username)
	if err != nil || !auth.IsAdmin(user) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		http.Error(w, "Error loading users", http.StatusInternalServerError)
		return
	}

	data := struct {
		User           *db.PdmUser
		ShowBackButton bool
		BackButtonLink string
		Users          []db.PdmUser
	}{
		User:           user,
		ShowBackButton: true,
		BackButtonLink: "/admin",
		Users:          users,
	}

	for _, u := range users {
		log.Printf("[DEBUG] User: ID=%d, UserName=%s, FullName=%s, PhotoPath=%s", u.ID, u.UserName, u.FullName, u.PhotoPath)
	}

	s.ExecuteTemplate(w, "admin-users.html", data)
}
