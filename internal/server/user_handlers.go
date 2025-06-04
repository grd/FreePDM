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
	"strconv"
	"strings"
	"time"

	"github.com/grd/FreePDM/internal/db"
)

func (s *Server) HandleShowPhoto(w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] URL.Path:", r.URL.Path)
	userIDStr := strings.TrimPrefix(r.URL.Path, "/admin/users/upload-photo/")
	log.Println(userIDStr)
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

	data := struct {
		User           *db.PdmUser
		ShowBackButton bool
		BackButtonLink string
	}{
		User:           user,
		ShowBackButton: false,
		BackButtonLink: "/admin/users",
	}

	s.ExecuteTemplate(w, "admin-upload-photo.html", data)
}

func (s *Server) HandleUploadPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/admin/users/show-photo/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(parts[0])
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

	filename := fmt.Sprintf("%d-%s%s", userID, time.Now().Format("2006-01-02"), ext)
	photoPath := filepath.Join("static", "uploads", filename)

	out, err := os.Create(photoPath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to store file", http.StatusInternalServerError)
		return
	}

	if err := s.UserRepo.UpdatePhotoPath(uint(userID), photoPath); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users/show-photo/"+strconv.Itoa(userID), http.StatusSeeOther)
}
