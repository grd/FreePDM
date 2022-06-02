"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
from itemdatamodel import ItemDataModel

# The directorymodel is a list of files that can be used inside a GUI

class DirectoryModel():
	def __init__(self, dir):
		self.dir = dir
		self.dirList = {}
		self.getDirList()

	def getDirList(self):
		dir_list = os.listdir(self.dir)
		for x in dir_list:
			# eleminate files and directories starting with "."
			if x.startswith("."):
				continue
			if ".FCStd" in x:
				self.dirList[x] = ItemDataModel(x) 
			else:
				self.dirList[x] = x



