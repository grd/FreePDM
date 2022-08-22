"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
from itemdatamodel import ItemDataModel
import config as conf

# the file types
types = ('Directory', 'Assembly A2P', 'Assembly A3', 'Assembly A4', 
         'Drawing', 'Part', 'Group', 'ConfTable', 'Object', 
         'FC document', 'FC old', 'Other')

# The directorymodel is a list of files that can be used inside a GUI
class DirectoryModel(object):
    def __init__(self, directory, withParentDirectory = True):  # '..'
        conf.read()
        self.directory = directory
        self.directoryList = [{"nr": "number",
                               "filename": "File name",
                               "description": "File description",
                               "state": "state of item",
                               "version": "0",
                               "type": "Dir or File",
                               "size": "file size"}]
        self.get_dir_list(withParentDirectory)

    def full_path(self, file):
        return os.path.join(self.directory, file)

    def get_dir_list(self, withParentDirectory):
        # get the filter types
        fc_files = conf.get_filter(conf.show_fc_files_only)
        versioned_files = conf.get_filter(conf.hide_versioned_fc_files)

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
                    if file[file.index('.FCStd')+6:].isdigit(): 
                        # versioned FCStd file
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

    def size(self):
        return(len(self.directoryList))

    #   What kind of FC files do you want to see?
    #   Do you want to see the history of the files?
    def filter(self):
        pass

    #   Purge stored versions of files.
    def purge(self, number):
        self.number = number
        pass

    def save_item_as(self, source, dest):
        self.source = source
        self.dest = dest
        pass
