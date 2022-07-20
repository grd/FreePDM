"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
from itemdatamodel import ItemDataModel

# The directorymodel is a list of files that can be used inside a GUI

class DirectoryModel(object):
	def __init__(self, dir):
		self.dir = dir
		self.dirList = [["  nr", "Dir/File ", "Filename"]]
		self.getDirList()
		

	def getDirList(self):
		dir_list = os.listdir(self.dir)
		dir_list.sort()
		directory_list = []
		file_list = []
		nr = 0
		for d in dir_list:
			full_path = self.dir + '/' + d
			# eliminate files and directories starting with "."
			if not d.startswith("."):
				if os.path.isdir(full_path):
					directory_list.append(d)
				if os.path.isfile(full_path):
					file_list.append(d)

		for dl in directory_list:
			self.dirList.append([str(nr), 'Directory', dl]) 
			nr = nr + 1
		for fl in file_list:
			if fl.endswith(".FCStd"):
				self.dirList.append([str(nr), 'FCStd', fl])
				nr = nr + 1
			else:
				self.dirList.append([str(nr), 'File', fl])
				nr = nr + 1
		self.dirList.pop(0)

	def size(self):
		len(self.dirList)

	# What kind of FC files do you want to see? 
	# Do you want to see the history of the files?
	# def filter(self):

	# Purge stored versions of files. All up to number 
	# def purge(number):

	def helloWorld():
		return("Hello World.")