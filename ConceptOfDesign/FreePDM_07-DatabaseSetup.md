# FreePDM
***Concept Of Design***


## Database Design

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

| FileID | (Item number?) | Item Path        | Filename| File type     |
|--------|----------------|------------------|---------|---------------|
| Number | ######         | /Project/######/ | Words   | db Attributes |
| Number | ######         | /Project/######/ | Words   | Part          |
| Number | ######         | /Project/######/ | Words   | Assembly      |
| Number | ######         | /Project/######/ | Words   | Drawing       |
| Number | ######         | /Project/######/ | Words   | Manual        |
| Number | ######         | /Project/######/ | Words   | Specification |
| Number | ######         | /Project/######/ | Words   | Calculation   |
| Number | ######         | /Project/######/ | Words   | Image         |
| Number | ######         | /Project/######/ | Words   | etc..         |

Material Attributes

| ... | File type     | Material   | Surface finish | Volume     | Mass       | Weight     |
|-----|---------------|------------|----------------|------------|------------|------------|
| ... | db Attributes | None       | None           | None       | None       | None       |
| ... | Part          | Yes        | Yes            | Yes        | Yes        | Yes        |
| ... | Assembly      | None       | Optionally     | Yes        | Optionally | Yes        |
| ... | Drawing       | None       | None           | None       | None       | None       |
| ... | Manual        | None       | Optionally     | Optionally | Optionally | Optionally |
| ... | Specification | Optionally | Optionally     | Optionally | Optionally | Optionally |
| ... | Calculation   | None       | None           | None       | None       | None       |
| ... | Image         | None       | None           | None       | None       | None       |
| ... | etc..         | None       | None           | None       | None       | None       |

_Note: The Attributes needs probably some mixing._  
_Note2: What about the units? Dimensionless, Iso by default? Optionally accept more dimensions later on, even if i don't like to add thousand different options like 103 pounds and 2.3 ounce..._

As long as every 

### Attributes by File

TODO:

Last Edit By
Created By
Checked-Out by
RevisionNumber

| Column 1 | Column 2 | Column 3|
|----------|----------|---------|
| item1    | item2    | item 3  |

_Note: Maybe remove this header and add All Info in multiple Tables Below Files_

### Attributes by Item

TODO:

| Column 1 | Column 2 | Column 3|
|----------|----------|---------|
| item1    | item2    | item 3  |
|----------|----------|---------|

_Note: Maybe remove this header and add All Info in multiple Tables Below Items_

### References

Difference

- [SO difference between different types of sql](https://stackoverflow.com/questions/1326318/difference-between-different-types-of-sql)
- [compare multiple sql flavours](https://www.altexsoft.com/blog/business/comparing-database-management-systems-mysql-postgresql-mssql-server-mongodb-elasticsearch-and-others/)
- [Hevodata MariaDB vs PostgreSQL](https://hevodata.com/learn/mariadb-vs-postgresql/)
- [IBM SQL vs NO-SQL](https://www.ibm.com/cloud/blog/sql-vs-nosql)
- [Freecodecamp relational vs non-relational](https://www.freecodecamp.org/news/relational-vs-nonrelational-databases-difference-between-sql-db-and-nosql-db/)
- [Techtarget SQL vs No-sql vs new-sql](https://www.techtarget.com/whatis/feature/SQL-vs-NoSQL-vs-NewSQL-How-do-they-compare)

Explnation

- [Databases in plain english](https://www.freecodecamp.org/news/sql-and-databases-explained-in-plain-english/)
- [Geek for Geeks SQL tutorial](https://www.geeksforgeeks.org/sql-tutorial/)
- [Geek for Geeks types of databases](https://www.geeksforgeeks.org/types-of-databases/)
- [Freecodecamp relational database](https://www.freecodecamp.org/learn/relational-database/)

[<< Previous Chapter]() | [Content Table](README.md) | [Next Chapter >>]()
