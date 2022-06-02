"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from directorymodel import DirectoryModel

# right now it should only show a directory that you specify
#  and print the values about that directory 

# TODO: Well, right now it only shows the files in one directory.

def main():
    dm = DirectoryModel("/home/user/temp")
    dm.getDirList()
  
if __name__=="__main__":
    main();