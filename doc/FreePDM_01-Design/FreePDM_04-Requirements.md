# FreePDM
***Concept Of Design***

## Requirements Document

The file below contains 4 tabs.
In the following paragraphs there are explanations of each tab.

[Spreadsheet with requirements](../FreePDM_CoD-Figures/FreePDM_04-Requirements.fods)

### Explanation on the Fourth Tab

The third Tab is added to mention all the assumptions done.
These should be tracked as well

### Explanation on the Third Tab

The second tab contains all the definitions.
**All definitions start also within the requirements with a capital letter**.
This is done to make aware there is something _special_ with it.
If a requirement contains a word that is part of the definitions but has no capital letter. please change it.

### Explanation on the Second Tab

The second Tab contains all the functions that are distilled out of the requirements of Tab one.
Every Line contains the following information

- Function id
- Function name
- Short Function description
- Refer to requirement (if None add a "-").
If multiple Requirements are applied duplicate function name

### Explanation on the First Tab

The first tab contain all the requirements.
Every line contain a 

- requirement id
- level
- traced from
- the requirements
- a note(optionally)

For example there is a requirement:
_The System provide capability to trace Revision history_.
This requirement has _requirement ID_:1 and is a system level requirement so _level_ 0.
and another requirement
_Every File Shall have a Revision_.
This requirement has _ID_ 1 as a parent and so it _trace back to_ requirement _ID_ 1.
Defining the (system )level is more difficult.  
Sometimes a line would be replaced by another one.
To track changes in a single overview the changes are done according the following method:

1. Copy the line you want to change to the end
2. Changed the original line by adding a Strike through.
3. Modify the line in the end to the new versions.

## Distilling Functions



## [Requirements Writing](FreePDM_04-1-RequirementsWriting.md)

[<< Previous Chapter](FreePDM_03-2-SVNProjectStructure.md) | [Content Table](README.md) | [Next Chapter >>](FreePDM_05-Architecture.md)
