// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"log"
	"net/http"
	"unicode"

	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/db"
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

	log.Printf("[DEBUG] Login attempt for user: %s", loginname)

	user, err := s.UserRepo.LoadUser(loginname)
	if err != nil {
		log.Printf("[DEBUG] User %s not found or error: %v", loginname, err)
		s.ExecuteTemplate(w, "login.html", map[string]string{
			"Error": "Invalid login credentials",
		})
		return
	}

	// Check password
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		log.Printf("[DEBUG] Invalid password for user: %s", loginname)
		s.ExecuteTemplate(w, "login.html", map[string]string{
			"Error": "Invalid login credentials",
		})
		return
	}

	// Login success → set session cookie
	shared.SetSessionCookie(w, loginname)
	log.Printf("[DEBUG] User %s logged in successfully", loginname)

	// Check if user must change password
	if user.MustChangePassword {
		log.Printf("[DEBUG] User %s must change password — redirecting to /change-password", loginname)
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	// Role-based redirection with priority
	switch {
	case user.HasRole(string(db.Admin)):
		log.Printf("[DEBUG] User %s has role Admin — redirecting to /admin", loginname)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	case auth.IsProjectLead(user):
		log.Printf("[DEBUG] User %s has role ProjectLead — redirecting to /project-lead-dashboard", loginname)
		http.Redirect(w, r, "/project-lead-dashboard", http.StatusSeeOther)
	case auth.IsSeniorDesigner(user):
		log.Printf("[DEBUG] User %s has role SeniorDesigner — redirecting to /senior-designer-dashboard", loginname)
		http.Redirect(w, r, "/senior-designer-dashboard", http.StatusSeeOther)
	default:
		log.Printf("[DEBUG] User %s has no special roles — redirecting to /dashboard", loginname)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
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
