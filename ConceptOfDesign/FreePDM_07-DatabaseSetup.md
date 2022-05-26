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

| Item number    | Path         | Part of Project |
|----------------|--------------|-----------------|
| N digit number | /Project/... | item 3          |

_Notes:_

- _Do we allow that there exist files that don't have(And never get) any related CAD file but are accepted in the PDM? For example a financial document?_
- _Do we allow non "Just the next number approach"? For example that (Module) assemblies get a certain range OR that a project get a certain range(of for example 1000 numbers)_

### File(s) belongs to Item

Are we creating one single database list for _ALL_ files.
OR Are we creating a database file for every Item?

| FileID | (Item number?) | Item Path        | Filename| File type      |
|--------|----------------|------------------|---------|----------------|
| Number | ######         | /Project/######/ | Words   | db Attributes  |
| Number | ######         | /Project/######/ | Words   | Part / Assembly|
| Number | ######         | /Project/######/ | Words   | Drawing        |
| Number | ######         | /Project/######/ | Words   | Manual         |
| Number | ######         | /Project/######/ | Words   | Specification  |
| Number | ######         | /Project/######/ | Words   | Calculation    |
| Number | ######         | /Project/######/ | Words   | Image          |
| Number | ######         | /Project/######/ | Words   | etc..          |

As long as every 

### Attributes by File

TODO:

| Column 1 | Column 2 | Column 3|
|----------|----------|---------|
| item1    | item2    | item 3  |


### Attributes by Item

TODO:

| Column 1 | Column 2 | Column 3|
|----------|----------|---------|
| item1    | item2    | item 3  |
|----------|----------|---------|

### References

- [Databases in plain english](https://www.freecodecamp.org/news/sql-and-databases-explained-in-plain-english/)
- [Geek for Geeks SQL tutorial](https://www.geeksforgeeks.org/sql-tutorial/)
- [Freecodecamp relational database](https://www.freecodecamp.org/learn/relational-database/)

[<< Previous Chapter]() | [Content Table](README.md) | [Next Chapter >>]()
