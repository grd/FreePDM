"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

# Getting file info from Document.xml that is stored within each FC file

# from fileinput import filename  # not used
from hashlib import blake2b
import defusedxml.ElementTree as et
# from tempfile import tempdir  # not used
import tempfile
import os
import zipfile
from typing import TypedDict

from PySide2 import QtGui

# linked item lists
item_list = list[str, str, str, str]
"""
The item list is the list of linked items into a FreeCAD file.

The arguments:
    first argument is the name
    the second is the type of the item
    the third is the location of the item file
    the fourth is the date stamp
"""

# The Class ItemDataModel is the main class for reading and writing
# one FCStd file. The file can be read without opening a FreeCAD file.


class ItemDataModel():
    def __init__(self, filename: str):  # filename is also imported so change
        self.file_name: str = filename
        print(self.file_name)
        self.xml_document: str = None
        self.thumbnail = None
        self.document_properties: TypedDict = {}
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
            # print(tree)
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

        self.document_properties["Assembly"] = None

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
                a4_assy = ItemAssemblyA4(root)
                self.document_properties["Objects"] = a4_assy.get_external_item_list()
                print(str(self.export_BOM()))


        self.document_properties["Objects"] = []

        # parsing the rest of the document
    
        # Objects:

        Objects = root.find("./Objects")

        # for att in Objects:
        #     if att.tag == "Object":
        #         type_name = att.attrib["type"] 
        #         if type_name == "PartDesign::Body": print("Object Type = " + type_name + " Name = " + att.attrib["name"])
        #         elif type_name == "App::Part": print("Object Type = " + type_name + " Name = " + att.attrib["name"])
        #         elif type_name == "Part::FeaturePython": print("Object Type = " + type_name + " Name = " + att.attrib["name"]) # A2P en A3 assy
        #         elif type_name == "App::DocumentObjectGroup": # A4 assy
        #             print("Object Type = " + type_name + " Name = " + att.attrib["name"])
        #             a4_assy = ItemAssemblyA4(self.xml_document)
        #             self.document_properties["Objects"] = a4_assy.get_external_item_list()
        #             print("Object list = " + str(self.document_properties["Objects"]))
        #             break
        #         elif type_name[0:3] == "App": continue
        #         elif type_name[0:10] == "PartDesign": continue
        #         elif type_name[0:8] == "Sketcher": continue
        #         else: print("Object Type = " + type_name + " Name = " + att.attrib["name"])
        

    def get_file_name(self) -> str:
        return(self.file_name)

    def get_file_size(self) -> str:
        return(str(self.fileSize))

    def get_document_propterties(self) -> TypedDict:
        return(self.document_properties)

    def save(self):
        pass

    def add_variable(self, location, name, variable):
        self.location = location
        self.name = name
        self.variable = variable
        raise NotImplementedError("The function 'add_variable' is not implemented yet.")

    def rename_item(self, source, dest):
        self.source = source
        self.dest = dest
        raise NotImplementedError("The function 'rename' is not implemented yet.")
        # COMMENT: Is renaming of an item something that should be part of the (SQL) backend?

    def is_assembly(self) -> bool: 
        return self.document_properties["Assembly"] != None

    def export_BOM(self) -> item_list:
        if self.is_assembly() == True:
            return self.document_properties["Objects"]
        else: return None




# The following classes contains methods that deal with
# the following items: Object, Part, Group, ConfigurationTable, 
#                      Assembly 2P, 3, 4 and drawings

class ItemTemplate:
    def __init__(self, xml_root: str):
        self.root = xml_root


    def get_internal_item_list(self): # returns the list of components inside a FC file, such as Objects, Parts, Configurations and Drawings
        pass

    def get_external_item_list(self) -> item_list: # returns the list of linked items if applicable
        pass

class ItemObject:
    pass

class ItemPart:
    pass

class ItemGroup:
    pass

class ItemConfigurationTable:
    pass

class ItemAssemblyA2:
    pass

class ItemAssemblyA3:
    pass

class ItemAssemblyA4(ItemTemplate):
    def __init__(self, tree: str):
        ItemTemplate.__init__(self, tree)

    def get_external_item_list(self) -> item_list:
        object_list = []
        xlist: item_list = []

        # Parsing the Objects
        Objects = self.root.find("./Objects")
        for att in Objects:
            if att.tag == "Object":
                type_name = att.attrib["type"] 
                if type_name == "App::Link":
                    # print("Object Type = " + type_name + " Name = " + att.attrib["name"])
                    object_list.append(att.attrib["name"])

        # Parsing the ObjectData
        ObjectData = self.root.find("./ObjectData")
        for xlink in ObjectData.findall(".//Property[@name='LinkedObject']/XLink"):
            xfile = xlink.attrib.get('file')
            xstamp = xlink.attrib.get('stamp')
            xname = xlink.attrib.get('name')
            xlist.append([xfile, xname, xfile, xstamp])
        return xlist



   


class ItemDrawing:
    pass
