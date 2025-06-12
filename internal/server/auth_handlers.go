// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"
	"unicode"

	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		BackButtonShow bool
		MenuButtonShow bool
	}{
		BackButtonShow: false,
		MenuButtonShow: false,
	}
	if err := s.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
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
	if r.Method == http.MethodGet {
		s.ExecuteTemplate(w, "login.html", nil)
		return
	}

	loginName := r.FormValue("login_name")
	password := r.FormValue("password")

	user, err := s.UserRepo.LoadUserByLoginName(loginName)
	if err != nil || user == nil {
		http.Error(w, "Invalid login name", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	sess, _ := s.SessionStore.Get(r, "pdm-session")
	sess.Values["user"] = loginName
	sess.Save(r, w)

	if user.HasRole("Admin") {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
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
