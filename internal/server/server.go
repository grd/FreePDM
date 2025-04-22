// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/grd/FreePDM/internal/db"
)

type Server struct {
	UserRepo *db.UserRepo
	// TODO: Add things such as Logger, Config etc.
}

// Constructor
func NewServer(userRepo *db.UserRepo) *Server {
	return &Server{
		UserRepo: userRepo,
	}
}

func (s *Server) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/register", s.handleRegister)
	mux.HandleFunc("/login", s.HandleLogin)
	mux.HandleFunc("/logout", s.handleLogout)
	mux.HandleFunc("/pdm", s.handlePdm)
	mux.HandleFunc("/handler", s.handler)
	mux.HandleFunc("/welcome", s.welcomeHandler)
}
