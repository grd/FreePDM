// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/shared"
)

func (s *Server) AdminDashboardGet(w http.ResponseWriter, r *http.Request) {
	loginname, err := shared.GetSessionLoginname(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := s.UserRepo.LoadUser(loginname)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !auth.IsAdmin(user) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	data := map[string]any{
		"User":            user,
		"ThemePreference": user.ThemePreference,
		"BackButtonShow":  false,
		"MenuButtonShow":  true,
		"MenuButtonLink":  "/admin/preferences",
	}

	if err := s.ExecuteTemplate(w, "admin-dashboard.html", data); err != nil {
		http.Error(w, "Failed to load admin dashboard page", http.StatusInternalServerError)
	}
}

// Handler update for colored log HTML + raw log output
func (s *Server) ShowLogsGet(w http.ResponseWriter, r *http.Request) {
	loginname, err := shared.GetSessionLoginname(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := s.UserRepo.LoadUser(loginname)
	if err != nil || !auth.IsAdmin(user) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	logDir := "logs"
	files, err := os.ReadDir(logDir)
	if err != nil {
		http.Error(w, "Could not read log directory", http.StatusInternalServerError)
		return
	}

	var availableDates []string
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".log") {
			datePart := strings.TrimSuffix(name, ".log")
			if _, err := time.Parse("2006-01-02", datePart); err == nil {
				availableDates = append(availableDates, datePart)
			}
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(availableDates)))

	date := r.URL.Query().Get("date")
	if date == "" && len(availableDates) > 0 {
		date = availableDates[0]
	}

	logPath := filepath.Join(logDir, date+".log")
	logData, err := os.ReadFile(logPath)
	if err != nil {
		logData = []byte("Log file not found: " + logPath)
	}

	rawLog := string(logData)

	// Preprocess log lines into colored HTML
	var logHTML strings.Builder
	lines := strings.Split(rawLog, "\n")
	for _, line := range lines {
		if strings.Contains(line, "[ERROR]") {
			logHTML.WriteString("<span class=\"text-red-400\">" + html.EscapeString(line) + "</span><br>")
		} else if strings.Contains(line, "[WARNING]") {
			logHTML.WriteString("<span class=\"text-yellow-400\">" + html.EscapeString(line) + "</span><br>")
		} else if strings.Contains(line, "[INFO]") {
			logHTML.WriteString("<span class=\"text-green-400\">" + html.EscapeString(line) + "</span><br>")
		} else if strings.Contains(line, "[DEBUG]") {
			logHTML.WriteString("<span class=\"text-blue-400\">" + html.EscapeString(line) + "</span><br>")
		} else {
			logHTML.WriteString(html.EscapeString(line) + "<br>")
		}
	}

	// Combine everything into template data
	data := map[string]any{
		"User":           user,
		"BackButtonShow": true,
		"BackButtonLink": "/admin",
		"MenuButtonShow": false,
		"AvailableDates": availableDates,
		"LogHTML":        template.HTML(logHTML.String()),
		"LogOutput":      rawLog,
	}

	if err = s.ExecuteTemplate(w, "admin-logs.html", data); err != nil {
		http.Error(w, "Failed to load admin logs page", http.StatusInternalServerError)
	}
}

func (s *Server) LoginGet(w http.ResponseWriter, r *http.Request) {
	if err := s.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, "Failed to load login page", http.StatusInternalServerError)
	}
}

func (s *Server) LogoutPost(w http.ResponseWriter, r *http.Request) {
	sess, _ := s.SessionStore.Get(r, shared.SessionName)
	sess.Options.MaxAge = -1
	sess.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) ShowLogFileGet(w http.ResponseWriter, r *http.Request) {
	day := chi.URLParam(r, "day")
	logPath := fmt.Sprintf("logs/%s.log", day)
	http.ServeFile(w, r, logPath)
}

// AdminPreferencesGet shows the Admin preferences
func (s *Server) AdminPreferencesGet(w http.ResponseWriter, r *http.Request) {
	user, err := s.getSessionUser(r)
	log.Printf("[DEBUG] Session user: %v, error: %v", user, err)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	data := map[string]interface{}{
		"User":            user,
		"ThemePreference": user.ThemePreference,
		"BackButtonShow":  true,
		"BackButtonLink":  "/admin",
		"MenuButtonShow":  false,
	}

	if err = s.ExecuteTemplate(w, "admin-preferences.html", data); err != nil {
		http.Error(w, "Failed to load admin preferences page", http.StatusInternalServerError)
	}
}

func (s *Server) ThemePreferencePatch(w http.ResponseWriter, r *http.Request) {
	user, err := s.getSessionUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	preference := r.FormValue("theme_preference")
	if preference != "light" && preference != "dark" && preference != "system" {
		http.Error(w, "Invalid theme preference", http.StatusBadRequest)
		return
	}

	user.ThemePreference = preference
	if err := s.UserRepo.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update theme", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204, no content needed
}

func (s *Server) getSessionUser(r *http.Request) (*db.PdmUser, error) {
	session, _ := s.SessionStore.Get(r, shared.SessionName)

	var userID uint
	switch id := session.Values["user_id"].(type) {
	case uint:
		userID = id
	case int:
		userID = uint(id)
	default:
		return nil, fmt.Errorf("no valid user_id in session")
	}

	return s.UserRepo.LoadUserByID(userID)
}
