#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from database import Base
from database import Session
from project import Project

# Tables
from pdm_tables import PdmUser
from pdm_tables import PdmRole
from pdm_tables import PdmProject
from pdm_tables import PdmItem
from pdm_tables import PdmModel
from pdm_tables import PdmDocument
from pdm_tables import PdmMaterial
from pdm_tables import PdmHistory
from pdm_tables import PdmPurchase
from pdm_tables import PdmManufacturer
from pdm_tables import PdmVendor
# from sqlalchemy import select
# from collections.abc import Sequence
from typing import Optional, Union

# Note:
# According to: https://docs.sqlalchemy.org/en/14/tutorial/engine.html
# 'The engine is typically a global object created just once for a particular database server, ...'
#
# So up till now i expected every user has it's own login but it looks like this is not possible using engines.
# Now there has to be some research in a dedicated login system!


class Item():
    """Item related Class"""

    def __init__(self):
        print("Generic Item")

    def create_item_number(self, number: Union[str, int], ndigits: Optional[int]) -> str:
        """
        Create new project number

        Parameters
        ----------

        number [str] : Last current number
            Number as string including leading zeros.

        ndigits [int] : number of digits
            Number of digits of the number length.
            If ndigits is -1 the length is just the length.

        Returns
        -------

        New number[str]
        """
        self.number = number
        self.ndigits = ndigits

        call_proj = Project()
        item_nr = call_item.create_number(number, ndigits)

        return(item_nr)

    # https://stackoverflow.com/questions/73887390/handle-multiple-users-login-database-with-sqlalchemy
    def create_item(self, project: str, path: str, number: Optional[str], name: Optional[str], description: Optional[str], full_description: Optional[str]) -> None:
        """
        Create new item

        Parameters
        ----------

        project [str] : Project number or name
            Project name or project number have to be a string for use whit leading zeros.

        path [str] : path
            Path to folder where related models and documents are stored.

        name [str] : Optional Name
            Name is optional for when no number is wanted.
            If name is None then automatically a number is generated, otherwise added number is used.

        Returns
        -------

        No return - new object created in SQL database ?
        """
        self.project = project  # User works on current project
        self.path = path  # TODO: create path automatically
        self.number = number
        self.name = name
        self.description = description
        self.full_description = full_description
        # TODO: How to handle other related properties

        if self.number is None:
            # TODO: get latest number and ndigits from conf / db
            self.number = self.create_number(last_number, ndigits)

        if self.name is None:
            self.name = ""

        if self.description is None:
            self.description = ""

        if self.full_description is None:
            self.full_description = ""

        proj = Project()
        self.project_id = proj.get_id(self.project)
        # self.project_id = Select()  # get project id based on project number / project name
        new_item = PdmItem(item_number=self.number, item_name=self.name, item_description=self.description, item_full_description=self.full_description, path=self.path, project_id=self.project_id)

        # TODO: Import Engine - From where?
        Session.configure(bind=engine, future=True)

        # https://docs.sqlalchemy.org/en/14/orm/session_basics.html#id1
        with Session() as session:
            try:
                session.add(new_item)
            except:
                Session.rollback()
            finally:
                Session.close()

        # raise NotImplementedError("Function create_item is not implemented yet")


    def remove_item(self):
        """Remove existing item"""
        # check if item is new (local == no state)
        # -> if True user can remove
        # -> if False Check if user == admin
        #   -> if True User can remove item
        #   -> if False warning message
        raise NotImplementedError("Function remove_item is not implemented yet")

    def update_item(self):
        """Update existing item"""
        raise NotImplementedError("Function update_item is not implemented yet")

    def add_item_image(self):
        """Update existing item"""
        # TODO: Auto generate image from models
        raise NotImplementedError("Function add_item_image is not implemented yet")


# When inheritance not everything. Do i need a base class?
class Model():
    """Model related Class"""

    def __init__(self):
        print("Generic model")


    def create_model(self):
        """Create new model"""
        # create copy with iter: 0
        # Create model for:
        # -> Existing item
        # -> For new Item
        #    -> With new item also create item
        raise NotImplementedError("Function create_model is not implemented yet")
        
    def remove_model(self):
        """Remove existing model"""
        # check if model is new (local == no state)
        # -> if True user can remove
        # -> if False Check if user == admin
        #   -> if True User can remove item
        #   -> if False warning message
        raise NotImplementedError("Function remove_model is not implemented yet")

    def update_model(self):
        """Update existing model"""
        # create copy with iter: N
        raise NotImplementedError("Function update_model is not implemented yet")

    def get_version(self, save_iter):
        """Get model that is not latest version"""
        # in UI set only release versions or all
        # Optional two versions to compare
        # is FC able to reload when model is changed?
        raise NotImplementedError("Function get_version is not implemented yet")


# When inheritance not everything. Do i need a base class?
class Document():
    """Document related Class"""
    # How much difference is there between a document and a model?

    def __init__(self):
        print("Generic Document")


# Probably a new file
class OwnerStates():  # Acces States
    """ Item / Model / Document Ownership states Class"""
    # Can all checkin options performed from a central class?

    def __init__(self):
        print("Generic OwnerStates")

    def check_in(self, objects):
        """Check in Items, Models, Documents"""
        # check if new item?
        # create copy (for Model, Document) and add copy to DataBase
        raise NotImplementedError("Function check_in is not implemented yet")

    def check_out(self, objects):
        """Check in Items, Models, Documents"""
        # check latest version (only for Models and Documents)
        # check if checked-out by other user
        raise NotImplementedError("Function check_out is not implemented yet")

    def check_in_check_out(self, objects):
        """Check in Items, Models, Documents"""
        # check if new item?
        # create copy (for Model, Document) and add copy to DataBase
        #
        # checkout checks are not needed
        raise NotImplementedError("Function check_in_check_out is not implemented yet")


class ReleaseStates():
    """Item / Model /Document release states Class"""

    def __init__(self):
        print("Generic Release States from ")

    def chnge_release_state(self):
        """new Item, Model, Document"""
        # All new items, models, documents have state new - untill they are checked in.
        raise NotImplementedError("Function new is not implemented yet")

    def new(self, objects):
        """new Item, Model, Document"""
        # All new items, models, documents have state new - untill they are checked in.
        raise NotImplementedError("Function new is not implemented yet")

    def prototype(self, objects):
        """prototype Item, Model, Document"""
        # All items, models, documents get state prototype on first checkin - untill they are released.
        raise NotImplementedError("Function prototype is not implemented yet")

    def release(self, objects):
        """Check in Items, Models, Documents"""
        # All items, models, documents get state release after  - untill they are released.
        raise NotImplementedError("Function check_in_check_out is not implemented yet")

    def not_for_new(self, objects):
        """Check in Items, Models, Documents"""
        # check latest version (only for Models and Documents)
        # check if checked-out by other user
        raise NotImplementedError("Function check_out is not implemented yet")

    def depreciated(self, objects):
        """Check in Items, Models, Documents"""
        # check latest version (only for Models and Documents)
        # check if checked-out by other user
        raise NotImplementedError("Function check_out is not implemented yet")
