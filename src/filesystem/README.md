# The Filesystem page

The Filesystem is a major part of FreePDM. It is primarily a storage of files that are immutable. The primary goals are to import / export files, moving files, the revision of files, renaming of files, and all with access control. To make this happening you want to check-out and check-in of files, so that others know that you checked out a file and can't check-out the file without notifying you. The filesystem use directories, just like ordinary directories, to organize the files. The storage of files allows that you can also open a previous file. For access control an admin account needs to be used. He or she needs extra privileges, for instance when someone left and that guy checked out a file, then the admin needs to be able to fix that. The admin also sets the rules for the (automatic) renaming of files for instance and the things thet need to be make when someone changed a revision, such as to automatically generate a PDF for a drawing or generate STEP/STL files.

## The Filesystem class
TheFilesystem class containes all the things that are needed to manipulate files inside the PDM. Essentially it is a storage of files.

## Import / Export of files
`Filesystem.import` allows to import a file or files inside the PDM. When you import a file the meta-data also gets imported. The local files remain untouched. When you import a file or files you need to set a directory and a description. The new file inside the PDM gets a revision number automatically.

`Filesystem.export` allows to export a file of files to a certain directory.

## Revision of files
`Filesystem.revision` sets a file of files to a new level that is according to the previous file or files. It allows to automatically generates PDF and STEP/STL files.

## Opening a revisioned file
`Filesystem.openrevisionend` allows to export a revisoned file (an older version of a file). We all make mistakes so opening an older file can be crucial.

## Check-out / Check-in
`Filesystem.checkoutfile` locks a file so that others can't accidentally check-in a different file.

`Filesystem.checkinfile`removes the locking but also uploads the file to the PDM. You need to write a description of what you did.

## Renaming of file
`Filesystem.renamefile` can be used to rename a file, for instance when he or she wants to use a file with a specified numbering system.

## Moving files
`Filesystem.movefile` can be used to move a file to a different directory.

## Access control
Not written yet.


## Todo:
- [ ] Add a filesystem (directories)
- [ ] Allow importing and exporting of files
- [ ] Hashing of files
- [ ] Partnumbering (both inside the Filesystem and in local files itself). This also means renaming of files and where used. The "where used" functionality can only be used inside the PDM because FreeCAD doesn't have such a feature.
- [ ] Dealing with other info (such as LibreOffice files and other media)
- [ ] Replacing files in assemblies. Both into the Filesystem and the local files. This also relies on renaming of files.
- [ ] Copying files (incl. it's dependencies. Everything below and drawings, all optional). When it's not put into the Filesystem.
- [ ] Dealing with versions and releases.
- [ ] Disable released files with a disable state.
- [ ] Dealing with states
- [ ] Dealing with checkin and checkout. (the Filesystem should be read-only)
- [ ] Importing and exporting files.
- [ ] Dealing with projects (creating, modifying, closing, re-opening).
- [ ] Dealing with users. Who did what?
