#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from sqlalchemy.orm import declarative_base
from sqlalchemy import Column, Float, Integer, String, Date, ForeignKey
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
    role_name = Column(String(32))
    # TODO: add privileges - Also how to


class SQLProject(Base):
    """Class with default SQL Role properties"""
    __tablename__ = 'projects'

    project_id = Column(Integer, primary_key=True)
    project_number = Column(Integer)  # this can come another source as the Db so not same as project_id
    project_name = Column(String(32))
    project_path = Column(String)


class SQLItem(Base):
    """Class with SQL Item properties"""
    __tablename__ = 'items'

    item_id = Column(Integer, primary_key=True)
    item_number = Column(Integer)
    item_name = Column(String(32))
    item_description = Column(String(32))
    item_full_description = Column(String)
    # item should be able to exist in multiple projects. but need a singel store location
    item_in_project = Column(Integer, ForeignKey('projects.project_number'), nullable=False)
    item_path = Column(String)


class SQLModel(Base):
    """Class with SQL Item properties"""
    __tablename__ = 'models'

    model_id = Column(Integer, primary_key=True)
    model_number = Column(Integer, ForeignKey('items.item_number'))
    model_name = Column(String(32))
    model_description = Column(String(32))
    model_full_description = Column(String)
    model_filename = Column(String(255))  # Limit of current file systems
    # Auto calculate extension
    model_ext = Column(String(255))  # Total limit of filename and extension is 255
    # model_path = Column(String)  # should belongs to same path as described in item


class SQLDocument(Base):
    """Class with SQL Item properties"""
    __tablename__ = 'documents'

    document_id = Column(Integer, primary_key=True)
    document_number = Column(Integer, ForeignKey('items.item_number'), nullable=False)
    document_name = Column(String(32))
    document_description = Column(String(32))
    document_full_description = Column(String)
    document_filename = Column(String(255))  # Limit of current file systems
    # Auto calculate extension
    document_ext = Column(String(255))  # Total limit of filename and extension is 255
    # document_path = Column(String)  # should belongs to same path as described in item


class SQLModelMaterial(Base):
    """Class with SQL Material properties"""
    __tablename__ = 'materials'

    modelmat_id = Column(Integer, primary_key=True)
    model_material = Column(String(32))
    model_finish = Column(String(32))
    model_mass = Column(Float)  # unit mass
    model_volume = Column(Float)
    model_weight = Column(Float)  # model_mass * model_volume
    model_surface_area = Column(Float)
    model_number = Column(Integer, ForeignKey('models.model_number'), nullable=False)


# TODO: A histric item / Model / Document can  have a foreign key to all of those tables.
class SQLHistory(Base):
    """Class with SQL History properties"""
    __tablename__ = 'history'

    hisory_id = Column(Integer, primary_key=True)
    history_date_created = Column(Date)
    history_created_by = Column(String)
    history_date_last_edit = Column(Date)
    history_last_edit_by = Column(String)
    history_checked_out_by = Column(String)
    history_revision_state = Column(String)
    history_revision_number = Column(Integer)  # Maybe other format
    # Create Complex revisions (Example: Date, major.minor, major.letter_minor)
    history_store_number = Column(Integer)  # last store version iterator
