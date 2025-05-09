// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/grd/FreePDM/internal/db"
)

type Server struct {
	UserRepo  *db.UserRepo
	Templates *template.Template

	// TODO: Add things such as Logger, Config etc.
}

// Constructor
func NewServer(userRepo *db.UserRepo) *Server {
	// templates := template.Must(template.ParseFiles(
	// 	"templates/base.html",
	// 	"templates/dashboard.html",
	// 	"templates/login.html",
	// 	"templates/index.html",
	// 	"templates/change-password.html",
	// ))

	templates := template.Must(template.ParseGlob("templates/*.html"))

	return &Server{
		UserRepo:  userRepo,
		Templates: templates,
	}
}

// func (s *Server) ExecuteTemplate(w http.ResponseWriter, name string, data any) {
// 	err := s.Templates.ExecuteTemplate(w, name, data)
// 	if err != nil {
// 		log.Println("Template execution error:", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 	}
// }

func (s *Server) ExecuteTemplate(w http.ResponseWriter, name string, data any) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/"+name)
	if err != nil {
		log.Println("Template parse error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println("Template execution error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
