# The Filesystem page

The Filesystem is a major part of FreePDM. It is primarily a storage of files that are stored into the server. The primary goals are to import / export files, moving files, the revision of files, renaming of files, and all with access control. To make this happen you want to check-out and check-in of files, so that others know that you checked out a file and can't check-out the file without notifying you. The filesystem use directories, just like ordinary directories, to organize the files. The storage of files allows that you can also open a previous file.

For an example you can look at the test about the filesystem [here](../../tests/fileserver_setup.md)

## The Filesystem class
TheFilesystem class containes all the things that are needed to manipulate files inside the PDM. Essentially it is a storage of files. The backend is [pysftp ](https://pysftp.readthedocs.io/en/release_0.2.9/index.html), this is done to get access of the files, but also to have [SSH[(https://en.wikipedia.org/wiki/Secure_Shell)] and you can even use [sshfs](https://en.wikipedia.org/wiki/SSHFS), so you can access all your files remotely (in the LAN but also on the internet).



## Todo:
- [x] Add a filesystem (directories)
- [ ] Allow importing and exporting of files
- [ ] Copying files (incl. it's dependencies. Everything below and drawings, all optional). When it's not put into the Filesystem.
- [ ] Dealing with versions.
- [ ] Dealing with and releases.
- [ ] Disable released files with a disable state.
- [ ] Dealing with states?
- [x] Dealing with checkin and checkout.
- [ ] Importing and exporting files.
- [x] Dealing with users.
- [ ] Logging the activities.
