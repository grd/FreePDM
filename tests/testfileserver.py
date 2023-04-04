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

sys.path.append(os.fspath(Path(__file__).resolve().parents[1]))

from src.filesystem.filesystem import FileSystem
from src.skeleton.config import conf


if __name__ == "__main__":
    fd_dir = os.fspath(Path(__file__).resolve().parents[1])
    test_dir = path.join(fd_dir, "ConceptOfDesign/TestFiles")
    file1 = path.join(test_dir, "0001.FCStd")
    file2 = path.join(test_dir, "0002.FCStd")
    file3 = path.join(test_dir, "0003.FCStd")
    file4 = path.join(test_dir, "0004.FCStd")
    file5 = path.join(test_dir, "0005.FCStd")
    file6 = path.join(test_dir, "0006.FCStd")
    
    fs = FileSystem()
    fs.connect("/home/user/mnt/vault1", "user1", "passwd1")
    print("Root directory of the vault: " + os.getcwd())
    fs.mkdir("Standard Parts")
    fs.mkdir("Projects")
    print(fs.listdir())
    fs.chdir("Projects")
    print(fs.listdir())
    fs.import_new_file(file1, "ConceptOfDesign/TestFiles", "")
    fs.import_new_file(file2, "ConceptOfDesign/TestFiles", "")
    fs.import_new_file(file3, "ConceptOfDesign/TestFiles", "")
    
    fs.mkdir("Temp")
    fs.chdir("Temp")

    fs.import_new_file(file4, "ConceptOfDesign/TestFiles", "")
    fs.import_new_file(file5, "ConceptOfDesign/TestFiles", "")
    fs.import_new_file(file6, "ConceptOfDesign/TestFiles", "")
    
    fs.chdir("..")
    print(fs.listdir())

    # rev = fs.check_latest_file_version("0003.FCStd")
    # new_file = "0003.FCStd" + "." + str(rev).zfill(3)
    # fs.revision_file(new_file)
    # print("checking file number: " + str(fs.check_latest_file_version("0003.FCStd")))

    # fs.rename("0003.FCStd.003", "0001.FCStd")