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
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/grd/FreePDM/internal/auth"
	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/shared"
	"golang.org/x/crypto/bcrypt"
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

// Helper: check if system user exists
func isSystemUserExists(username string) bool {
	cmd := exec.Command("id", username)
	err := cmd.Run()
	return err == nil
}
func (s *Server) HandleAdminNewUser(w http.ResponseWriter, r *http.Request) {
	username, err := shared.GetSessionLoginname(r)
	if err != nil || username == "" {
		log.Printf("[DEBUG] No valid session — redirecting to /login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		log.Printf("[DEBUG] Serving New User page")
		s.ExecuteTemplate(w, "admin-new-user.html", map[string]interface{}{
			"AvailableRoles": db.GetAvailableRoles(),
			"ShowBackButton": true,
			"BackButtonLink": "/admin/users",
		})
		return
	}

	// POST processing
	loginname := r.FormValue("loginname")
	fullname := r.FormValue("fullname")
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")
	dateOfBirth := r.FormValue("date_of_birth")
	sex := r.FormValue("sex")
	phone := r.FormValue("phone_number")
	department := r.FormValue("department")
	roles := r.Form["roles"]

	log.Printf("[DEBUG] Processing New User form for loginname: %s", loginname)

	// Basic validations
	if loginname == "" || fullname == "" || email == "" || password == "" || len(roles) == 0 {
		s.ExecuteTemplate(w, "admin-new-user.html", map[string]interface{}{
			"Error":          "Please fill in all required fields.",
			"AvailableRoles": db.GetAvailableRoles(),
		})
		return
	}
	if password != confirmPassword {
		s.ExecuteTemplate(w, "admin-new-user.html", map[string]interface{}{
			"Error":          "Passwords do not match.",
			"AvailableRoles": db.GetAvailableRoles(),
		})
		return
	}

	// Check if user exists in FreePDM
	_, err = s.UserRepo.LoadUser(loginname)
	if err == nil {
		s.ExecuteTemplate(w, "admin-new-user.html", map[string]interface{}{
			"Error":          "A user with this login name already exists in FreePDM.",
			"AvailableRoles": db.GetAvailableRoles(),
		})
		return
	}

	// Check if system user exists (Linux)
	if !isSystemUserExists(loginname) {
		s.ExecuteTemplate(w, "admin-new-user.html", map[string]interface{}{
			"Error":          "This login name does not exist on the system.",
			"AvailableRoles": db.GetAvailableRoles(),
		})
		return
	}

	// Parse date
	dob, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		s.ExecuteTemplate(w, "admin-new-user.html", map[string]interface{}{
			"Error":          "Invalid date format.",
			"AvailableRoles": db.GetAvailableRoles(),
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.ExecuteTemplate(w, "admin-new-user.html", map[string]interface{}{
			"Error":          "Failed to hash password.",
			"AvailableRoles": db.GetAvailableRoles(),
		})
		return
	}

	// Handle photo upload
	file, handler, err := r.FormFile("photo")
	var photoPath string
	if err == nil {
		defer file.Close()

		photoFilename := fmt.Sprintf("%s_%s", loginname, handler.Filename)
		photoPath = filepath.Join("static/uploads", photoFilename)
		dst, err := os.Create(photoPath)
		if err != nil {
			s.ExecuteTemplate(w, "admin-new-user.html", map[string]interface{}{
				"Error":          "Failed to save photo.",
				"AvailableRoles": db.GetAvailableRoles(),
			})
			return
		}
		defer dst.Close()
		io.Copy(dst, file)
	}

	// Create new user in DB
	newUser := db.PdmUser{
		LoginName:          loginname,
		FullName:           fullname,
		FirstName:          firstname,
		LastName:           lastname,
		EmailAddress:       email,
		PasswordHash:       string(hashedPassword),
		MustChangePassword: true,
		DateOfBirth:        dob,
		Sex:                sex,
		PhoneNumber:        phone,
		Department:         department,
		PhotoPath:          photoPath,
		Roles:              roles,
	}

	if err := s.UserRepo.CreateUser(&newUser); err != nil {
		log.Printf("[ERROR] Failed to create user: %v", err)
		s.ExecuteTemplate(w, "admin-new-user.html", map[string]interface{}{
			"Error":          "Failed to create user. Please try again.",
			"AvailableRoles": db.GetAvailableRoles(),
		})
		return
	}

	log.Printf("[DEBUG] User %s successfully created by admin %s", loginname, username)
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (s *Server) HandleAdminEditUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := strings.TrimPrefix(r.URL.Path, "/admin/users/edit/")
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

	availableRoles := db.GetAvailableRoles()
	roleChecks := make(map[string]bool)
	for _, role := range availableRoles {
		roleChecks[role] = slices.Contains(user.Roles, role)
	}

	if r.Method == http.MethodGet {
		s.ExecuteTemplate(w, "admin-edit-user.html", map[string]interface{}{
			"User":           user,
			"AvailableRoles": availableRoles,
			"RoleChecks":     roleChecks,
			"ShowBackButton": true,
			"BackButtonLink": "/admin/users",
		})
		return
	}

	// POST: Update user details
	fullName := r.FormValue("fullname")
	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	email := r.FormValue("email")
	dateOfBirthStr := r.FormValue("date_of_birth")
	sex := r.FormValue("sex")
	phone := r.FormValue("phone_number")
	department := r.FormValue("department")
	roles := r.Form["roles"]

	dob, err := time.Parse("2006-01-02", dateOfBirthStr)
	if err != nil {
		s.ExecuteTemplate(w, "admin-edit-user.html", map[string]interface{}{
			"User":           user,
			"AvailableRoles": availableRoles,
			"RoleChecks":     roleChecks,
			"ShowBackButton": true,
			"BackButtonLink": "/admin/users",
			"Error":          "Invalid date format",
		})
		return
	}

	// Update fields
	user.FullName = fullName
	user.FirstName = firstName
	user.LastName = lastName
	user.EmailAddress = email
	user.DateOfBirth = dob
	user.Sex = sex
	user.PhoneNumber = phone
	user.Department = department
	user.Roles = roles

	// Save changes
	if err := s.UserRepo.UpdateUser(user); err != nil {
		log.Printf("[ERROR] Failed to update user: %v", err)
		s.ExecuteTemplate(w, "admin-edit-user.html", map[string]interface{}{
			"User":           user,
			"AvailableRoles": availableRoles,
			"RoleChecks":     roleChecks,
			"ShowBackButton": true,
			"BackButtonLink": "/admin/users",
			"Error":          "Failed to update user. Please try again.",
		})
		return
	}

	log.Printf("[INFO] Updated user %s (ID %d)", user.LoginName, user.ID)
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
