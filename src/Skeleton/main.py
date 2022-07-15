"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from numpy import true_divide
from directorymodel import DirectoryModel
from itemdatamodel import ItemDataModel
import os
import sys

# right now it should only show a directory that you specify
#  and print the values of the FC file about that directory. 

def handleDirectory(dir):
    dm = DirectoryModel(dir)
    print("  nr Dir/File  Filename")
    for list in dm.dirList:
        print(list[0].rjust(4, ' ') + ' ' + list[1].ljust(9, ' ') + ' ' + list[2])
  
    ip = input("Press a number or'q' to quit, or '-1' to go the higher directory ") 
    if ip == 'q':
        exit()
    num = int(ip) # Just don't press anything. It doesn't work like that...
    if num == -1:
        dir = dir[0:dir.rfind("/")]
        return(dir)
    if num >= 0 and num <= len(dm.dirList): # We have a number here...
        item = dm.dirList[num]
        if item[1] == 'Directory':
            dir = dir + '/' + dm.dirList[num][2]
            return(dir)
        elif item[1] == 'FCStd':
            fcFile = ItemDataModel(dir + '/' + item[2])
            idm = fcFile.documentProperties
            print("File Size: " + fcFile.getFileSize())
            print("Data from FCStd model")
            for x, y in idm.items():
                print(x.ljust(18, ' ') + y)
            input("Press Enter to return")
            return(dir)
        else:
            return(dir)


def main():
    dir = os.path.expanduser('~')
    if len(sys.argv) == 2:
        dir = sys.argv[1]
    
    while True:
        dir = handleDirectory(dir)

        

if __name__=="__main__":
    main();