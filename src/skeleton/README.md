# FreePDM Skeleton
A PDM for FreeCAD


## Skeleton
- [ ] Reading the files for meta-data and storing the meta-data in the DB
- [ ] Modifying the meta-data
- [ ] Select material and surface finish and calculate the weight of the part
- [ ] Dealing with users. Who did what? (this can be dealth with the meta-data DB)
- [ ] Partnumbering (both inside the DB and in FreeCAD itself)
- [ ] Adding an owner, checker, approver
- [ ] Dealing with the history of the parts (this can be dealth with the meta-data DB)
- [ ] Dealing with the -three- assembly types. I am sure that they all have their own quirks. The answer of that question is Yes.
- [ ] Dealing with other info (such as LidbreOffice files and other media)
- [ ] Adding, Renaming and deleting files. Renaming can be tricky when the file is used in assemblies.
- [ ] Replacing files in assemblies
- [ ] Copying files (incl. it's dependencies. Everything below and drawings, all optional)
- [ ] Dealing with versions and releases
- [ ] Disable files with a disable state (again in the data DB)
- [ ] Dealing with states (temporary ones, before putting into a release), again in the data DB
- [ ] Dealing with checkin and checkout
- [ ] Importing and exporting files (in the data DB and meta-data DB)
- [ ] Dealing with projects (creating, modifying, closing, re-opening) in the data DB
- [ ] Dealing with a database (PostgreSQL, MariaDB, SQLite). The preferred one is PostgreSQL.
- [ ] Use SQLAlchemy https://en.wikipedia.org/wiki/SQLAlchemy