"""
    :copyright: Copyright 2022 by Gerard van de Schoot.
    :license:   MIT License.
"""

# import OS module
import os
 
# Get the list of all files and directories
path = "C://Users//Vanshi//Desktop//gfg"
dir_list = os.listdir(path)
 
print("Files and directories in '", path, "' :")
 
# prints all files
print(dir_list)
