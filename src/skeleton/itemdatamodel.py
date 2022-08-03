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
from PySide2 import QtGui

# The Class ItemDataModel is the main class for reading and writing
# one FCStd file. The file can be read without opening a FreeCAD file.


class ItemDataModel():
    def __init__(self, filename):  # filename is also imported so change
        self.fileName = filename
        self.xmldocument = None
        self.thumbnail = None
        self.documentProperties = {}
        self.read()

        # check out file size
        self.fileSize = os.path.getsize(filename)

    def read(self):
        with tempfile.TemporaryDirectory() as tmpdirname:
            # print('created temporary directory', tmpdirname)

            # extracting the FCStd zip file
            zippedfile = zipfile.ZipFile(self.fileName)  # zip is python internal function
            zippedfile.extractall(tmpdirname)

            # reading xml file
            tree = tmpdirname + "/Document.xml"
#            print(tree)
            self.readXML(tree)

            # Check whether there is a thumbnail
            if os.path.isdir(tmpdirname + "/thumbnails"):
                self.documentProperties["thumbnail"] = tmpdirname + "/thumbnails/Thumbnail.png"
                self.thumbnail = QtGui.QImage(self.documentProperties['thumbnail'])

    def readXML(self, data):
        with open(data) as docxml:
            self.xmldocument = docxml.read()

        root = et.fromstring(self.xmldocument)

        element = root.find(".//Property[@name='Comment']/String")
        self.documentProperties["Comment"] = element.attrib.get('value')

        element = root.find(".//Property[@name='Company']/String")
        self.documentProperties["Company"] = element.attrib.get('value')

        element = root.find(".//Property[@name='CreatedBy']/String")
        self.documentProperties["CreatedBy"] = element.attrib.get('value')

        element = root.find(".//Property[@name='CreationDate']/String")
        self.documentProperties["CreationDate"] = element.attrib.get('value')

        element = root.find(".//Property[@name='Id']/String")
        self.documentProperties["Id"] = element.attrib.get('value')

        element = root.find(".//Property[@name='Label']/String")
        self.documentProperties["Label"] = element.attrib.get('value')

        element = root.find(".//Property[@name='LastModifiedBy']/String")
        self.documentProperties["LastModifiedBy"] = element.attrib.get('value')

        element = root.find(".//Property[@name='LastModifiedDate']/String")
        self.documentProperties["LastModifiedDate"] = element.attrib.get('value')

        element = root.find(".//Property[@name='Uid']/Uuid")
        self.documentProperties["Uid"] = element.attrib.get('value')

        element = root.find(".//Property[@name='a2p_Version']/String")
        if element is not None:
            self.documentProperties["Assembly-A2P"] = element.attrib.get('value')

        element = root.find(".//Property[@name='Proxy']/Python")
        if element is not None:
            a3assy = element.attrib.get('module')
            if a3assy == 'freecad.asm3.assembly':
                self.documentProperties["Assembly-A3"] = 'A3-Assy'

        element = root.find(".//Property[@name='SolverId']/String")
        if element is not None:
            a4assy = element.attrib.get('value')
            if a4assy == 'Asm4EE':
                self.documentProperties["Assembly-A4"] = 'A4-Assy'

    def getFileName(self):
        return(self.fileName)

    def getFileSize(self):
        return(str(self.fileSize))

    def getDocumentPropterties(self):
        return(self.documentProperties)