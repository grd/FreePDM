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

sys.path.append(os.fspath(Path(__file__).resolve().parents[1]))
from skeleton.config import conf
from filesystem import file_index
import filesystem.file_directory as fd
from filesystem.extras import *

class FileSystem():
    """File System related Class"""

    def __init__(self):
        self.conf = conf()
        self.conf.read()
        self._index = file_index.FileIndex()
        self._user: str
        self._passwd: str
        self._vault_uid = self.conf.vault_uid
        self._main_pdm_dir: str
        self.current_working_dir: str

    def connect(self, vault_dir, user, passwd: str):
        os.chdir(vault_dir)
        print("Vault dir: " + os.getcwd())
        self._user = user
        self._user_uid = self.conf.get_user_uid(self._user)
        self._passwd = passwd
        self._vault_dir = vault_dir
        self._main_pdm_dir = path.join(self._vault_dir, "PDM")
        self.current_working_dir = self._main_pdm_dir
        self._index.init(self._vault_dir, self._user_uid, self._vault_uid)
        os.chdir(self.current_working_dir)
        self._locked_txt = path.join(self._vault_dir, "Locked.txt")
        # self._locked_index: List[int, str] = get_key_value_list(self._locked_txt)
        self._locked_index: List[int, str] = []
        self.read()


    def read(self):
        """read updates the locked index by reading from the locked file."""
        self._locked_index = get_key_value_list(self._locked_txt)


    def cred(self):
        """cred returns the credentials that are being used inside a FileDirectory."""

        idx = len(self._main_pdm_dir)
        dir_name = self.current_working_dir[idx:]
        return dir_name, self._user, self._user_uid, self._vault_uid


    def import_new_file(self, fname, descr, long_descr: str) -> int:
        """import a file inside the PDM. When you import a
        file the meta-data also gets imported, which means uploaded to the server.
        When you import a file or files you are placing the new file in the current directory.
        The new file inside the PDM gets a revision number automatically.
        The function returns the number of the imported file."""

        if not path.isfile(fname):
            raise FileNotFoundError("File " + fname + " could not be found.")

        index = self._index.add_item(fname, self.current_working_dir)

        fd.new_directory(index, self.cred()).new_version(fname, descr, long_descr)

        return index


    def export_file(self, fname, dest_dir):
        """export a file to a local directory."""
        raise NotImplementedError("Function export_file is not implemented yet")


    def mkdir(self, dname: str):
        """Creates a new directory inside the current directory, with the correct uid and gid."""
        os.mkdir(dname)
        # os.chown(remotepath=dname, uid=self._user_uid, gid=self._vault_uid)
        os.chown(dname, self._user_uid, self._vault_uid)


    def chdir(self, dir: str):
        self.current_working_dir = path.join(self.current_working_dir, dir)
        os.chdir(self.current_working_dir)


    def _latest_file_version(self, dir: str) -> str:
        """Returns the latest file name from directory "dir" """

        prevdir = os.getcwd()

        version: str
        version_list = []
        file_name: str

        with open("VER.txt", "r") as file:
            version_list = file.readlines()
        if len(version_list) >= 2:
            version = version_list[-2].strip()

        os.chdir(version)

        with open('File.txt', "r") as file:
            file_name = file.readline()

        os.chdir(prevdir)
        
        return file_name


    def _latest_file_information_version(self, dir: str) -> List[str]:
        """Returns the latest file information from directory "dir" """
        version: str
        version_list = []
        file_information: List[str]
        with open(path.johnpath(dir, "VER.txt")):
            version_list = IOBase.readlines()
        if len(version_list >= 2):
            version = version_list[-2]
        with open(path.joinpath(dir, version, version + '.txt')):
            file_information = IOBase.readlines()
        return file_information


    def _all_file_versions(self, dir: str) -> List[str]:
        """Returns all file versions name from directory "dir" """
        version: str
        version_list = []
        file_names: List[str]
        with open(path.johnpath(dir, "VER.txt")):
            version_list = IOBase.readlines()
        if len(version_list >= 2):
            # TODO; Implement this. ATM I don't know how.
            version = version_list[-2]
        with open(path.joinpath(dir, version, version + ' File.txt')):
            file_names = IOBase.readline()
        return file_names


    def listdir(self) -> List[str]:
        """list the sorted directories and files of the current working directory"""
        dir_list = os.listdir(self.current_working_dir)
        directory_list = []
        file_list = []
        sub_dir_list: List[str] = []

        for sub_dir in dir_list:
            if path.isdir(sub_dir):
                os.chdir(sub_dir)
                # check wheter a file named "VER.txt" exists
                if path.isfile("VER.txt"):
                    file_list.append(self._latest_file_version(sub_dir))
                else: # ordinary directory
                    directory_list.append(sub_dir)
                os.chdir('..')

        if path.abspath(self.current_working_dir) != self._main_pdm_dir:
            sub_dir_list.append("d: ..")

        for sub_dir in directory_list:
            sub_dir_list.append("d: " + sub_dir)

        for file in file_list:
            sub_dir_list.append(file)
        
        return sub_dir_list


    def check_latest_file_version(self, fname: str) -> int:
        """ returns the latest version number of a file in the current
        directory or -1 when the file doesn't exist."""
        file_list = os.listdir()
        result = -1
        for file in file_list:
            if os.isdir(file):
                continue
            file1, ext1 = os.path.splitext(file)
            if fname == file1:
                result = int(ext1[1:])
        return result


    # def revision_file(self, fname):
    #     """copy a file and increments revision number."""
    #     # TODO; make it work for older files and not simply overwrite the files.
    #     # TODO: we should have to deal with version ".999"
    #     self.sftp.get(fname)
    #     file, ext = os.path.splitext(fname)
    #     ext_int = int(ext[1:])
    #     ext_int += 1
    #     new_file = file + "." + str(ext_int).zfill(3)
    #     os.rename(fname, new_file)
    #     self.sftp.put(new_file)
    #     self.sftp.chown(remotepath=new_file, uid=self._user_uid, gid=self._vault_uid)
    #     os.remove(new_file)


    def is_locked(self, itemnr: int) -> bool:
        """Check whether the itemnr is locked."""

        self.read() # update the index

        for item in self._locked_index:
            if int(item[0]) == itemnr:
                return True

        return False


    def checkout(self, itemnr: int):
        """Checkout means locking a itemnr so that only you can use it."""

        self.read() # update the index

        # check whether the itemnr is locked
        if self.is_locked(itemnr):
            raise Exception("file is already locked.")

        self._locked_index.append([itemnr, self._user])

        with open(self._locked_txt, "a") as file:
            file.write(str(itemnr) + "=" + self._user + "\n")

        os.chown(self._locked_txt, self._user_uid, self._vault_uid)


    def checkout_status(self, itemnr: int) -> str:
        """Returns the status of a check-out."""

        self.read() # updates the index

        for item in self._locked_index:
            if int(item[0]) == itemnr:
                return item[1]

        return "The file is not checked out."


    def checkin(self, itemnr: int):
        """Checkin means unlocking a itemnr."""

        self.read() # updates the index

        num = -1
        for idx, item in enumerate(self._locked_index):
            if int(item[0]) == itemnr:
                num = idx

        if num != -1:
            self._locked_index.pop(num)
            store_key_value_list(self._locked_txt, self._locked_index)
            os.chown(self._locked_txt, self._user_uid, self._vault_uid)


    def rename(self, src, dest: str) -> bool:
        """rename a directory or a file, for instance when the user wants to use a file 
        with a specified numbering system."""

        # TODO: When renaming a directory, they may also have parts inside!
        # So I also need to modify the file locations.

        # TODO: This code still belongs to the old filesystem version,
        # so I need to rewrite it completely.

        if path.isdir(src):
            try:
                os.rename(src, dest)
            except IOError:
                # The new directory already exist. Can't rename the directory.
                return False
        else:
            file, ext = os.path.splitext(src)
            file_list = os.listdir()
            for item in file_list:
                file1, ext1 = os.path.splitextitem(item)
                if file == file1:
                    try:
                        os.rename(item, dest + ext1)
                    except IOError:
                        # The new file already exist. Can't rename the file.
                        return False
        # The file(s) is / are successfully renamed
        return True


    def move_file(self, fname, dest_dir: str):
        """moves a file to a different directory."""
        raise NotImplementedError("Function move_file is not implemented yet")
        # TODO: This can be a bit tricky because of the multiple files.
        # Also look into file_index.py function move_item !!!
