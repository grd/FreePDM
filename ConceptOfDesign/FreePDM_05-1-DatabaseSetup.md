# FreePDM
***Concept Of Design***


## Database Design

The database stores data from CAD files, including links to where the files are stored.
Besides this it contains Part, user, project attributes etc.
For more information about the attributes or a shortlist see [Attributes](FreePDM_03-1-Attributes.md).

<!--
| Column 1 | Column 2 | Column 3|
|----------|----------|---------|
| item1    | item2    | item 3  |
-->

### user

| User Id         | User Name | Name  | Last name | Email             | Phone number   | Department    | Role(s) | Aliases |
|-----------------|-----------|-------|-----------|-------------------|----------------|---------------|---------|---------|
| N digit number? | Words     | words | words     | name@something.st | N digit Number | Word or words | List(?) | List(?) |

### Projects

| Project ID | Project Number | Path         | Column 3|
|------------|----------------|--------------|---------|
| Number     |N digit number  | /Project/... | item 3  |

### Item 

List of all item numbers

| Item number    | Path         | Part of Project | Number of Files Attached |
|----------------|--------------|-----------------|--------------------------|
| N digit number | /Project/... | item 3          | N Digit Number           |

_Notes:_

- _Do we allow that there exist files that don't have(And never get) any related CAD file but are accepted in the PDM? For example a financial document?_
- _Do we allow non "Just the next number approach"? For example that (Module) assemblies get a certain range OR that a project get a certain range(of for example 1000 numbers)_

### File(s) belongs to Item

Are we creating one single database list for _ALL_ files.
OR Are we creating a database file for every Item?

General attributes

| FileID | Item No. | Item Name  | Description | Description long | File type     |Item Path         | Filename|
|--------|----------|------------|-------------|------------------|---------------|------------------|---------|
| Number | ######   | Words / Nr | Short term  | Multiple terms   | db Attributes | /Project/######/ | Words   |
| Number | ######   | Words / Nr | Short term  | Multiple terms   | Part          | /Project/######/ | Words   |
| Number | ######   | Words / Nr | Short term  | Multiple terms   | Assembly      | /Project/######/ | Words   |
| Number | ######   | Words / Nr | Short term  | Multiple terms   | Drawing       | /Project/######/ | Words   |
| Number | ######   | Words / Nr | Short term  | Multiple terms   | Manual        | /Project/######/ | Words   |
| Number | ######   | Words / Nr | Short term  | Multiple terms   | Specification | /Project/######/ | Words   |
| Number | ######   | Words / Nr | Short term  | Multiple terms   | Calculation   | /Project/######/ | Words   |
| Number | ######   | Words / Nr | Short term  | Multiple terms   | Image         | /Project/######/ | Words   |
| Number | ######   | Words / Nr | Short term  | Multiple terms   | etc..         | /Project/######/ | Words   |

Example of a part can be:  

- Item number is: 123456
- Item name is often same as Item number so: 123456
- Description is: Screw
- Description long: ISO 4762 M8x32, A2 - Hexagon Socket Head Cap Screw
- Filetype: Library

Material Attributes

| ... | File type     | Material   | Surface finish | Volume     | Mass       | Weight     | Surface Area |
|-----|---------------|------------|----------------|------------|------------|------------|--------------|
| ... | db Attributes | None       | None           | None       | None       | None       | None         |
| ... | Part          | Yes        | Yes            | Yes        | Yes        | Yes        | Yes          |
| ... | Assembly      | None       | Optionally     | Yes        | Optionally | Yes        | Yes          |
| ... | Drawing       | None       | None           | None       | None       | None       | None         |
| ... | Manual        | None       | Optionally     | Optionally | Optionally | Optionally | Optionally   |
| ... | Specification | Optionally | Optionally     | Optionally | Optionally | Optionally | Optionally   |
| ... | Calculation   | None       | None           | None       | None       | None       | None         |
| ... | Image         | None       | None           | None       | None       | None       | None         |
| ... | etc..         | None       | None           | None       | None       | None       | None         |

_Note: The Attributes needs probably some mixing._  
_Note2: What about the units? Dimensionless, Iso by default? Optionally accept more dimensions later on, even if i don't like to add thousand different options like 103 pounds and 2.3 ounce..._

Historical Attributes

| ... | File type     | Date created | Created by | Last Edit  | Edited by | Checked-out by | Revision state| Revision nr. |
|-----|---------------|--------------|------------|------------|-----------|----------------|----------------|-------------|
| ... | db Attributes | Date         | Name       | Date       | Name      | Name           | List of States | Number      |
| ... | Part          | Date         | Name       | Date       | Name      | Name           | List of States | Number      |
| ... | Assembly      | Date         | Name       | Date       | Name      | Name           | List of States | Number      |
| ... | Drawing       | Date         | Name       | Date       | Name      | Name           | List of States | Number      |
| ... | Manual        | Date         | Name       | Date       | Name      | Name           | List of States | Number      |
| ... | Specification | Date         | Name       | Date       | Name      | Name           | List of States | Number      |
| ... | Calculation   | Date         | Name       | Date       | Name      | Name           | List of States | Number      |
| ... | Image         | Date         | Name       | Date       | Name      | Name           | List of States | Number      |
| ... | etc..         | Date         | Name       | Date       | Name      | Name           | List of States | Number      |

Revision state:

- New( Item)
- Concept
- Revised / Revision
- work
- Not for new Projects
- Depreciated

Buy / Manufacturing(Purchasing)

| ... | File type     | Traceability | Source | Manufacturer ID | Manufacturer Name | Vendor ID | Vendor Name |
|-----|---------------|--------------|--------|-----------------|-------------------|-----------|-------------|
| ... | db Attributes | None         | List   | Number          | Name              | Number    | Name        |
| ... | Part          | List         | List   | Number          | Name              | Number    | Name        |
| ... | Assembly      | List         | List   | Number          | Name              | Number    | Name        |
| ... | Drawing       | List         | List   | Number          | Name              | Number    | Name        |
| ... | Manual        | List         | List   | Number          | Name              | Number    | Name        |
| ... | Specification | List         | List   | Number          | Name              | Number    | Name        |
| ... | Calculation   | List         | List   | Number          | Name              | Number    | Name        |
| ... | Image         | none         | List   | Number          | Name              | Number    | Name        |
| ... | etc..         | List         | List   | Number          | Name              | Number    | Name        |

Traceability state:

- Lot
- Lot And Serial Number
- Serial Number
- Not traced

Source state:

- Buy
- Make

### References

Difference

- [SO difference between different types of sql](https://stackoverflow.com/questions/1326318/difference-between-different-types-of-sql)
- [compare multiple sql flavours](https://www.altexsoft.com/blog/business/comparing-database-management-systems-mysql-postgresql-mssql-server-mongodb-elasticsearch-and-others/)
- [Hevodata MariaDB vs PostgreSQL](https://hevodata.com/learn/mariadb-vs-postgresql/)
- [IBM SQL vs NO-SQL](https://www.ibm.com/cloud/blog/sql-vs-nosql)
- [Freecodecamp relational vs non-relational](https://www.freecodecamp.org/news/relational-vs-nonrelational-databases-difference-between-sql-db-and-nosql-db/)
- [Techtarget SQL vs No-sql vs new-sql](https://www.techtarget.com/whatis/feature/SQL-vs-NoSQL-vs-NewSQL-How-do-they-compare)

Explanation

- [Databases in plain english](https://www.freecodecamp.org/news/sql-and-databases-explained-in-plain-english/)
- [Geek for Geeks SQL tutorial](https://www.geeksforgeeks.org/sql-tutorial/)
- [Geek for Geeks types of databases](https://www.geeksforgeeks.org/types-of-databases/)
- [Freecodecamp relational database](https://www.freecodecamp.org/learn/relational-database/)

[<< Previous Chapter](FreePDM_05-Architecture.md) | [Content Table](README.md) | [Next Chapter >>](FreePDM_06-Roadmap.md)
