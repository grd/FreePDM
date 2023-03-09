"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
from skeleton.itemdatamodel import ItemDataModel
from skeleton.config import *

# the file types
types = ('Directory', 'Assembly A2P', 'Assembly A3', 'Assembly A4', 
         'Drawing', 'Part', 'Group', 'ConfTable', 'Object', 
         'FC document', 'FC old', 'Other')

# The directorymodel is a list of files that can be used inside a GUI
class DirectoryModel(object):
    def __init__(self, directory, withParentDirectory = True):  # '..'
        self.conf = conf()
        self.conf.read()
        self.purge_list = list()
        self.directory = directory
        self.directoryList = [{"nr": "number",
                               "filename": "File name",
                               "description": "File description",
                               "state": "state of item",
                               "version": "0",
                               "type": "Dir or File",
                               "size": "file size"}]
        self.set_dir_list(withParentDirectory)
 
    def full_path(self, file: str) -> str:
        return os.path.join(self.directory, file)

    def set_dir_list(self, withParentDirectory):
        # get the filter types
        fc_files = self.conf.get_filter(show_fc_files_only)
        versioned_files = self.conf.get_filter(hide_versioned_fc_files)

        dir_list = os.listdir(self.directory)
        dir_list.sort()
        directory_list = []
        file_list = []
        nr = 0
        for directory in dir_list:
            full_item_path = self.full_path(directory)
            # eliminate files and directories starting with "."
            if not directory.startswith("."):
                if os.path.isdir(full_item_path):
                    directory_list.append(directory)
                if os.path.isfile(full_item_path):
                    file_list.append(directory)

        # Keep this in mind:

        # self.directoryList = [{"nr": "number",
        #                        "filename": "File name",
        #                        "description": "File description",
        #                        "state": "state of item",
        #                        "version": "0",
        #                        "type": "Dir or File",
        #                        "size": "file size"}]

        for directory in directory_list:
            self.directoryList.append({'nr': str(nr), 'filename': directory, 'type': 'Directory', 'size': ''})
            nr += 1
        for file in file_list:
            size = str(os.path.getsize(self.full_path(file)))
            if ".FCStd" in file:
                if file.endswith(".FCStd"):
                    item = ItemDataModel(self.full_path(file))
                    if 'Assembly' in item.document_properties:
                        type = item.document_properties['Assembly']
                    else:
                        type = 'FCStd'
                else: 
                    # versioned FCStd file
                    if file[file.index('.FCStd')+6:].isdigit(): 
                        self.purge_list.append(self.full_path(file))
                        if versioned_files: # Filter out the versioned files
                            continue
                        type = 'FC old'
                    else: # TODO: Check for other versioned file types
                        type = 'Other'
                self.directoryList.append({'nr': str(nr), 'filename': file, 'type': type, 'size': size})
                nr += 1
            else:
                # Non FC files:
                if fc_files: # non-fc files out
                    continue
                self.directoryList.append({'nr': str(nr), 'filename': file, 'type': 'File', 'size': size})
                nr += 1

        # dealing with the first item
        if withParentDirectory is True:
            self.directoryList[0] = ({'nr': str(0), 'filename': '..', 'type': 'Directory', 'size': ''})
        else:
            self.directoryList.pop(0)

    def size(self) -> int:
        return(len(self.directoryList))

    # Purge stored versions of FreeCAD files.
    def purge(self):
        for i in self.purge_list:
            os.remove(i)

    def save_item_as(self, source, dest):
        self.source = source
        self.dest = dest
        raise NotImplementedError("The function 'save_item_as' is not implemented yet.")
