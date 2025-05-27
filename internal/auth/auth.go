// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"github.com/grd/FreePDM/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	HashedPassword string
	SessionToke    string
	CSRFToke       string
}

// temporary Key is loginname
var Users = map[string]Login{}

func IsValidUser(loginname, password string, repo *db.UserRepo) bool {
	user, err := repo.LoadUser(loginname)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil
}

func IsValidSession(sessionValue string, repo *db.UserRepo) bool {
	user, err := repo.LoadUser(sessionValue)
	return err == nil && user != nil
}
