// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
)

// Basically searching in the SQL database requires a coonetcion to the database too.
// Is it than a sub-class or better to request access to a class or has this class a interface too?

// Search
// https://docs.sqlalchemy.org/en/14/orm/session_basics.html
type Search struct {
}

func NewSearch() *Search {
	fmt.Println("Search in Databases")
	return &Search{}
}

// Search on number
func (s Search) SearchNumber(number int) {
	// raise NotImplementedError("Function search_number is not implemented yet")
}

// Search on description
func (s Search) SearchDescription(description string) {
	// raise NotImplementedError("Function search_description is not implemented yet")
}

// Search on something else
func (s Search) SearchSomething_else(something string) {
	// raise NotImplementedError("Function search_something_else is not implemented yet")
}

// help function
func (s Search) SearchHelp() string {
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

func NewSearchItem() *SearchItem {
	fmt.Println("Search in items")
	return &SearchItem{}
}

// Search on project number
func (s SearchItem) ItemNumber(user_number string) {
	// raise NotImplementedError("Function project_number is not implemented yet")
}

// Search on project description
func (s SearchItem) ItemDescription(description string) {
	// raise NotImplementedError("Function project_description is not implemented yet")
}

// Search for projects
type SearchProject struct {
}

func NewSearchProject() *SearchProject {
	fmt.Println("Search in projects")
	return &SearchProject{}
}

// Search on project number
func (s SearchProject) ProjectNumber(user_number string) {
	// raise NotImplementedError("Function project_number is not implemented yet")
}

// Search on project description
func (s SearchProject) ProjectDescription(description string) {
	// raise NotImplementedError("Function project_description is not implemented yet")
}

// Search for projects
type SearchUser struct {
}

func NewSearchUser() *SearchUser {
	fmt.Println("Search in Users")
	return &SearchUser{}
}

// Search on user number
func (s SearchUser) UserNumber(user_number string) {
	// raise NotImplementedError("Function user_number is not implemented yet")
}

// Search on user name
func (s SearchUser) UserName(user_name string) {
	// raise NotImplementedError("Function user_name is not implemented yet")
}

// Search on user first name
func (s SearchUser) UserFirstFame(user_first_name string) {
	// raise NotImplementedError("Function user_first_name is not implemented yet")
}

// Search on user last name
func (s SearchUser) UserLastName(user_last_name string) {
	// raise NotImplementedError("Function user_last_name is not implemented yet")
}

// Search on user role
func (s SearchUser) UserRole(user_role string) {
	// raise NotImplementedError("Function user_role is not implemented yet")
}
