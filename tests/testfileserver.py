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
    fs.connect("/mnt/test/vault1", "user1", "passwd1")
    os.chdir("/vault/TestFiles3")
    print(fs.listdir())
    print("checking file number: " + str(fs.check_latest_file_version("0003.FCStd")))
    print("checking file number: " + str(fs.check_latest_file_version("v0.FCStd")))
    print("checking file number: " + str(fs.check_latest_file_version("non_existing_file.FCStd")))
    rev = fs.check_latest_file_version("0003.FCStd")
    new_file = "0003.FCStd" + "." + str(rev).zfill(3)
    fs.revision_file(new_file)
    print("checking file number: " + str(fs.check_latest_file_version("0003.FCStd")))

    fs.rename("0003.FCStd.003", "0001.FCStd")