#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from sqlalchemy.orm import declarative_base
from sqlalchemy import Column, ForeignKey
from sqlalchemy import Boolean, Integer, Float, String, Date, Enum, LargeBinary
from typing import NewType, Optional

Base = declarative_base()


class SQLUser(Base):
    """Class with default SQL User properties"""
    __tablename__ = 'user_account'

    user_id = Column(Integer, primary_key=True)
    user_name = Column(String(30))
    user_first_name = Column(String(30))
    user_last_name = Column(String(30))
    user_full_name = Column(String)
    user_email_adress = Column(String, nullable=False)  # TODO: change to mail address
    user_phonenumber = Column(Integer)  # TODO: add phonenumber property
    # user_role = []  # TODO: add List of roles
    user_department = Column(String)  # Enum optionally
    # user_aliases = []  # TODO: What to do with aliases


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
    project_status = Column(Enum)  # TODO: Create Enum list
    Project_date_start = Column(Date)
    Project_date_finish = Column(Date)  # TODO check finsh date is after start date
    project_path = Column(String)


class SQLItem(Base):
    """Class with SQL Item properties"""
    __tablename__ = 'items'

    item_id = Column(Integer, primary_key=True)
    item_number = Column(Integer)
    item_name = Column(String(32))
    item_description = Column(String(32))
    item_full_description = Column(String)
    item_number_linked_files = Column(Integer)
    # item should be able to exist in multiple projects. but need a singel store location
    item_in_project = Column(Integer, ForeignKey('projects.project_number'), nullable=False)
    item_path = Column(String)
    # TODO: get image from Model. If there is no fileimage add default empty image.
    item_preview = Column(LargeBinary)  # Change when no image is available
    purchasing_id = Column(Integer, ForeignKey('purchasing.purchasing_id'))


class SQLModel(Base):
    """Class with SQL Item properties"""
    __tablename__ = 'models'

    model_id = Column(Integer, primary_key=True)
    model_number = Column(Integer, ForeignKey('items.item_number'))
    model_name = Column(String(32))
    model_description = Column(String(32))
    model_full_description = Column(String)
    model_filename = Column(String(253), nullable=False)  # Limit of current file systems
    # Auto calculate extension
    model_ext = Column(String(253), nullable=False)  # Total limit of filename and extension is 255
    # model_path = Column(String)  # should belongs to same path as described in item
    model_preview = Column(LargeBinary)  # Change when no image is available


class SQLDocument(Base):
    """Class with SQL Item properties"""
    __tablename__ = 'documents'

    document_id = Column(Integer, primary_key=True)
    document_number = Column(Integer, ForeignKey('items.item_number'), nullable=False)
    document_name = Column(String(32))
    document_description = Column(String(32))
    document_full_description = Column(String)
    document_filename = Column(String(253), nullable=False)  # Limit of current file systems
    # Auto calculate extension
    document_ext = Column(String(253), nullable=False)  # Total limit of filename and extension is 255
    # document_path = Column(String)  # should belongs to same path as described in item


class SQLMaterial(Base):
    """Class with SQL Material properties"""
    __tablename__ = 'materials'

    material_id = Column(Integer, primary_key=True)
    material_name= Column(String(32))
    material_finish = Column(String(32))
    material_mass = Column(Float)
    material_mass_unit = Column(Enum)  # TODO: Create Enum list
    material_volume = Column(Float)  # TODO: From CAD file
    material_volume_unit = Column(Enum)  # TODO: Create Enum list
    material_weight = Column(Float)  # material_mass * material_volume
    material_weight_unit = Column(Enum)  # TODO: Create Enum list
    material_surface_area = Column(Float)
    material_surface_area_unit = Column(Enum)  # TODO: Create Enum list
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
    history_revision_state = Column(Enum)  # TODO: Create Enum list
    history_revision_number = Column(Integer)  # Maybe other format
    # TODO: Create Complex revisions (Example: Date, major.minor, major.letter_minor)
    history_stored_number = Column(Integer)  # last store version iterator
    # TODO: Every stored version in the database should be traceble


class SQLPurchase(Base):
    """Class with SQL purchasing properties"""
    __tablename__ = 'purchasing'

    purchasing_id = Column(Integer, primary_key=True)
    purchasing_source = Column(Boolean)  # Repesent Buy OR Manufacture
    # Optionally Enum
    purchasing_tracebility = Column(Enum)  # TODO: Create Enum list
    # Represent list: Lot, Lot And Serial Number, Serial Number, Not traced
    manufacturer_id = Column(Integer, ForeignKey('manufacturers.manufacturer_id'))
    vendor_id = Column(Integer, ForeignKey('vendors.vendor_id'))


class SQLManufacturer(Base):
    """Class with SQL Manufacturing properties"""
    __tablename__ = 'manufacturers'

    manufacturer_id = Column(Integer, primary_key=True)
    manufacturer_name = Column(String(32))
    # TODO: Add manufacturer address


class SQLVendor(Base):
    """Class with SQL Vendor properties"""
    __tablename__ = 'vendors'

    vendor_id = Column(Integer, primary_key=True)
    vendor_name = Column(String(32))
    # TODO: Add manufacturer address
