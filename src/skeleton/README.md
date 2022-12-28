# FreePDM Skeleton
A PDM for FreeCAD


## Skeleton
- [ ] Reading the files for meta-data and storing the meta-data in the DB
- [ ] Modifying the meta-data. This can be a bit tricky.
- [ ] Select material and surface finish and calculate the weight of the part
- [ ] What belongs into the Filesystem?
    - [ ] Partnumbering (both inside the Filesystem and in local files itself)
    - [ ] Adding an owner, checker, approver
    - [ ] Dealing with the history of the parts
    - [ ] Dealing with the -three- assembly types. I am sure that they all have their own quirks. The answer of that question is Yes.
    - [ ] Dealing with other info (such as LidbreOffice files and other media)
    - [ ] Adding, Renaming and deleting files. Renaming can be tricky when the file is used in assemblies.
    - [ ] Replacing files in assemblies. Both into the Filesystem and the local files.
    - [ ] Copying files (incl. it's dependencies. Everything below and drawings, all optional). When it's not put into the Filesystem.
    - [ ] Dealing with versions and releases.
    - [ ] Disable released files with a disable state.
    - [ ] Dealing with states
    - [ ] Dealing with checkin and checkout.
    - [ ] Importing and exporting files.
    - [ ] Dealing with projects (creating, modifying, closing, re-opening).
    - [ ] Dealing with users. Who did what?

