#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
import sys
from pathlib import Path
import pysftp
from PySide2.QtWidgets import QDialog  # type: ignore
from PySide2.QtCore import QFile  # type: ignore
from PySide2.QtUiTools import QUiLoader  # type: ignore

sys.path.append(os.fspath(Path(__file__).resolve().parents[1]))
from skeleton.config import conf


class FileSystem():
    """File System related Class"""

    def __init__(self):
        print("Generic File System")
        self.conf = conf()
        self.conf.read()
        self._user = ""
        self._passwd = ""


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

    def connect(self, server, user, passwd): # TODO: also make it work with logging in with a passkey
        """Create a new connection to the server with SSHFS"""
        # print("server = " + server + ", user = " + user + ", passwd = " + passwd)
        self.server_pdm_name = server
        self._user = user
        self._passwd = passwd

        try:
            self.sftp = pysftp.Connection(
                self.server_pdm_name, username=self._user, password=self._passwd
        )
        except:
            # TODO: fill the exeptions in. They are located at: 
            # https://pysftp.readthedocs.io/en/release_0.2.9/pysftp.html
            print("something went wrong")

        self.sftp.cwd("/vault")
        return self.sftp


    def import_new_file(self, fname, dest_dir, descr, long_descr=None):
        """import a file inside the PDM. When you import a 
        file the meta-data also gets imported. The local files remain untouched. 
        When you import a file or files you need to set a directory and a description. 
        The new file inside the PDM gets a revision number automatically."""
        
        raise NotImplementedError("Function import_file is not implemented yet")

    def export_file(self, fname, dest_dir):
        """export a file to a local directory."""
        raise NotImplementedError("Function export_file is not implemented yet")
    
    def ls(self, dir=""):
        """list the sorted directories and latest files only"""
        prevdir = self.sftp.pwd
        self.sftp.chdir(dir)
        file_list = self.sftp.listdir(dir)
        ret_file_list = []
        for idx, file in enumerate(file_list):
            if self.sftp.isdir(file):
                ret_file_list.append(file)
            if self.sftp.isfile(file):
                if idx < len(file_list) + 1:
                    file1, ext1 = os.path.splitext(file)
                    file2, ext2 = os.path.splitext(file_list[idx+1])
                    if file1 == file2:
                        i = int(ext1[1:])
                        j = int(ext2[1:])
                        if i + 1 != j:
                            ret_file_list.append(file)
                    else:
                        ret_file_list.append(file) # Adding latest file in dir
        self.sftp.chdir(prevdir)
        return ret_file_list

    def check_latest_file_version(self, fname, dir):
        """ returns the latest version number of a file or -1 when the file doesn't exist."""
        pass

    def revision_file(self, fname):
        """increments a file revision number."""
        raise NotImplementedError("Function revision is not implemented yet")

    def checkout_file(self, fname):
        """locks a file so that others can't accidentally check-in a different file."""
        raise NotImplementedError("Function checkout_file is not implemented yet")

    def checkin_file(self, fname, descr, long_descr=None):
        """removes the locking but also uploads the file to the PDM. 
        You need to write a description of what you did."""
        raise NotImplementedError("Function checkin_file is not implemented yet")


    def rename_file(self, fname):
        """rename a file, for instance when he or she wants to use a file 
        with a specified numbering system."""
        raise NotImplementedError("Function rename_file is not implemented yet")
        # TODO: This can be a bit tricky because of the multiple files and users that are involved, besides the file name extensions that are ending with '.FCStd.#x'

    def move_file(self, fname, dest_dir):
        """moves a file to a different directory."""
        raise NotImplementedError("Function move_file is not implemented yet")
        # TODO: This can be a bit tricky because of the multiple files and users that are involved, besides the file name extensions that are ending with '.FCStd.#x'

if __name__ == "__main__":
    fs = FileSystem()
    fs.connect("10.0.0.11", "user1", "passwd1")
    fs.sftp.cwd("/vault/TestFiles2")
    print(fs.ls())
