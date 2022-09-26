#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from sqlalchemy.orm import sessionmaker
from database import Base
import typing


# Basically searching in the SQL database requires a coonetcion to the database too.
# Is it than a sub-class or better to reqest acces to a class or has this class a interface too?
class Project():
    """Search class"""
    # https://docs.sqlalchemy.org/en/14/orm/session_basics.html

    def __init__(self):
        print("Search in Databases")

    def search_number(self, number: int):
        """Search on number"""
        raise NotImplementedError("Function search_number is not implemented yet")

    def search_description(self, description: str):
        """Search on description"""
        raise NotImplementedError("Function search_description is not implemented yet")

    def search_something_else(self, something: str):
        """Search on something else"""
        raise NotImplementedError("Function search_something_else is not implemented yet")

    def search_help(self) -> str:
        """Print help function"""
        help_text = """
        ad some text
        - Modifiers
        - Search keys etc
        """
        print(help_text)
