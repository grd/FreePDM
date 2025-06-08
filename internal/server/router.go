// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) Routes(mux *http.ServeMux) {
	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)

	// ✅ Public routes (géén login nodig!)
	r.Group(func(r chi.Router) {
		r.Get("/", s.HandleHomePage)
		r.Get("/login", s.HandleLoginPage)
		r.Post("/login", s.HandleLogin)
		r.Get("/logout", s.HandleLogout)
	})

	// ✅ Alle routes hieronder vereisen login
	r.Group(func(r chi.Router) {
		r.Use(s.RequireLoginChi)

		// ✅ Dashboard
		r.With(s.RequireRoleChi("Admin")).Get("/admin", s.HandleAdminDashboard)
		r.Get("/dashboard", s.HandleDashboard)

		// ✅ Logs
		r.With(s.RequireRoleChi("Admin")).Get("/admin/logs", s.HandleShowLogs)
		r.With(s.RequireRoleChi("Admin")).Get("/admin/logs/{day}", s.HandleShowLogFile)

		// ✅ User management (Admin only)
		r.With(s.RequireAdminChi).Group(func(r chi.Router) {
			r.Get("/admin/users", s.HandleAdminUsers)
			r.Get("/admin/users/new", s.HandleAdminNewUser)
			r.Post("/admin/users/new", s.HandleAdminNewUser)
			r.Get("/admin/users/edit/{userID}", s.HandleAdminEditUser)
			r.Post("/admin/users/edit/{userID}", s.HandleAdminEditUser)
			r.Post("/admin/users/reset-password/{userID}", s.HandleAdminResetPassword)
			r.Get("/admin/users/show-photo/{userID}", s.HandleShowPhoto)
			r.Post("/admin/users/upload-photo/{userID}", s.HandleUploadPhoto)
		})

		// 	// ✅ Vault routes (Admin only)
		// 	r.With(s.RequireAdminChi).Group(func(r chi.Router) {
		// 		r.Get("/admin/vault", s.HandleVaultIndex)
		// 		r.Get("/admin/vault/{vaultID}", s.HandleVaultView)
		// 		r.Post("/admin/vault/{vaultID}/upload", s.HandleVaultUpload)
		// 		r.Post("/admin/vault/{vaultID}/delete", s.HandleVaultDelete)
		// 	})
	})

	// ✅ Mount op de standaard mux
	mux.Handle("/", r)
}
