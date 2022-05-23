# FreePDM
***Concept Of Design***

## Update Database Version Methods

### Introduction

There is a lot of research needed to know what the best way is to handle location, version management and database management.
This work here is an extension on the work done [here](FreePDM_03-2-SVNProjectStructure.md).


### Tree way system

#### General

First an explanation about the context see [diagram](BDD-Item-3Way) below.
The _Item_ is the base building block for everything.
The _Item_ connect to the database, to the versioning system.
At the same time the _Item_ holds all , like CAD files, Documents etc.
(See all the 'blocks' that hang below the item).  
_Note: A  Triangle arrow means general case.
So a document is a general case and a docx file is a special case of a document._

![BDD Item 3Way](FreePDM_CoD-Figures/BDD_Item-3Way.png)

The next diagram is to create a better understanding what happens inside the _Item_ 'block'.

![IBD Item 3Way](FreePDM_CoD-Figures/IBD_Item-3Way.png)

In this diagram everything is a bit more general.
So all the special cases are ignored.
There are _Connections_ between the different files and the _meta-data_.
This are data transfers between the files / based on the files and the _Meta-Data_('local' database) of the _Item_.  
From the _Meta-Data_ there is a _Connection_ to the right side of the _Item_.
All Other Files have a _connection_ to the left side of the _Item_.
What this mean are explained in [Interfaces](#interface).

#### Interface

Basically an _Interface_ show what systems talk to each other / Need information from each other.
In the diagram below show the _Interfaces_ between the _Item_, _Database_ and _Versioning System_.

![BDD ItemInterface 3Way](FreePDM_CoD-Figures/BDD_ItemInterface-3Way.png)

If we go back to the _Item_ block there we have two _Interfaces_. These two _Interfaces_ correspondent with the diagram directly above.  
As al the _Rationale_ and _Problems_ explain there are some difficulties to overcome.  
The main issue is split between the _Version system_ and the _database_. 
In this split lies ho to handle binary files where most version system are not good at.
Binary files are in such a system basically all copies and not delta files what it is with text files.

### two way system

#### General

There are two alternatives in a two way system that walk side by side.
The first alternative that is part of this whole research.
The second alternative drops in by the [_Interface_](#Interface).  
Below a diagram that is almost equal to the equivalent diagram in the [tree-way](#Tree-way-system) method.
The difference is the missing _Version System_.

![BDD Item](FreePDM_CoD-Figures/BDD_Item.png)

The next diagram is to create a better understanding what happens inside the _Item_ 'block'.

![IBD ItemInterface](FreePDM_CoD-Figures/IBD_ItemInterface.png)

But as you can see there is now a single interface.
Both the files and the _Meta-data_ goes to the the same 

#### Interface

Here is the diverge between the two alternatives.
The divergence have also change inside the _Item_ block as result.
But what that is follows later.

![BDD ItemInterface](FreePDM_CoD-Figures/BDD_ItemInterface.png)

First we ignore the _LDAP-server_.
What do we see then:  
All data goes to the _Database_, Some part directly(The binary Files) and some via the _Meta-Data_.

It is possible to add binary data to a sql server. Normally a database have a known width(call it that way) and grows only length wise.
In this case we have an unknown amount of versions (of every CAD file) that are stored.
With this method the database expand in both directions.
For comparison with preview images.
A preview image in the _Database_ is fine since every item get's a single image.
The history of images is not important so only the last version needs acces tot the database. 
If the preview change the existing file is overwritten.
So the size of the table don't expand.

Now the second alternative.
The _LDAP-server_ is now taken into account.

What is a LDAP server?
According to Okta (See What is LDAP) LDAP is the following:

> LDAP is an open, vendor-neutral application protocol for accessing and maintaining that data. 
> LDAP can also tackle authentication, so users can sign on just once and access many different files on the server.

This is one of the big problems.
Who have acces to what file(s).
If we go now to the internals of an _Item_.
In the situation without _LDAP-server_ (as shown)all files are synchronised with the _Database_. 
In the version with _LDAP-server_ is much more easy. Only the _Meta-Data_ is synchronised with the _Database_ because the data itself is part of the _LDAP-server_.
This makes this implementation less difficult.  
Another pro of this implementation is that it is widely used by other PLD / PDM systems. see for reference by the [Additional information](#additional-information).

The versioning is not handled by this system.
Since it are binary files it is equally problematic as with a versioning system (that also create copies).  
FreeCAD saves backup files as _FCStd#_ where the _#_ is a number.
It is possible make use of this same scheme.
If preferred Al versions that are before a release can be removed at the moment of a revision.

#### Additional information

LDAP server in general

- [LDAP wikipedia](https://en.wikipedia.org/wiki/Lightweight_Directory_Access_Protocol)
- [What is LDAP](https://www.okta.com/identity-101/what-is-ldap/)
- [Open source LDAP server](https://opensource.com/business/14/5/four-open-source-alternatives-LDAP)
- [Difference between SQL and LDAP](https://stackoverflow.com/questions/5075394/difference-between-sql-and-ldap#5075461)

It looks like That Windchill, Teamcenter AND SolidWorks PLM use a combination of SQL server && LDAP server.

- [Creo LDAP](https://community.ptc.com/t5/Windchill/what-is-the-purpose-of-Method-server-Server-Manager-and-LDAP/td-p/211635
)
- [Teamcenter LDAP question](https://community.sw.siemens.com/s/question/0D54O000061xshJSAQ/teamcenter-sso-setup-unknown-ldap-exception)
- [Solid Works LDAP](https://help.solidworks.com/2022/english/enterprisepdm/archiveserver/c_Login_Settings.htm)

Databases and storing binary data.

- [binary data types in sql server](https://www.thoughtco.com/binary-data-types-in-sql-server-1019807)
- [storing binary data types in sql server](https://codingsight.com/storing-binary-data-types-in-sql-server/)
- [my-SQL binary varbinary](https://dev.mysql.com/doc/refman/8.0/en/binary-varbinary.html)

### Conclusion

From the three Alternatives looks it like the system version without _Versioning System_ and with _LDAP-server_ is implementation wise the best option.  
Is the _LDAP-server_ alway required?
Maybe not for small companies where _ALl_ users have acces to everything. As long a file and folder structure can be maintained.

[<< Previous Chapter](FreePDM_03-2-SVNProjectStructure.md) | [Content Table](README.md) | [Next Chapter >>](FreePDM_04-Requirements.md)
