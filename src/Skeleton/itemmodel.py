"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from fileinput import filename
from hashlib import blake2b
import tempfile
from defusedxml.ElementTree import parse
import os
from tempfile import tempdir
import zipfile

# The Class ItemModel is the main class for reading and writing
# one FCStd file. The file can only be read without opening a FreeCAD file.

class ItemModel():
	def __init__(self, filename):
		self.filename = filename

	def read(self):
		with tempfile.TemporaryDirectory() as tmpdirname:
			print('created temporary directory', tmpdirname)

			# extracting the FCStd zip file
			zip=zipfile.ZipFile(self.filename)
			zip.extractall(tmpdirname)

			# reading xml file
			tree = parse(tmpdirname + "/Document.xml")
			print(tree)
			self.readXML(tree)
			print(tree.getroot())
			# def readpicture():
	
	def readXML(self, xmlParseFile):


fcfile = ItemModel("/home/user/temp/freecad.FCStd")
fcfile.read()
