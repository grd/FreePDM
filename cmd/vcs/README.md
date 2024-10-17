# The vcs tool

The vcs (version control system) is a tool that interacts with the vfs (virtual filesystem). It is roughly based onto https://github.com/Abdulsametileri/vX for the "status" part, but for the rest this will be compeletly different.

## Commands

`vcs status file`
`vcs newversion file`
`vcs update file`
`vcs checkout file`
`vcs checkin file`
`vcs history file`

Notes:
- The 'file' can be any file, it doesn't need to be a FreeCAD file.
- I don't know exactly how to implement `vcs checkin file` yet. This probably will look rather different.
- Update means updating the file info from inside the file

## Todo:

- [ ] Status of a file
- [ ] New version of a file
- [ ] Check-Out and Check-In of a file

## Future idea:

- Maybe also add a search in file fields