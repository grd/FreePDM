// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"log"
	"net/http"
	usr "os/user"
	"slices"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/shared"
)

func (s *Server) HandleAdminUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		http.Error(w, "Error loading users", http.StatusInternalServerError)
		return
	}

	formattedDOB := ""
	if !user.DateOfBirth.IsZero() {
		formattedDOB = user.DateOfBirth.Format("2006-01-02")
	}

	data := struct {
		User           *db.PdmUser
		ShowBackButton bool
		BackButtonLink string
		Users          []db.PdmUser
		FormattedDOB   string
	}{
		User:           user,
		ShowBackButton: true,
		BackButtonLink: "/admin",
		Users:          users,
		FormattedDOB:   formattedDOB,
	}

	for _, u := range users {
		log.Printf("[DEBUG] User: ID=%d, LoginName=%s, FullName=%s, PhotoPath=%s", u.ID, u.LoginName, u.FullName, u.PhotoPath)
	}

	s.ExecuteTemplate(w, "admin-users.html", data)
}

func (s *Server) HandleNewUserForm(w http.ResponseWriter, r *http.Request) {
	user := &db.PdmUser{}

	availableRoles := db.GetAvailableRoles()

	roleChecks := make(map[string]bool)
	for _, role := range user.Roles {
		roleChecks[role] = true
	}

	data := map[string]interface{}{
		"User":           user,
		"RoleChecks":     roleChecks,
		"AvailableRoles": availableRoles,
		"ShowBackButton": true,
		"BackButtonLink": "/admin/users",
	}

	s.ExecuteTemplate(w, "admin-new-user.html", data)
}

func (s *Server) HandleCreateNewUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form parsing error", http.StatusBadRequest)
		return
	}

	user := &db.PdmUser{
		FullName:     r.FormValue("full_name"),
		FirstName:    r.FormValue("first_name"),
		LastName:     r.FormValue("last_name"),
		EmailAddress: r.FormValue("email_address"),
		Sex:          r.FormValue("sex"),
		PhoneNumber:  r.FormValue("phone_number"),
		Department:   r.FormValue("department"),
		Roles:        r.Form["roles"],
	}

	dobStr := r.FormValue("date_of_birth")
	if dobStr != "" {
		dob, err := time.Parse("2006-01-02", dobStr)
		if err == nil {
			user.DateOfBirth = dob
		}
	}

	// Check previously created login name
	existingUser, _ := s.UserRepo.LoadUserByLoginName(user.LoginName)
	if existingUser != nil {
		http.Error(w, "Login name already exists", http.StatusBadRequest)
		return
	}

	// Check user exists in Linux
	_, err := usr.Lookup(user.LoginName)
	if err != nil {
		http.Error(w, "Login name does not exist as a Linux user", http.StatusBadRequest)
		return
	}

	// Simpel wachtwoord bij aanmaak (voor later te resetten)
	// In praktijk moet je hier een gegenereerd wachtwoord + e-mail reset link gebruiken
	defaultPassword := "changeme123"
	hash, err := auth.HashPassword(defaultPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = hash

	if err := s.UserRepo.CreateUser(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Created new user %s with roles %v", user.LoginName, user.Roles)
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (s *Server) HandleEditUserForm(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := s.UserRepo.LoadUserByID(uint(userID))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	roleChecks := make(map[string]bool)
	for _, role := range user.Roles {
		roleChecks[role] = true
	}

	availableRoles := db.GetAvailableRoles()

	data := map[string]interface{}{
		"User":           user,
		"RoleChecks":     roleChecks,
		"AvailableRoles": availableRoles,
		"ShowBackButton": true,
		"BackButtonLink": "/admin/users",
	}

	s.ExecuteTemplate(w, "admin-edit-user.html", data)
}

func (s *Server) HandleSaveEditedUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := s.UserRepo.LoadUserByID(uint(userID))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form parsing error", http.StatusBadRequest)
		return
	}

	user.FullName = r.FormValue("full_name")
	user.FirstName = r.FormValue("first_name")
	user.LastName = r.FormValue("last_name")
	user.EmailAddress = r.FormValue("email_address")
	user.Sex = r.FormValue("sex")
	user.PhoneNumber = r.FormValue("phone_number")
	user.Department = r.FormValue("department")

	dobStr := r.FormValue("date_of_birth")
	if dobStr != "" {
		dob, err := time.Parse("2006-01-02", dobStr)
		if err == nil {
			user.DateOfBirth = dob
		}
	}

	roles := r.Form["roles"]
	if roles == nil {
		roles = []string{}
	}
	// Beveiliging: Admin ID mag niet zijn Admin-rechten verliezen
	if user.ID == 1 && !slices.Contains(roles, "Admin") {
		roles = append(roles, "Admin")
	}
	user.Roles = roles

	if err := s.UserRepo.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
