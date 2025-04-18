// Copyright 2024 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/shared"
	"github.com/grd/FreePDM/internal/util"
	fsm "github.com/grd/FreePDM/internal/vaults"
)

func CommandHandler(w http.ResponseWriter, r *http.Request) {
	var req shared.CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJsonError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println(req)

	// Using the right command
	switch req.Command {
	case "root":
		handleRoot(w)
	case "list":
		handleList(w)
	case "direxists":
		path, ok := req.Params["path"]
		if !ok {
			writeJsonError(w, "Missing parameters", http.StatusBadRequest)
		}
		handleDirexists(w, req.User, req.Vault, path)
	case "ls":
		path, ok := req.Params["path"]
		if !ok {
			writeJsonError(w, "Missing parameters", http.StatusBadRequest)
		}
		handleLs(w, req.User, req.Vault, path)
	case "allocate":
		path, ok := req.Params["path"]
		if !ok {
			writeJsonError(w, "Missing parameters", http.StatusBadRequest)
		}
		handleAllocate(w, req.User, req.Vault, path)

	// case "rename":
	// 	// Get 'vault', 'src' en 'dst' out of params map
	// 	src := req.Params["src"]
	// 	dst := req.Params["dst"]
	// 	handleRename(w, src, dst)

	default:
		writeJsonError(w, "Unknown command: "+req.Command, http.StatusBadRequest)
	}
}

func writeJsonError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func handleRoot(w http.ResponseWriter) {
	var resp shared.CommandResponse

	root := fsm.Root()
	resp = shared.CommandResponse{
		Data: []string{root},
	}
	json.NewEncoder(w).Encode(resp)
}

// Shows the existing vaults
func handleList(w http.ResponseWriter) {
	var resp shared.CommandResponse

	list, err := fsm.ListVaults()
	if err != nil {
		resp = shared.CommandResponse{
			Error: "Failed to show the list of vaults",
		}
	} else {
		resp = shared.CommandResponse{
			Data: list,
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func handleDirexists(w http.ResponseWriter, user, vault, dir string) {
	var resp shared.CommandResponse

	fs, err := fsm.NewFileSystem(vault, user)
	if err != nil {
		log.Fatalf("unable to access the filesystem : %s", err)
	}

	if ok := fs.DirExists(dir); !ok {
		resp = shared.CommandResponse{
			Error: "Directory " + dir + " does not exists",
		}
	} else {
		resp = shared.CommandResponse{
			Error: "Directory " + dir + " exists",
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func handleLs(w http.ResponseWriter, user, vault, path string) {
	var resp shared.CommandResponse

	fs, err := fsm.NewFileSystem(vault, user)
	if err != nil {
		log.Fatalf("unable to access the filesystem : %s", err)
	}

	fmt.Printf("path = %s\n", path)

	list, err := fs.ListDir(path)
	if err != nil {
		resp = shared.CommandResponse{
			Error: "Failed to show directory",
		}
	} else {
		files := make([]string, len(list))
		for i, item := range list {
			files[i] = item.Name()
		}
		resp = shared.CommandResponse{
			Data: files,
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func handleAllocate(w http.ResponseWriter, user, vault, path string) {
	var resp shared.CommandResponse

	fs, err := fsm.NewFileSystem(vault, user)
	if err != nil {
		log.Fatalf("unable to access the filesystem : %s", err)
	}

	bla, err := fs.Allocate(path)
	if err != nil {
		resp = shared.CommandResponse{
			Error: "Failed to allocate a container",
		}
	} else {
		resp = shared.CommandResponse{
			Data: util.StringToSlice(bla.ContainerNumber),
		}
	}
	json.NewEncoder(w).Encode(resp)
}

// func handleRename(w http.ResponseWriter, src, dst string) {

// 	// Rename the file
// 	if err := os.Rename(src, dst); err != nil {
// 		resp := shared.CommandResponse{
// 			Status:  "error",
// 			Message: "Failed to rename file",
// 			Data:    map[string]interface{}{"error": err.Error()},
// 		}
// 		json.NewEncoder(w).Encode(resp)
// 		return
// 	}

// 	// Send success response
// 	resp := shared.CommandResponse{
// 		Status:  "success",
// 		Message: "File renamed successfully",
// 	}
// 	json.NewEncoder(w).Encode(resp)
// }

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if isValidUser(username, password) {
		sessionID := uuid.NewString()
		expiration := time.Now().Add(30 * time.Minute)

		db.Sessions[sessionID] = db.Session{
			Username:   username,
			Expiration: expiration,
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionID,
			Expires: expiration,
		})

		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	} else {
		http.Error(w, "Ongeldige gebruikersnaam of wachtwoord", http.StatusUnauthorized)
	}
}

func isValidUser(username, password string) bool {
	users := map[string]string{
		"user1": "password1",
		"user2": "password2",
	}

	if pass, ok := users[username]; ok {
		return pass == password
	}
	return false
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Niet geautoriseerd", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		sessionID := cookie.Value
		session, exists := db.Sessions[sessionID]

		if !exists || session.Expiration.Before(time.Now()) {
			http.Error(w, "Niet geautoriseerd", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", session.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	fmt.Fprintf(w, "Welkom, %s!", username)
}
