#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

# https://stackoverflow.com/questions/54118182/sqlalchemy-not-creating-tables
from database import Base
from sqlalchemy.orm import relationship
from sqlalchemy import Column, ForeignKey
from sqlalchemy import Boolean, Integer, Float, String, Date, Enum, LargeBinary
import pdm_enum

# https://dataedo.com/kb/data-glossary/what-is-metadata
# https://www.geeksforgeeks.org/difference-between-data-and-metadata/
# https://www.geeksforgeeks.org/metadata-in-dbms-and-its-types/


class PdmUser(Base):
    """Class with default SQL User properties"""
    __tablename__ = 'users'

    user_id = Column(Integer, primary_key=True)
    user_name = Column(String(30))
    user_first_name = Column(String(30))
    user_last_name = Column(String(30))
    user_full_name = Column(String)
    user_email_adress = Column(String, nullable=False)  # TODO: change to mail address
    user_phonenumber = Column(Integer)  # TODO: add phonenumber property
    user_department = Column(String)  # Enum optionally
    # user_aliases = []  # TODO: What to do with aliases

    # relationships with other tables
    roles = relationship("PdmRole", secondary="user_role_link", back_populates="users")
    projects = relationship("PdmProjects", secondary="user_project_link", back_populates="users")
    items = relationship("PdmItem", back_populates="user")
    models = relationship("PdmModel", back_populates="user")
    documents = relationship("PdmDocument", back_populates="user")

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return f"PdmUser(user_id={self.user_id!r}, user_name={self.user_name!r}, user_first_name={self.user_first_name!r}, user_last_name={self.user_last_name!r}, user_full_name={self.user_full_name!r}, user_email_adress={self.user_email_adress!r}, user_phonenumber={self.user_phonenumber!r}, user_department={self.user_department!r})"


class PdmRole(Base):
    """Class with default SQL Role properties"""
    __tablename__ = 'roles'

    role_id = Column(Integer, primary_key=True)
    role_name = Column(String(32))
    # TODO: add privileges - Also how to

    # relationships with other tables
    user_id = Column(Integer, ForeignKey("user_accounts.user_id"))
    users = relationship("PdmUser", secondary="user_role_link", back_populates="roles")

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"PdmRole(role_id={self.role_id!r}, role_name={self.role_name!r})")


class PdmUserRoleLink(Base):
    """Association Table between user and project"""
    # https://www.pythoncentral.io/sqlalchemy-association-tables/
    __tablename__ = "user_role_link"

    user_id = Column("user_id", ForeignKey("user_accounts.user_id"), primary_key=True)
    role_id = Column("role_id", ForeignKey("roles.role_id"), primary_key=True)

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"Pdm User Role Association Table(user_id={self.user_id!r}, role_id={self.role_id!r})")


class PdmProject(Base):
    """Class with default SQL Role properties"""
    __tablename__ = 'projects'

    project_id = Column(Integer, primary_key=True)
    project_number = Column(String(16), nullable=False)  # this can come another source as the Db so not same as project_id
    project_name = Column(String(32))
    project_status = Column(Enum(pdm_enum.ProjectState))
    Project_date_start = Column(Date)
    Project_date_finish = Column(Date)  # TODO check finish date is after start date
    project_path = Column(String)

    # relationships with other tables
    users = relationship("PdmUsers", secondary="user_project_link", back_populates="projects")

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"PdmProject(project_id={self.project_id!r}, project_number={self.project_number!r}, project_name={self.project_name!r}, project_status={self.project_status!r}, project_date_start={self.project_date_start!r}, project_date_finish={self.project_date_finish!r}, project_path={self.project_path!r})")


class PdmUserProjectLink(Base):
    """Association Table between User and Project"""
    # https://www.pythoncentral.io/sqlalchemy-association-tables/
    __tablename__ = "user_project_link"

    user_id = Column("user_id", ForeignKey("user_accounts.user_id"), primary_key=True)
    project_id = Column("project_id", ForeignKey("projects.project_id"), primary_key=True)

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"Pdm User Project Association Table(user_id={self.user_id!r}, project_id={self.project_id!r})")


class PdmItem(Base):
    """Class with SQL Item properties"""
    __tablename__ = 'items'

    item_id = Column(Integer, primary_key=True)
    item_number = Column(String(16))  # Don't expect that numbers longer than 16 digits are needed
    item_name = Column(String(32))
    item_description = Column(String(32))
    item_full_description = Column(String)
    item_number_linked_files = Column(Integer)
    item_path = Column(String, nullable=False)
    # TODO: get image from Model. If there is no fileimage add default empty image.
    item_preview = Column(LargeBinary)  # Change when no image is available
    # item should be able to exist in multiple projects. but need a single store location

    # relationships with other tables
    user_id = Column(Integer, ForeignKey("user_accounts.user_id"))
    user = relationship("PdmUser", back_populates="items")
    project_id = Column(Integer, ForeignKey('projects.project_id'), nullable=False)
    models = relationship("PdmModel", back_populates="item")
    documents = relationship("PdmDocument", back_populates="item")
    Material = relationship("PdmMaterial", back_populates="item", uselist=False)
    purchasing = relationship("PdmPurchasing", back_populates="item", uselist=False)

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        # ignore cross columns
        return(f"PdmItem(item_id={self.item_id!r}, tem_number={self.item_number!r}, item_name={self.item_name!r}, item_description={self.item_description!r}, item_full_description={self.item_full_description!r}, item_number_linked_files={self.item_number_linked_files!r}, item_path={self.item_path!r}, item_preview={self.item_preview!r})")


class PdmProjectItemLink(Base):
    """Association Table between Project and Item"""
    __tablename__ = "Project_item_link"

    project_id = Column("project_id", ForeignKey("projects.project_id"), primary_key=True)
    item_id = Column("item_id", ForeignKey("items.item_id"), primary_key=True)

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"Pdm Project Item Association Table(project_id={self.project_id!r}, item_id={self.item_id!r})")


class PdmModel(Base):
    """Class with SQL Item properties"""
    __tablename__ = 'models'

    model_id = Column(Integer, primary_key=True)
    model_number = Column(Integer)  # or should this be a string?
    model_name = Column(String(32))
    model_description = Column(String(32))
    model_full_description = Column(String)
    model_filename = Column(String(253), nullable=False)  # Limit of current file systems
    # Auto calculate extension
    model_ext = Column(String(253), nullable=False)  # Total limit of filename and extension is 255
    # model_path = Column(String)  # should belongs to same path as described in item
    model_preview = Column(LargeBinary)  # Change when no image is available

    # relationships with other tables
    user_id = Column(Integer, ForeignKey("user_accounts.user_id"))
    user = relationship("PdmUser", back_populates="models")
    item_id = Column(Integer, ForeignKey("items.item_id"))
    item = relationship("PdmItem", back_populates="models")
    material = relationship("PdmMaterial", back_populates="model", uselist=False)

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"PdmModel(model_id={self.model_id!r}, model_number={self.model_number!r}, model_name={self.model_name!r}, model_description={self.model_description!r}, model_full_description={self.model_full_description!r}, model_filename={self.model_filename!r}, model_ext={self.model_ext!r}, model_preview={self.model_preview!r})")


class PdmDocument(Base):
    """Class with SQL Item properties"""
    __tablename__ = 'documents'

    document_id = Column(Integer, primary_key=True)
    document_number = Column(Integer)  # or should this be a string?
    document_name = Column(String(32))
    document_description = Column(String(32))
    document_full_description = Column(String)
    document_filename = Column(String(253), nullable=False)  # Limit of current file systems
    # Auto calculate extension
    document_ext = Column(String(253), nullable=False)  # Total limit of filename and extension is 255
    # document_path = Column(String)  # should belongs to same path as described in item

    # relationships with other tables
    user_id = Column(Integer, ForeignKey("user_accounts.user_id"))
    user = relationship("PdmUser", back_populates="documents")
    item_id = Column(Integer, ForeignKey("items.item_id"))
    item = relationship("PdmItem", back_populates="models")

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"PdmDocument(document_id={self.document_id!r}, document_number={self.document_number!r}, document_name={self.document_name!r}, document_description={self.document_description!r}, document_full_description={self.document_full_description!r}, document_filename={self.document_filename!r}, document_ext={self.document_ext!r})")


class PdmMaterial(Base):
    """Class with SQL Material properties"""
    __tablename__ = 'materials'

    material_id = Column(Integer, primary_key=True)
    material_name = Column(String(32))
    material_finish = Column(String(32))
    material_density = Column(Float)
    material_density_unit = Column(Enum(pdm_enum.DensityUnit))
    material_volume = Column(Float)  # TODO: From CAD file
    material_volume_unit = Column(Enum(pdm_enum.VolumeUnit))
    material_weight = Column(Float)  # material_density * material_volume
    material_weight_unit = Column(Enum(pdm_enum.WeightUnit))
    material_surface_area = Column(Float)
    material_surface_area_unit = Column(Enum(pdm_enum.AreaUnit))

    # relationships with other tables
    model_id = Column(Integer, ForeignKey('models.model_id'))
    model = relationship("PdmModel", back_populates="material")
    item_id = Column(Integer, ForeignKey('items.item_id'))
    item = relationship("PdmItem", back_populates="material")

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"PdmMaterial(material_id={self.material_id!r}, material_name={self.material_name!r}, material_finish={self.material_finsh!r}, material_mass={self.material_mass!r}, material_mass_unit={self.material_mass_unit!r}, material_volume={self.material_volume!r}, material_volume_unit={self.material_volume_unit!r}, material_weight={self.material_weight!r}, material_weight_unit={self.material_weight_unit!r}, material_surface_area={self.material_surface_area!r}, material_surface_area_unit={self.material_surface_area_unit})")


# TODO: A histric item / Model / Document can  have a foreign key to all of those tables.
class PdmHistory(Base):
    """Class with SQL History properties"""
    __tablename__ = 'history'

    history_id = Column(Integer, primary_key=True)
    history_date_created = Column(Date)
    history_created_by = Column(String)
    history_date_last_edit = Column(Date)
    history_last_edit_by = Column(String)
    history_checked_out_by = Column(String)
    # TODO: Create Complex revisions (Example: Date, major.minor, major.letter_minor)
    history_revision_state = Column(Enum(pdm_enum.RevisionState))
    history_revision_number = Column(Integer)  # Maybe other format
    history_stored_number = Column(Integer)  # last store version iterator
    # TODO: Every stored version in the database should be traceble

    # relationships with other tables

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"PdmHistory(history_id={self.history_id!r}, history_date_created={self.history_date_created!r}, history_created_by={self.history_created_by!r}, history_date_last_edit={self.history_date_last_edit!r}, history_last_edit_by={self.history_last_edit_by!r}, history_checked_out_by={self.history_checked_out_by!r}, history_revision_state={self.history_revision_state!r}, history_revision_number={self.history_revision_number!r}, history_stored_number={self.history_stored_number!r})")


class PdmPurchase(Base):
    """Class with SQL purchasing properties"""
    __tablename__ = 'purchasing'

    purchasing_id = Column(Integer, primary_key=True)
    purchasing_source = Column(Boolean)  # Represent Buy OR Manufacture
    # Optionally Enum
    purchasing_tracebility = Column(Enum(pdm_enum.TracebilityState))
    # Represent list: Lot, Lot And Serial Number, Serial Number, Not traced

    # relationships with other tables
    item_id = Column(Integer, ForeignKey("items.item_id"))
    item = relationship("PdmItems", back_populates="purchasing")
    manufacturer_id = Column(Integer, ForeignKey('manufacturers.manufacturer_id'))
    manufacturers = relationship("PdmManufacturer", back_populates="purchasing")
    vendor_id = Column(Integer, ForeignKey('vendors.vendor_id'))
    vendors = relationship("PdmVendors", back_populates="purchasing")

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"PdmPurchase(purchasing_id={self.purchasing_id!r}, purchasing_source={self.purchasing_source!r}, purchasing_tracebility={self.purchasing_tracebility!r})")


class PdmManufacturer(Base):
    """Class with SQL Manufacturing properties"""
    __tablename__ = 'manufacturers'

    manufacturer_id = Column(Integer, primary_key=True)
    manufacturer_name = Column(String(32))
    # TODO: Add manufacturer address

    # relationships with other tables
    purchasing_id = Column(Integer, ForeignKey("purchasing.purchasing_id"))
    purchasing = relationship("PdmPurchase", back_populates="manufacturers")

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"PdmManufacturer(manufacturer_id={self.manufacturer_id!r}, manufacturer_name={self.manufacturer_name!r})")


class PdmVendor(Base):
    """Class with SQL Vendor properties"""
    __tablename__ = 'vendors'

    vendor_id = Column(Integer, primary_key=True)
    vendor_name = Column(String(32))
    # TODO: Add manufacturer address

    # relationships with other tables
    vpurchasing_id = Column(Integer, ForeignKey("purchasing.purchasing_id"))
    vpurchasing = relationship("PdmPurchase", back_populates="vendors")

    def __repr__(self):
        # Fstrings are better ?
        # https://realpython.com/python-f-strings/
        return(f"PdmVendor(vendor_id={self.vendor_id!r}, vendor_name={self.vendor_name!r})")
