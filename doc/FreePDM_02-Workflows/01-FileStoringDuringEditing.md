# FreePDM
***Concept Of Design***


## file storing during editing


### Assumptions

- Part already exists
- Search function exists and is implemented


### Workflow 1

The user goes to the PDM and searches for the required file. Subsequently, the file is downloaded to a local folder. 
- If this is a part, they can open the part
- If it is an assembly, it downloads all part(s) / assembly(ies) that are required for opening the assembly. There should also be checkboxes for downloading the individual drawings of the items
- If it is a drawing... it automatically downloads the required part(s) / assembly(ies) dependent on what type of drawing it is
After editing the part / assembly / drawing it can be saved in the local folder.
When completing edits to the file(s), it/they be checked in and then all the data is uploaded to the PDM. 

### questions / Comments 1

- What do we do with local work? Are multiple branches / workspaces an idea to implement? If yes what state? (Probably state >4)...

### Workflow 2

The user goes to the PDM and searches for the required file. Then the file is downloaded and opened within FreeCAD.
- If this is a part, they can open the part
- If it is an assembly, it downloads all part(s) / assembly(ies) that are required for opening the assembly
- If it is a drawing... it automatically downloads the required part(s) / assembly(ies) dependent on what type of drawing it is   

When saving, this is done directly to the server. 

#### questions / Comments 2

- Since the risk of a crash is probable, we have to check how FreeCAD: 
    1. handles temporary files
    1. how this can interact with the PDM
- The other possibility (close all files first), is to: 
    1. Extract and uncompress the contents of a FreeCAD archive (note: the file format [FCStd](https://wiki.freecad.org/File_Format_FCStd) is actually a compressed gzipped file)
    1. Read/Write the `Document.xml`  file
    1. Re-compress everything back in to an .FCStd again. 
    
    This is easy because we know the name `Document.xml` and therefore we can leverage a script to automate this. See [this stackexchange](https://superuser.com/questions/647674/is-there-a-way-to-edit-files-inside-of-a-zip-file-without-explicitly-extracting) post. But, we need to be sure that the file is in fact closed. See question #3

#### questions / Comments 3

How can you check whether a file is open or closed?

[<< Previous Chapter](../FreePDM_02-Workflows.md) | [Content Table](README.md) | [Next Chapter >>](02-CheckoutFile.md)