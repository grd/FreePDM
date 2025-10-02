// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/grd/FreePDM/internal/config"
	"github.com/grd/FreePDM/internal/db"
	vfs "github.com/grd/FreePDM/internal/vault/localfs"
)

type Server struct {
	UserRepo     *db.UserRepo
	Templates    *template.Template
	SessionStore *sessions.CookieStore
	FS           *vfs.FileSystem

	// TODO: Add things such as Logger, Config etc.
}

var sessionKey = []byte("your-secret-session-key")

// Constructor
func NewServer(userRepo *db.UserRepo) *Server {
	templatePath := filepath.Join(config.AppDir(), "templates", "*.html")
	templates := template.Must(template.ParseGlob(templatePath))

	return &Server{
		UserRepo:     userRepo,
		Templates:    templates,
		SessionStore: sessions.NewCookieStore(sessionKey),
	}
}

func (s *Server) ExecuteTemplate(w http.ResponseWriter, name string, data any) error {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/"+name)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return fmt.Errorf("template %s not found", name)
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("[ERROR] Executing template %s: %v", name, err)
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		return err
	}

	return nil
}
