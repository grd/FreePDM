// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// https://reintech.io/blog/implementing-authentication-authorization-go
// https://dev.to/tanmayvaish/how-to-implement-authentication-and-authorization-in-golang-20of

var ErrUserNotFound = errors.New("user not found")

// Ease of handling
type UserRepo struct {
	DB *gorm.DB
}

// Constructor
func NewUserRepo(db *gorm.DB) UserRepo {
	return UserRepo{DB: db}
}

// LoadUser search a user based on user name.
func (r *UserRepo) LoadUser(username string) (*PdmUser, error) {
	var user PdmUser
	result := r.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepo) AddUserToSql(username string) {
	fmt.Println("This is basically the interface")
}

// Delete existing user
func (r *UserRepo) RemoveUserFromSql(user_id int, username string) {
	fmt.Println("existing user deleted")
}

func (r *UserRepo) AddUserToLdap(username string) {
	fmt.Println("This is basically the interface")
}

// Delete existing user
func (r *UserRepo) RemoveUserFromLdap(user_id int, username string) {
	fmt.Println("existing user deleted")
}
