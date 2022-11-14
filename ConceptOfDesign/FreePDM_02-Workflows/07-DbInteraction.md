# FreePDM
***Concept Of Design***


## User / File / Database / Version management interaction 

### Introduction

There is a complex interaction between the User, the file(s), the database and the version management system needed.  
(Most of) these interactions are noted in the workflows below.
_All_ functionalities interact with each other.

### Assumptions

What is happening with the version handling and the database.

- The database that is part of the item can be locked
- The database that has knowledge of all available items not
- The database that is related to the CAD files is part of the versioning system
- If a CAD file is released, it's related (database) item is also released
- From the database interaction the version system gives the user write access
- From the database interaction the version system releases the user write access.

### Workflow 1  <!-- Edit a file -->

A user wants to edit a file, here is the logic that occurs: 

1. The database needs to be checked if the file is `writable()`
1. The database needs to be checked if the user is allowed to check-out the file
    * Not in release state.
    * Not Checked-Out by another user
1. Give the user write access in the version system AND block all other users.

After finishing the edits: 

1. Save the file to the server
1. (Automatically) edit the database Attributes regarding changes
1. Check-in the file
1. Release the write access AND release the block of other users

### questions / Comments 1

- Within which databases is set who checked-out a file? Is that on an object level or at main level?

### Workflow 2  <!-- Edit an Attribute -->

A user wants to edit an Attribute (what is part (not CAD part) of the database), the logic is involved: 

1. The database needs to be checked if the item where the attributes belongs to are writable (AKA not locked).
    * Not in release state.
    * Not Checked-Out by another user
1. Give the user write-access in the version system AND block all other users.

After finishing the edits: 

1. Save the changed attributes to the server
1. (Automatically) edit the database Attributes regarding changes
1. Check-in the (database) item
1. Release the write-access AND release the block of other users

#### questions / Comments 2

- This step should also work for multiple attributes from multiple file(s) (and filetypes).

### Workflow 3  <!-- Release an Item -->

A user wants to release an Item (basically the CAD model, CAD drawing, database Item, and other documentation): 

1. The database needs to checked what files are already released
1. Fulfill the Engineerings Change Process (not described, but basically a set of checks)
3. Database Change all Attributes that are affected to a higher revision number.
4. Database Change all Attributes that are affected to a _Release_ state.
5. The affected database items are locked AND the affected files (file locations) are locked for _All_ users.

#### questions / Comments 3

- This step should also work for a whole module / assembly structure.


### Workflow 4  <!-- Change an Item -->

A user wants to change an previously released Item. For example, a CAD file.

1. The user requests a change to the needed file(s) AND the (database) item
1. (Optionally) accepted by supervisor
1. (Optionally) a copy is made from the released version so a changed version can always be reverted.
1. Database Change all Attributes that are affected to an _InWork_ state
1. [Workflow 1](#workflow-1)
1. Optionally [Workflow 2](#workflow-2)
1. [Workflow 3](#workflow-3).

#### questions / Comments 4

- This step should also work for a whole module / assembly structure.

[<< Previous Chapter](06-DbShape.md) | [Content Table](README.md) | [Next Chapter >>](08-DevelopmentBuilds.md)
