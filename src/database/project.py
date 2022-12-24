#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from base import Base
from base import Session
from base import import metadata_obj
from default_tables import PdmProject
from default_enum import ProjectState
from sqlalchemy import Date
from typing import Optional, Union


# Basically searching in the SQL database requires a coonetcion to the database too.
# Is it than a sub-class or better to request access to a class or has this class a interface too?
class Project():
    """Project related Class"""
    # https://docs.sqlalchemy.org/en/14/orm/session_basics.html

    def __init__(self):
        print("Generic Project")

    def get_id(self, number: str) -> int:
        """Get id with Project number"""
        raise NotImplementedError("Function get_id is not implemented yet")

    def create_number(self, number: Union[str, int], ndigits: Optional[int]) -> str:
        """
        Create new Project number

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

        # TODO: load defaults for ndigits
        if self.ndigits is None:
            raise ValueError("Value for 'ndigits' can't be None.")

        self.number = int(self.number)
        self.number += 1

        # https://thecodingbot.com/count-the-number-of-digits-in-an-integer-python4-ways/
        if self.ndigits != -1:
            counter = 0
            num = self.number
            while (num):
                counter += 1
                num = int(num / 10)

            leading_zeros = ""
            for digit in range(self.ndigits - counter):
                leading_zeros += "0"

            return(leading_zeros + str(self.number))
        else:
            return(str(self.number))

    def create_project(self, number: Optional[str], name: Optional[str], status: Optional[ProjectState], date_start: Optional[Date], date_finish: Optional[Date], path: str):
        """
        Create new project

        Parameters
        ----------

        project [str] : Project number or name
            Project name or project number have to be a string for use with leading zeros.

        path [str] : path
            Path to folder where related models and documents are stored.

        name [str] : Optional Name
            Name is optional for when no number is wanted.
            If name is None then automatically a number is generated, otherwise added number is used.

        Returns
        -------

        No return - new object created in SQL database ?
        """
        self.number = number
        self.name = name
        self.status = status
        self.date_start = date_start
        self.date_finish = date_finish
        self.path = path  # TODO: create path automatically
        # TODO: How to handle other related properties

        if self.number is None:
            # TODO: get latest number and ndigits from conf / db
            self.number = self.create_number(last_number, ndigits)

        if self.name is None:
            self.name = ""

        if self.status is None:
            self.status = ProjectState.new

        if self.date_start is None:
            self.date_start = ""  # TODO: replace with date

        if self.date_finish is None:
            self.date_finish = ""  # TODO: replace with date

        if self.path is None:
            self.path = ""

        new_project = PdmProject(self.number, self.name, self.status, self.date_start, self.date_finish, self.path)

        # TODO: Import Engine - From where?
        Session.configure(bind=engine, future=True)

        # https://docs.sqlalchemy.org/en/14/orm/session_basics.html#id1
        with Session() as session:
            try:
                session.add(new_project)
            except:
                Session.rollback()
            finally:
                Session.close()

        # raise NotImplementedError("Function create_project is not implemented yet")

    def remove_project(self):
        """Remove existing project"""
        raise NotImplementedError("Function remove_project is not implemented yet")

    def update_project(self):
        """Update existing project"""
        raise NotImplementedError("Function update_model is not implemented yet")

    def add_user_to_project(self):
        """Add user to project"""
        raise NotImplementedError("Function add_user_to_project is not implemented yet")

    def remove_user_from_project(self):
        """Remove user from project"""
        raise NotImplementedError("Function remove_user_from_project is not implemented yet")
