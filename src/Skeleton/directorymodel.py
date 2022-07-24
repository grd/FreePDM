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
		self.dirList = [{"nr": "number", 
						"dirOrFile": "Dir or File", 
						"filename": "File name", 
						"size": "file size"}]
		self.getDirList()
		

	def full_path(self, file):
		return (self.dir + '/' + file)

	def getDirList(self):
		dir_list = os.listdir(self.dir)
		dir_list.sort()
		directory_list = []
		file_list = []
		nr = 0
		for d in dir_list:
			full_item_path = self.full_path(d)
			# eliminate files and directories starting with "."
			if not d.startswith("."):
				if os.path.isdir(full_item_path):
					directory_list.append(d)
				if os.path.isfile(full_item_path):
					file_list.append(d)

		# Keep this in mind:

		# self.dirList = [{"nr": "number", 
		# 		"dirOrFile": "Dir or File", 
		# 		"filename": "File name", 
		# 		"size": "file size"}]

		for dl in directory_list:
			self.dirList.append({'nr': str(nr), 'dirOrFile': 'Directory', 'filename': dl, 'size': ''}) 
			nr = nr + 1
		for fl in file_list:
			size = str(os.path.getsize(self.full_path(fl)))
			if fl.endswith(".FCStd"):
				self.dirList.append({'nr': str(nr), 'dirOrFile': 'FCStd', 'filename': fl, 'size': size})
				nr = nr + 1
			else:
				self.dirList.append({'nr':str(nr), 'dirOrFile': 'File', 'filename': fl, 'size': size})
				nr = nr + 1
		# getting rid of first item		
		self.dirList.pop(0)

	def size(self):
		return(len(self.dirList))

#TODO:
	# What kind of FC files do you want to see? 
	# Do you want to see the history of the files?
	# def filter(self):

#TODO:
 	# Purge stored versions of files. All up to number 
	# def purge(number):

