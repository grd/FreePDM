# The Filesystem page

The Filesystem is a major part of FreePDM. It is primarily a storage of files that are stored into the server. The primary goals are to import / export files, the revision of files, moving -, copying and renaming of files, and all with access control. To make this happen you want to check-out and check-in of files, so that others know that you checked out a file and can't modify / edit the file without notifying you. The filesystem use directories, just like ordinary directories, to organize the files. The storage of files allows that you can also open a previous file.


## Todo:

- [x] Add a filesystem (directories)
- [x] Allow importing of files
- [x] Allow showing and editing files with SMB
- [x] Renaming files
- [x] Copying files
- [x] Moving files
- [ ] Renaming directories
- [ ] Moving directories
- [ ] Copying directories
- [x] Dealing with versions.
- [ ] Disable / enable versioned files (mode 0700)
- [x] Dealing with checkin and checkout.
- [x] Dealing with users.
- [ ] Reporting with json (dir list, properties, etc.)
- [x] Logging the activities.
- [ ] Hovering over an item (show list of description, checked out, etc.)