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
	"time"

	"github.com/grd/FreePDM/internal/config"
)

// HandleUploadPhoto handles the file upload, saves it, and updates the user record.
func (s *Server) HandleUploadPhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(5 << 20); err != nil { // max 5MB
		log.Printf("[ERROR] Failed to parse multipart form: %v", err)
		http.Error(w, "Invalid upload", http.StatusBadRequest)
		return
	}

	userID := r.FormValue("user_id")
	if userID == "" {
		log.Println("[ERROR] No user_id provided in form")
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	idUint64, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		log.Printf("[ERROR] Invalid user_id format: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint64)
	log.Printf("id = %d", id)

	file, handler, err := r.FormFile("photo")
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve file: %v", err)
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		log.Printf("[ERROR] Unsupported file type: %s", ext)
		http.Error(w, "Only PNG and JPEG files are allowed", http.StatusBadRequest)
		return
	}

	filename := fmt.Sprintf("%s-%s%s", userID, time.Now().Format("2006-01-02"), ext)
	fullPath := filepath.Join(config.AppDir(), "static", "uploads", filename)
	log.Printf("[DEBUG] Saving photo to: %s", fullPath)

	dst, err := os.Create(fullPath)
	if err != nil {
		log.Printf("[ERROR] Failed to create file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("[ERROR] Failed to write file contents: %v", err)
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	// Update PhotoPath in DB
	if err := s.UserRepo.UpdatePhotoPath(id, filepath.Join("static", "uploads", filename)); err != nil {
		log.Printf("[ERROR] Failed to update photo path in DB: %v", err)
		http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Successfully uploaded photo for user %s", userID)
	http.Redirect(w, r, "/admin/users/edit/"+userID, http.StatusSeeOther)
}
