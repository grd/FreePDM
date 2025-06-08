// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grd/FreePDM/internal/db"
)

// Middleware adapted for chi

type ctxKey string

const ctxCurrentUser ctxKey = "currentUser"

func (s *Server) RequireLoginChi(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := s.SessionStore.Get(r, "pdm-session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		loginNameRaw := sess.Values["user"]
		loginName, ok := loginNameRaw.(string)
		if !ok || loginName == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := s.UserRepo.LoadUserByLoginName(loginName)
		if err != nil || user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), ctxCurrentUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// In server/middleware.go (in jouw `server` package)

// func (s *Server) RequireAdminChi(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		session, err := s.SessionStore.Get(r, "pdm-session")
// 		if err != nil {
// 			http.Error(w, "Session error", http.StatusUnauthorized)
// 			return
// 		}

// 		loginNameRaw := session.Values["login"]
// 		loginName, ok := loginNameRaw.(string)
// 		if !ok || loginName == "" {
// 			http.Redirect(w, r, "/login", http.StatusFound)
// 			return
// 		}

// 		user, err := s.UserRepo.LoadUserByLoginName(loginName)
// 		if err != nil {
// 			http.Error(w, "User not found", http.StatusUnauthorized)
// 			return
// 		}

// 		if !user.HasRole("Admin") {
// 			http.Redirect(w, r, "/dashboard", http.StatusFound)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), "currentUser", user)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func (s *Server) RequireAdminChi(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := s.SessionStore.Get(r, "pdm-session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		loginNameRaw := sess.Values["user"]
		loginName, ok := loginNameRaw.(string)
		if !ok || loginName == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := s.UserRepo.LoadUserByLoginName(loginName)
		if err != nil || user == nil || !user.HasRole(string(db.Admin)) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), ctxCurrentUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) GetUserFromSession(r *http.Request) (*db.PdmUser, error) {
	session, err := s.SessionStore.Get(r, "pdm-session")
	if err != nil {
		return nil, err
	}

	loginName, ok := session.Values["user"].(string)
	if !ok || loginName == "" {
		return nil, fmt.Errorf("user not found in session")
	}

	return s.UserRepo.LoadUserByLoginName(loginName)
}

// func (s *Server) RequireRoleChi(requiredRole string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			sess, err := s.SessionStore.Get(r, "pdm-session")
// 			if err != nil {
// 				http.Redirect(w, r, "/login", http.StatusSeeOther)
// 				return
// 			}

// 			loginNameRaw := sess.Values["user"]
// 			loginName, ok := loginNameRaw.(string)
// 			if !ok || loginName == "" {
// 				http.Redirect(w, r, "/login", http.StatusSeeOther)
// 				return
// 			}

// 			user, err := s.UserRepo.LoadUserByLoginName(loginName)
// 			if err != nil || user == nil || !user.HasRole(requiredRole) {
// 				http.Error(w, "Forbidden", http.StatusForbidden)
// 				return
// 			}

// 			ctx := context.WithValue(r.Context(), ctxCurrentUser, user)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

// In server.go of middleware.go (in package server)

func (s *Server) RequireRoleChi(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := s.SessionStore.Get(r, "pdm-session")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			loginName, ok := sess.Values["user"].(string)
			if !ok || loginName == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			user, err := s.UserRepo.LoadUserByLoginName(loginName)
			if err != nil || user == nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
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
