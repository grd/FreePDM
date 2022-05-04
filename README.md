# FreePDM
A PDM for FreeCAD

This is going to be my first attempt to make a PDM for FreeCAD. And it is going to be in Python.

Yeah. It's a bit scary.

The initial idea is to make a Skeleton (model), GUI, a MVC and an Admin module. The GUI is based on pyQT.

This is the whish list:

## Admin
* Add a path to the location of the database
* Allow to make different numbering scheme's
* Allow to make different version types

## Skeleton
* Reading the files for meta-data and storing the meta-data
* Modifying the meta-data
* Dealing with the three assembly types
* Adding, Renaming and deleting files
* Replacing files
* Dealing with versions

## GUI
* Show a list of files
* Add buttons
* Edit text strings
* Show versions
* Edit versions
* Releasing
* Rename components (so that the assy knows that the component has been renamed)
* Replace components (so that the assy knows that the component has been renamed)
*

## Controller

## Custom made scripts
* For what you can do after releasing a file or files
