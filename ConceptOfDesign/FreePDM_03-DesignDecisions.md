# FreePDM
***Concept Of Design***

## Design decisions

- Version management see:
- User interface
- Database
- Used programming language
- Attributes

### General setup
![Block diagram FreePDM general](FreePDM_CoD-Figures/BDD_FreePMD-design.png)

### Version management

Version management is a difficult topic. A main issue is that it should be able to create releases of every individual part, assembly (and optional drawing.) 

There are serveral different (open source) version control systems available.
The main difficulty is that these software are most targetting software development.
The difference is that, where software is released as a whole project, that is not always the case within the mechanical world. see: @ref to story 
_Note Grd: But there are also closed source versions. I think that we should also accept those because in a company they might only allow this kind of software._
This has to be tackled. 
Currently the most promising implementation is the SVN version system (a Client-server). 
The most well known Distributed version systems are [git](https://git-scm.com/) and [mercurial](https://www.mercurial-scm.org/).
_Note Grd: Why does it need to be tackled? You can also allow them all. But in that case we need to have an admin page in where we can select which one we use. They all work roughly the same, but we need to make all work of course. But we are gonna start with SVN._

_Note: We have to decide if the model(part / assembly) gets an update, the drawing get also an update. The same issue applied in revers order._
Note grd: This is never gonna work. When a part is changed you aso need to update the drawing of that part, the assy (or assies) in where it belongs and the drawings of those too.

### Database

Within databases there are two main types of databases.
[SQL](https://en.wikipedia.org/wiki/SQL) based databases. Examples of SQL like databases are [SQLite](https://sqlite.org/index.html), [MariaDB](https://mariadb.org/), [PostgreSQL](https://www.postgresql.org/). 
<!--On https://sqlite.org/fileformat2.html is written:-->
<!--The main database file consists of one or more pages. The size of a page is a power of two between 512 and 65536 inclusive. All pages within the same database are the same size. The page size for a database file is determined by the 2-byte integer located at an offset of 16 bytes from the beginning of the database file.--> 
<!--Does this mean that it basically a big spreadsheet?-->

The other is [noSQL](https://en.wikipedia.org/wiki/NoSQL). An example of this database is [MongoDB](https://www.mongodb.com/)

_Note: Add What the difference is and how such a database is build up_

_Note Grd: Why do we need such a database? That is the question. We need it because we can't extract the relevant information out of VCS and that is because a FC file attributes are stored somewhere into that file. And you want to be able to search that metadata information. To be honest, I don't care about which kind of software you use (SQL / noSQL), but [SQLAlchemy](https://en.wikipedia.org/wiki/SQLAlchemy) helps a lot when you use SQL data. That is why I think that we better use SQL. Starting with sqlite._


### User interface
![Block diagram interface design](FreePDM_CoD-Figures/BDD_UI-design.png)

For the user interface there are two main issues but both are related. the issue about the programming language is described below.
What i mean with the figure<!--(report me(==Jee-Bee) when i'm wrong also explain why)--> above is that there are two main options of interfaces.

1. Web interface
2. Local stored interface

Looking at the web interface has some pro's and cons

Pro's                      | Con's
-------------------------- | --------------------------
Centralised updates        | Security is more difficult(specially .js)
More easy to add other SW  | More data transport
                           | 
(Maintainability?)           | (Maintainability?)

See also Comment of user1234 in the related [topic](https://forum.freecad.org/viewtopic.php?f=8&t=68350&p=594331#p594252)

Looking at the local storage there are two tastes available(for now).

- An independent tool
- Inside FreeCAD (for example [_Add-on manager_](https://wiki.freecad.org/Std_AddonMgr))
  - If Using the _Add-on Manager_ see: the [Workebench start kit](https://github.com/FreeCAD/freecad.workbench_starterkit)
  - Another important Thread about the [_Add-on Manager_ Redesign](https://forum.freecad.org/viewtopic.php?f=9&t=64628)

Specially the Add-on manager can help for at least (semi-)centralised updates. 

_Note: How does the GUI interact with the user?_

### Used programming language

In the case of the web interface, it can be programmed with Javascript or Python (there are more but these are most well known).

In the case of a local program, Python can be a good option too.
This is also due to its interoperating well with FreeCAD and the FreeCAD [Add-on Manager](https://wiki.freecad.org/Std_AddonMgr). 

What are pros/cons of the Add-on manager: 

Pros                           | Cons
-------------------------------| --------------------------
Semi Centralised updates       | Not direct integration
Fast(er) update speed⁰         | python only
Translation to other languages | More difficult user handling
 Maintainability               | 

⁰ Faster is compared to FreeCAD when every chance goes through checks(Except when WMayer fix it). <!-- unicode superscript see: https://stackoverflow.com/questions/15155778/superscript-in-markdown-github-flavored -->
 
_What i meant with user handling see story later on_

By using the add-on manager it is possible to change quickly when a change / update is needed.
_Note: We have to decide to:
- Make a FreeCAD Add-on
- ?_


### [Attributes](FreePDM_03-1-Attributes.md)

### [SVN extension](FreePDM_03-2-SVNProjectStructure.md)

### [DB Versioning Update](FreePDM_03-3-DBVersioningUpd.md)

[<< Previous Chapter](FreePDM_02-Workflows.md) | [Content Table](README.md) | [Next Chapter >>](FreePDM_03-1-Attributes.md)
