# The Filesystem page

The Filesystem is a major part of FreePDM. It is primarily a storage of files that are stored into the server. The primary goals are to import / export files, moving files, the revision of files, renaming of files, and all with access control. To make this happen you want to check-out and check-in of files, so that others know that you checked out a file and can't check-out the file without notifying you. The filesystem use directories, just like ordinary directories, to organize the files. The storage of files allows that you can also open a previous file.

For an example you can look at the test about the filesystem [here](../../tests/fileserver_setup.md)

## The Filesystem class
TheFilesystem class containes all the things that are needed to manipulate files inside the PDM. Essentially it is a storage of files. The backend is [pysftp ](https://pysftp.readthedocs.io/en/release_0.2.9/index.html), this is done to get access of the files, but also to have [SSH[(https://en.wikipedia.org/wiki/Secure_Shell)] and you can even use [sshfs](https://en.wikipedia.org/wiki/SSHFS), so you can access all your files remotely (in the LAN but also on the internet).

## Connecting and disconnect
`connect(server, user, passwd)` connects to the remote Filesystem and returns the SFTP filesystem or an error.

## Import / Export of files
`import_file(fname, dest_dir, descr, long_descr=None)` import a file or files inside the PDM. When you import a file the meta-data also gets imported, which means updated to the server. When you import a file or files you are placing the new file in the current directory. The new file inside the PDM gets a revision number automatically.

`export_file(fname, dest_dir)` export a file to a local directory.

## Create a directory
`mkdir(dname)` Creates a new directory inside the current directory, with the correct uid and gid.

## List a directory
`ls(dir="")` list a sorted directory of only the latest files.

## Check the latest version of a file (for adding another one)
`check_latest_file_version(fname)` Returns the number of the file, or -1 when the file doesn't exist.

## Revision of files
`revision_file(fname)` copy a file and increments revision number.

## Check-out / Check-in
`checkout_file(fname)` locks a file so that others can't accidentally check-in a different file.

`checkin_file(fname, descr, long_descr=None)`removes the locking but also uploads the file to the PDM. You need to write a description of what you did.

## Renaming of file
`rename_file(fname)` rename a file, for instance when he or she wants to use a file with a specified numbering system.

## Moving files
`move_file(fname, dest_dir)` moves a file to a different directory.


## Todo:
- [x] Add a filesystem (directories)
- [ ] Allow importing and exporting of files
- [ ] Copying files (incl. it's dependencies. Everything below and drawings, all optional). When it's not put into the Filesystem.
- [x] Dealing with versions.
- [ ] Dealing with and releases.
- [ ] Disable released files with a disable state.
- [ ] Dealing with states?
- [ ] Dealing with checkin and checkout.
- [ ] Importing and exporting files.
- [x] Dealing with users.
- [ ] Logging the activities.
