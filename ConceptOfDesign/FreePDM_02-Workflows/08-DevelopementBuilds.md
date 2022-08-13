# FreePDM
***Concept Of Design***


## Building the project

### Introduction

During the development (and maybe also during its use) it would be great if there is one or more methods to setup an environment.  
Be aware that in such case it would be an idea to create a [test database](06-DbShape.md) as well.  
**Edit:** building a server, virtual machine or Docker container is not only necessary during build, but also during deployment.

### Assumptions

- 

### Workflow 1  <!-- Bash script for setting up environment-->

Create a bash file for doing a set of actions.
This works well on Linux-like systems.
This could be a virtual machine or a Raspberry Pi.
The actions the script can do are:

- installing database software(SQL etc.)
- installing the requested version management system (especially if multiple are applicable).
- creating a database
- downloading and filling database(s) with example dataset

Additional features can be added later.

### questions / Comments 1

- If using a virtual machine be aware that you probably need port-forwarding.

### Workflow 2  <!-- Cmd OR Powershell script for setting up environment-->

For actions the same as [Workflow 1](#workflow-1), but now with _Power Shell_ or _cmd prompt_.

### Workflow 3  <!-- Docker container script for setting up environment-->

Another option is to use a Docker container or a Docker script to create an environment where everything is available.
Since Docker containers don't have history it has some advantages for rapidly testing.
On the other hand, if something needs to be installed then every iteration has to be done again.
So it is more important to setup.

#### questions / Comments 3

- [BeCFG makes active use of docker containers](https://docs.becpg.fr/en/installation/server_installation.html#software-prerequisites)

[<< Previous Chapter](07-DbInteraction.md) | [Content Table](README.md) | [Next Chapter >>](../FreePDM_03-DesignDecisions.md)
