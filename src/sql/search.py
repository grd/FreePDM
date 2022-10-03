#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

# from sqlalchemy.orm import sessionmaker
from database import Base
# import typing


# Basically searching in the SQL database requires a coonetcion to the database too.
# Is it than a sub-class or better to reqest acces to a class or has this class a interface too?
class Search():
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
        return(help_text)


class SearchItem():
    """Search for projects"""

    def __init__(self):
        print("Search in items")

    def item_number(self, user_number: str):
        """Search on project number"""
        raise NotImplementedError("Function project_number is not implemented yet")

    def item_description(self, description: str):
        """Search on project description"""
        raise NotImplementedError("Function project_description is not implemented yet")


class SearchProject():
    """Search for projects"""

    def __init__(self):
        print("Search in projects")

    def project_number(self, user_number: str):
        """Search on project number"""
        raise NotImplementedError("Function project_number is not implemented yet")

    def project_description(self, description: str):
        """Search on project description"""
        raise NotImplementedError("Function project_description is not implemented yet")


class SearchUser():
    """Search for projects"""

    def __init__(self):
        print("Search in Users")

    def user_number(self, user_number: str):
        """Search on user number"""
        raise NotImplementedError("Function user_number is not implemented yet")

    def user_name(self, user_name: str):
        """Search on user name"""
        raise NotImplementedError("Function user_name is not implemented yet")

    def user_first_name(self, user_first_name: str):
        """Search on user first name"""
        raise NotImplementedError("Function user_first_name is not implemented yet")

    def user_last_name(self, user_last_name: str):
        """Search on user last name"""
        raise NotImplementedError("Function user_last_name is not implemented yet")

    def user_role(self, user_role: str):
        """Search on user role"""
        raise NotImplementedError("Function user_role is not implemented yet")
