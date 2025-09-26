// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grd/FreePDM/internal/config"
)

// Routes handling
func (s *Server) Routes(mux *http.ServeMux) {
	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)

	staticPath := path.Join(config.AppDir(), "static")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	// ✅ Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", s.HomeGet)
		r.Get("/login", s.LoginGet)
		r.Post("/login", s.LoginPost)
		r.Post("/logout", s.LogoutPost)
	})

	// ✅ Require login
	r.Group(func(r chi.Router) {
		r.Use(s.RequireLoginChi)

		// ✅ Dashboard
		r.With(s.RequireRoleChi("Admin")).Get("/admin", s.AdminDashboardGet)
		r.Get("/dashboard", s.DashboardGet)
		r.Get("/admin/preferences", s.AdminPreferencesGet)
		r.Patch("/preferences/theme", s.ThemePreferencePatch)

		// ✅ Logs
		r.With(s.RequireRoleChi("Admin")).Get("/admin/logs", s.ShowLogsGet)
		r.With(s.RequireRoleChi("Admin")).Get("/admin/logs/{day}", s.ShowLogFileGet)

		// ✅ User management (Admin only)
		r.With(s.RequireAdminChi).Group(func(r chi.Router) {
			r.Get("/admin/users", s.AdminUsersGet)
			r.Get("/admin/users/reset-password/{userID}", s.AdminUserResetPasswordGet)
			r.Post("/admin/users/reset-password/{userID}", s.AdminUserResetPasswordPost)
			r.Get("/admin/users/new", s.NewUserGet)
			r.Post("/admin/users/new", s.NewUserPost)
			r.Get("/admin/users/edit/{userID}", s.EditUserGet)
			r.Post("/admin/users/edit/{userID}", s.EditUserPost)
			r.Get("/admin/users/show-photo/{userID}", s.UserPhotoGet)
			r.Post("/admin/users/upload-photo/{userID}", s.UserPhotoPost)
			r.Get("/admin/users/change-status/{userID}", s.UserChangeStatusGet)
			r.Post("/admin/users/change-status/{userID}", s.UserChangeStatusPost)
		})

		// ✅ Vault routes (Admin only)
		r.With(s.RequireAdminChi).Group(func(r chi.Router) {
			r.Get("/vaults/list", s.VaultsListGet)
			r.Get("/vaults/{vaultName}", s.VaultBrowseGet)
			r.Get("/vaults/{vaultName}/*", s.VaultPathBrowseGet)

			// r.Get("/admin/vault/{vaultID}", s.VaultViewGet)
			// r.Post("/admin/vault/{vaultID}/upload", s.VaultUploadPost)
			// r.Post("/admin/vault/{vaultID}/delete", s.VaultDelete)
		})
	})

	// ✅ Mount on the standard mux
	mux.Handle("/", r)
}
