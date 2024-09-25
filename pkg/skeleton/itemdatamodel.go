// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package skeleton

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/grd/FreePDM/pkg/util"
)

// Getting file info from Document.xml that is stored within each FC file

// linked item lists
// item_list = list[str, str, str, str]

// The item list is the list of linked items into a FreeCAD file.

// The arguments:
//     first argument is the name
//     the second is the type of the item
//     the third is the location of the item file
//     the fourth is the date stamp

// The overall document XML structure
type Document struct {
	XMLName        xml.Name `xml:"Document" json:"-"`
	SchemaVersion  int      `xml:"SchemaVersion,attr"`
	ProgramVersion string   `xml:"ProgramVersion,attr"`
	FileVersion    int      `xml:"FileVersion,attr"`
	Properties     struct {
		XMLName        xml.Name    `xml:"Properties" json:"-"`
		Count          int         `xml:"Count,attr"`
		TransientCount int         `xml:"TransientCount,attr"`
		UProperty      []UProperty `xml:"_Property"`
		Property       []Property  `xml:"Property"`
	}
	Objects struct {
		XMLName      xml.Name     `xml:"Objects" json:"-"`
		Count        int          `xml:"Count,attr"`
		Dependencies int          `xml:"Dependencies,attr"`
		ObjectDeps   []ObjectDeps `xml:"ObjectDeps"`
		Object       []Object     `xml:"Object"`
	}
}

type Properties struct {
	XMLName        xml.Name    `xml:"Properties" json:"-"`
	Count          int         `xml:"Count,attr"`
	TransientCount int         `xml:"TransientCount,attr"`
	PProperty      []UProperty `xml:"_Property"`
	Property       []Property  `xml:"Property"`
}

type UProperty struct {
	XMLName xml.Name `xml:"_Property" json:"-"`
	Name    string   `xml:"name,attr" json:"name"`
	Type    string   `xml:"type,attr" json:"type"`
	Status  int      `xml:"status,attr" json:"status"`
}

type Property struct {
	XMLName xml.Name `xml:"Property" json:"-"`
	Name    string   `xml:"name,attr" json:"name"`
	Type    string   `xml:"type,attr" json:"type"`
	Status  *int     `xml:"status,attr" json:"status,omitempty"`
	String  *struct {
		Value string `xml:"value,attr" json:"value,omitempty"`
	} `xml:"String" json:"String,omitempty"`
	Bool *struct {
		Value string `xml:"value,attr" json:"value,omitempty"`
	} `xml:"Bool" json:"Bool,omitempty"`
	Uuid *struct {
		Value string `xml:"value,attr" json:"value,omitempty"`
	} `xml:"Uuid" json:"Uuid,omitempty"`
	Map *struct {
		Count int `xml:"count,attr" json:"value,omitempty"`
	} `xml:"Map" json:"Map,omitempty"`
}

type ObjectDeps struct {
	XMLName xml.Name `xml:"ObjectDeps" json:"-"`
	Name    string   `xml:"Name,attr" json:"Name"`
	Count   int      `xml:"Count,attr" json:"Count"`
	Dep     *[]struct {
		Name string `xml:"Name,attr" json:"Name,omitempty"`
	} `xml:"Dep" json:"Dep,omitempty"`
}

type Object struct {
	XMLName    xml.Name    `xml:"Object" json:"-"`
	Name       string      `xml:"name,attr" json:"name"`
	Type       *string     `xml:"type,attr" json:"type,omitempty"`
	Id         *int        `xml:"id,attr" json:"id,omitempty"`
	Extension  *bool       `xml:"Extensions,attr" json:"extensions,omitempty"`
	Extensions *Extensions `xml:"Extensions" json:"Extensions,omitempty"`
}

type Extensions struct {
	XMLName   xml.Name `xml:"Extensions" json:"-"`
	Count     int      `xml:"Count,attr" json:"Count"`
	Extension []struct {
		Name string `xml:"Name,attr" json:"Name,omitempty"`
		Type string `xml:"Type,attr" json:"Type,omitempty"`
	}
}

func read() {
	prop := `
	<?xml version='1.0' encoding='utf-8'?>
	<!--
	 FreeCAD Document, see https://www.freecadweb.org for more information...
	-->
	<Document SchemaVersion="4" ProgramVersion="0.20R28647 (Git)" FileVersion="1">
		<Properties Count="15" TransientCount="3">
			<_Property name="FileName" type="App::PropertyString" status="50331649"/>
			<_Property name="Tip" type="App::PropertyLink" status="33554433"/>
			<_Property name="TransientDir" type="App::PropertyString" status="50331649"/>
			<Property name="Comment" type="App::PropertyString">
				<String value=""/>
			</Property>
			<Property name="Company" type="App::PropertyString">
				<String value=""/>
			</Property>
			<Property name="CreatedBy" type="App::PropertyString">
				<String value=""/>
			</Property>
			<Property name="CreationDate" type="App::PropertyString" status="16777217">
				<String value="2022-05-19T06:43:58Z"/>
			</Property>
			<Property name="Id" type="App::PropertyString">
				<String value=""/>
			</Property>
			<Property name="Label" type="App::PropertyString" status="1">
				<String value="0002"/>
			</Property>
			<Property name="LastModifiedBy" type="App::PropertyString">
				<String value=""/>
			</Property>
			<Property name="LastModifiedDate" type="App::PropertyString" status="16777217">
				<String value="2022-05-19T15:22:47Z"/>
			</Property>
			<Property name="License" type="App::PropertyString" status="1">
				<String value="All rights reserved"/>
			</Property>
			<Property name="LicenseURL" type="App::PropertyString" status="1">
				<String value="http://en.wikipedia.org/wiki/All_rights_reserved"/>
			</Property>
			<Property name="Material" type="App::PropertyMap" status="1">
				<Map count="0">
				</Map>
			</Property>
			<Property name="Meta" type="App::PropertyMap" status="1">
				<Map count="0">
				</Map>
			</Property>
			<Property name="ShowHidden" type="App::PropertyBool" status="1">
				<Bool value="false"/>
			</Property>
			<Property name="TipName" type="App::PropertyString" status="83886080">
				<String value=""/>
			</Property>
			<Property name="Uid" type="App::PropertyUUID" status="16777217">
				<Uuid value="ec21dc52-98d0-4d18-b29c-02cfe3cf78bf"/>
			</Property>
		</Properties>
		<Objects Count="17" Dependencies="1">
			<ObjectDeps Name="Body" Count="11">
				<Dep Name="Chamfer001"/>
				<Dep Name="Origin"/>
				<Dep Name="Sketch"/>
				<Dep Name="Pad"/>
				<Dep Name="Sketch001"/>
				<Dep Name="Hole"/>
				<Dep Name="Chamfer"/>
				<Dep Name="Chamfer001"/>
				<Dep Name="placement"/>
				<Dep Name="LCS_1"/>
				<Dep Name="LCS_2"/>
			</ObjectDeps>
			<ObjectDeps Name="Origin" Count="6">
				<Dep Name="X_Axis"/>
				<Dep Name="Y_Axis"/>
				<Dep Name="Z_Axis"/>
				<Dep Name="XY_Plane"/>
				<Dep Name="XZ_Plane"/>
				<Dep Name="YZ_Plane"/>
			</ObjectDeps>
			<ObjectDeps Name="X_Axis" Count="0"/>
			<ObjectDeps Name="Y_Axis" Count="0"/>
			<ObjectDeps Name="Z_Axis" Count="0"/>
			<ObjectDeps Name="XY_Plane" Count="0"/>
			<ObjectDeps Name="XZ_Plane" Count="0"/>
			<ObjectDeps Name="YZ_Plane" Count="0"/>
			<ObjectDeps Name="Sketch" Count="1">
				<Dep Name="XZ_Plane"/>
			</ObjectDeps>
			<ObjectDeps Name="Pad" Count="3">
				<Dep Name="Sketch"/>
				<Dep Name="Sketch"/>
				<Dep Name="Body"/>
			</ObjectDeps>
			<ObjectDeps Name="Sketch001" Count="1">
				<Dep Name="Pad"/>
			</ObjectDeps>
			<ObjectDeps Name="Hole" Count="3">
				<Dep Name="Sketch001"/>
				<Dep Name="Pad"/>
				<Dep Name="Body"/>
			</ObjectDeps>
			<ObjectDeps Name="Chamfer" Count="3">
				<Dep Name="Hole"/>
				<Dep Name="Hole"/>
				<Dep Name="Body"/>
			</ObjectDeps>
			<ObjectDeps Name="Chamfer001" Count="3">
				<Dep Name="Chamfer"/>
				<Dep Name="Chamfer"/>
				<Dep Name="Body"/>
			</ObjectDeps>
			<ObjectDeps Name="placement" Count="1">
				<Dep Name="Chamfer001"/>
			</ObjectDeps>
			<ObjectDeps Name="LCS_1" Count="1">
				<Dep Name="Chamfer001"/>
			</ObjectDeps>
			<ObjectDeps Name="LCS_2" Count="1">
				<Dep Name="Chamfer001"/>
			</ObjectDeps>
			<Object type="PartDesign::Body" name="Body" id="4258" />
			<Object type="App::Origin" name="Origin" id="4259" />
			<Object type="App::Line" name="X_Axis" id="4260" />
			<Object type="App::Line" name="Y_Axis" id="4261" />
			<Object type="App::Line" name="Z_Axis" id="4262" />
			<Object type="App::Plane" name="XY_Plane" id="4263" />
			<Object type="App::Plane" name="XZ_Plane" id="4264" />
			<Object type="App::Plane" name="YZ_Plane" id="4265" />
			<Object type="Sketcher::SketchObject" name="Sketch" id="4266" />
			<Object type="PartDesign::Pad" name="Pad" id="4267" />
			<Object type="Sketcher::SketchObject" name="Sketch001" id="4269" />
			<Object type="PartDesign::Hole" name="Hole" id="4270" />
			<Object type="PartDesign::Chamfer" name="Chamfer" id="4271" />
			<Object type="PartDesign::Chamfer" name="Chamfer001" id="4272" />
			<Object type="PartDesign::CoordinateSystem" name="placement" id="4278" />
			<Object type="PartDesign::CoordinateSystem" name="LCS_1" id="4279" />
			<Object type="PartDesign::CoordinateSystem" name="LCS_2" id="4280" />
		</Objects>
	</Document>
	`

	var document Document
	if err := xml.Unmarshal([]byte(prop), &document); err != nil {
		log.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(document, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

// The Class ItemDataModel is the main class for reading and writing
// one FCStd file. The file can be read without opening a FreeCAD file.

type StringHolder struct {
	Name  string
	Value string
}

type ItemDataModel struct {
	FileName string
	// xml_document        XmlNode
	thumbnail           []byte
	document_properties []byte
}

func InitItemDataModel(filename string) (ret ItemDataModel) {
	ret.FileName = filename
	ret.ReadFcFile()
	return ret
}

func (idm *ItemDataModel) ReadFcFile() {
	// Note that this temp-directory and all it's contents is automatically being deleted
	tempDir, err := os.MkdirTemp("/etc", "*-FC-File")
	util.CheckErr(err)
	defer os.RemoveAll(tempDir) // clean up
	log.Printf("Created temporary directory with path: %s", tempDir)

	// Unzipping the file to the directory
	err = unzipSource(idm.FileName, tempDir)
	util.CheckErr(err)

	dirList, err := os.ReadDir(tempDir)
	util.CheckErr(err)
	for _, item := range dirList {
		fmt.Println(item.Name())
	}

	// reading xml file
	idm.readXml(path.Join(tempDir, "Document.xml"))

	// Check whether there is a thumbnail
	if util.DirExists(path.Join(tempDir, "/thumbnails")) {
		idm.thumbnail, err = os.ReadFile(path.Join(tempDir, "/thumbnails/Thumbnail.png"))
		util.CheckErr(err)
	}
}

func (idm *ItemDataModel) readXml(fileName string) {
	// text, _ := os.ReadFile(fileName)
	// self.xml_document = loadXml(text)

	// // echo data.xml_document

	// var properties = findAll(self.xml_document, "Properties")

	// // echo "Length of properties" & $properties.len

	// root = et.fromstring(self.xml_document)

	// element = root.find(".//Property[@name='Comment']/String")
	// self.document_properties["Comment"] = element.attrib.get("value")

	// element = root.find(".//Property[@name='Company']/String")
	// self.document_properties["Company"] = element.attrib.get("value")

	// element = root.find(".//Property[@name='CreatedBy']/String")
	// self.document_properties["CreatedBy"] = element.attrib.get("value")

	// element = root.find(".//Property[@name='CreationDate']/String")
	// self.document_properties["CreationDate"] = element.attrib.get("value")

	// element = root.find(".//Property[@name='Id']/String")
	// self.document_properties["Id"] = element.attrib.get("value")

	// element = root.find(".//Property[@name='Label']/String")
	// self.document_properties["Label"] = element.attrib.get("value")

	// element = root.find(".//Property[@name='LastModifiedBy']/String")
	// self.document_properties["LastModifiedBy"] = element.attrib.get("value")

	// element = root.find(".//Property[@name='LastModifiedDate']/String")
	// self.document_properties["LastModifiedDate"] = element.attrib.get("value")

	// element = root.find(".//Property[@name='Uid']/Uuid")
	// self.document_properties["Uid"] = element.attrib.get("value")

	// self.document_properties["Assembly"] = None

	// element = root.find(".//Property[@name='a2p_Version']/String")
	// if element != nil {
	// 	self.document_properties["Assembly"] = "A2-Assy"
	// }
	// element = root.find(".//Property[@name='Proxy']/Python")
	// if element != nil {
	// 	a3assy = element.attrib.get("module")
	// 	if a3assy == "freecad.asm3.assembly" {
	// 		self.document_properties["Assembly"] = "A3-Assy"
	// 	}
	// }

	// element = root.find(".//Property[@name='SolverId']/String")
	// if element != nil {
	// 	a4assy = element.attrib.get("value")
	// 	if a4assy == "Asm4EE" {
	// 		self.document_properties["Assembly"] = "A4-Assy"
	// 		a4_assy = ItemAssemblyA4(root)
	// 		self.document_properties["Objects"] = a4_assy.get_external_item_list()
	// 		print(str(self.export_BOM()))
	// 	}
	// }

	// self.document_properties["Objects"] = []
}

// parsing the rest of the document

// Objects:

// Objects = root.find("./Objects")

// for att in Objects:
// 	if att.tag == "Object":
// 		type_name = att.attrib["type"]
// 		if type_name == "App::FeaturePython":
// 			if att.attrib["name"] == "Properties":
// 				self.read_properties()

// element = root.find(".//Property[@name='Comment']/String")
// self.document_properties["Comment"] = element.attrib.get('value')

//   <Property name="LastModifiedDate" type="App::PropertyString" status="16777217">
//     <String value="2022-11-21T13:03:46Z"/>
// </Property>

// element = root.find(".//Property[@name='a2p_Version']/String")
// if element is not None:
//     self.document_properties["Assembly"] = 'A2-Assy'

// for att in Objects:
//     if att.tag == "Object":
//         type_name = att.attrib["type"]
//         if type_name == "PartDesign::Body": print("Object Type = " + type_name + " Name = " + att.attrib["name"])
//         elif type_name == "App::Part": print("Object Type = " + type_name + " Name = " + att.attrib["name"])
//         elif type_name == "Part::FeaturePython": print("Object Type = " + type_name + " Name = " + att.attrib["name"]) // A2P en A3 assy
//         elif type_name == "App::DocumentObjectGroup": // A4 assy
//             print("Object Type = " + type_name + " Name = " + att.attrib["name"])
//             a4_assy = ItemAssemblyA4(self.xml_document)
//             self.document_properties["Objects"] = a4_assy.get_external_item_list()
//             print("Object list = " + str(self.document_properties["Objects"]))
//             break
//         elif type_name[0:3] == "App": continue
//         elif type_name[0:10] == "PartDesign": continue
//         elif type_name[0:8] == "Sketcher": continue
//         else: print("Object Type = " + type_name + " Name = " + att.attrib["name"])
// }

// func (self ItemDataModel) readProperties() {
// 	// Parsing the ObjectData
// 	ObjectData = self.root.find("./ObjectData")
// 	props = ObjectData.find("./Properties")
// 	count = props.attrib["Count"]
// 	print("Count = " + str(count))

// 	for att := range props {
// 		if att.tag == "Object" {
// 			type_name = att.attrib["type"]
// 		}
// 	}
// 	for xlink := ObjectData.findall(".//Property[@name='LinkedObject']/XLink") {
// 		xfile = xlink.attrib.get("file")
// 		xstamp = xlink.attrib.get("stamp")
// 		xname = xlink.attrib.get("name")
// 		xlist = append(xlist, [xfile, xname, xfile, xstamp])
// 	}
// 	return xlist
// }

func (idm ItemDataModel) GetFileName() string {
	return idm.FileName
}

// func (self ItemDataModel) get_document_propterties() TypedDict {
// 	return self.document_properties
// }

func (idm ItemDataModel) Save() {
	// TODO
}

// func (self ItemDataModel) add_variable(location, name, variable string){
// 	self.location = location
// 	self.name = name
// 	self.variable = variable
// 	// raise NotImplementedError("The function 'add_variable' is not implemented yet.")
// }

// func (self ItemDataModel) rename_item(source, dest) {
// 	self.source = source
// 	self.dest = dest
// 	// COMMENT: Is renaming of an item something that should be part of the (SQL) backend?
// }

// func (self ItemDataModel) is_assembly() bool {
// 	return self.document_properties["Assembly"] == true
// }

// func (self ItemDataModel) export_BOM() item_list {
// 	if self.is_assembly() == true {
// 		return self.document_properties["Objects"]
// 	} else {
// 		 return nil
// 	}
// }

// // The following classes contains methods that deal with
// // the following items: Object, Part, Group, ConfigurationTable,
// //                      Assembly 2P, 3, 4 and drawings

// class ItemTemplate:
//     func __init__(self, xml_root: str):
//         self.root = xml_root

//     func get_internal_item_list(self): // returns the list of components inside a FC file, such as Objects, Parts, Configurations and Drawings
//         pass

//     func get_external_item_list(self) -> item_list: // returns the list of linked items if applicable
//         pass

// class ItemObject:
//     pass

// class ItemPart:
//     pass

// class ItemGroup:
//     pass

// class ItemConfigurationTable:
//     pass

// class ItemAssemblyA2:
//     pass

// class ItemAssemblyA3:
//     pass

// class ItemAssemblyA4(ItemTemplate):
//     func __init__(self, tree: str):
//         ItemTemplate.__init__(self, tree)

//     func get_external_item_list(self) -> item_list:
//         object_list = []
//         xlist: item_list = []

//         // Parsing the Objects
//         Objects = self.root.find("./Objects")
//         for att in Objects:
//             if att.tag == "Object":
//                 type_name = att.attrib["type"]
//                 if type_name == "App::Link":
//                     // print("Object Type = " + type_name + " Name = " + att.attrib["name"])
//                     object_list.append(att.attrib["name"])

//         // Parsing the ObjectData
//         ObjectData = self.root.find("./ObjectData")
//         for xlink in ObjectData.findall(".//Property[@name='LinkedObject']/XLink"):
//             xfile = xlink.attrib.get('file')
//             xstamp = xlink.attrib.get('stamp')
//             xname = xlink.attrib.get('name')
//             xlist.append([xfile, xname, xfile, xstamp])
//         return xlist

type ItemDrawing struct {
	// generate a DXF ;-)
}
