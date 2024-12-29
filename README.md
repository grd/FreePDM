<center> 

![FreePDM-logo](assets/FreePDM-logo.png)

</center> 

# FreePDM

A PDM for FreeCAD. A PDM is a Product Data Management system. Per [wikipedia](https://en.wikipedia.org/wiki/Product_data_management): 

> Product data management (PDM) should not be confused with product information management (PIM). PDM is the name of a business function within product lifecycle management (PLM) that is denotes the management and publication of product data. 

>In software engineering, this is known as version control. The goals of product data management include ensuring all stakeholders share a common understanding, that confusion during the execution of the processes is minimized, and that the highest standards of quality controls are maintained. 

[FreeCAD](https://www.freecad.org) is a free libre opensource cross-platform Computer Aided Design (CAD) software. 

## Background
The initial idea is to make a Skeleton (model), GUI and an Admin module. The GUI is based on Fyne, and will be a multi platform app.

### Previous attempts made at creating a FOSS PDM / PLM
* [OpenPLM](https://github.com/cadracks-project/openplm) (abandoned)
* German users on the FreeCAD forum attempted a PDM (see FreeCAD [forum thread](https://forum.freecad.org/viewtopic.php?f=22&t=63794))
* [Taack PLM](https://github.com/Taack/taack-plm-freecad)
* [Ondsel](https://ondsel.com/) A very promising design. Unfortunately it shut down.
* [NanoPLM](https://github.com/alekssadowski95/nanoPLM)

Relevant:
* The FreeCAD [Reporting workbench](https://github.com/furti/FreeCAD-Reporting) addon that uses SQL to extract information out of a FreeCAD document.
* The [fcinfo](https://wiki.freecad.org/Macro_FCInfo) macro for measuring the weight of a model.

## Install
Proposed [Install](ConceptOfDesign/FreePDM_Install/README.md).

## Concept

Proposed [concept](ConceptOfDesign/FreePDM_01-Design/README.md) of design.

## Workflow

Proposed [workflow](ConceptOfDesign/FreePDM_02-Workflows/README.md).

## Licence
MIT [LICENSE](LICENSE)
