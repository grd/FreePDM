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
from io import IOBase
from typing import List
import shutil

sys.path.append(os.fspath(Path(__file__).resolve().parents[1]))
from filesystem import extras

class FileDirectory():
    """
    File Directory related Class.

    This is a helper class. It helps with setting up directories inside a directory.
    These directories are versions of a file, which means that it helps with with version control.
    It is designed to work with method chaining, such as:

    new_directory(index).new_version(fname, descr, long_descr)

    open_directory(index).checkout(user)

    open_directory(index).checkout_status()

    open_directory(index).checkin()

    open_directory(index).release()

    open_directory(index).version_info()

    open_directory(index).edit_version_info()

    open_directory(index).delete_version()
    """

    def __init__(self, cred):
        dir, user_name, user_uid, vault_uid = cred
        self._dir: str = dir
        self._user_name: str = user_name
        self._user_uid: int = user_uid
        self._vault_uid: int = vault_uid
        self._locked_file = "Locked.txt"


    def new_directory(self, dname):
        os.mkdir(dname)
        os.chown(dname, self._user_uid, self._vault_uid)
        self._dir = dname
        return self


    def new_version(self, fname, descr, long_descr):
        """Creates a new file version."""

        # check whether the file is locked
        if self.is_locked():
            raise Exception("File is locked.")
        
        print("_dir = " + self._dir)
        # going into the directory...
        os.chdir(self._dir)

        # create a new version string
        new_version = self.latest_version_string() + 1
        version = "VER" + str(new_version).zfill(2)
        to_day = extras.today()
        with open("VER.txt", 'a') as file:
            file.write(version + "\n")
            file.write(to_day + "\n")

        # create a new version dir
        os.mkdir(version)
        os.chown(version, self._user_uid, self._vault_uid)
        os.chdir(version)

        # create a new properties text
        with open("Properties.txt", 'w') as file:
            for item in long_descr:
                file.write(item + "\n")

        os.chown("Properties.txt", self._user_uid, self._vault_uid)

        # create a new file reference text
        split_file = path.split(fname)
        copied_file = split_file[1]

        with open("File.txt", "w") as file:
            file.write(copied_file)

        os.chown("File.txt", self._user_uid, self._vault_uid)

        # copy the file inside the new version
        shutil.copyfile(fname, copied_file)

        print(self._dir)
        os.chdir("../..")
        return self


    def latest_version_string(self):
        """Returns the latest version string, or nothing when VER.txt doesn't exist."""

        if not path.isfile("VER.txt"):
            return -1

        with open("VER.txt", "r") as file:
            my_list = file.readlines()
            verstr = my_list[-2]
            ver = int(verstr[3:])
            return ver


    def is_locked(self):
        """Check whether the file is locked by someone else."""
        if path.isfile(self._locked_file):
            # check locking user
            with open(self._locked_file, "r") as file:
                line = file.readline()
                _, user = line.split("=")
                return user != self._user_name
        else:
            return False


    def checkout(self):
        """Checkout means locking a file so that only you can use it."""

        # check whether the file is locked
        if self.is_locked():
            raise Exception("File is locked.")

        with open(self._locked_file, "w") as file:
            file.write("Locked=" + self._user_name)

        os.chown(self._locked_file, self._user_uid, self._vault_uid)


    def checkout_status(self):
        """Returns the status of a check-out."""

        if self.is_locked():
            with open(self._locked_file, "r") as file:
                line = file.readline()
                _, user = line.split("=")
                return "The file is checked out by " + user
        else:
            return "The file is not checked out."


    def checkin(self):
        """Checkin means unlocking a file."""

        os.remove(self._locked_file)
        return self


    def release():
        pass

    def version_info():
        pass

    def edit_version_info():
        pass

    def delete_version():
        pass



def new_directory(item: int, cred) -> FileDirectory:
    """ Creates a new directory with a version string inside.
        The item number is the name of the directory.
        The cred are the necessary credentials."""
    dir = FileDirectory(cred)
    dir.new_directory(str(item))
    return dir


def open_directory(item: int) -> FileDirectory:
    pass
