#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from base import Base
from base import Session
from base import metadata_obj
from main import start_your_engine
from sqlalchemy.engine import Engine
from sqlalchemy import Table
from typing import Optional
# Table classes - Test if it still works if all tables are commented out
# from default_tables import PdmUser
# from default_tables import PdmRole
# from default_tables import PdmProject
# from default_tables import PdmItem
# from default_tables import PdmModel
# from default_tables import PdmDocument
# from default_tables import PdmMaterial
# from default_tables import PdmHistory
# from default_tables import PdmPurchase
# from default_tables import PdmManufacturer
# from default_tables import PdmVendor


def create_default_tables(engine: Engine):
    """
    Create default set of tables
    Tables are defined in sql_tables.py

    Parameters
    ----------

    Engine [Str]:
        Used Engine for creating tables

    TODO: Add additional classes later on
    """
    # https://stackoverflow.com/questions/54118182/sqlalchemy-not-creating-tables
    Session.configure(bind=engine)
    Base.metadata.create_all(engine)

    user_table = Table("user_accounts", metadata_obj, autoload_with=engine)
    role_table = Table("user_roles", metadata_obj, autoload_with=engine)
    project_table = Table("projects", metadata_obj, autoload_with=engine)
    item_table = Table("items", metadata_obj, autoload_with=engine)
    model_table = Table("models", metadata_obj, autoload_with=engine)
    document_table = Table("documents", metadata_obj, autoload_with=engine)
    material_table = Table("materials", metadata_obj, autoload_with=engine)
    history_table = Table("history", metadata_obj, autoload_with=engine)
    purchase_table = Table("purchasing", metadata_obj, autoload_with=engine)
    manufacturer_table = Table("manufacturers", metadata_obj, autoload_with=engine)
    vendor_table = Table("vendors", metadata_obj, autoload_with=engine)
    return(user_table, role_table, project_table, item_table, model_table, document_table, material_table, history_table, purchase_table, manufacturer_table, vendor_table)


if __name__ == "__main__":
    import sys

    print(len(sys.argv))

    # at least two variables required: filename; url_to_database.
    # if more variables need to be added it is always: filename; database_type; url_to_database; **var
    if len(sys.argv) == 1:
        # raise ValueError("Not enough parameters added")
        # means run from IDE
        username = "freepdm"
        password = "PsqlPassword123!"  # remove password. this one only for development purposes in VM!
        dbname = "freepdm"

        url = "postgresql://" + username + ":" + password + "@localhost/" + dbname

        # sqldb = start_your_engine(sys.argv[1], "postgresql")
        sqldb = start_your_engine(url, "postgresql")
    elif len(sys.argv) == 2:
        # default SQL engine chosen: PostgreSQL
        sqldb = start_your_engine(sys.argv[1], "postgresql")
    elif len(sys.argv) == 3:
        # Choose own SQL Engine
        sqldb = start_your_engine(sys.argv[1], sys.argv[2])
    else:
        # pass all parameters through
        raise NotImplementedError("Parsing argumenets is not implemented yet.")
        # sqldb = start_your_engine(sys.argv[1],sys.argv[2], sys.argv[2:])
        # This could be used for creation of additional tables

    tables = create_default_tables(sqldb)
    for table in tables:
        print(table)
        for key in table.c.keys():
            print(key)
