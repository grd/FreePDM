#!/usr/bin/python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2023 by the FreePDM team
    :license:   MIT License.
"""

# Script for creating a new PDM vault.

import os
from os import path
import sys
from pathlib import Path

sys.path.append(os.fspath(Path(__file__).resolve().parents[1]))

from src.filesystem.filesystem import FileSystem
from src.skeleton.config import conf


def clear():
    if os.name == 'nt':
        os.system('cls')
    else:
        os.system('clear')


if __name__ == "__main__":
    clear()
    print("")
    print("Welcome to the create new vault program.")
    print("")
    print("Now a small questionaire is coming for the creation of the new vault.")
    print("For more info look at https://github.com/grd/FreePDM/tests/README.md")
    print("")
    print("The necessary information is the user uid, the vault uid,")
    print("the directory of the mounted sshfs drive and the name of the new vault.")
    print("")
    print("Are you ready? [Y/n]")
    drive = input()
    if drive.lower() == "n":
        exit("Program terminated.")
    
    print("Enter the user uid")
    user_uid = int(input())

    print("Enter the vault uid")
    vault_uid = int(input())

    print("Enter the path to the sshfs directory, e.g. '/mnt/test'")

    while True:
        drive = input()
        if not path.isdir(drive):
            print("Path doesn't exist. Try again.")
        else:
            break

    os.chdir(drive)
    
    if len(os.listdir()) > 0:
        print("")
        print("This is a list of exisisting vaults:")
        for file in os.listdir():
            print(file)

    print("")
    print("Input your new vault...")
    vault_dir = input()
    
    if path.isdir(vault_dir):
        sys.exit("The vault " + vault_dir + " already exist.")
    
    os.mkdir(vault_dir)
    os.chown(vault_dir, user_uid, vault_uid)

    os.chdir(vault_dir)

    # creating three files: 'All Files.txt', 'FileLocation.txt' and 'IndexNumber.txt'

    all_files = "All Files.txt"
    file_location = "FileLocation.txt"
    index_number = "IndexNumber.txt"
    locked_file = "Locked.txt"

    open(all_files, 'a').close()
    os.chown(all_files, user_uid, vault_uid)

    open(file_location, 'a').close()
    os.chown(file_location, user_uid, vault_uid)

    open(locked_file, 'a').close()
    os.chown(locked_file, user_uid, vault_uid)

    with open(index_number, 'w') as file:
        file.write("0")
    os.chown(index_number, user_uid, vault_uid)

    os.mkdir("PDM")
    os.chown('PDM', user_uid, vault_uid)

    print("Three files have been created: 'All Files.txt', 'FileLocation.txt' and the directory 'PDM")
    print("")
    print(os.listdir())
    print("")
    print("If that is correct, then the vault has been created.")
    print("Program ended successfully.")