"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
from pathlib import Path
import sys
sys.path.append(os.fspath(Path(__file__).resolve().parents[1]))
from filesystem.filesystem import FileSystem

from skeleton.itemdatamodel import ItemDataModel
from skeleton.config import *

# the file types
types = ('Directory', 'Assembly A2P', 'Assembly A3', 'Assembly A4', 
         'Drawing', 'Part', 'Group', 'ConfTable', 'Object', 
         'FC document', 'FC old', 'Other')

# The directorymodel is a list of files that can be used inside a GUI
class DirectoryModel(object):
    def __init__(self, fs: FileSystem, directory: str):  # '..'
        self.conf = conf()
        self.fs = fs
        self.conf.read()
        self.directory: str = directory
        self.directoryList = [{"nr": "number",
                               "filename": "File name",
                               "description": "File description",
                               "state": "state of item",
                               "version": "0",
                               "type": "Dir or File",
                               "size": "file size"}]
        self.set_dir_list()
 
    def full_path(self, file: str) -> str:
        return os.path.join(self.directory, file)

    def set_dir_list(self):
        dir_list = self.fs.listdir()
        directory_list = []
        file_list = []
        print("dir_list = " + str(dir_list))
        print("ls = " + str(os.listdir()))
        nr = 0
        for directory in dir_list:
            if directory.startswith("d: "):
                directory = directory[3:]
                self.directoryList.append({'nr': str(nr), 'filename': directory, 'type': 'Directory', 'size': ''})
            else:
                self.directoryList.append({'nr': str(nr), 'filename': directory})
                nr += 1

    def size(self) -> int:
        return(len(self.directoryList))

