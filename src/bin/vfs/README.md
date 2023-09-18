# The Virtual Filesystem page

The Virtual Filesystem (vfs) is a representation of the file structures. I think that I am using [this](https://github.com/C2FO/vfs) to start with. In the end it is an executable vfs that can be mounted in the server. The vfs is read-writeable. It consists of multiple levels that are highly complicated.

### Level 1 (and two)

![Level 1](images/Level%201.png)
![Level 2](images/Level%202.png)

This are an ordinary directory structures. Users or admins should have write access here. Possible actions: Create, move, delete, rename and copy.

### Level 3

![Level 3](images/Level%203.png)

This is the "file" level, but you are only seeing there directories and that is because the files are stored at level 5. You can notice that the FreeCAD or other file has some extra text and that text is the number of versions and if the file is checked out then it also shows who checked the file out.

Possible actions: Import a file (with drag and drop), rename, move, copy, delete. With an extra (command line) tool : create a new version, check-out and check-in.

### Level 4

![Level 4](images/Level%204.png)

This shows the versions.

### Level 5

![Level 5](images/Level%205.png)

This level shows the files. The files are read-only by default (and some of the version info files are always read-only) but when you check-out a file then you can modify the file, for instance when you open a FreeCAD file you are allowed to save it. After that you can check the file in.

## Todo:

- [ ] Show the entire filesystem
- [ ] Level 3: Import a file with drag and drop
- [ ] Level 3: Perform file actions: Move, rename, copy and delete
- [ ] Level 3: Show the version numbers and file checkout status
- [ ] Level 3: Show the correct filename
- [ ] Test the VFS
