// Copyright 2024 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/grd/FreePDM/internal/config"
	"github.com/grd/FreePDM/internal/shared"
	"github.com/grd/FreePDM/internal/util"
	vfs "github.com/grd/FreePDM/internal/vault/localfs"
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
	// 	// Get 'vault', 'src' and 'dst' out of params map
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

	root := vfs.Root()
	resp = shared.CommandResponse{
		Data: []string{root},
	}
	json.NewEncoder(w).Encode(resp)
}

// Shows the existing vaults
func handleList(w http.ResponseWriter) {
	var resp shared.CommandResponse

	list, err := vfs.ListVaults()
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

	fs, err := vfs.NewFileSystem(vault, user)
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

	fs, err := vfs.NewFileSystem(vault, user)
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

	fs, err := vfs.NewFileSystem(vault, user)
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

// VaultsListGet shows all vaults inside the filesystem
func (s *Server) VaultsListGet(w http.ResponseWriter, r *http.Request) {
	vaultRoot := config.VaultsDir()
	dirs, err := os.ReadDir(vaultRoot)
	if err != nil {
		http.Error(w, "Unable to read vaults directory", http.StatusInternalServerError)
		log.Printf("[ERROR] Unable to read the root vault: %v", err)
		return
	}

	var vaults []string
	for _, d := range dirs {
		if d.IsDir() && !strings.HasPrefix(d.Name(), ".") {
			vaults = append(vaults, d.Name())
		}
	}

	data := map[string]any{
		"Vaults":          vaults,
		"Title":           "Vaults",
		"BackButtonShow":  true,
		"BackButtonLink":  "/dashboard",
		"MenuButtonShow":  false,
		"ThemePreference": "system",
	}

	err = s.ExecuteTemplate(w, "vaults-list.html", data)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func (s *Server) VaultBrowseGet(w http.ResponseWriter, r *http.Request) {
	vaultName := chi.URLParam(r, "vaultName")
	subPath := chi.URLParam(r, "*") // alles ná de vaultName
	subPath = strings.TrimPrefix(subPath, "/")

	user, err := s.getSessionUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fs, err := vfs.NewFileSystem(vaultName, user.LoginName)
	if err != nil {
		log.Printf("[ERROR] Failed to create FS for vault=%s user=%s: %v", vaultName, user.LoginName, err)
		http.Error(w, "Vault init error", http.StatusInternalServerError)
		return
	}

	fullPath := path.Join(fs.VaultDir(), subPath)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		http.Error(w, "Read error", http.StatusInternalServerError)
		return
	}

	var results []VaultEntry
	for _, entry := range entries {
		results = append(results, VaultEntry{
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			NextURL: path.Join("/vaults", vaultName, subPath, entry.Name()),
		})
	}

	data := map[string]any{
		"VaultName":      vaultName,
		"SubPath":        subPath,
		"Entries":        results,
		"BackButtonShow": true,
		"BackButtonLink": "/vaults/list",
		"MenuButtonShow": false,
	}

	if err := s.ExecuteTemplate(w, "vaults-browse.html", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func (s *Server) VaultPathBrowseGet(w http.ResponseWriter, r *http.Request) {
	vaultName := chi.URLParam(r, "vaultName")
	subPath := chi.URLParam(r, "*")

	// Clean vault path
	cleanPath := path.Join("vaults", vaultName, subPath)
	fullPath := filepath.Join(config.Conf.VaultsDirectory, cleanPath)

	entries, err := fs.ReadDir(*s.FS, fullPath)
	if err != nil {
		http.Error(w, "Unable to read vault path", http.StatusNotFound)
		log.Printf("[ERROR] Cannot read vault path %q: %v", fullPath, err)
		return
	}

	var files []VaultEntry
	for _, entry := range entries {
		name := entry.Name()
		entryPath := path.Join("/vaults", vaultName, subPath, name)

		// Ensure directories end with a slash for UX clarity (optional)
		if entry.IsDir() {
			entryPath += "/"
		}

		files = append(files, VaultEntry{
			Name:    name,
			IsDir:   entry.IsDir(),
			NextURL: entryPath,
		})
	}

	data := map[string]any{
		"VaultName":      vaultName,
		"SubPath":        subPath,
		"Entries":        files,
		"BackButtonShow": true,
		"BackButtonLink": "/vaults/list",
		"MenuButtonShow": false,
	}

	if err := s.ExecuteTemplate(w, "vaults-browse.html", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

type VaultEntry struct {
	Name    string
	IsDir   bool
	NextURL string
}
