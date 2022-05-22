# FreePDM
***Concept Of Design***


## Developments builds

TODO:

- Bash script for creating everything on a server / VM
- Docker environment

### Introduction

During the developments(and maybe also during it's use) it would be great if there is one or more methods to setup a environment.  
Be aware that in such case it would be an idea to create a [test database](06-DbShape.md) as well.

### Assumptions

- 

### Workflow 1  <!-- Bash script for setting up environment-->

Create a bash file for doing a set of actions.
This works well on linux like systems.
this could be a virtual machine or a Raspberry Pi.
The action the script can do are:

- Installing database software(SQL etc.)
- Installing the requested(specially if multiple are applicable) version management system.
- creating a database
- Downloading and Filling database with example dataset

Other thing can be added later on.

### questions / Comments 1

- If using a virtual machine be aware that you probably need port forwarding.

### Workflow 2  <!-- Docker container script for setting up environment-->

Another option is to create a Docker container or a docker script to create a environment where everything is available. Since docker containers don't have history it has some advantages to test quickly.
On the other hand if something needs to install every iteration has be done again. So it is more important to setup.

#### questions / Comments 2

[<< Previous Chapter](07-DbInteraction.md) | [Content Table](README.md) | [Next Chapter >>](../FreePDM_03-DesignDecisions.md)
