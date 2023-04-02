#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2023 by the FreePDM team
    :license:   MIT License.
"""

import os
from os import path
from typing import List
from datetime import datetime

""" TODO: 1. Combine both indexes into one, but only after the index
             is working completely and is tested!!! """

""" TODO: 2. Optimize the index by making it a tuple and
             reading from / writing to that tuple binary."""


class FileIndex():
    """File Index Files in the root"""

    def init(self, vault_dir: str, user_uid, vault_uid: int):
        self.vault_dir = vault_dir
        self._user_uid = user_uid
        self._vault_uid = vault_uid
        if not path.isdir(vault_dir):
            raise FileNotFoundError("Directory " + vault_dir + " does not exist.")

        """
            The self.file_index is a List with three fields:
            0. The number. This is a integer. This is the name of the directory inside the PDM.
            1. The file name.
            2. The previous name of the file including the date when the file was renamed.
        """
        self.file_index: List[int, str, str] = []

        """
            The self.file_location_list is a List with two fields:
            0. The number. This is a integer. This is the name of the directory.
            1. The location, in what directory the file was stored, offset from '<VaultDir>/PDM/'.
        """
        self.file_location_list: List[int, str] = []

        self.all_files_txt = path.join(self.vault_dir, "All Files.txt")
        self.file_location_txt = path.join(self.vault_dir, "FileLocation.txt")

        self.index_number: int
        self.index_number_txt = path.join(self.vault_dir, "IndexNumber.txt")

        if not path.isfile(self.all_files_txt):
            raise FileNotFoundError("File " + self.all_files_txt + " does not exist.")
        if not path.isfile(self.file_location_txt):
            raise FileNotFoundError("File " + self.file_location_txt + " does not exist.")

        if not path.isfile(self.index_number_txt):
            raise FileNotFoundError("File " + self.index_number_txt + " does not exist.")

        self.read() # read the indexes


    def read(self):
        """ Reads the values from both "All Files.txt" and "FileLocation.txt",
            and stores the data into memory. """

        with open(self.all_files_txt, "r") as file:
            while (line := file.readline()):
                index, complete_file_name = line.split("=")
                if complete_file_name.find("<") != -1:
                    file_name, rename_name = complete_file_name.split("<")
                    if len(rename_name) > 0:
                        rename_name.removesuffix(">")
                    self.file_index.append([int(index), file_name, rename_name])
                else:
                    self.file_index.append([int(index), complete_file_name, ""])


        with open(self.file_location_txt, "r") as file:
            while (line := file.readline()):
                index, path_name = line.split("=")
                self.file_location_list.append([int(index), path_name])


    def _increase_index_number(self) -> int:
        """Increases the index number, that is stored in the file 'IndexNumber.txt' 
           in the root directory of the vault."""

        with open(self.index_number_txt, "r") as file:
            self.index_number = int(file.readline())

        with open(self.index_number_txt, "w") as file:
            file.write(str(self.index_number + 1))

        os.chown(self.index_number_txt, self._user_uid, self._vault_uid)

        print("The index number = " + str(self.index_number))
 
        return self.index_number

 
    def add_item(self, name, dir: str) -> int:
        """ Adds the item from both self.file_index and self.file_location_list
            and store the information on disk. Returns the index number.

            It does not add a file on disk. """

        self.read() # refreshing the index

        index_len = self._increase_index_number()
        self.file_index.append([index_len, name, ""])
        self.file_location_list.append([index_len, dir])

        with open(self.all_files_txt, "a") as file:
            file.write(str(index_len) + "=" + name + "\n")
        os.chown(self.all_files_txt, self._user_uid, self._vault_uid)

        with open(self.file_location_txt, "a") as file:
            file.write(str(index_len) + "=" + dir + "\n")
        os.chown(self.file_location_txt, self._user_uid, self._vault_uid)

        self.read() # again refreshing the index

        return index_len


    def rename_item(self, index: int, new_name: str):
        """ Renames the filename to new_name in self.file_index
            and store the information on disk.

            It does not rename a file on disk. """

        self.read() # refreshing index

        # check whether new name already exist
        for item in self.file_index:
            if new_name == item[1]:
                raise IndexError("Duplicate file in index: " + new_name)

        for item in self.file_index:
            if index == item[0]:
                old_filename = item[1]
                item[1] = new_name
                today = today()
                item[2] = "<" + old_filename + "," + today + ">"
                break

        with open(self.all_files_txt, "w") as file:
            for item in self.file_index:
                file.write(str(item[0]) + "=" + item[1] + item[2] + "\n")

        os.chown(self.all_files_txt, self._user_uid, self._vault_uid)


    def remove_item(self, filename: str) -> int:
        """ Removes the filename both self.file_index and self.file_location_list
            and store the information on disk.

            It does not remove a file on disk. """

        self.read() # refreshing the index

        index = -1

        for idx, item in enumerate(self.file_index):
            if filename == item[1]:
                index = idx
                break

        if index == -1:
            return -1

        self.file_index.pop(index)
        self.file_location_list.pop(index)

        with open(self.all_files_txt, "w") as file:
            for item in self.file_index:
                file.write(str(item[0]) + "=" + item[1] + item[2] + "\n")
        os.chown(self.all_files_txt, self._user_uid, self._vault_uid)

        with open(self.file_location_txt, "a") as file:
            for item in self.file_location_list:
                file.write(str(item[0]) + "=" + item[1] + "\n")
        os.chown(self.file_location_txt, self._user_uid, self._vault_uid)

        return index


    def get_filename_index(self, name: str) -> int:
        """ Returns the index number of the file name,
            or -1 when the file is not found. """

        for item in self.file_index:
            if name == item[1]:
                return item[0]

        return -1  # not found the name


    def move_item(self, filename, new_filename: str):
        """ Moves the filename to new_filename in both self.file_index and self.file_location_list
            and store the information on disk.

            It does not move a file on disk. """
        pass


def today():
    # helper function
    return datetime.today().isoformat()
