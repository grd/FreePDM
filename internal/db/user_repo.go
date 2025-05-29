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
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

// UpdateUser saves the updated user record into the database.
// It overwrites the existing user entry based on its primary key.
func (r *UserRepo) UpdateUser(user *PdmUser) error {
	return r.DB.Save(user).Error
}

// LoadUser search by ID on user name.
func (r *UserRepo) LoadUserByID(id uint) (*PdmUser, error) {
	var user PdmUser
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// LoadUser search a user based on user name.
func (r *UserRepo) LoadUser(loginname string) (*PdmUser, error) {
	var user PdmUser
	result := r.DB.Where("loginname = ?", loginname).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// Renames password and also sets the MustChangePassword flag to false
func (r *UserRepo) UpdatePassword(loginname, hash string) error {
	var user PdmUser
	if err := r.DB.Where("loginname = ?", loginname).First(&user).Error; err != nil {
		return err
	}

	user.PasswordHash = hash
	user.MustChangePassword = false

	return r.DB.Save(&user).Error
}

// Resets the MustChangePassword flag to false
func (r *UserRepo) ClearMustChangePassword(loginname string) error {
	return r.DB.Model(&PdmUser{}).
		Where("loginname = ?", loginname).
		Update("must_change_password", false).
		Error
}

func (r *UserRepo) GetAllUsers() ([]PdmUser, error) {
	var users []PdmUser
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *UserRepo) CreateUser(user *PdmUser) error {
	return r.DB.Create(user).Error
}

func (r *UserRepo) AddUserToLdap(loginname string) {
	fmt.Println("This is basically the interface")
}

// Delete existing user
func (r *UserRepo) RemoveUserFromLdap(user_id int, loginname string) {
	fmt.Println("existing user deleted")
}
