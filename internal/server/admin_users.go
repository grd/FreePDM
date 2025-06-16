// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	usr "os/user"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/config"
	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/shared"
)

func (s *Server) AdminUsersGet(w http.ResponseWriter, r *http.Request) {
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

	data := map[string]any{
		"User":           user,
		"BackButtonShow": true,
		"BackButtonLink": "/admin",
		"MenuButtonShow": false,
		"Users":          users,
		"FormattedDOB":   formattedDOB,
	}

	for _, u := range users {
		log.Printf("[DEBUG] User: ID=%d, LoginName=%s, FullName=%s, PhotoPath=%s", u.ID, u.LoginName, u.FullName, u.PhotoPath)
	}

	s.ExecuteTemplate(w, "admin-users.html", data)
}

func (s *Server) NewUserGet(w http.ResponseWriter, r *http.Request) {
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
		"BackButtonShow": true,
		"BackButtonLink": "/admin/users",
		"MenuButtonShow": false,
	}

	s.ExecuteTemplate(w, "admin-new-user.html", data)
}

func (s *Server) NewUserPost(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) EditUserGet(w http.ResponseWriter, r *http.Request) {
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
	availableStatuses := db.GetAvailableStatuses()

	data := map[string]interface{}{
		"User":              user,
		"RoleChecks":        roleChecks,
		"AvailableRoles":    availableRoles,
		"AvailableStatuses": availableStatuses,
		"BackButtonShow":    true,
		"BackButtonLink":    "/admin/users",
		"MenuButtonShow":    false,
	}

	s.ExecuteTemplate(w, "admin-edit-user.html", data)
}

func (s *Server) EditUserPost(w http.ResponseWriter, r *http.Request) {
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

	// dealing with FormValues

	user.FullName = r.FormValue("full_name")
	user.FirstName = r.FormValue("first_name")
	user.LastName = r.FormValue("last_name")
	user.Sex = r.FormValue("sex")
	user.PhoneNumber = r.FormValue("phone_number")
	user.Department = r.FormValue("department")
	user.AccountStatus = r.FormValue("account_status")

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

	email := r.FormValue("email_address")
	if email == "" {
		http.Error(w, "Email address is required", http.StatusBadRequest)
		return
	}
	user.EmailAddress = email

	// Security: Admin ID can't lose it's Admin rights
	if user.ID == 1 && !slices.Contains(roles, "Admin") {
		roles = append(roles, "Admin")
	}
	user.Roles = roles

	log.Printf("[DEBUG] User update failed for ID %d: %+v", user.ID, user)

	if err := s.UserRepo.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprint("/admin/users/edit/", userIDStr)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (s *Server) UserPhotoGet(w http.ResponseWriter, r *http.Request) {
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

	data := map[string]interface{}{
		"User":           user,
		"BackButtonShow": true,
		"BackButtonLink": "/admin/users",
		"MenuButtonShow": false,
	}

	if err := s.ExecuteTemplate(w, "show-photo.html", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("[ERROR] Template execution failed: %v", err)
	}
}

func (s *Server) UserPhotoPost(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > 5*1024*1024 {
		http.Error(w, "File too large (max 5MB)", http.StatusBadRequest)
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		http.Error(w, "Only JPG and PNG allowed", http.StatusBadRequest)
		return
	}

	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%d-%s%s", userID, timestamp, ext)
	savePath := filepath.Join(config.AppDir(), "static", "uploads", filename)

	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "Failed to store file", http.StatusInternalServerError)
		return
	}

	photoPath := filename
	if err := s.UserRepo.UpdatePhotoPath(uint(userID), photoPath); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Updated photo for user ID %d => %s", userID, photoPath)
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// AdminUserResetPasswordGet shows the reset password form
func (s *Server) AdminUserResetPasswordGet(w http.ResponseWriter, r *http.Request) {
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

	data := map[string]interface{}{
		"User": user,
	}
	s.ExecuteTemplate(w, "admin-reset-password.html", data)
}

// AdminUserResetPasswordPost updates the password
func (s *Server) AdminUserResetPasswordPost(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := s.UserRepo.LoadUserByID(uint(userID))
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	newPassword := r.FormValue("new_password")
	if newPassword == "" {
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	if err := s.UserRepo.UpdatePasswordByID(uint(userID), hash); err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		log.Printf("[ERROR] Failed to update password for user %d: %v", userID, err)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (s *Server) UserChangeStatusGet(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")
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

	data := map[string]interface{}{
		"User":              user,
		"AvailableStatuses": db.GetAvailableStatuses(),
	}

	if err := s.ExecuteTemplate(w, "admin-change-status.html", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func (s *Server) UserChangeStatusPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	status := r.FormValue("account_status")

	if err := s.UserRepo.UpdateAccountStatus(uint(userID), status); err != nil {
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	redirectStr := path.Join("/admin/users", idStr)
	http.Redirect(w, r, redirectStr, http.StatusSeeOther)
}

func StatusTooltip(status string) string {
	switch status {
	case "Active":
		return "The user account is active and fully functional."
	case "Disabled":
		return "The account is disabled and cannot be used until re-enabled."
	case "Locked":
		return "The account is locked due to security reasons."
	case "Pending":
		return "The account is pending approval or verification."
	case "Suspended":
		return "The account is temporarily suspended."
	case "Expired":
		return "The account access has expired."
	case "Deleted":
		return "The account is marked as deleted."
	case "Invited":
		return "The user has been invited but has not yet accepted."
	default:
		return "Unknown status"
	}
}
