// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Basically searching in the SQL database requires a connection to the database too.
// Is it than a sub-class or better to request access to a class or has this class a interface too?

// Project represents the projects table.
type Project struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	ProjectName string `gorm:"size:100"`
	number      string
	name        string
	status      ProjectState
	path        string
}

// Get id with Project number
func (p Project) GetId(number string) int {
	// raise NotImplementedError("Function get_id is not implemented yet")

	return -1 // BUG: Watch out !!!
}

// var number int

// Create new Project number
// number is only added when it has a value.
// CreateNumber generates a new project number based on the last current number.
// number: Last current number as a string, including leading zeros.
// ndigits: Number of digits for the new number length. If -1, the length is just the length of the number.
func CreateNumber(number string, ndigits *int) (string, error) {
	if ndigits == nil {
		return "", fmt.Errorf("value for 'ndigits' can't be nil")
	}

	num, err := strconv.Atoi(number)
	if err != nil {
		return "", err
	}
	num++

	if *ndigits != -1 {
		// Count the current digits in the number
		counter := len(strconv.Itoa(num))

		// Create leading zeros
		leadingZeros := strings.Repeat("0", *ndigits-counter)

		return leadingZeros + strconv.Itoa(num), nil
	}
	return strconv.Itoa(num), nil
}

// Create new project
func NewProject(number string, name string, status *string, path string) (*Project, error) {
	proj := new(Project)
	proj.name = name
	proj.path = path // TODO: create path automatically
	// TODO: How to handle other related properties

	if number == "" {
		// TODO: get latest number and ndigits from conf / db
		// var err error
		// // proj.number, err = CreateNumber(last_number, ndigits)
		// if err != nil {
		// 	return nil, err
		// }
	} else {
		proj.number = number
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
