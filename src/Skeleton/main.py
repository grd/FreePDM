"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from directorymodel import DirectoryModel
import os
import sys

# right now it should only show a directory that you specify
#  and print the values about that directory. 

def main():
    dir = os.path.expanduser('~')
    if len(sys.argv) == 2:
        dir = sys.argv[1]

    dm = DirectoryModel(dir)
    print("  nr Dir/File  Filename")
    for list in dm.dirList:
        print(list[0].rjust(4, ' ') + ' ' + list[1] + ' ' + list[2])
  

if __name__=="__main__":
    main();