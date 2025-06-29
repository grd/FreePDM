// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package shared

import (
	"errors"
	"net/http"
)

const SessionName = "pdm_session"

var (
	ErrNoSessionCookie = errors.New("no session cookie found")
	SessionCookieName  = "PDM_Session"
)

// GetSessionLoginname retrieves the loginname stored in the session cookie.
func GetSessionLoginname(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return "", ErrNoSessionCookie
	}
	return cookie.Value, nil
}

// SetSessionCookie sets the session cookie for the user.
func SetSessionCookie(w http.ResponseWriter, loginname string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    loginname,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to false if not using HTTPS in dev
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie deletes the session cookie.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
