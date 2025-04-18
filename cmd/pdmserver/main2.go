// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/server"
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

func handler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register")
}

func requireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("PDM_Session")
		if err != nil || !isValidUser(cookie.Value) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := db.LoadUser(username)
	if err != nil {
		if err == db.ErrUserNotFound {
			http.Error(w, "Invalid login", http.StatusUnauthorized)
			return

		}
	}
	if !err || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
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

// func handleLogin(w http.ResponseWriter, r *http.Request) {
// 	if err := r.ParseForm(); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	username := r.FormValue("admin")
// 	password := r.FormValue("password")

// 	if username == "admin" && password == "password" {
// 		renderTemplate(w, "login")
// 	} else {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 	}
// }

func handleLogout(w http.ResponseWriter, r *http.Request) {
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

func handlePdm(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "pdm")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
	fmt.Println("Endpoint Hit: homePage")
}

func RunServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/register", handleRegister)
	mux.HandleFunc("/login", handleLogin)
	mux.HandleFunc("/logout", handleLogout)
	mux.HandleFunc("/pdm", handlePdm)
	mux.HandleFunc("/handler", handler)
	mux.HandleFunc("/welcome", server.WelcomeHandler)

	log.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
