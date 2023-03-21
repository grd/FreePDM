#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2023 by the FreePDM team
    :license:   MIT License.
"""

import os
import sys
from pathlib import Path

sys.path.append(os.fspath(Path(__file__).resolve().parents[1]))

from src.filesystem.filesystem import FileSystem
from src.skeleton.config import conf


if __name__ == "__main__":
    fs = FileSystem()
    fs.connect("10.0.0.11", "user1", "passwd1")
    fs.sftp.cwd("/vault/TestFiles2")
    print(fs.ls("/vault/TestFiles2"))
    print("checking file number: " + str(fs.check_latest_file_version("0003.FCStd")))
    print("checking file number: " + str(fs.check_latest_file_version("v0.FCStd")))
    print("checking file number: " + str(fs.check_latest_file_version("non_existing_file.FCStd")))
    rev = fs.check_latest_file_version("0003.FCStd")
    new_file = "0003.FCStd" + "." + str(rev).zfill(3)
    fs.revision_file(new_file)
    print("checking file number: " + str(fs.check_latest_file_version("0003.FCStd")))

    # s1 = "023"
    # i = int(s1)
    # print("str to int check: " + str(i))
    # i2 = 0
    # s2 = str(i2).zfill(3)
    # print("s2 = " + s2)

    fs.sftp.close()