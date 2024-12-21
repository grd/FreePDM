// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package skeleton

import (
	"path"
)

// the file types

const (
	Directory    = "Directory"
	FCStd        = "FCStd"
	Assembly_A2P = "Assembly_A2P"
	Assembly_A3  = "Assembly_A3"
	Assembly_A4  = "Assembly_A4"
	Drawing      = "Drawing"
	Part         = "Part"
	Group        = "Group"
	ConfTable    = "ConfTable"
	FC_Document  = "FC_Document"
	FC_Old       = "FC_Old"
	Other        = "Other"
)

type DirectoryList struct {
	Nr          int // number
	Filename    string
	Description string // file description
	State       string // state of item
	Version     string
	// Filetype FileType   // Dir or File
	Filesize int64 // file size
}

type DirectoryModel struct {
	PurgeList     []string
	Directory     string
	DirectoryList []DirectoryList
}

func fullPath(directory, file string) string {
	return path.Join(directory, file)
}

// func (self *DirectoryModel) SetDirList(directory string, withParentDirectory bool) {
//     // get the filter types
//     // fc_files = self.conf.get_filter(config.show_fc_files_only)
//     // versioned_files = self.conf.get_filter(config.hide_versioned_fc_files)

//   var directory_list []string
//   var file_list []string

//   // TODO: Dit wordt een "normale" Go variant van ls
//    dirs, err := os.ReadDir(directory)
//    ex.CheckErr(err)

//   case kind:
//     of pcFile:
//       if path[0] != '.':
//         file_list.add(path)
//     of pcDir:
//       if path[0] != '.':
//         directory_list.add(path)
//     else: continue
//     // of pcLinkToFile:
//     //   echo "Link to file: ", path
//     // of pcLinkToDir:
//     //   echo "Link to dir: ", path
//     }
//   directory_list.sort()
//   file_list.sort()

//   // for e in directory_list: echo e
//   // for f in file_list: echo f

//   var
//     nr = 0
//     filetype FileType
//     filesize int64

//   // dealing with the first item
//   if withParentDirectory == true:
//     self.directoryList.add(DirectoryList(nr: nr, filename: "..", filetype: Directory, filesize: 0))
//     nr += 1

//   for directory in directory_list:
//     self.directoryList.add(DirectoryList(nr: nr, filename: directory, filetype: Directory, filesize: 0))
//     nr += 1

//   for item in file_list:
//     filesize = getfilesize(full_path(directory, item))
//     if ".FCStd" in item:
//       if item.endswith(".FCStd"):
//         var data = initItemDataModel(item)
//         // item = ItemDataModel(self.full_path(file))
//         // if "Assembly" in item.document_properties:
//         //   filetype = item.document_properties["Assembly"]
//         // else:
//         filetype = FCStd
//       // else:
//         // // versioned FCStd file
//         // if item[item.index(".FCStd")+6:].isdigit():
//         //   self.purge_list.add(self.full_path(item))
//         //   if versioned_files: // Filter out the versioned files
//         //     continue
//         //   filetype = FC Old
//         // else: // TODO: Check for other versioned file types
//         //   filetype = Other
//         self.directoryList.add(DirectoryList(nr: nr, filename: item, filetype: filetype, filesize: filesize))
//         nr += 1
//     else:
//       // Non FC files:
//       // if fc_files: // non-fc files out
//       //   continue
//       self.directoryList.add(DirectoryList(nr: nr, filename: item, filetype: Other, filesize: filesize))
//       nr += 1
// }

// func (self: DirectoryModel) Size() int {
//   return len(self.directoryList)
// }

// // func save_item_as(self: DirectoryModel, source, dest: string):
// //     self.source = source
// //     self.dest = dest
// //     raise NotImplementedError("The function 'save_item_as' is not implemented yet.")

// // The directorymodel is a list of files that can be used inside a GUI
// func initDirectoryModel(directory string, withParentDirectory bool) DirectoryModel {
//   result.directory = directory
//   readConfig()
//   result.set_dir_list(directory, withParentDirectory)
// }
