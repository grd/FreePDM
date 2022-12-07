#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from general import GeneralDb
# from database import Base  # uncemment if needed
import typing


# https://www.osohq.com/post/sqlalchemy-role-rbac-basics
class Role(GeneralDb):
    """Class for generating Roles"""

    def __init__(self):
        print("creating role")

    def add_role(self):
        """Create new role"""
        print("new role created")

    def remove_role(self):
        """Delete existing role"""
        print("existing role deleted")


class User(GeneralDb):
    """Class for generating users
    Users are Aliases for roles in SQL see: https://www.postgresql.org/docs/14/sql-createuser.html
    """

    def __init__(self):
        print("creating users")

    def add_user_to_sql(self, username: str):
        print("This is basically the interface")

    def remove_user_from_sql(self, user_id: int, username: str):
        """Delete existing user"""
        print("existing user deleted")

    def add_user_to_ldap(self, username: str):
        print("This is basically the interface")

    def remove_user_from_ldap(self, user_id: int, username: str):
        """Delete existing user"""
        print("existing user deleted")
