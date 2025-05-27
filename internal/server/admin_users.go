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
	"path/filepath"
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

	data := struct {
		User           *db.PdmUser
		ShowBackButton bool
		BackButtonLink string
		Users          []db.PdmUser
	}{
		User:           user,
		ShowBackButton: true,
		BackButtonLink: "/admin",
		Users:          users,
	}

	for _, u := range users {
		log.Printf("[DEBUG] User: ID=%d, LoginName=%s, FullName=%s, PhotoPath=%s", u.ID, u.LoginName, u.FullName, u.PhotoPath)
	}

	s.ExecuteTemplate(w, "admin-users.html", data)
}

func (s *Server) HandleAdminNewUser(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == http.MethodGet {
		// Render form
		s.ExecuteTemplate(w, "admin-new-user.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseMultipartForm(10 << 20) // max 10MB upload
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		loginname := r.FormValue("loginname")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		dateOfBirth := r.FormValue("date_of_birth")
		sex := r.FormValue("sex")
		email := r.FormValue("email_address")
		phone := r.FormValue("phone_number")
		department := r.FormValue("department")
		roles := r.Form["roles"]

		// Basic validation
		if password != confirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}
		if len(password) < 10 {
			http.Error(w, "Password must be at least 10 characters", http.StatusBadRequest)
			return
		}

		// Parse date
		dob, err := time.Parse("2006-01-02", dateOfBirth)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Handle photo upload
		file, handler, err := r.FormFile("photo")
		var photoPath string
		if err == nil {
			defer file.Close()

			// Save photo (e.g., static/uploads/userID_photo.ext)
			photoFilename := fmt.Sprintf("%s_%s", loginname, handler.Filename)
			photoPath = filepath.Join("static/uploads", photoFilename)
			dst, err := os.Create(photoPath)
			if err != nil {
				http.Error(w, "Failed to save photo", http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			io.Copy(dst, file)
		}

		// Create user object
		newUser := db.PdmUser{
			LoginName:          loginname,
			PasswordHash:       string(hashedPassword),
			MustChangePassword: true,
			FirstName:          firstName,
			LastName:           lastName,
			FullName:           fmt.Sprintf("%s %s", firstName, lastName),
			DateOfBirth:        dob,
			Sex:                sex,
			EmailAddress:       email,
			PhoneNumber:        phone,
			Department:         department,
			PhotoPath:          photoPath,
			Roles:              roles,
		}

		// Insert into database
		err = s.UserRepo.CreateUser(&newUser)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Redirect back to user list
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	}
}
