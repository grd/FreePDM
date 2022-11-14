# FreePDM
***Concept Of Design***

## Workflows

### Administrators
![Usecase Diagram users](FreePDM_CoD-Figures/UC_Admin.png)


- Create new users
  - Edit user permissions
- Create new projects
  - Edit project Attributes
- Can override Part attributes(related to its state)
- Edit Server side settings

### Users

![Usecase Diagram users](FreePDM_CoD-Figures/UC_User.png)

- Open File from pdm
- Get write Access (Check-out)
  - List all files with write access
- Save file locally
- Save file to server (Check-in)
- Edit Metadata
  - Add (new) Metadata
- Create new part.

### Workflow Stories

- [File storing during editing](FreePDM_02-Workflows/01-FileStoringDuringEditing.md)
- [Check-Out File](FreePDM_02-Workflows/02-CheckoutFile.md)
- [Revision Procedures](FreePDM_02-Workflows/03-RevisionProcedures.md)
- [User Interface Functions](FreePDM_02-Workflows/04-UIFunctions.md)
- [User Interface Interaction FC](FreePDM_02-Workflows/05-UIInteractionFC.md)
- [Database Shape](FreePDM_02-Workflows/06-DbShape.md)
- [Database Interaction](FreePDM_02-Workflows/07-DbInteraction.md)
- [Development Builds](FreePDM_02-Workflows/08-DevelopmentBuilds.md)


[<< Previous Chapter](FreePDM_01-RequestedInformation.md) | [Content Table](README.md) | [Next Chapter >>](FreePDM_02-Workflows/01-FileStoringDuringEditing.md)
