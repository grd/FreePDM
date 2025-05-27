// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	s.ExecuteTemplate(w, "index.html", nil)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.HandleLogin(w, r)
	} else if r.Method == http.MethodGet {
		s.ServeLoginPage(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	s.ExecuteTemplate(w, "login.html", nil)
}

// login handler
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	loginname := r.FormValue("loginname")
	password := r.FormValue("password")

	user, err := s.UserRepo.LoadUser(loginname)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		s.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Invalid login credentials"})
		return
	}

	shared.SetSessionCookie(w, loginname)

	// Explicit debug to stdout
	fmt.Printf("[DEBUG] User %s roles: %v\n", loginname, user.Roles)
	log.Printf("[DEBUG] User %s roles: %v", loginname, user.Roles)

	// EXTRA CHECK: show exact roles to be sure
	for i, role := range user.Roles {
		fmt.Printf("[DEBUG] Role %d: %s\n", i, role)
	}

	// Robust admin detection (case-insensitive match)
	for _, role := range user.Roles {
		if strings.EqualFold(role, "Admin") {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
	}

	// Continue with other priority roles
	if auth.IsProjectLead(user) {
		http.Redirect(w, r, "/project-lead-dashboard", http.StatusSeeOther)
		return
	} else if auth.IsSeniorDesigner(user) {
		http.Redirect(w, r, "/senior-designer-dashboard", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (s *Server) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.ExecuteTemplate(w, "change-password.html", nil)
	}

	loginname, err := shared.GetSessionLoginname(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Form data
	oldPassword := r.FormValue("old_password")
	newPassword := r.FormValue("new_password")
	repeatPassword := r.FormValue("repeat_password")

	user, err := s.UserRepo.LoadUser(loginname)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)) != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Current password is incorrect")
		return
	}

	// Basic checks
	if newPassword != repeatPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}
	if len(newPassword) < 10 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Password must be at least 10 characters long")
		return
	}
	if !containsUppercase(newPassword) {
		http.Error(w, "Password must contain at least one uppercase letter", http.StatusBadRequest)
		return
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// Update DB
	err = s.UserRepo.UpdatePassword(loginname, string(hash))
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Set MustChangePassword = false
	err = s.UserRepo.ClearMustChangePassword(loginname)
	if err != nil {
		http.Error(w, "Failed to finalize update", http.StatusInternalServerError)
		return
	}

	// Return redirect via HX
	w.Header().Set("HX-Redirect", "/dashboard")
}

func containsUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	shared.ClearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	loginname, err := shared.GetSessionLoginname(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := struct {
		LoginName string
	}{
		LoginName: loginname,
	}

	s.ExecuteTemplate(w, "dashboard.html", data)
}
