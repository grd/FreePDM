// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/grd/FreePDM/internal/db"
)

// RequireLogin checks for a valid PDM_Session cookie and valid user in DB
func RequireLogin(userRepo db.UserRepo, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("PDM_Session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, err = userRepo.LoadUser(cookie.Value)
		if err != nil {
			// Optionally: clear cookie
			http.SetCookie(w, &http.Cookie{
				Name:   "PDM_Session",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Later: add user to context for downstream handlers
		next(w, r)
	}
}
