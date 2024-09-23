// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"fmt"
)

// Basically searching in the SQL database requires a coonetcion to the database too.
// Is it than a sub-class or better to request access to a class or has this class a interface too?

// Search
// https://docs.sqlalchemy.org/en/14/orm/session_basics.html
type Search struct {
}

func (self Search) init() {
	fmt.Println("Search in Databases")
}

// Search on number
func (self Search) search_number(number int) {
	// raise NotImplementedError("Function search_number is not implemented yet")
}

// Search on description
func (self Search) search_description(description string) {
	// raise NotImplementedError("Function search_description is not implemented yet")
}

// Search on something else
func (self Search) search_something_else(something string) {
	// raise NotImplementedError("Function search_something_else is not implemented yet")
}

// help function
func (self Search) search_help() string {
	help_text := `
        ad some text
        - Modifiers
        - Search keys etc
        `
	fmt.Println(help_text)
	return help_text
}

// Search for projects
type SearchItem struct {
}

func (self SearchItem) init() {
	fmt.Println("Search in items")
}

// Search on project number
func (self SearchItem) item_number(user_number string) {
	// raise NotImplementedError("Function project_number is not implemented yet")
}

// Search on project description
func (self SearchItem) item_description(description string) {
	// raise NotImplementedError("Function project_description is not implemented yet")
}

// Search for projects
type SearchProject struct {
}

func (self SearchProject) init() {
	fmt.Println("Search in projects")
}

// Search on project number
func (self SearchProject) project_number(user_number string) {
	// raise NotImplementedError("Function project_number is not implemented yet")
}

// Search on project description
func (self SearchProject) project_description(description string) {
	// raise NotImplementedError("Function project_description is not implemented yet")
}

// Search for projects
type SearchUser struct {
}

func (self SearchUser) init() {
	fmt.Println("Search in Users")
}

// Search on user number
func (self SearchUser) user_number(user_number string) {
	// raise NotImplementedError("Function user_number is not implemented yet")
}

// Search on user name
func (self SearchUser) user_name(user_name string) {
	// raise NotImplementedError("Function user_name is not implemented yet")
}

// Search on user first name
func (self SearchUser) user_first_name(user_first_name string) {
	// raise NotImplementedError("Function user_first_name is not implemented yet")
}

// Search on user last name
func (self SearchUser) user_last_name(user_last_name string) {
	// raise NotImplementedError("Function user_last_name is not implemented yet")
}

// Search on user role
func (self SearchUser) user_role(user_role string) {
	// raise NotImplementedError("Function user_role is not implemented yet")
}
