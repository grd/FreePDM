# FreePDM
***Concept Of Design***


## file storing during editing


### Assumptions

- Part already exist
- search function exist and is implemented


### Workflow 1

The user goes to the pdm and searched for the required file. After this the file is downloaded to a local folder. 
- If this is a part he can open the part
- If it is an assembly it download all part(s) / Assemblie(s) that are required for opening the assembly
- If it is a drawing... it automatically download the required part(s) / Assemblie(s) dependent on what type of drawing it is
After editing the part / assembly / drawing it can be saved in the local folder.
When finished editing the file(s) can be checked in and then all the data is uploaded to the pdm. 

### questions / Comments 1

- What do we do with local work? are multiple branches / workspaces an idea to implememnt? If yes what state?(Probably state >4)...

### Workflow 2

The user goes to the pdm and searched for the required file. Then the file is downloaded and opened in FC.
- If this is a part he can open the part
- If it is an assembly it download all part(s) / Assemblie(s) that are required for opening the assembly
- If it is a drawing... it automatically download the required part(s) / Assemblie(s) dependent on what type of drawing it is
When saving this is done directly to the server. 

#### questions / Comments 2

- Since the risk on a crash is bigger we have to check how FC handle temporary files and how this can interact with the pdm
