# The Filesystem page

The Filesystem is a major part of FreePDM. It is primarily a storage of files that are stored into the server. The primary goals are to import / export files, the revision of files, moving -, copying and renaming of files, and all with access control. To make this happen you want to check-out and check-in of files, so that others know that you checked out a file and can't modify / edit the file without notifying you. The filesystem use directories, just like ordinary directories, to organize the files. The storage of files allows that you can also open a previous file.


## Todo:

- [x] Add a filesystem (directories)
- [x] Allow importing of files
- [x] Allow showing and editing files with SMB
- [x] Renaming files
- [x] Copying files
- [x] Moving files
- [x] Renaming directories
- [x] Copying directories
- [x] Dealing with versions.
- [x] Disable / enable versioned files (mode 0700)
- [x] Dealing with checkin and checkout.
- [x] Dealing with users.
- [ ] Reporting with json (dir list, properties, etc.)
- [x] Logging the activities.
- [x] Checks about checkout during file operations (rename, copy) and new versions.
- [ ] Make it working for multi-user, multi-vaults
- [ ] Checks about CheckIn comments (descr and longDescr)
- [ ] Checks about the VER.txt in file versions
- [ ] Change directory and file structure to Read Only for security. That way no accidental issues can happen.
