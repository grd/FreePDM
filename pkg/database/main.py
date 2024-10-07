#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from base import Base
from base import Session
from base import metadata_obj
from general import MySQLDb, PostgreSQLDb, SQLiteDb
from sqlalchemy.engine import Engine
# from sqlalchemy import MetaData
from sqlalchemy import Table
from typing import Optional
# Table classes - Test if it still works if all tables are commented out
# from pdm_tables import PdmUser
# from pdm_tables import PdmRole
# from pdm_tables import PdmProject
# from pdm_tables import PdmItem
# from pdm_tables import PdmModel
# from pdm_tables import PdmDocument
# from pdm_tables import PdmMaterial
# from pdm_tables import PdmHistory
# from pdm_tables import PdmPurchase
# from pdm_tables import PdmManufacturer
# from pdm_tables import PdmVendor


# https://dataedo.com/kb/data-glossary/what-is-metadata
# https://www.geeksforgeeks.org/difference-between-data-and-metadata/
# https://www.geeksforgeeks.org/metadata-in-dbms-and-its-types/
metadata_obj = MetaData()


# https://stackoverflow.com/questions/73887390/handle-multiple-users-login-database-with-sqlalchemy
def start_your_engine(url_string: str, db_type: Optional[str], split: Optional[str] = ',', **vargs) -> Engine:
    """
    Start your chosen engine

    Parameters
    ----------

    prefered_engine [str] :
        What sql engine is preferred to use.

    url [str]:
        url for database. This can be a single string but also a list of parameters (comma-separated).
        This list contain the following information: drivername; username; password; host; port; database_name.
        TODO: add dialect as optional value

    split [str]:
        split string in list (Optional). Default value is ','.

    **vargs:
        Other parameters are parsed through
    """
    url_list = url_string.split(split)

    if (db_type == "mysql") or (db_type == "MySQL"):
        msql = MySQLDb()

        if len(url_list) == 1:
            print("Complete url received.")
            new_url = url_list[0]
            dialect = None
        elif len(url_list) == 6:
            print("Url shall be created")
            new_url = msql.make_url(url_list[0], url_list[1], url_list[2], url_list[3], url_list[4], int(url_list[5]))
            dialect = None
            # if dialect is part of the url...
        elif len(url_list) == 7:
            # list including dialect
            url_dialect = url_list[0] + '+' + url_list[6]
            new_url = msql.make_url(url_dialect, url_list[1], url_list[2], url_list[3], int(url_list[4]), url_list[5])
            dialect = url_list[6]
        else:
            raise ValueError("{} is not the right amount of values for the url. [1, 6 or 7]\n".format(len(url_list)))

        msql_engine = msql.start_engine(new_url, dialect=dialect, encoding='utf-8', echo=False, future=True)
        return(msql_engine)

    elif (db_type is None) or (db_type == "postgresql") or (db_type == "PostgresSQL"):
        psql = PostgreSQLDb()

        if len(url_list) == 1:
            print("Complete url received.")
            new_url = url_list[0]
            dialect = None
        elif len(url_list) == 6:
            print("Url shall be created")
            new_url = psql.make_url(url_list[0], url_list[1], url_list[2], url_list[3], url_list[4], int(url_list[5]))
            dialect = None
            # if dialect is part of the url...
        elif len(url_list) == 7:
            # list including dialect
            url_dialect = url_list[0] + '+' + url_list[6]
            new_url = psql.make_url(url_dialect, url_list[1], url_list[2], url_list[3], int(url_list[4]), url_list[5])
            dialect = url_list[6]
        else:
            raise ValueError("{} is not the right amount of values for the url. [1, 6 or 7]\n".format(len(url_list)))

        psql_engine = psql.start_engine(new_url, dialect=dialect, encoding='utf-8', echo=False, future=True)
        return(psql_engine)  # Not sure if returning this is required

    elif (db_type == "sqlite") or (db_type == "SQLite"):
        sqli = SQLiteDb()

        if len(url_list) == 1:
            print("Complete url received.")
            new_url = url_list[0]
        elif len(url_list) == 6:
            print("Url shall be created")
            new_url = sqli.make_url(url_list[0], url_list[1], url_list[2], url_list[3], url_list[4], int(url_list[5]))
            # dialect = None  # not needed? remove later
            # if dialect is part of the url...
        elif len(url_list) == 7:
            # list including dialect
            url_dialect = url_list[0] + '+' + url_list[6]
            new_url = sqli.make_url(url_dialect, url_list[1], url_list[2], url_list[3], int(url_list[4]), url_list[5])
            # dialect = url_list[6]  # not needed ? remove later
        else:
            raise ValueError("{} is not the right amount of values for the url. [1, 6 or 7]\n".format(len(url_list)))

        sqli_engine = sqli.start_engine(new_url, encoding='utf-8', echo=False, future=True)
        return(sqli_engine)

    else:
        raise ValueError("{} Is not a Valid input for 'db_type'.".format(db_type))


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

    # get existing tables?
