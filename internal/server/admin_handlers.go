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
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/shared"
	"golang.org/x/crypto/bcrypt"
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
		BackButtonShow bool
		BackButtonLink string
		MenuButtonShow bool
		MenuButtonLink string
	}{
		User:           user,
		BackButtonShow: false,
		BackButtonLink: "",
		MenuButtonShow: true,
		MenuButtonLink: "/admin/edit",
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
		BackButtonShow bool
		BackButtonLink string
		MenuButtonShow bool
		AvailableDates []string
		LogHTML        template.HTML
		LogOutput      string
	}{
		User:           user,
		BackButtonShow: true,
		BackButtonLink: "/admin",
		MenuButtonShow: false,
		AvailableDates: availableDates,
		LogHTML:        template.HTML(logHTML.String()),
		LogOutput:      rawLog,
	}

	s.ExecuteTemplate(w, "admin-logs.html", data)
}

// Handler for resetting password
func (s *Server) HandleAdminResetPassword(w http.ResponseWriter, r *http.Request) {
	username, err := shared.GetSessionLoginname(r)
	if err != nil || username == "" {
		log.Printf("[DEBUG] No valid session — redirecting to /login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/admin/users/reset-password/")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := s.UserRepo.LoadUserByID(uint(userID))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Generate a temporary password (or use a fixed default)
	tempPassword := "changeme123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = string(hashedPassword)
	user.MustChangePassword = true

	if err := s.UserRepo.UpdateUser(user); err != nil {
		log.Printf("[ERROR] Failed to reset password for user %s: %v", user.LoginName, err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	log.Printf("[DEBUG] Admin %s reset password for user %s", username, user.LoginName)
	http.Redirect(w, r, fmt.Sprintf("/admin/users/edit/%d", userID), http.StatusSeeOther)
}

func (s *Server) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := s.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, "Failed to load login page", http.StatusInternalServerError)
	}
}

func (s *Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sess, _ := s.SessionStore.Get(r, "pdm-session")
	sess.Values["user"] = ""
	sess.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) HandleShowLogFile(w http.ResponseWriter, r *http.Request) {
	day := chi.URLParam(r, "day")
	logPath := fmt.Sprintf("logs/%s.log", day)
	http.ServeFile(w, r, logPath)
}

func (s *Server) HandleAdminEdit(w http.ResponseWriter, r *http.Request) {
	err := s.ExecuteTemplate(w, "admin-edit.html", nil)
	if err != nil {
		http.Error(w, "Failed to load admin edit page", http.StatusInternalServerError)
		log.Printf("[ERROR] Rendering admin edit page: %v", err)
	}
}
