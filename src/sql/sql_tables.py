#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

# https://stackoverflow.com/questions/54118182/sqlalchemy-not-creating-tables
from sql import Base
from sqlalchemy import Column, ForeignKey
from sqlalchemy import Boolean, Integer, Float, String, Date, Enum, LargeBinary
import sql_enum

# https://dataedo.com/kb/data-glossary/what-is-metadata
# https://www.geeksforgeeks.org/difference-between-data-and-metadata/
# https://www.geeksforgeeks.org/metadata-in-dbms-and-its-types/


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

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return f"SQLUser(user_id={self.user_id!r}, user_name={self.user_name!r}, user_first_name={self.user_first_name!r}, user_last_name={self.user_last_name!r}, user_full_name={self.user_full_name!r}, user_email_adress={self.user_email_adress!r}, user_phonenumber={self.user_phonenumber!r}, user_department={self.user_department!r})"


class SQLRole(Base):
    """Class with default SQL Role properties"""
    __tablename__ = 'roles'

    role_id = Column(Integer, primary_key=True)
    role_name = Column(String(32))
    # TODO: add privileges - Also how to

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"SQLRole(role_id={self.role_id!r}, role_name={self.role_name!r})")


class SQLProject(Base):
    """Class with default SQL Role properties"""
    __tablename__ = 'projects'

    project_id = Column(Integer, primary_key=True)
    project_number = Column(Integer)  # this can come another source as the Db so not same as project_id
    project_name = Column(String(32))
    project_status = Column(Enum(sql_enum.ProjectState))
    Project_date_start = Column(Date)
    Project_date_finish = Column(Date)  # TODO check finish date is after start date
    project_path = Column(String)

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"SQLProject(project_id={self.project_id!r}, project_number={self.project_number!r}, project_name={self.project_name!r}, project_status={self.project_status!r}, project_date_start={self.project_date_start!r}, project_date_finish={self.project_date_finish!r}, project_path={self.project_path!r})")


class SQLItem(Base):
    """Class with SQL Item properties"""
    __tablename__ = 'items'

    item_id = Column(Integer, primary_key=True)
    item_number = Column(Integer)
    item_name = Column(String(32))
    item_description = Column(String(32))
    item_full_description = Column(String)
    item_number_linked_files = Column(Integer)
    item_path = Column(String)
    # TODO: get image from Model. If there is no fileimage add default empty image.
    item_preview = Column(LargeBinary)  # Change when no image is available
    # item should be able to exist in multiple projects. but need a single store location
    project_id = Column(Integer, ForeignKey('projects.project_number'), nullable=False)
    purchasing_id = Column(Integer, ForeignKey('purchasing.purchasing_id'))

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        # ignore cross columns
        return(f"SQLItem(item_id={self.item_id!r}, tem_number={self.item_number!r}, item_name={self.item_name!r}, item_description={self.item_description!r}, item_full_description={self.item_full_description!r}, item_number_linked_files={self.item_number_linked_files!r}, item_path={self.item_path!r}, item_preview={self.item_preview!r})")


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

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"SQLModel(model_id={self.model_id!r}, model_number={self.model_number!r}, model_name={self.model_name!r}, model_description={self.model_description!r}, model_full_description={self.model_full_description!r}, model_filename={self.model_filename!r}, model_ext={self.model_ext!r}, model_preview={self.model_preview!r})")


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

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"SQLDocument(document_id={self.document_id!r}, document_number={self.document_number!r}, document_name={self.document_name!r}, document_description={self.document_description!r}, document_full_description={self.document_full_description!r}, document_filename={self.document_filename!r}, document_ext={self.document_ext!r})")


class SQLMaterial(Base):
    """Class with SQL Material properties"""
    __tablename__ = 'materials'

    material_id = Column(Integer, primary_key=True)
    material_name = Column(String(32))
    material_finish = Column(String(32))
    material_density = Column(Float)
    material_density_unit = Column(Enum(sql_enum.DensityUnit))
    material_volume = Column(Float)  # TODO: From CAD file
    material_volume_unit = Column(Enum(sql_enum.VolumeUnit))
    material_weight = Column(Float)  # material_density * material_volume
    material_weight_unit = Column(Enum(sql_enum.WeightUnit))
    material_surface_area = Column(Float)
    material_surface_area_unit = Column(Enum(sql_enum.AreaUnit))
    model_number = Column(Integer, ForeignKey('models.model_number'), nullable=False)

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"SQLMaterial(material_id={self.material_id!r}, material_name={self.material_name!r}, material_finish={self.material_finsh!r}, material_mass={self.material_mass!r}, material_mass_unit={self.material_mass_unit!r}, material_volume={self.material_volume!r}, material_volume_unit={self.material_volume_unit!r}, material_weight={self.material_weight!r}, material_weight_unit={self.material_weight_unit!r}, material_surface_area={self.material_surface_area!r}, material_surface_area_unit={self.material_surface_area_unit})")


# TODO: A histric item / Model / Document can  have a foreign key to all of those tables.
class SQLHistory(Base):
    """Class with SQL History properties"""
    __tablename__ = 'history'

    history_id = Column(Integer, primary_key=True)
    history_date_created = Column(Date)
    history_created_by = Column(String)
    history_date_last_edit = Column(Date)
    history_last_edit_by = Column(String)
    history_checked_out_by = Column(String)
    history_revision_state = Column(Enum(sql_enum.RevisionState))
    history_revision_number = Column(Integer)  # Maybe other format
    # TODO: Create Complex revisions (Example: Date, major.minor, major.letter_minor)
    history_stored_number = Column(Integer)  # last store version iterator
    # TODO: Every stored version in the database should be traceble

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"SQLHistory( history_id={self.history_id!r}, history_date_created={self.history_date_created!r}, history_created_by={self.history_created_by!r}, history_date_last_edit={self.history_date_last_edit!r}, history_last_edit_by={self.history_last_edit_by!r}, history_checked_out_by={self.history_checked_out_by!r}, history_revision_state={self.history_revision_state!r}, history_revision_number={self.history_revision_number!r}, history_stored_number={self.history_stored_number!r})")


class SQLPurchase(Base):
    """Class with SQL purchasing properties"""
    __tablename__ = 'purchasing'

    purchasing_id = Column(Integer, primary_key=True)
    purchasing_source = Column(Boolean)  # Represent Buy OR Manufacture
    # Optionally Enum
    purchasing_tracebility = Column(Enum(sql_enum.TracebilityState))
    # Represent list: Lot, Lot And Serial Number, Serial Number, Not traced
    manufacturer_id = Column(Integer, ForeignKey('manufacturers.manufacturer_id'))
    vendor_id = Column(Integer, ForeignKey('vendors.vendor_id'))

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"SQLPurchase(purchasing_id={self.purchasing_id!r}, purchasing_source={self.purchasing_source!r}, purchasing_tracebility={self.purchasing_tracebility!r})")


class SQLManufacturer(Base):
    """Class with SQL Manufacturing properties"""
    __tablename__ = 'manufacturers'

    manufacturer_id = Column(Integer, primary_key=True)
    manufacturer_name = Column(String(32))
    # TODO: Add manufacturer address

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"SQLManufacturer(manufacturer_id={self.manufacturer_id!r}, manufacturer_name={self.manufacturer_name!r})")


class SQLVendor(Base):
    """Class with SQL Vendor properties"""
    __tablename__ = 'vendors'

    vendor_id = Column(Integer, primary_key=True)
    vendor_name = Column(String(32))
    # TODO: Add manufacturer address

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"SQLVendor(vendor_id={self.vendor_id!r}, vendor_name={self.vendor_name!r})")
