#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2023 by the FreePDM team
    :license:   MIT License.
"""

import os
from os import path
import sys
from pathlib import Path
from typing import List

class FileIndex():
    """File Index Files in the root"""

    def __init__(self, vault_dir: str):
        self.vault_dir = vault_dir
        if not path.isdir(vault_dir):
            raise FileNotFoundError("Directory " + vault_dir + " does not exist.")
        self.file_index: List[int, str]
        self.file_location_list: List[int, str]
        self.all_files_txt = path.join(self.vault_dir, "All Files.txt")
        self.file_location_txt = path.join(self.vault_dir, "FileLocation.txt")

        if not path.isfile(self.all_files_txt):
            raise FileNotFoundError("File " + self.all_files_txt + " does not exist.")
        if not path.isfile(self.file_location_txt):
            raise FileNotFoundError("File " + self.file_location_txt + " does not exist.")

        self.read() # read the indexes


    def read(self):
        with open(self.all_files_txt, "r") as file:
            while (line := file.readline()):
                index, file_name = line.split("=")
                self.file_index.append([int(index), file_name])

        with open(self.file_location_txt, "r") as file:
            while (line := file.readline()):
                index, path_name = line.split("=")
                self.file_location_list.append([int(index), path_name])


    def add_item(self, name, dir: str):
        """ This code only adds items to the index."""
        self.read() # refreshing the index
        index_len = len(self.file_index)

        with open(self.all_files_txt, "a") as file:
            file.write(str(index_len) + "=" + name + "\n")

        with open(self.file_location_txt, "a") as file:
            file.write(str(index_len) + "=" + dir + "\n")

        self.read() # again refreshing the index


    def rename_item(self, item, new_item, date: str):
        pass

    def remove_item(self, item: str):
        pass

    def get_filename_index(self, name: str) -> int:
        pass

