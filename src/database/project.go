// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package database

import "gorm.io/gorm"

// Basically searching in the SQL database requires a connection to the database too.
// Is it than a sub-class or better to request access to a class or has this class a interface too?

// Project related struct
type Project struct {
	gorm.Model
	number int
	name   string
	status ProjectState
	path   string
}

// Get id with Project number
func (self Project) GetId(number string) int {
	// raise NotImplementedError("Function get_id is not implemented yet")

	return -1 // BUG: Watch out !!!
}

// Create new Project number
// number is only added when it has a value.
func (self *Project) CreateNumber() int {

	// increase the number
	self.number += 1

	return self.number
}

// Create new project
func (self *Project) CreateProject(number *int, name string, status *string, path string) {
	self.name = name
	self.path = path // TODO: create path automatically
	// TODO: How to handle other related properties

	if number == nil {
		// TODO: get latest number and ndigits from conf / db
		self.number = self.CreateNumber()
	} else {
		self.number = *number
	}

	if status == nil {
		self.status = New
	}

	// new_project = PdmProject(self.number, self.name, self.status, self.date_start, self.date_finish, self.path)

	// // TODO: Import Engine - From where?
	// Session.configure(bind=engine, future=True)

	// // https://docs.sqlalchemy.org/en/14/orm/session_basics.html#id1
	// with Session() as session:
	//     try:
	//         session.add(new_project)
	//     except:
	//         Session.rollback()
	//     finally:
	//         Session.close()
}

// Remove existing project
func (self *Project) RemovePproject() {
	// raise NotImplementedError("Function remove_project is not implemented yet")
}

// Update existing project
func (self *Project) UpdateProject() {
	// raise NotImplementedError("Function update_model is not implemented yet")
}

// Add user to project
func (self *Project) AddUserToProject() {
	// raise NotImplementedError("Function add_user_to_project is not implemented yet")
}

// Remove user from project
func (self *Project) RemoveUserFromProject() {
	// raise NotImplementedError("Function remove_user_from_project is not implemented yet")
}
