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
		self.dirList = []
		self.getDirList()

	def getDirList(self):
		dir_list = os.listdir(self.dir)
		dir_list.sort()
		directory_list = []
		file_list = []
		for d in dir_list:
			full_path = self.dir + '/' + d
			# eleminate files and directories starting with "."
			if not d.startswith("."):
				if os.path.isdir(full_path):
					directory_list.append(d)
				if os.path.isfile(full_path):
					file_list.append(d)

		for dl in directory_list:
			print('[ ' + dl + ' ]')
		for fl in file_list:
			if ".FCStd" in fl:
				idm = ItemDataModel(self.dir + '/' + fl).get()
				print(os.path.basename(idm))
			else:
				print(fl)



