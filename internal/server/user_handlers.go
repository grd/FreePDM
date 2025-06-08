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

	"github.com/go-chi/chi/v5"
)

func (s *Server) HandleShowPhoto(w http.ResponseWriter, r *http.Request) {
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
		"ShowBackButton": true,
		"BackButtonLink": "/admin/users",
	}

	if err := s.ExecuteTemplate(w, "show-photo.html", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("[ERROR] Template execution failed: %v", err)
	}
}

func (s *Server) HandleUploadPhoto(w http.ResponseWriter, r *http.Request) {
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
	savePath := filepath.Join("static", "uploads", filename)

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

	photoPath := filepath.Join("uploads", filename)
	if err := s.UserRepo.UpdatePhotoPath(uint(userID), photoPath); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Updated photo for user ID %d => %s", userID, photoPath)
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
