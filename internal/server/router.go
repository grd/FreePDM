// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/grd/FreePDM/internal/middleware"
)

func (s *Server) Routes(mux *http.ServeMux) {
	// Without auth
	mux.HandleFunc("/", s.HandleHomePage)
	mux.HandleFunc("/login", s.Login)
	mux.HandleFunc("/change-password", s.HandleChangePassword)

	// With auth
	mux.HandleFunc("/dashboard", middleware.RequireLogin(*s.UserRepo, s.HandleDashboard))
	mux.HandleFunc("/projects/manage", middleware.RequireRoleWithLogin(*s.UserRepo, s.HandleProjectManagement, "Admin", "ProjectLead"))
	// mux.HandleFunc("/vaults/manage", middleware.RequireRoleWithLogin(*s.UserRepo, s.HandleVaultManagement, "Admin", "SeniorDesigner"))
	mux.HandleFunc("/logout", middleware.RequireLogin(*s.UserRepo, s.handleLogout))
	mux.HandleFunc("/pdm", s.handlePdm)
}
