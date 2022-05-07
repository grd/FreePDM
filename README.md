# FreePDM
A PDM for FreeCAD

This is going to be my first attempt to make a PDM. And it is going to be in Python.

Yeah. It's a bit scary.

The initial idea is to make a Skeleton (model), GUI, a MVC and an Admin module. The GUI is based on pySide. I need to make a WB.

Looking at previous attempts to make a PDM:
* https://github.com/cadracks-project/openplm (but FreePDM is going to be a "poor man's" PDM. You can compare it with SolidWorks PDM. Nothing fancy. That all brokes down sooner or later)
* 

Also interesting:
* https://github.com/furti/FreeCAD-Reporting
* The fcinfo macro for measure weight of a model

This is the whish list:

## Admin
* Add buttons
* Add roles for the type of use, such as Architecture, mechanical, etc) 
* Add a path to the location of the files or data database (SVN, Git, MS Azure, etc)
* Add a path to the meta-data database (A SQL server such as PostgreSQL, MySQL, etc)
* Allow to make different numbering scheme's
* Allow to make different version types (this also depends on the data database)

## Skeleton
* Reading the files for meta-data and storing the meta-data in the DB
* Modifying the meta-data
* Select material and surface finish and calculate the weight of the part
* Dealing with users. Who did what? (this can be dealth with the meta-data DB)
* Partnumbering
* Adding an owner, checker, approver
* Dealing with the history of the parts (this can be dealth with the meta-data DB)
* Dealing with the -three- assembly types. I am sure that they all have their own quirks
* Dealing with other info (such as LidbreOffice files and other media)
* Adding, Renaming and deleting files. Renaming can be tricky when the file is used in assemblies.
* Replacing files in assemblies
* Dealing with versions (IMO this can be dealt with the data DB)
* Disable files with a disable state (again in the data DB)
* Dealing with states (temperary ones, before putting into a release), again in the data DB
* Dealing with checkin and checkout (in the data DB and meta-data DB)
* Importing and exporting files (in the data DB and meta-data DB)
* Dealing with projects (creating, modifying, closing, re-opening) in the data DB
* Dealing with a database (PostgreSQL, MongoDB, trython, sqlite ?) I think about PostgreSQL since I now know about SQLAlchemy. But sqlite to start with.
* Use SQLAlchemy https://en.wikipedia.org/wiki/SQLAlchemy

## GUI
* Add buttons
* Show a list of files
* Dealing with projects (creating, modifying, closing, re-opening)
* Add a search button
* Add all the things done in the skeleton (such as edit text strings, show versions, edit versions etc.)
* Releasing (by who?)
*

## The Web?
* I think that it is inevitable at some time point but from now, I don't think that is is necessary ATM
* qtwebkit ?
* The big problem is that the web is a big pile of BS, and securing the web can become a nightmare.


## Controller
* The glue between the skeleton and the GUI

## Custom made scripts
* For what you can do after releasing a file or files
* Probably a thousand more

