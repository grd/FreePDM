#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
""" 

import database
import typing


# Basically searching in the SQL database requires a coonetcion to the database too.
# Is it than a sub-class or better to reqest acces to a class or has this class a interface too?
class search():
    """Search class"""
    # https://docs.sqlalchemy.org/en/14/orm/session_basics.html

    def __init__(self):
        print("Search in Databases")

    def search_number(self, number: int):
        """Search on number"""
        pass

    def search_description(self, description):
        """Search on description"""
        pass

    def search_something_else(self, something):
        """Search on number"""
        pass
