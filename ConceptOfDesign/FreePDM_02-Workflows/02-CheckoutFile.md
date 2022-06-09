# FreePDM
***Concept Of Design***

## update a FreeCAD part


### Assumptions

- Part already exist

### Workflow 1

The user has opened a CAD file from the pdm(see: file storing during editing). While busy the need to edit a file that isn't checked out yet. He goes to that particular file in the assembly tree / goes to the TAB of that file. Then RBM on the part end press check-out to get write access.

#### questions / Comments

- I assume this is not something we implement in one of the first states.(Probably state >4)...
- For check-in: The same workflow but than with checkin

### Workflow 2

The user want to open AND edit a CAD file from the pdm(see: file storing during editing). The user select the part in the pdm and end press check-out to get write access.

#### questions / Comments

- We have to decide how to do it also Right Mouse Button, Via Menu or both
- For check-in: The same workflow but than with checkin

_Note grd: I think more about adding a button [button](../FreePDM_CoD-Figures/check-in-out.png) with the name "Check Out", "Check In" and at the left of each item a selection button. The procedure is that you first select which items you want to check out. Then you press the Check Out button and the items are checked out and you can do your work. An indicator will show which user has checked out an item. When you are ready with your work then you can select those files again and press Check In._  
_Note Jee-Bee: What you describe is workflow 2. The Symbol that that it is checked out is a good addition!_


### Extra
There is a [sequence diagram](../FreePDM_CoD-Figures/SEQ_CheckIn-CheckOut.png) added with a more graphical representation of the process.

[<< Previous Chapter](FreePDM_01-RequestedInformation.md) | [Content Table](README.md) | [Next Chapter >>](FreePDM_03-DesignDecisions.md)