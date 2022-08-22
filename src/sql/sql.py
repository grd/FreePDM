#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from sqlalchemy.orm import declarative_base
# from typing import Optional

Base = declarative_base()


class SQL(Base):
    """Generic SQL"""
    # https://docs.python.org/3/library/typing.html#typing.NamedTuple

    def __init__(self):
        # self.Base = declarative_base()
        print("Generic SQL")

    def connect_db(self):
        """Connect to database"""
        # https://docs.sqlalchemy.org/en/14/tutorial/dbapi_transactions.html
        print("connect to database")
