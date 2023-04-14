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
import appdirs

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
    mounting_point = input()
    if mounting_point.lower() == "n":
        exit("Program terminated.")
    
    print("Enter the user name")
    user_name = input()

    print("Enter the user uid")
    user_uid = int(input())

    print("Enter the vault uid")
    vault_uid = int(input())

    print("Enter mounting point of the sshfs directory, e.g. '/mnt/test'")

    while True:
        mounting_point = input()
        if not path.isdir(mounting_point):
            print("Path doesn't exist. Try again.")
        else:
            break

    os.chdir(mounting_point)
    
    if len(os.listdir()) > 0:
        print("")
        print("This is a list of exisisting vaults:")
        for file in os.listdir():
            print(file)

    print("")
    print("Input the directory of your new vault...")
    vault_dir = input()
    
    if path.isdir(vault_dir):
        sys.exit("The vault " + vault_dir + " already exist.")
    
    os.mkdir(vault_dir)
    os.chown(vault_dir, user_uid, vault_uid)

    os.chdir(vault_dir)

    # creating four files: 'All Files.txt', 'FileLocation.txt',
    # 'IndexNumber.txt' and 'Locked.txt'

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

    print("Four files have been created: 'All Files.txt', 'FileLocation.txt', ")
    print("'IndexNumber.txt' and 'Locked.txt', and the directory 'PDM")
    print("")
    print(os.listdir())
    print("")
    print("If that is correct, then the vault has been created.\n\n")
    print("Creating a configuration file...\n\n")

    appname = 'FreePDM'
    config_dir = appdirs.user_config_dir(appname)
    config_name = os.path.join(config_dir, 'FreePDM.conf')

    if path.isfile(config_name):
        print("Configuration file already exist.\n\n")
        print("Changing the vault directory...")
        item_list = []
        with open(config_name, "r") as file:
            item_list = file.readlines()

        idx = len(item_list) - 1
        while idx > 0:
            item = item_list[idx]
            if item.startswith("vault"):
                item_list[idx] = "vault = " + "/vault/" + vault_dir + "\n"
            idx -= 1

        with open(config_name, "w") as file:
            for item in item_list:
                file.write(item)
    else:
        os.mkdir(config_dir)
        with open(config_name, "w") as file:
            file.write("[DEFAULT]\n")
            file.write("log_file = \n")
            file.write("logging_is_on = False\n")
            file.write("mounting_point = \"" + mounting_point + "\"\n")
            file.write("server_path = \"" + "/vault/" + vault_dir + "\"\n\n")
            file.write("[user]\n")
            file.write("vault = " + vault_uid + "\n")
            file.write(user_name + " = " + str(user_uid) + "\n")



    print("Program ended successfully.")

