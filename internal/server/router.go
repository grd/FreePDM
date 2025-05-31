// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"
	"path/filepath"

	"github.com/grd/FreePDM/internal/config"
	"github.com/grd/FreePDM/internal/middleware"
)

func (s *Server) Routes(mux *http.ServeMux) {
	staticPath := filepath.Join(config.AppDir(), "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	// Without auth
	mux.HandleFunc("/", s.HandleHomePage)
	mux.HandleFunc("/login", s.Login)
	mux.HandleFunc("/change-password", s.HandleChangePassword)

	// With auth
	mux.HandleFunc("/admin/logs", middleware.RequireRoleWithLogin(s.HandleShowLogs, "Admin"))
	mux.HandleFunc("/admin/users", middleware.RequireRoleWithLogin(s.HandleAdminUsers, "Admin"))
	mux.HandleFunc("/admin/users/new", middleware.RequireAdmin(s.HandleAdminNewUser))
	mux.HandleFunc("/admin/users/edit/", middleware.RequireAdmin(s.HandleAdminEditUser))
	mux.HandleFunc("/admin/users/reset-password/", middleware.RequireAdmin(s.HandleAdminResetPassword))

	mux.HandleFunc("/admin/users/upload-photo", middleware.RequireAdmin(s.HandleUploadPhoto))

	mux.HandleFunc("/dashboard", middleware.RequireLogin(s.HandleDashboard))
	mux.HandleFunc("/projects/manage", middleware.RequireRoleWithLogin(s.HandleProjectManagement, "Admin", "ProjectLead"))
	mux.HandleFunc("/admin", middleware.RequireRoleWithLogin(s.HandleAdminDashboard, "Admin"))
	mux.HandleFunc("/logout", middleware.RequireLogin(s.handleLogout))
	mux.HandleFunc("/pdm", s.handlePdm)
}
