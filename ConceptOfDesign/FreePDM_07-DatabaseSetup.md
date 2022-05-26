# FreePDM
***Concept Of Design***


## Database Design

<!--
| Column 1 | Column 2 | Column 3|
|----------|----------|---------|
| item1    | item2    | item 3  |
-->

### user


| User Id         | Name | Last name | Email             | Phone number    | Department    |
|-----------------|------|-----------|-------------------|-----------------|---------------|
| N digit number? | Word | words     | name@something.st | 12 digit Number | Word or words |



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

| FileID | (Item number?) | Item Path        | Filename| File type     | Material   |
|--------|----------------|------------------|---------|---------------|------------|
| Number | ######         | /Project/######/ | Words   | db Attributes | None       |
| Number | ######         | /Project/######/ | Words   | Part          | Yes        |
| Number | ######         | /Project/######/ | Words   | Assembly      | None       |
| Number | ######         | /Project/######/ | Words   | Drawing       | None       |
| Number | ######         | /Project/######/ | Words   | Manual        | None       |
| Number | ######         | /Project/######/ | Words   | Specification | Optionally |
| Number | ######         | /Project/######/ | Words   | Calculation   | None       |
| Number | ######         | /Project/######/ | Words   | Image         | None       |
| Number | ######         | /Project/######/ | Words   | etc..         | None       |

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

- [Databases in plain english](https://www.freecodecamp.org/news/sql-and-databases-explained-in-plain-english/)
- [Geek for Geeks SQL tutorial](https://www.geeksforgeeks.org/sql-tutorial/)
- [Freecodecamp relational database](https://www.freecodecamp.org/learn/relational-database/)

[<< Previous Chapter]() | [Content Table](README.md) | [Next Chapter >>]()
