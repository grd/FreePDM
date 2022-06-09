# FreePDM
***Concept Of Design***


## User / File / Database / Version management interaction 

### Introduction

There is a complex interaction between the User, the file(s), the database and the version management system needed.  
(Most of )these interactions are written down in the workflows below.
_All_ functionality interact with each other.

### Assumptions

What is happening with the version handling and the database.

- The database that is part of the item can be locked
- The database that have knowledge of all available items not
- The database that is related to the CAD files is part of the versioning system
- If a CAD file is released it's related (database )item is also released
- From the database interaction the version system give the user write access
- From the database interaction the version system release the user write access.

### Workflow 1  <!-- Edit a file -->

An user want to edit a file, what has to be done: 

1. The database need to be checked if the file is writable()
2. The database need to be checked if the user is aloud to check-out the file
  -   Not in release state.
  -   Not Checked-Out by another user
3. Give the user write access in the version system AND block all other users.

After finishing the edits 

1. Save the file to the server
2. (Automatically )edit the database Attributes regarding changes
3. Check-in the file
4. Release the write access and release the block of other users

### questions / Comments 1

- In what databases is set who checked-out a file? Is that on object level or at main level.

### Workflow 2  <!-- Edit an Attribute -->

An user want to edit a Attribute(what is part(not CAD part) of the database), what has to be done: 

1. The database need to be checked if the item where the attributes belongs to is writable(not locked).
  -   Not in release state.
  -   Not Checked-Out by another user
3. Give the user write access in the version system AND block all other users.

After finishing the edits 

1. Save the changed attributes to the server
2. (Automatically )edit the database Attributes regarding changes
3. Check-in the (database )item
4. Release the write access and release the block of other users

#### questions / Comments 2

- This stap should also work for multiple attributes from multiple File(s)( and filetypes).

### Workflow 3  <!-- Release an Item -->

An user want to release a Item(basically the CAD model, Cad drawing, database Item, and other documentation): 

1. The database need to be checked what files are already released
2. Full-fill the Engineerings Change Process(Not described, but basically a set of checks)
3. Database Change all Attributes that are affected to a higher revision number.
4. Database Change all Attributes that are affected to a _Release_ state.
5. The affect database items are locked AND the affected files(file locations) are locked for _all_ users.

#### questions / Comments 3

- This stap should also work for a whole module / assembly structure.


### Workflow 4  <!-- Change an Item -->

An user want to change a earlier released Item. For example a CAD file.

1. The user request a change to the needed file(s) AND the (database )item
2. (Optionally) Accept by supervisor
3. (Optionally) a copy is made from the released version so a changed version can always be reverted.
4. Database Change all Attributes that are affected to an _InWork_ state
5. [Workflow 1](#workflow-1)
6. Optionally [Workflow 2](#workflow-2)
7. [Workflow 3](#workflow-3).

#### questions / Comments 4

- This stap should also work for a whole module / assembly structure.

[<< Previous Chapter](06-DbShape.md) | [Content Table](README.md) | [Next Chapter >>](08-DevelopementBuilds.md)
