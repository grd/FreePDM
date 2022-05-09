# FreePDM
***Concept Of Design***


## Revision procedures


### Assumptions

- All model types(like part, assemly, drawing, but also (external )documentation) are called item


### Workflow 1

When a new item is needed it is possible / required to create this from the pdm system.
This can / have to be done by go to the pdm system, than _Create new item_.  
From now on the model can be _CheckedOut_, edited and _CheckedIn_ until the model is ready for release.  
Now can start a procedure(Can be self defined) that check the item(s) etc if it is release ready. If all checks are good than the item go to the release state.

### questions / Comments 1

- in [this FC forum thread](https://forum.freecadweb.org/viewtopic.php?f=8&t=68350) there is spoken about the usage of svn. Is this possible with this
- If this procedure is start within FC it should be possible to directly open the created part. This could also be two feature(_Create new item_ AND _Create & Open new item_).
- There should be a method added that there are one or more people who can check the 

### Workflow 2

If a part has the **revision** state _NoOne_ has the right to edit this item.
If the item needs a change: The only way is through [workflow 3](#workflow-3)

#### questions / Comments 2


### Workflow 3

An item that have a revision state needs an update. This can be done through change the items from **revision** to a **work**(define proper name) state. From now on [workflow 1](#workflow-1) start again. except that revision is added 1 up.

#### questions / Comments 3

- Which roles have the permission to do this?
- Let's say a part and the coresponding drawing are in the revision state. But later on there is noticed that something is forgotten that has no influence on the model(for example a typo). Is it valid that the drawing get's a new revision but the part holds the existing one?

