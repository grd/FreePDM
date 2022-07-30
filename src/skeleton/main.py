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


def handleDirectory(directory):  # dir is a build in python function
    dirmodel = DirectoryModel(directory)
    print("  nr Dir/File  Filename              size")
    for list in dirmodel.directoryList:
        print(list['nr'].rjust(4, ' ') + ' ' + list['dirOrFile'].ljust(9, ' ') + ' ' + list['filename'].ljust(25, ' ') + list['size'])
#        print(list['nr'].rjust(4, ' ') + ' ' + list['filename'].ljust(25, ' ') + list['size'])

    keyinput = input("Press a number or'q' to quit, or '-1' to go the higher directory ")
    if keyinput == 'q':
        exit()
    num = int(keyinput)  # Just don't press anything. It doesn't work like that...
    if num == -1:
        directory = directory[0:directory.rfind("/")]
        return(directory)
    if num >= 0 and num <= len(dirmodel.directoryList):  # We have a number here...
        item = dirmodel.directoryList[num]
        if item['dirOrFile'] == 'Directory':
            directory = directory + '/' + dirmodel.directoryList[num]['filename']
            return(directory)
        elif item['dirOrFile'] == 'FCStd':
            fcFile = ItemDataModel(directory + '/' + item['filename'])
            idm = fcFile.documentProperties
            print("File Size: " + fcFile.getFileSize())
            print("Data from FCStd model")
            for x, y in idm.items():
                print(x.ljust(18, ' ') + y)
            input("Press Enter to return")
            return(directory)
        else:
            return(directory)


def main():
    directory = os.path.expanduser('~')
    if len(sys.argv) == 2:
        directory = sys.argv[1]

    while True:
        directory = handleDirectory(directory)


if __name__ == "__main__":
    main()
