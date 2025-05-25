// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"
	"path/filepath"

	"github.com/grd/FreePDM/internal/middleware"
	"github.com/grd/FreePDM/internal/server/config"
)

func (s *Server) Routes(mux *http.ServeMux) {
	staticPath := filepath.Join(config.AppDir, "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	// Without auth
	mux.HandleFunc("/", s.HandleHomePage)
	mux.HandleFunc("/login", s.Login)
	mux.HandleFunc("/change-password", s.HandleChangePassword)

	// With auth
	mux.HandleFunc("/admin", middleware.RequireRoleWithLogin(*s.UserRepo, s.HandleAdminDashboard, "Admin"))
	mux.HandleFunc("/admin/logs", middleware.RequireRoleWithLogin(*s.UserRepo, s.HandleShowLogs, "Admin"))

	mux.HandleFunc("/dashboard", middleware.RequireLogin(*s.UserRepo, s.HandleDashboard))
	mux.HandleFunc("/projects/manage", middleware.RequireRoleWithLogin(*s.UserRepo, s.HandleProjectManagement, "Admin", "ProjectLead"))
	// mux.HandleFunc("/vaults/manage", middleware.RequireRoleWithLogin(*s.UserRepo, s.HandleVaultManagement, "Admin", "SeniorDesigner"))
	mux.HandleFunc("/logout", middleware.RequireLogin(*s.UserRepo, s.handleLogout))
	mux.HandleFunc("/pdm", s.handlePdm)
}
