# FreePDM
***Concept Of Design***

## Design decisions

- Version management see:
- User interface
- Database
- Used programming language
- Attributes

### General setup
![Block diagram FreePDM general](./FreePDM_CoD-Figures/BDD_FreePMD-design.png)

### Version management

Version management is a difficult topic. A main issue is that it should be able to create releases of every individual part, assembly (and optional drawing.)  

There are serval different (Open source)version control systems available.
The main difficulty is that those are most written for software.  
Where software is released as a whole that is not always the case within the mechanical world. see: @ref to story

This has to be tackeled.  
Currently the most promissing implementation is the svn version system (a Client-server).  
The most wel known Distributed version systems are [git](https://git-scm.com/) and [mercurial](https://www.mercurial-scm.org/).

_Note: We have to decide if the model(part / assembly) get's an update,  the drawing get also an update. The same issue applied in revers order._

### Database

Within databases there are two main types of databases.
[SQL](https://en.wikipedia.org/wiki/SQL) based databases. Examples of SQL like databases are [SQLite](https://sqlite.org/index.html), [MariaDB](https://mariadb.org/), [PostgreSQL](https://www.postgresql.org/). 
<!--On https://sqlite.org/fileformat2.html is written:-->
<!--The main database file consists of one or more pages. The size of a page is a power of two between 512 and 65536 inclusive. All pages within the same database are the same size. The page size for a database file is determined by the 2-byte integer located at an offset of 16 bytes from the beginning of the database file.-->  
<!--Does this mean that it basically a big spreadsheet?-->

The other is [noSQL](https://en.wikipedia.org/wiki/NoSQL). An example of this database is [MongoDB](https://www.mongodb.com/)

_Note: Add What the difference is and how such a database is build up_

### User interface
![Block diagram interface design](./FreePDM_CoD-Figures/BDD_UI-design.png)

For the user interface there are two main issues buth both are related. the issue about the programming language is described below.


### Used programming language



### [Attributes](Attributes.md)


[<< Previous Chapter](FreePDM_Workflows.md) | [Content Table](FreePDM_CoD.md) | [Next Chapter >>](Attributes.md)
