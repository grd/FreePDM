// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

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

var number int

// Create new Project number
// number is only added when it has a value.
func createNumber() int {

	// increase the number
	number += 1

	return number
}

// Create new project
func NewProject(number *int, name string, status *string, path string) (*Project, error) {
	proj := new(Project)
	proj.name = name
	proj.path = path // TODO: create path automatically
	// TODO: How to handle other related properties

	if number == nil {
		// TODO: get latest number and ndigits from conf / db
		proj.number = createNumber()
	} else {
		proj.number = *number
	}

	if status == nil {
		proj.status = New
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

	return proj, nil
}

// Remove existing project
func (p *Project) RemovePproject() {
	// raise NotImplementedError("Function remove_project is not implemented yet")
}

// Update existing project
func (p *Project) UpdateProject() {
	// raise NotImplementedError("Function update_model is not implemented yet")
}

// Add user to project
func (p *Project) AddUserToProject() {
	// raise NotImplementedError("Function add_user_to_project is not implemented yet")
}

// Remove user from project
func (p *Project) RemoveUserFromProject() {
	// raise NotImplementedError("Function remove_user_from_project is not implemented yet")
}
