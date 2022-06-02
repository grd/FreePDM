import os

# The directorymodel is a list of files that can be used inside a GUI

from itemdatamodel import ItemDataModel

class DirectoryModel():
	def __init__(self, dir):
		self.dir = dir
		self.dirList = {}

	def getList(self):
		dir_list = os.listdir(self.dir)
		for x in dir_list:
			# eleminate files and directories starting with "."
			if x.startswith("."):
				continue
			if x.contains(".FCStd"):
				self.dirList[x] = ItemDataModel(x) 


