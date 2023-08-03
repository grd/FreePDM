# The Filesystem page

The Filesystem is a major part of FreePDM. It is primarily a storage of files that are stored into the server. The primary goals are to import / export files, moving files, the revision of files, renaming of files, and all with access control. To make this happen you want to check-out and check-in of files, so that others know that you checked out a file and can't check-out the file without notifying you. The filesystem use directories, just like ordinary directories, to organize the files. The storage of files allows that you can also open a previous file.

~~ For an example you can look at the test about the filesystem [here](../../tests/fileserver_setup.md) ~~ 

## Todo:

- [x] Add a filesystem (directories)
- [x] Allow importing of files
- [ ] Allow showing and editing files with SMB
- [ ] Copying files (latest only).
- [x] Moving files
- [x] Renaming files
- [x] Dealing with versions.
- [ ] Dealing with releases.
- [ ] Disable versioned files with a disable state (with a negative number)
- [x] Dealing with checkin and checkout.
- [x] Dealing with users.
- [ ] Reporting with json (dir list, properties, etc.)
- [ ] Logging the activities.
- [ ] Hovering over an item (show list of description, checked out, etc.)