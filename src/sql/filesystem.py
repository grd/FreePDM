#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""


class FileSystem():
    """File System related Class"""

    def init(self):
        print("Generic File System")

    def create_folder(self):
        """Create new folder"""
        # Create folder for project
        #  -> How to handle project tree?
        #  -> If on the same lever are projects and items: How to distinguish the difference?
        # create folder for item
        raise NotImplementedError("Function create_folder is not implemented yet")

    def create_file_idx(self):
        """Create file copy"""
        # create copy of file with index: for example 12345678.FCStd.2
        raise NotImplementedError("Function create_file_copy is not implemented yet")
