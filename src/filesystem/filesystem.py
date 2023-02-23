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

sys.path.append('../')
import skeleton.config as conf


class FileSystem():
    """File System related Class"""

    _server_pdm_name: str = None
    _server_pdm_username: str = None
    _server_pdm_path: str = None
    _fs = None
    _user: str = None
    _passwd: str = None

    def init(self):
        print("Generic File System")
        self._server_pdm_name = conf.server_pdm_name
        self._server_pdm_user = conf.server_pdm_user_name
        self._server_pdm_path = conf.server_pdm_path


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

    def connect(self): # TODO: also make it work with logging in with a passkey
        loader = QUiLoader()
        path = os.fspath(Path(__file__).resolve().parents[1] / "gui/authenticate.ui")
        print(path)
        ui_file = QFile(path)
        ui_file.open(QFile.ReadOnly)
        self.ui = loader.load(ui_file, self)
        self.ui.exec()

        # """Create a new connection to the server with SSHFS"""
        # self._fs = pysftp.Connection(
        #     self._server_pdm_name,
        #     username=self._user,
        #     password=self._passwd
        # )
        # self._fs.cd(self._server_pdm_path)

    def disconnect(self):
        """Disconnects the connection"""
        self._fs.close()

    def import_file(self, fname, dest_dir, descr):
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

    def checkin_file(self, fname, descr):
        """removes the locking but also uploads the file to the PDM. 
        You need to write a description of what you did."""
        raise NotImplementedError("Function checkin_file is not implemented yet")


    def rename_file(self, fname):
        """rename a file, for instance when he or she wants to use a file 
        with a specified numbering system."""
        raise NotImplementedError("Function rename_file is not implemented yet")
        # TODO: This can be a bit tricky

    def move_file(self, fname, dest_dir):
        """moves a file to a different directory."""
        raise NotImplementedError("Function move_file is not implemented yet")
        # TODO: This can be a bit tricky

    def create_directory(self, dir_name):
        """creates a directory."""
        self._fs.mkdir(dir_name)
