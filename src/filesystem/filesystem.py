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

sys.path.append(os.fspath(Path(__file__).resolve().parents[1]))
from skeleton.config import conf

class FileSystem():
    """File System related Class"""

    def __init__(self):
        self.conf = conf()
        self.conf.read()

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
        self._user_uid = self.conf.get_user_uid(self._user)
        self._vault_uid = self.conf.vault_uid

        try:
            self.sftp = pysftp.Connection(
                self.server_pdm_name, username=self._user, password=self._passwd
        )
        except:
            # TODO: fill the exeptions in. They are located at: 
            # https://pysftp.readthedocs.io/en/release_0.2.9/pysftp.html
            print("something went wrong")

        self.sftp.cwd("/vault")
        print("Connected to: " + self.conf.server_name)
        print("Vault directory: " + self.conf.server_path)
        return self.sftp


    def import_new_file(self, fname, dest_dir, descr, long_descr=None):
        """import a file inside the PDM. When you import a 
        file the meta-data also gets imported, which means uploaded to the server.
        When you import a file or files you are placing the new file in the current directory. 
        The new file inside the PDM gets a revision number automatically."""
        raise NotImplementedError("Function import_file is not implemented yet")

    def export_file(self, fname, dest_dir):
        """export a file to a local directory."""
        raise NotImplementedError("Function export_file is not implemented yet")

    def mkdir(self, dname):
        """Creates a new directory inside the current directory, with the correct uid and gid."""
        self.sftp.mkdir(dname)
        self.sftp.chown(remotepath=dname, uid=self._user_uid, gid=self._vault_uid)

    def ls(self, dir):
        """list the sorted directories and filtered the latest files only"""
        prevdir = self.sftp.pwd
        self.sftp.chdir(dir)
        file_list = self.sftp.listdir(dir)
        index = len(file_list)
        while index > 0:
            index = index - 1
            if self.sftp.isdir(file_list[index]):
                continue
            file, ext = os.path.splitext(file_list[index])
            number = int(ext[1:])
            while number > 1:
                index = index - 1
                file1, ext1 = os.path.splitext(file_list[index])
                if file == file1:
                    file_list.pop(index)
                    number = number - 1
        self.sftp.chdir(prevdir)
        return file_list

    def check_latest_file_version(self, fname):
        """ returns the latest version number of a file in the current
        directory or -1 when the file doesn't exist."""
        file_list = self.sftp.listdir()
        result = -1
        for file in file_list:
            if self.sftp.isdir(file):
                continue
            file1, ext1 = os.path.splitext(file)
            if fname == file1:
                result = int(ext1[1:])
        return result

    def revision_file(self, fname):
        """copy a file and increments revision number."""
        # TODO; make it work for older files and not simply overwrite the files.
        # TODO: we should have to deal with version ".999"
        self.sftp.get(fname)
        file, ext = os.path.splitext(fname)
        ext_int = int(ext[1:])
        ext_int += 1
        new_file = file + "." + str(ext_int).zfill(3)
        os.rename(fname, new_file)
        self.sftp.put(new_file)
        self.sftp.chown(remotepath=new_file, uid=self._user_uid, gid=self._vault_uid)
        os.remove(new_file)



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
        # TODO: This can be a bit tricky because of the multiple files and users that are involved.

    def move_file(self, fname, dest_dir):
        """moves a file to a different directory."""
        raise NotImplementedError("Function move_file is not implemented yet")
        # TODO: This can be a bit tricky because of the multiple files and users that are involved

