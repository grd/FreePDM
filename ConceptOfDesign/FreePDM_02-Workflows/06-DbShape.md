# FreePDM
***Concept Of Design***


## Database Shape

### Introduction

How to setup a database makes a lot of impact on how the pdm system should interact.
To explain the two workflows here is an example.  

Let's assume there is need for a gearbox. The chosen gearbox is a [Apex Dynamics _AFH ###_](https://www.apexdyna.com/AFH_pro.aspx).
To handle that part there is some info required like

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

#### TODO:

- Create a database setup! we need to know
  - how to represent data,
  - how added data looks like
  - how to create a structure

### Assumptions

- This gearbox can be used in multiple projects.

### Workflow 1  <!-- a single big database -->

The first workflow uses a single database with all information.
A databases with all information is fine when the expected amount of information known.
But as shown in the example is not all information always available OR extra information is given.
The extra given information result that _ALL_ items in the databases get an extra field(s) so the database grow real quick.  
When for example an extra document is added does this need version control?

Another question is how to handle attributes for an item(like the apex gearbox). When the apex gearbox has a revision it should be impossible to create edit the attributes while other items should still be able to edit.
There are methods needed to iterate over the available attributes and ignore the information that is not needed.

### questions / Comments 1

- When iterating over the available attributes when the databases has more optional attributes it takes longer

### Workflow 2  <!-- a small databases for every item -->

This workflow uses a single (micro)database per item.
In these databases is stored what information is available.
So only the fields that are available in that part are there and information that is not added is just not there in the database.
Because the databases are much more compact iteration goes probably much faster.
With this method _ALl_ databases can be individually locked( with the svn versioning system).  
What this makes it more difficult is that with scattered databases(basically inside their [part environment](../FreePDM_03-2-SVNProjectStructure.md)) that there is no overall information what parts are available. This can be tackled with a databases with a link to all microdatabases.

A counter point for all small databases is that for the [Item view](04-UIFunctions.md) it is more difficult to implement because there are more unknowns.

#### questions / Comments 2

### Workflow 3  <!-- a layered databases for every item -->

A third option is try to create a 3D database.
For easy comparison think about the database as a spreadsheet. Where every item has it's own tab and the attributes are stored in columns and rows.  
Partly there is still the problem of locking (a part) of the database.
Even if this is probably less difficult since it locks a tab and not a few lines..

#### questions / Comments 3

[<< Previous Chapter](05-UIInteractionFC.md) | [Content Table](README.md) | [Next Chapter >>](07-DbInteraction.md)