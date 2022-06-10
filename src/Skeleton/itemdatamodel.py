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
		self.read()

	def read(self):
		with tempfile.TemporaryDirectory() as tmpdirname:
#			print('created temporary directory', tmpdirname)

			# extracting the FCStd zip file
			zip=zipfile.ZipFile(self.filename)
			zip.extractall(tmpdirname)

			# reading xml file
			tree = tmpdirname + "/Document.xml"
#			print(tree)
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

		el = root.find(".//Property[@name='a2p_Version']/String")
		if el != None:
			self.documentProperties["Assembly-A2P"] = el.attrib.get('value')

		el = root.find(".//Property[@name='Proxy']/Python")
		if el != None:
			a3assy = el.attrib.get('module')
			if a3assy == 'freecad.asm3.assembly':
				self.documentProperties["Assembly-A3"] = 'A3-Assy'

		el = root.find(".//Property[@name='SolverId']/String")
		if el != None:
			a4assy = el.attrib.get('value')
			if a4assy == 'Asm4EE':
				self.documentProperties["Assembly-A4"] = 'A4-Assy'


	def getFileName(self):
		return(self.filename)	                            

	def getDocumentPropterties(self):
		return(self.documentProperties)
		