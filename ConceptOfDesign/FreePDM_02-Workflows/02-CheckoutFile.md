# FreePDM
***Concept Of Design***

## update a FreeCAD part


### Assumptions

- Part already exist
- 


### Workflow 1

The user has opened a CAD file from the pdm(see: file storing during editing). While busy the need to edit a file that isn't checked out yet. He goes to that particular file in the assembly tree / goes to the TAB of that file. Then RBM on the part end press check-out to get write acces.

#### questions / Comments

- I assume this is not something we implement in one of the first states.(Probably state >4)...

### Workflow 2

The user want to open AND edit a CAD file from the pdm(see: file storing during editing). The user select the part in the pdm and end press check-out to get write acces.

#### questions / Comments

- We have to decide how to do it also Right Mouse Button, Via Menu or both
