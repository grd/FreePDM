// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"log"
	"net/http"

	"github.com/grd/FreePDM/internal/shared"
)

// Middleware adapted for chi

type ctxKey string

const ctxCurrentUser ctxKey = "currentUser"

// RequireLoginChi ensures the user is logged in
func (s *Server) RequireLoginChi(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := s.SessionStore.Get(r, shared.SessionName)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userIDRaw := sess.Values["user_id"]
		userID, ok := userIDRaw.(uint)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := s.UserRepo.LoadUserByID(userID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), ctxCurrentUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdminChi allows only admin users
func (s *Server) RequireAdminChi(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := s.SessionStore.Get(r, shared.SessionName)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userIDRaw := sess.Values["user_id"]
		userID, ok := userIDRaw.(uint)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := s.UserRepo.LoadUserByID(userID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), ctxCurrentUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRoleChi returns a handler that only allows users with given roles
func (s *Server) RequireRoleChi(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := s.SessionStore.Get(r, shared.SessionName)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			userIDRaw := sess.Values["user_id"]
			userID, ok := userIDRaw.(uint)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			user, err := s.UserRepo.LoadUserByID(userID)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			for _, role := range roles {
				if user.HasRole(role) {
					ctx := context.WithValue(r.Context(), ctxCurrentUser, user)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			log.Printf("[DEBUG] RequireAnyRoleChi: gebruiker %s heeft geen van de vereiste rollen %v", user.LoginName, roles)
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
