# FreePDM
***Concept Of Design***


## Revision procedures


### Assumptions

- All model types (like part, assembly, drawing, but also (external) documentation) are called item

### Workflow 1

When a new item is needed it is possible / required to create this from within PDM system.
Within the PDM system, find `Create new item`.
From now on the model can be `CheckedOut`, edited, and `CheckedIn` until the model is ready for release. 
Now we can start a procedure (can be self-defined) that checks the item(s) etc... if it/they is/are release ready. If all checks are good, than the item is pushed to the release state.

### questions / Comments 1

- in [this FC forum thread](https://forum.freecad.org/viewtopic.php?f=8&t=68350) it mentions the usage of SVN. Is this possible with this? 
  - grd: To be honest, I don't know but I think that is possible. I looked at some videos that showed that it's possible. But I also talked to Yorik and he said that he liked `git` as well. So database should be independent of which kind of database you use.
- If this procedure is initiated within FreeCAD it should be possible to directly open the created part. This could also be two feature (`Create new item` AND `Create & Open new item`).
  - grd: I don't know. We are gonna find out.
- There should be a method added that there are multiple people who can check out simultaneously 
   - grd: I don't think so. Check out yes. But when the file is checked out the file is locked for everyone else. And checking in, that should be done with the person who checked it out. 

### Workflow 2

If a part has the **revision** state `NoOne` has the right to edit this item.
If the item needs a change: the only way is through [workflow 3](#workflow-3)

#### questions / Comments 2


### Workflow 3

An item that has a revision state needs an update. This can be done through changing the item's state from **revision** to **work** (define proper name). From now on [workflow 1](#workflow-1) starts again. Except that revision is iterated by 1.

#### questions / Comments 3

- Which roles have the permission to do this?
- Let's say a part and its corresponding drawing are in the revision state. But later on we notice that something is changed that has no influence on the model. Is it valid that the drawing gets a new revision but the part holds the existing one? See example [forum post](https://forum.freecad.org/viewtopic.php?f=8&t=68350&start=60#p594331) 
<!--I wrote it the wrong way around. Of course this change if the drawing is inside the related part / assembly. 
Let's assume i created a assembly and a drawing. everything is released and there has to be a change for example i described in the notes a type of glue that is not strong enough. is it valid to change the drawing without releasing the model. (So the drawing get release version 2, but the assembly hold release state 1)-->

[<< Previous Chapter](02-CheckoutFile.md) | [Content Table](README.md) | [Next Chapter >>](04-UIFunctions.md)