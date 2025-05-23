# FreePDM
***Concept Of Design***


## Database Shape

### Introduction

How to setup a database has a lot of impact on how the PDM system should interact.
To explain the two workflows here is an example:  

Let's assume there is need for a gearbox. The chosen gearbox is a [Apex Dynamics _AFH ###_](https://www.apexdyna.com/AFH_pro.aspx).
To handle that part there is some info required such as:

- A part / assembly file
  - optionally step file that is imported
- specifications
- (one or more )user manual(s)
- optionally a (one or more )repair manual(s)
- (one or more )calculation(s) for the required size
- optionally a drawing
- attributes about
  - filename
  - creation date
  - last edit date
  - buy / create /...
  - revision
  - reseller
  - etc

#### TODO

- Create a database setup! We need to know:
  - how to represent data,
  - how added data looks like
  - how to create a structure

### Assumptions

- This gearbox can be used in multiple projects.

### Workflow 1  <!-- a single big database -->

The first workflow uses a single database with all information.
A databases with all information is fine when the expected amount of information is known.
But as shown in the example, not all information is always available OR extra information is given.
In the case of extra information, results that _ALL_ items in the databases get an extra field(s) so subsequently the database grows very quickly.  
When, for example, an extra document is added does this need version control?

Another question is how to handle attributes for an item (like the apex gearbox). When the apex gearbox has a revision, it should be impossible to create/edit the attributes while other items should still be able to edit.
There are methods needed to iterate over the available attributes and ignore the information that is not needed.

### questions / Comments 1

- When iterating over the available attributes when the databases has more optional attributes it takes longer

### Workflow 2  <!-- a small databases for every item -->

This workflow uses a single (micro)database per item.
In these databases are stored, what information is available.
So only the fields that are available in that part are there and information that is not added is just not there in the database.
Because the databases are much more compact - iteration probably happens much faster.
With this method _ALL_ databases can be individually locked (with the svn versioning system).  
What adds more complexity is that with scattered databases (basically inside their [part environment](../FreePDM_03-2-SVNProjectStructure.md)) that there is no overall information what parts are available. This can be tackled with a databases with a link to all microdatabases.

A counter point for all small databases is that for the [Item view](04-UIFunctions.md) it is more difficult to implement because there are more unknowns.

#### questions / Comments 2

### Workflow 3  <!-- a layered databases for every item -->

A third option is try to create a 3D database.
For an easy comparison think about the database as a spreadsheet. Where every item has its own tab and the attributes are stored in columns and rows.  
Partly there is still the problem of locking (a part) of the database.
Even if this is probably less difficult since it locks a tab and not a few lines..

#### questions / Comments 3

[<< Previous Chapter](05-UIInteractionFC.md) | [Content Table](README.md) | [Next Chapter >>](07-DbInteraction.md)