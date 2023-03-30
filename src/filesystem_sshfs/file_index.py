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
from datetime import date

class FileIndex():
    """File Index Files in the root"""

    def init(self, vault_dir: str, user_uid, vault_uid: int):
        self.vault_dir = vault_dir
        self._user_uid = user_uid
        self._vault_uid = vault_uid
        if not path.isdir(vault_dir):
            raise FileNotFoundError("Directory " + vault_dir + " does not exist.")
        self.file_index: List[int, str, str]
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
                index, complete_file_name = line.split("=")
                file_name, rename_name = complete_file_name.split("<")
                if len(rename_name) > 0:
                    rename_name.removesuffix(">")
                self.file_index.append([int(index), file_name, rename_name])

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
        os.chown(self.all_files_txt, self._user_uid, self._vault_uid)

        with open(self.file_location_txt, "a") as file:
            file.write(str(index_len) + "=" + dir + "\n")
        os.chown(self.file_location_txt, self._user_uid, self._vault_uid)

        self.read() # again refreshing the index


    def rename_item(self, index: int, new_name: str):
        self.read() # refreshing index

        for item in self.file_index:
            if index == item[0]:
                old_filename = item[1]
                item[1] = new_name
                today = date.today()
                item[2] = "<" + old_filename + "," + today + ">"
                break

        with open(self.all_files_txt, "w") as file:
            for item in self.file_index:
                file.write(str(item[0]) + "=" + item[1] + item[2] + "\n")

        os.chown(self.all_files_txt, self._user_uid, self._vault_uid)


    def remove_item(self, item: str):
        pass

    def get_filename_index(self, name: str) -> int:
        pass

