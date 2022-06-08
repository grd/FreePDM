"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from numpy import true_divide
from directorymodel import DirectoryModel
import os
import sys

# right now it should only show a directory that you specify
#  and print the values about that directory. 

dir =""

def handleDirectory():
    global dir
    dm = DirectoryModel(dir)
    print("  nr Dir/File  Filename")
    for list in dm.dirList:
        print(list[0].rjust(4, ' ') + ' ' + list[1].ljust(9, ' ') + ' ' + list[2])
  
    ip = input("Press a number or'q' to quit, or '-1' to go back in directory ") 
    if ip == 'q':
        exit()
    num = int(ip)
    if num == -1:
        dir = dir[0:dir.rfind("/")]
        return(True)
    if num >= 0 and num <= len(dm.dirList): # We have a number here...
        item = dm.dirList[num]
        print(item)
        if item[1] == 'Directory':
            dir = dir + '/' + dm.dirList[num][2]
            return(True)
        if item[1] == 'FCStd':
            print(item[2])

def main():
    global dir
    dir = os.path.expanduser('~')
    if len(sys.argv) == 2:
        dir = sys.argv[1]
    
    okay = True
    while okay == True:
        okay = handleDirectory()

        

if __name__=="__main__":
    main();