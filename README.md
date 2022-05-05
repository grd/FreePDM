# FreePDM
A PDM for FreeCAD

This is going to be my first attempt to make a PDM. And it is going to be in Python.

Yeah. It's a bit scary.

The initial idea is to make a Skeleton (model), GUI, a MVC and an Admin module. The GUI is based on pySide. I need to make a WB.

Looking at previous attempts to make a PDM:
* https://github.com/cadracks-project/openplm (but this is going to be a "poor man's" PDM. You can compare it with SolidWorks PDM)
* 
This is the whish list:

## Admin
* Add buttons
* Add a path to the location of the files or database
* Allow to make different numbering scheme's
* Allow to make different version types

## Skeleton
* Reading the files for meta-data and storing the meta-data
* Modifying the meta-data
* Dealing with users. Who did what?
* Partnumbering
* Adding an owner, checker, approver
* Dealing with the history of the parts
* Dealing with the -three- assembly types
* Dealing with other info (such as LidbreOffice files and other media)
* Adding, Renaming and deleting files
* Replacing files
* Dealing with versions
* Disable files with a disable state
* Dealing with states (temperary ones, before putting into a release)
* Dealing with checkin and checkout
* Importing and exporting files
* Dealing with projects (creating, modifying, closing, re-opening)
* Dealing with a database (PostgreSQL, MongoDB, trython, sqlite ?) I think about PostgreSQL since I now know about SQLAlchemy. But sqlite to start with.
* Use SQLAlchemy https://en.wikipedia.org/wiki/SQLAlchemy

## GUI
* Show a list of files
* Add buttons
* Dealing with projects (creating, modifying, closing, re-opening)
* Add a search button
* Add all the things done in the skeleton (such as edit text strings, show versions, edit versions etc.)
* Releasing (by who?)
*

## The Web?
* I think that it is inevitable at some time point but from now, I don't think that is is necessary
* qtwebkit ?


## Controller

## Custom made scripts
* For what you can do after releasing a file or files
