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
        # self.server_pdm_name = conf.get_pdm_name()
        # self.server_pdm_path = conf.get_pdm_path()
        self._sftp = None
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
            self._sftp = pysftp.Connection(
                self.server_pdm_name, username=self._user, password=self._passwd
        )
        except:
            # TODO: fill the exeptions in. They are located at: 
            # https://pysftp.readthedocs.io/en/release_0.2.9/pysftp.html
            print("something went wrong")

        print(self._sftp.listdir())
        # self._sftp.cd(self._server_pdm_path)
        # print(self._sftp.getcwd())


    def close(self):
        """Disconnects the connection"""
        self._sftp.close()

    def import_new_file(self, fname, dest_dir, descr, long_descr=None):
        """import a file inside the PDM. When you import a 
        file the meta-data also gets imported. The local files remain untouched. 
        When you import a file or files you need to set a directory and a description. 
        The new file inside the PDM gets a revision number automatically."""
        raise NotImplementedError("Function import_file is not implemented yet")

    def export_file(self, fname, dest_dir):
        """export a file to a local directory."""
        raise NotImplementedError("Function export_file is not implemented yet")

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

    def create_directory(self, dir_name):
        """creates a directory."""
        self._sftp.mkdir(dir_name)

    def exists(self):
        """ check wheter there is a connection."""
        return self._sftp.exists(self.server_pdm_name)

