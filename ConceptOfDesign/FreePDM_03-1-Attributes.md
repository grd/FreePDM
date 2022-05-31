# FreePDM
***Concept Of Design***

## Attributes

### Stored general attributes
What database are the databases convertible?
See for set of example attributes in the [Links to documents](#links-to-documents)

The set of properties of a FreeCAD file:

- [FreeCAD file Data](https://wiki.freecadweb.org/Std_Part#Data)
- [Body Data](https://wiki.freecadweb.org/PartDesign_Body#Hidden_properties_Data)
- [Assembly A2 Plus](https://wiki.freecadweb.org/A2plus_Workbench#Assembly_Structure)
- Assembly A3 and A4 uses Std LinkMate for their data.
[Std LinkMate Data](https://wiki.freecadweb.org/Std_LinkMake#Data)

### Parts vs Bodies
A Part is a container and it can contain multiple bodies and as so it can become an assembly (of Bodies). A Body is a 3D geometry that can either be a solid or a mesh object.

### Stored body attributes

- Name / Number (is the number an unique identifier?)
- Description
- Date of creation
- User Name
- Material (and mass)
- surface finish
- Surface treatment
- Weight(or volume and calculate weight)
- Key words
- Unit? (FreeCAD only allows the Metric System)

### Stored Assembly Attributes

- Name / Number (is this the unique identifier?)
- Description
- Date of creation
- User Name
- BOM (Read out file?)
- Surface treatment (for example for weldings)
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

### Links to documents

Links in alphabetic order:

- [NewUserAttributes pdf](NewUserAttributes.pdf) original see [Windchill New User Attributes](https://support.ptc.com/help/windchill/wnc/r11.1_m020/whc_en/index.html#page/Windchill_Help_Center/ParticipantAdminUserCreateAttributesTab.html)
- [PartAttributes pdf](PartAttributes.pdf) original see [Windchill Part Attributes](https://support.ptc.com/help/windchill/wc110_hc/whc_en/index.html#page/Windchill_Help_Center%2FPMPartAttributesRef.html)
- [ProjectAttributes pdf](ProjectAttributes.pdf) original see [Windchill Project Attributes](https://support.ptc.com/help/wnc/r12.0.2.0/en/index.html#page/Windchill_Help_Center/ProjMgmtProjectAttributes.html)
- [DefaultSystemAttributes pdf](DefaultSystemAttributes.pdf) original see [Windchill Default System Attributes](https://support.ptc.com/help/windchill/wc110_hc/whc_en/index.html#page/Windchill_Help_Center/WGMCATIAV5AdmConfigWCDefaultSysAtt.html)

Check also [database setup](FreePDM_07-DatabaseSetup.md) for more (and better split out )attributes.

[<< Previous Chapter](FreePDM_03-DesignDecisions.md) | [Content Table](README.md) | [Next Chapter >>](FreePDM_03-2-SVNProjectStructure.md)
