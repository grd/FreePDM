"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

# Getting file info from Document.xml that is stored within each FC file

from fileinput import filename
from hashlib import blake2b
import defusedxml.ElementTree as et 
from tempfile import tempdir
import tempfile
import os
import imageio
import zipfile

# The Class ItemDataModel is the main class for reading and writing
# one FCStd file. The file can be read without opening a FreeCAD file.

class ItemDataModel():
	def __init__(self, filename):
		self.filename = filename
		self.thumbnail = None
		self.documentProperties = {}

	def read(self):
		with tempfile.TemporaryDirectory() as tmpdirname:
			print('created temporary directory', tmpdirname)

			# extracting the FCStd zip file
			zip=zipfile.ZipFile(self.filename)
			zip.extractall(tmpdirname)

			# reading xml file
			tree = tmpdirname + "/Document.xml"
			print(tree)
			self.readXML(tree)
		
			# Check whether there is a thumbnail
			if os.path.isdir(tmpdirname + "/thumbnails"):
				self.thumbnail = imageio.v3.imread(tmpdirname + "/thumbnails/Thumbnail.png")

	def readXML(self, data):
		with open(data) as docxml:
			docx = docxml.read()

		root = et.fromstring(docx)

		el = root.find(".//Property[@name='Comment']/String")
		self.documentProperties["Comment"] = el.attrib.get('value')

		el = root.find(".//Property[@name='Company']/String")
		self.documentProperties["Company"] = el.attrib.get('value')

		el = root.find(".//Property[@name='CreatedBy']/String")
		self.documentProperties["CreatedBy"] = el.attrib.get('value')

		el = root.find(".//Property[@name='CreationDate']/String")
		self.documentProperties["CreationDate"] = el.attrib.get('value')

		el = root.find(".//Property[@name='Id']/String")
		self.documentProperties["Id"] = el.attrib.get('value')

		el = root.find(".//Property[@name='Label']/String")
		self.documentProperties["Label"] = el.attrib.get('value')

		el = root.find(".//Property[@name='LastModifiedBy']/String")
		self.documentProperties["LastModifiedBy"] = el.attrib.get('value')

		el = root.find(".//Property[@name='LastModifiedDate']/String")
		self.documentProperties["LastModifiedDate"] = el.attrib.get('value')

		el = root.find(".//Property[@name='Uid']/Uuid")
		self.documentProperties["Uid"] = el.attrib.get('value')


#fcfile = ItemDataModel("/home/user/temp/freecad3.FCStd")
#fcfile.read()
 