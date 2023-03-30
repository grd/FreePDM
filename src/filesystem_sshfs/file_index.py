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

class FileIndex():
    """File Index Files in the root"""

    def __init__(self, vault_dir: str):
        self.vault_dir = vault_dir
        self.file_index = None
        self.file_list = None
        self.read()

    def get_index(self, name: str) -> int:
        pass

    def read(self):
        with open(path.join(self.vault_dir, "All Files.txt"), "r") as file:
            while (line := file.readline()):
                print(line.split("="))
        with open(path.join(self.vault_dir, "FileLocation.txt"), "r") as file:
            while (line := file.readline()):
                print(line.split("="))

    def add_item(self, name, dir: str):
        pass

    def rename_item(self, item, new_item, date: str):
        pass

    def remove_item(self, item: str):
        pass