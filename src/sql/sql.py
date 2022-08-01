#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from sqlalchemy.orm import declarative_base
from sqlalchemy import Column, Integer, String
from typing import NewType, Optional

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


class SQLUser(Base):
    """Class with default SQL User properties"""
    __tablename__ = 'user_account'
    user_id = Column(Integer, primary_key=True)
    user_name = Column(String(30))
    first_name = Column(String(30))
    last_name = Column(String(30))
    full_name = Column(String)
    email_adress = Column(String, nullable=False)
    # phonenumber =   # TODO: add phonenumber property
    # role = []  # TODO: add List of roles


class SQLRole(Base):
    """Class with default SQL Role properties"""
    __tablename__ = 'roles'
    role_id = Column(Integer, primary_key=True)
    role_name = Column(String(30))
    # TODO: add privileges - Also how to


class SQLProject(Base):
    """Class with default SQL Role properties"""
    __tablename__ = 'roles'
    project_id = Column(Integer, primary_key=True)
    Project_number = Column(Integer)
    project_name = Column(String(30))
    project_path = Column(String)
