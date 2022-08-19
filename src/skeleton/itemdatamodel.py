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
import zipfile
from PySide2 import QtGui

# The Class ItemDataModel is the main class for reading and writing
# one FCStd file. The file can be read without opening a FreeCAD file.


class ItemDataModel():
    def __init__(self, filename):  # filename is also imported so change
        self.file_name = filename
        self.xml_document = None
        self.thumbnail = None
        self.document_properties = {}
        self.read()

        # check out file size
        self.fileSize = os.path.getsize(filename)

    def read(self):
        with tempfile.TemporaryDirectory() as tmpdirname:
            # print('created temporary directory', tmpdirname)

            # extracting the FCStd zip file
            zippedfile = zipfile.ZipFile(self.file_name)  # zip is python internal function
            zippedfile.extractall(tmpdirname)

            # reading xml file
            tree = tmpdirname + "/Document.xml"
#            print(tree)
            self.read_xml(tree)

            # Check whether there is a thumbnail
            if os.path.isdir(tmpdirname + "/thumbnails"):
                self.document_properties["thumbnail"] = tmpdirname + "/thumbnails/Thumbnail.png"
                self.thumbnail = QtGui.QImage(self.document_properties['thumbnail'])

    def read_xml(self, data):
        with open(data) as docxml:
            self.xml_document = docxml.read()

        root = et.fromstring(self.xml_document)

        element = root.find(".//Property[@name='Comment']/String")
        self.document_properties["Comment"] = element.attrib.get('value')

        element = root.find(".//Property[@name='Company']/String")
        self.document_properties["Company"] = element.attrib.get('value')

        element = root.find(".//Property[@name='CreatedBy']/String")
        self.document_properties["CreatedBy"] = element.attrib.get('value')

        element = root.find(".//Property[@name='CreationDate']/String")
        self.document_properties["CreationDate"] = element.attrib.get('value')

        element = root.find(".//Property[@name='Id']/String")
        self.document_properties["Id"] = element.attrib.get('value')

        element = root.find(".//Property[@name='Label']/String")
        self.document_properties["Label"] = element.attrib.get('value')

        element = root.find(".//Property[@name='LastModifiedBy']/String")
        self.document_properties["LastModifiedBy"] = element.attrib.get('value')

        element = root.find(".//Property[@name='LastModifiedDate']/String")
        self.document_properties["LastModifiedDate"] = element.attrib.get('value')

        element = root.find(".//Property[@name='Uid']/Uuid")
        self.document_properties["Uid"] = element.attrib.get('value')

        element = root.find(".//Property[@name='a2p_Version']/String")
        if element is not None:
            self.document_properties["Assembly"] = 'A2-Assy'

        element = root.find(".//Property[@name='Proxy']/Python")
        if element is not None:
            a3assy = element.attrib.get('module')
            if a3assy == 'freecad.asm3.assembly':
                self.document_properties["Assembly"] = 'A3-Assy'

        element = root.find(".//Property[@name='SolverId']/String")
        if element is not None:
            a4assy = element.attrib.get('value')
            if a4assy == 'Asm4EE':
                self.document_properties["Assembly"] = 'A4-Assy'

    def get_file_name(self):
        return(self.file_name)

    def get_file_size(self):
        return(str(self.fileSize))

    def get_document_propterties(self):
        return(self.document_properties)

    def save(self):
        pass

    def add_variable(location, name, variable):
        pass

    def rename_item(source, dest):
        pass

# The following classes contains methods that deal with
# the following items: Object, Part, Group, ConfigurationTable, 
#                      Assembly 2P, 3, 4 and drawings

class ItemTemplate:
    pass

class ItemObject:
    pass

class ItemPart:
    pass

class ItemGroup:
    pass

class ItemConfigurationTable:
    pass

class ItemAssemblyA4:
    pass

class ItemAssemblyA3:
    pass

class ItemAssemblyA2:
    pass

class ItemDrawing:
    pass

