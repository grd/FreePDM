#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from sql import Base
from sqlalchemy import create_engine
from sqlalchemy import MetaData
from sqlalchemy import Table
from sqlalchemy.orm import sessionmaker
from sqlalchemy.engine import URL
from typing import Optional, Union
# Table classes
from sql_tables import SQLHistory
# from sql_tables import SQLUser, SQLRole, SQLProject, SQLItem, SQLModel, SQLDocument, SQLMaterial, SQLHistory, SQLPurchase, SQLManufacturer, SQLVendor



# Base = declarative_base()
# https://dataedo.com/kb/data-glossary/what-is-metadata
# https://www.geeksforgeeks.org/difference-between-data-and-metadata/
# https://www.geeksforgeeks.org/metadata-in-dbms-and-its-types/
metadata_obj = MetaData()

class CreateDb():
    """Create database"""

    def __init__(self):
        pass

    def make_url(self, drivername: str, username: Optional[str], password:  Optional[str], host: Optional[str], port: Optional[int], database_name: Optional[str]):
        """
        Create new url

        Parameters
        ----------

        drivername [str] :
            drivername. For example 'postgresql+psycopg2'

        username [str] :
            username

        password [str] :
            password.

        host [str] :
            host address. For example localhost

        port [int] :
            SQL port.

        database_name [str] :
            name of the database

        Returns
        -------

        url [str]
        """
        self.drivername = drivername
        self.username = username
        self.password = password
        self.host = host
        self.port = port
        self.database_name = database_name
        new_url = URL.create(self.drivername, self.username, self.password, self.host, self.port, self.database_name)
        return(new_url)


class CreateMySQLDb(CreateDb):  # Everything in a file or better to split it?
    """Feed Forward of generic SQL functions to MySQL"""
    # https://docs.sqlalchemy.org/en/14/core/engines.html#mysql

    def __init__(self):
        print("MySQL")
        super(CreateMySQLDb, self).__init__()

    def start_engine(self, url: Union[str, URL], encoding: Optional[str], echo: Optional[bool], future: Optional[bool], dialect: Optional[str]):
        """
        Start MySQL engine.
        Note: MySQL engine is not default development database.

        Parameters
        ----------

        url [str] : url
            Path to MySQL database.

        encoding [str] : Set text encoding for database.
            The default encoding is `utf-8`

        echo [bool] : Set logging On / Off
            If True Engine will log all statements.

        future [bool] : Future proof style
            SQLAlchemy 2.0 up style Engine, Connection (Introduced in SQLAlchemy 1.4).

        dialect [string] : If other SQL python libraries are used this can be set.
            Optional parameter.
        """
        self.url = url
        self.encoding = encoding
        self.echo = echo
        self.future = future
        self.dialect = dialect
        # # https://docs.sqlalchemy.org/en/14/core/engines.html#database-urls
        if (self.dialect == "default") or (self.dialect is None):
            # Installing via `FreePDM-ServerInstaller.sh` installs default engine
            # default
            self.engine = create_engine(self.url, echo=self.echo, future=self.future)
            return(self.engine)
        elif (self.dialect == "mysqlclient") or (self.dialect == "mysqldb"):
            # mysqlclient (a maintained fork of MySQL-Python)
            self.engine = create_engine(self.url, echo=self.echo, future=self.future)
            return(self.engine)
        elif (self.dialect == "PyMySQL") or (self.dialect == "pymysql"):
            # PyMySQL
            self.engine = create_engine(self.url, echo=self.echo, future=self.future)
            return(self.engine)
        else:
            pass


class CreatePostgreSQLDb(CreateDb):
    """Feed Forward of generic SQL functions to PostgeSQL"""
    # https://docs.sqlalchemy.org/en/14/core/engines.html#postgresql

    def __init__(self):
        print("PostgreSQL")
        super(CreatePostgreSQLDb, self).__init__()

    def start_engine(self, url: Union[str, URL], echo: Optional[bool], encoding: Optional[str], future: Optional[bool], dialect: Optional[str]):
        """
        Start PostgreSQL engine.

        Parameters
        ----------

        url [str] : url
            Path to PostgreSQL database.

        encoding [str] : Set text encoding for database.
            The default encoding is `utf-8`

        echo [bool] : Set logging On / Off
            If True Engine will log all statements.

        future [bool] : Future proof style
            SQLAlchemy 2.0 up style Engine, Connection (Introduced in SQLAlchemy 1.4).

        dialect [string] : If other SQL python libraries are used this can be set.
            Optional parameter.
        """
        self.url = url
        self.encoding = encoding
        self.echo = echo
        self.future = future
        self.dialect = dialect
        # https://docs.sqlalchemy.org/en/14/core/engines.html#database-urls
        if (self.dialect == "default") or (self.dialect is None):
            # default
            self.engine = create_engine(self.url, echo=self.echo, encoding=self.encoding, future=self.future)
            return(self.engine)
        elif self.dialect == "psycopg2":
            # psycopg2
            self.engine = create_engine(self.url, echo=self.echo, encoding=self.encoding, future=self.future)
            return(self.engine)
        elif self.dialect == "pg8000":
            # pg8000
            self.engine = create_engine(self.url, echo=self.echo, encoding=self.encoding, future=self.future)
            return(self.engine)
        else:
            pass


class CreateSQLiteDb(CreateDb):  # Everything in a file or better to split it?
    """Feed Forward of generic SQL functions to SQLite"""
    # https://docs.sqlalchemy.org/en/14/core/engines.html#sqlite

    def __init__(self):
        super(CreateSQLiteDb, self).__init__()
        print("SQLite")

    def start_engine(self, url: Union[str, URL], encoding: Optional[str], echo: Optional[bool], future: Optional[bool]):
        """
        Start SQLite engine.
        Note: SQLite engine is not default development database.

        Parameters
        ----------

        url [str] : url
            Path to SQLite database.

        encoding [str] : Set text encoding for database.
            The default encoding is `utf-8`

        echo [bool] : Set logging On / Off
            If True Engine will log all statements.

        future [bool] : Future proof style
            SQLAlchemy 2.0 up style Engine, Connection (Introduced in SQLAlchemy 1.4).
        """
        self.url = url
        self.encoding = encoding
        self.echo = echo
        self.future = future
        # https://docs.sqlalchemy.org/en/14/core/engines.html#database-urls
        # exampleurl: "sqlite+pysqlite:///:memory:"
        self.engine = create_engine(self.url, echo=self.echo, encoding=self.encoding,future=self.future)  # start from memory
        return(self.engine)


def start_your_engine(url_string: str, db_type: Optional[str], split: Optional[str] = ',', **vargs):
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
        msql = CreateMySQLDb()
        if len(url_list) == 1:
            print("Complete url received.")
            new_url = url_list[0]
            dialect = None
        elif len(url_list) == 6:
            print("Url shall be created")
            new_url = msql_engine.make_url(url_list[0], url_list[1], url_list[2], url_list[3], url_list[4], int(url_list[5]))
            dialect = None
            # if dialect is part of the url...
        elif len(url_list) == 7:
            # list including dialect
            url_dialect = url_list[0] + '+' + url_list[6]
            new_url = msql_engine.make_url(url_dialect, url_list[1], url_list[2], url_list[3], int(url_list[4]), url_list[5])
            dialect = url_list[6]
        else:
            raise ValueError("{} is not the right amount of values for the url. [1, 6 or 7]\n".format(len(url_list)))

        msql_engine = msql.start_engine(new_url, dialect=dialect, encoding='utf-8', echo=False, future=True)
        return(msql_engine)
    elif (db_type is None) or (db_type == "postgresql") or (db_type == "PostgresSQL"):
        psql = CreatePostgreSQLDb()
        if len(url_list) == 1:
            print("Complete url received.")
            new_url = url_list[0]
            dialect = None
        elif len(url_list) == 6:
            print("Url shall be created")
            new_url = psql_engine.make_url(url_list[0], url_list[1], url_list[2], url_list[3], url_list[4], int(url_list[5]))
            dialect = None
            # if dialect is part of the url...
        elif len(url_list) == 7:
            # list including dialect
            url_dialect = url_list[0] + '+' + url_list[6]
            new_url = psql_engine.make_url(url_dialect, url_list[1], url_list[2], url_list[3], int(url_list[4]), url_list[5])
            dialect = url_list[6]
        else:
            raise ValueError("{} is not the right amount of values for the url. [1, 6 or 7]\n".format(len(url_list)))

        psql_engine = psql.start_engine(new_url, dialect=dialect, encoding='utf-8', echo=False, future=True)
        return(psql_engine)  # Not sure if returning this is required
    elif (db_type == "sqlite") or (db_type == "SQLite"):
        sqli = CreateSQLiteDb()
        if len(url_list) == 1:
            print("Complete url received.")
            new_url = url_list[0]
        elif len(url_list) == 6:
            print("Url shall be created")
            new_url = sqli_engine.make_url(url_list[0], url_list[1], url_list[2], url_list[3], url_list[4], int(url_list[5]))
            # dialect = None  # not needed? remove later
            # if dialect is part of the url...
        elif len(url_list) == 7:
            # list including dialect
            url_dialect = url_list[0] + '+' + url_list[6]
            new_url = sqli_engine.make_url(url_dialect, url_list[1], url_list[2], url_list[3], int(url_list[4]), url_list[5])
            # dialect = url_list[6]  # not needed ? remove later
        else:
            raise ValueError("{} is not the right amount of values for the url. [1, 6 or 7]\n".format(len(url_list)))

        sqli_engine = sqli.start_engine(new_url, encoding='utf-8', echo=False, future=True)
        return(sqli_engine)
    else:
        raise ValueError("{} Is not a Valid input for 'db_type'.".format(db_type))


def create_default_tables(engine: str):
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
    session = sessionmaker()
    session.configure(bind=engine)
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
        password = "PsqlPassword123!"  # remove password. this one only for developement purposes in VM!
        dbname = "freepdm"

        url = "postgresql://" + username + ":" + password +"@localhost/" + dbname

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
