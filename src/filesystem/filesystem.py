#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from sshfs import SSHFileSystem


class FileSystem():
    """File System related Class"""

    def init(self):
        print("Generic File System")

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

    def connect(self):
        """Create a new connection to the server"""
        # create a new connection to the server with SSHFS
        # 
        # # Connect with a password
        # fs = SSHFileSystem(
        #     '127.0.0.1',
        #     username='sam',
        #     password='fishing'
        # )
        raise NotImplementedError("Function connect is not implemented yet")

    def disconnect(self):
        """Disconnects the connection"""
        raise NotImplementedError("Function disconnect is not implemented yet")

    def import_file(self, fname, dest_dir):
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

    def move_file(self, fname, dest_dir):
        """moves a file to a different directory."""
        raise NotImplementedError("Function move_file is not implemented yet")

    def create_directory(self, dir_name):
        """creates a directory."""
        raise NotImplementedError("Function create_directory is not implemented yet")
