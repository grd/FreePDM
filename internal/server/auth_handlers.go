// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func renderTemplate(w http.ResponseWriter, tmpl string) {
	content, err := os.ReadFile("templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	t, err := template.New(tmpl).Parse(string(content))
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register")
}

func (s *Server) requireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("PDM_Session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if !auth.IsValidSession(cookie.Value, s.UserRepo) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := s.UserRepo.LoadUser(username)
	if err != nil {
		if err == db.ErrUserNotFound {
			http.Error(w, "Invalid login", http.StatusUnauthorized)
			return

		}
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		http.Error(w, "Invalid login", http.StatusUnauthorized)
		return
	}

	// login success → set cookie
	cookie := http.Cookie{
		Name:     "PDM_Session",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // True with HTTPS
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "admin" && password == "password" {
		renderTemplate(w, "dashboard")
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
	fmt.Println("Endpoint Hit: homePage")
}
