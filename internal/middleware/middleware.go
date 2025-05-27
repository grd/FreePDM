// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/shared"
)

// RequireLogin ensures the user is logged in
func RequireLogin(repo db.UserRepo, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginame, err := shared.GetSessionLoginname(r)
		if err != nil || loginame == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// RequireRole returns a handler that only allows users with given roles
func RequireRole(repo db.UserRepo, next http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginname, err := shared.GetSessionLoginname(r)
		if err != nil || loginname == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := repo.LoadUser(loginname)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !auth.HasAnyRole(user, roles...) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// RequireRoleWithLogin ensures the user is logged in and has one of the specified roles
func RequireRoleWithLogin(repo db.UserRepo, next http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginname, err := shared.GetSessionLoginname(r)
		if err != nil || loginname == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := repo.LoadUser(loginname)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !auth.HasAnyRole(user, roles...) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func RequireAdmin(userRepo db.UserRepo, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginname, err := shared.GetSessionLoginname(r)
		if err != nil || loginname == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := userRepo.LoadUser(loginname)
		if err != nil || !user.HasRole(string(db.Admin)) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
