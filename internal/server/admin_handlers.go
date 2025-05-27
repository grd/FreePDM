// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"html"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/shared"
)

func (s *Server) HandleAdminDashboard(w http.ResponseWriter, r *http.Request) {
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

	data := struct {
		User           *db.PdmUser
		ShowBackButton bool
		BackButtonLink string
	}{
		User:           user,
		ShowBackButton: false,
		BackButtonLink: "",
	}

	s.ExecuteTemplate(w, "admin-dashboard.html", data)
}

// Handler update for colored log HTML + raw log output
func (s *Server) HandleShowLogs(w http.ResponseWriter, r *http.Request) {
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
	data := struct {
		User           *db.PdmUser
		ShowBackButton bool
		BackButtonLink string
		AvailableDates []string
		LogHTML        template.HTML
		LogOutput      string
	}{
		User:           user,
		ShowBackButton: true,
		BackButtonLink: "/admin",
		AvailableDates: availableDates,
		LogHTML:        template.HTML(logHTML.String()),
		LogOutput:      rawLog,
	}

	s.ExecuteTemplate(w, "admin-logs.html", data)
}
