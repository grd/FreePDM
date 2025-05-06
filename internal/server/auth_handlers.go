// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

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
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println("template error:", err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := s.UserRepo.LoadUser(username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		// Template opnieuw renderen met foutmelding
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Println("Template error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		data := map[string]string{
			"Error": "Invalid login credentials",
		}
		tmpl.Execute(w, data)
		return
	}

	// Login success → cookie en redirect
	cookie := http.Cookie{
		Name:     "PDM_Session",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)

	if user.MustChangePassword {
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (s *Server) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	// Check cookie always, even on GET
	cookie, err := r.Cookie("PDM_Session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		log.Println("Rendering change-password page for user:", cookie.Value)
		tmpl, err := template.ParseFiles("templates/change-password.html")
		if err != nil {
			log.Println("Template error:", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		tmpl.Execute(w, map[string]string{
			"Username": cookie.Value,
		})
		return
	}

	// POST
	username := cookie.Value
	newPassword := r.FormValue("new_password")

	hashed, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	err = s.UserRepo.UpdatePassword(username, string(hashed))
	if err != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Remove session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "PDM_Session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
}

func (s *Server) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("PDM_Session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Loading HTML template
	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, "Error loading dashboard", http.StatusInternalServerError)
		return
	}

	// Returns template (or data)
	tmpl.Execute(w, nil)
}
