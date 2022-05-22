# FreePDM
***Concept Of Design***

## Attributes

### Stored general attributes
What database are the databases convertible?
See for set of example attributes Part [Attributes pdf](PartAttributes.pdf) original see [Windchill Attributes](https://support.ptc.com/help/windchill/wc110_hc/whc_en/index.html#page/Windchill_Help_Center%2FPMPartAttributesRef.html)

The properties of a FreeCAD file: 
- [FreeCAD file Data](https://wiki.freecadweb.org/Std_Part#Data)
- [Body Data](https://wiki.freecadweb.org/PartDesign_Body#Hidden_properties_Data)
- [Assembly A2 Plus](https://wiki.freecadweb.org/A2plus_Workbench#Assembly_Structure)
- Assembly A3 and A4 uses Std LinkMate for their data.
[Std LinkMate Data](https://wiki.freecadweb.org/Std_LinkMake#Data)

### Parts vs Bodies
A Part is a container and it can contain multiple bodies and as so it can become an assembly (of Bodies). A Body is a 3D geometry that can either be a solid or a mesh object.

### Stored body attributes

- Name / Number (is this the unique identifier?)
- Description
- Date of creation
- User Name
- Material
- surface finish
- Surface treathment
- Weight(or volume and calculate weight)
- Key words
- Unit? (FreeCAD only allows the Metric System)

### Stored Assembly Attributes

- Name / Number (is this the unique identifier?)
- Description
- Date of creation
- User Name
- BOM (Read out file?)
- Surface treathment (for example for weldings)
- Weight
- keywords

### Stored Drawing? Attributes

- Name / Number (is this the unique identifier?)
- Description
- Date of creation
- User Name
- Drawing standards
- keywords
- Revision data? Revision Text, Name, Date, Revision

### User Attributes
- Name / Number (is this the unique identifier?)
- Access to projects / top level Systems
- Role(s)

[<< Previous Chapter](FreePDM_03-DesignDecisions.md) | [Content Table](README.md) | [Next Chapter >>](FreePDM_04-Requirements.md)
