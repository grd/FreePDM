// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"gorm.io/gorm"
)

// Note:
// According to: https://docs.sqlalchemy.org/en/14/tutorial/engine.html
// 'The engine is typically a global object created just once for a particular database server, ...'
//
// So up till now i expected every user has it's own login,
// but it looks like this is not possible using engines.
// Now there has to be some research in a dedicated login system!

// Item related
type Item struct {
	gorm.Model
	number          int
	Project         string
	Path            string
	Name            string
	Description     string
	FullDescription string
}

// Create new project number
// number [str] : Last current number
//
//	Number as string including leading zeros.
//
// ndigits [int] : number of digits
//
//	Number of digits of the number length.
//	If ndigits is -1 the length is just the length.
func (self Item) CreateItemNumber() int {
	call_proj := Project{} // BUG: This needs to be something that is running because of the number which would be zero.
	item_nr := call_proj.CreateNumber()

	return item_nr
}

// https://stackoverflow.com/questions/73887390/handle-multiple-users-login-database-with-sqlalchemy

func (self *Item) CreateItem(project, path string, number *int, name, description, full_description string) {
	self.Project = project // User works on current project
	self.Path = path       // TODO: create path automatically
	self.Name = name
	self.Description = description
	self.FullDescription = full_description
	// TODO: How to handle other related properties

	if number == nil {
		// TODO: get latest number and ndigits from conf / db
		self.number = self.CreateItemNumber()
	}

	// proj := new(Project)
	// self.project_id = proj.get_id(self.project)
	// // self.project_id = Select()  // get project id based on project number / project name
	// new_item = PdmItem(item_number=self.number, item_name=self.name, item_description=self.description, item_full_description=self.full_description, path=self.path, project_id=self.project_id)

	// // TODO: Import Engine - From where?
	// Session.configure(bind=engine, future=True)

	// // https://docs.sqlalchemy.org/en/14/orm/session_basics.html#id1
	// with Session() as session:
	//     try:
	//         session.add(new_item)
	//     except:
	//         Session.rollback()
	//     finally:
	//         Session.close()
}

// Remove existing item
func (self *Item) RemoveItem() {
	// check if item is new (local == no state)
	// -> if True user can remove
	// -> if False Check if user == admin
	//   -> if True User can remove item
	//   -> if False warning message
}

// Update existing item
func (self *Item) UpdateItem() {
	// raise NotImplementedError("Function update_item is not implemented yet")
}

// Update existing item
func (self *Item) AddItemImage() {
	// TODO: Auto generate image from models

	// raise NotImplementedError("Function add_item_image is not implemented yet")
}

// When inheritance not everything. Do i need a base class?

// Model related
type Model struct {
}

// Create new model
func (self Model) CreateModel() {
	// create copy with iter: 0
	// Create model for:
	// -> Existing item
	// -> For new Item
	//    -> With new item also create item

	// raise NotImplementedError("Function create_model is not implemented yet")
}

// Remove existing model
func (self Model) RemoveModel() {
	// check if model is new (local == no state)
	// -> if True user can remove
	// -> if False Check if user == admin
	//   -> if True User can remove item
	//   -> if False warning message

	// raise NotImplementedError("Function remove_model is not implemented yet")
}

// Update existing model
func (self Model) UpdateModel() {
	// create copy with iter: N

	// raise NotImplementedError("Function update_model is not implemented yet")
}

// Get model that is not latest version
func (self Model) GetVersion( /*save_iter*/ ) {
	// in UI set only release versions or all
	// Optional two versions to compare
	// is FC able to reload when model is changed?

	// raise NotImplementedError("Function get_version is not implemented yet")
}

// Document related
type Document struct {
	// How much difference is there between a document and a model?
}

// Item / Model / Document Ownership states
type OwnerStates struct { // Access States
	// Can all checkin options performed from a central class?
}

// Check in Items, Models, Documents
func (self OwnerStates) CheckIn( /*objects*/ ) {
	// check if new item?
	// create copy (for Model, Document) and add copy to DataBase

	// raise NotImplementedError("Function check_in is not implemented yet")
}

// Check in Items, Models, Documents
func (self OwnerStates) CheckOut( /*objects*/ ) {
	// check latest version (only for Models and Documents)
	// check if checked-out by other user

	// raise NotImplementedError("Function check_out is not implemented yet")
}

// Check in Items, Models, Documents
func (self OwnerStates) CheckInCheckOut( /*objects*/ ) {
	// check if new item?
	// create copy (for Model, Document) and add copy to DataBase
	//
	// checkout checks are not needed

	// raise NotImplementedError("Function check_in_check_out is not implemented yet")
}

// Item / Model /Document release states struct
type ReleaseStates struct {
}

// new Item, Model, Document
func (self ReleaseStates) ChangeReleaseState() {
	// All new items, models, documents have state new - until they are checked in.

	// raise NotImplementedError("Function new is not implemented yet")
}

// New Item, Model, Document
func (self ReleaseStates) New( /*objects*/ ) {
	// All new items, models, documents have state new - until they are checked in.

	// raise NotImplementedError("Function new is not implemented yet")
}

// Prototype Item, Model, Document
func (self ReleaseStates) Prototype( /*objects*/ ) {
	// All items, models, documents get state prototype on first checkin - until they are released.

	// raise NotImplementedError("Function prototype is not implemented yet")
}

// Check in Items, Models, Documents
func (self ReleaseStates) Release( /*objects*/ ) {
	// All items, models, documents get state release after  - until they are released.

	// raise NotImplementedError("Function check_in_check_out is not implemented yet")
}

// Check in Items, Models, Documents
func (self ReleaseStates) NotForNew( /*objects*/ ) {
	// check latest version (only for Models and Documents)
	// check if checked-out by other user

	// raise NotImplementedError("Function check_out is not implemented yet")
}

// Check in Items, Models, Documents
func (self ReleaseStates) Depreciated( /*objects*/ ) {
	// check latest version (only for Models and Documents)
	// check if checked-out by other user

	// raise NotImplementedError("Function check_out is not implemented yet")
}
