"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
from itemdatamodel import ItemDataModel

# the file types
types = ('Directory', 'Assembly A2P', 'Assembly A3', 'Assembly A4', 
         'Drawing', 'Part', 'Object', 'FC document', 'FC old', 'Other')

# The directorymodel is a list of files that can be used inside a GUI
class DirectoryModel(object):
    def __init__(self, directory, withParentDirectory = True):  # '..'
        self.directory = directory
        self.directoryList = [{"nr": "number",
                               "filename": "File name",
                               "description": "File description",
                               "state": "state of item",
                               "version": "0",
                               "type": "Dir or File",
                               "size": "file size"}]
        self.getDirList(withParentDirectory)

    def full_path(self, file):
        return (self.directory + '/' + file)

    def getDirList(self, withParentDirectory):
        dir_list = os.listdir(self.directory)
        dir_list.sort()
        directory_list = []
        file_list = []
        nr = 0
        for d in dir_list:
            full_item_path = self.full_path(d)
            # eliminate files and directories starting with "."
            if not d.startswith("."):
                if os.path.isdir(full_item_path):
                    directory_list.append(d)
                if os.path.isfile(full_item_path):
                    file_list.append(d)

        # Keep this in mind:

        # self.directoryList = [{"nr": "number",
        #                        "filename": "File name",
        #                        "description": "File description",
        #                        "state": "state of item",
        #                        "version": "0",
        #                        "type": "Dir or File",
        #                        "size": "file size"}]

        for dl in directory_list:
            self.directoryList.append({'nr': str(nr), 'filename': dl, 'type': 'Directory', 'size': ''})
            nr = nr + 1
        for fl in file_list:
            size = str(os.path.getsize(self.full_path(fl)))
            if fl.endswith(".FCStd"):
                self.directoryList.append({'nr': str(nr), 'filename': fl, 'type': 'FCStd', 'size': size})
                nr = nr + 1
            else:
                self.directoryList.append({'nr': str(nr), 'filename': fl, 'type': 'File', 'size': size})
                nr = nr + 1
 
        # dealing with the first item
        if withParentDirectory == True:
            self.directoryList[0] = ({'nr': str(0), 'filename': '..', 'type': 'Directory', 'size': ''})
        else:
            self.directoryList.pop(0)

    def size(self):
        return(len(self.directoryList))

# TODO:
#   What kind of FC files do you want to see?
#   Do you want to see the history of the files?
#   def filter(self):

# TODO:
#   Purge stored versions of files. All up to number
#   def purge(number):
