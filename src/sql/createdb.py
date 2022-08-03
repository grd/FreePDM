#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from sql import SQLUser, SQLRole, SQLProject, SQLItem, SQLModel, SQLDocument, SQLModelMaterial, SQLHistory, SQLPurchase, SQLManufacturer, SQLVendor

from sqlalchemy.orm import declarative_base
from sqlalchemy import create_engine
from sqlalchemy.engine import URL
from typing import NewType, Optional

Base = declarative_base()


class CreateDb(Base):
    """Create database"""

    def __init__(self):
        pass

    def make_url(self, drivername: str, username: str, password: str, host: str, port: int, database_name: str) -> str:
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
            host adress. For example localhost

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
        create_url = URL()
        new_url = create_url.create(self.drivername, self.username, self.password, self.host, self.port, self.database_name)
        return(new_url)

    # https://docs.sqlalchemy.org/en/14/tutorial/dbapi_transactions.html#committing-changes
    def create_db(self):
        """create database"""
        pass

    def create_tables(self):
        """create tables for database"""
        pass

    def add_columns_to_table(self):
        """add columns to table"""
        pass


class CreateMySQLDb(CreateDb):  # Alles in een file of beter te splitsen?
    """Feed Forward of generic SQL functions to MySQL"""
    # https://docs.sqlalchemy.org/en/14/core/engines.html#mysql

    def __init__(self):
        print("MySQL")
        super(CreateMySQLDb, self).__init__()

    def start_engine(self, url: str, encoding: str, echo: bool, future: bool, dialect: Optional[str]):
        """
        Start MySQL engine.
        Note: MySQL engine is not default developement database. 

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

        dialect [string] : If other SQL python libraies are used this can be set.
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
            engine = create_engine(self.url, echo=self.echo, future=self.future)
            pass
        elif (self.dialect == "mysqlclient") or (self.dialect == "mysqldb"):
            # mysqlclient (a maintained fork of MySQL-Python)
            engine = create_engine(self.url, echo=self.echo, future=self.future)
            pass
        elif (self.dialect == "PyMySQL") or (self.dialect == "pymysql"):
            # PyMySQL
            engine = create_engine(self.url, echo=self.echo, future=self.future)
        else:
            pass


class CreatePostgreSQLDb(CreateDb):
    """Feed Forward of generic SQL functions to PostgeSQL"""
    # https://docs.sqlalchemy.org/en/14/core/engines.html#postgresql

    def __init__(self):
        print("PostgreSQL")
        super(CreatePostgreSQLDb, self).__init__()

    def start_engine(self, url: str, encoding: str, echo: bool, future: bool, dialect: Optional[str]):
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

        dialect [string] : If other SQL python libraies are used this can be set.
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
            engine = create_engine(self.url, echo=self.echo, future=self.future)
            pass
        elif self.dialect == "psycopg2":
            # psycopg2
            engine = create_engine(self.url, echo=self.echo, future=self.future)
            pass
        elif self.dialect == "pg8000":
            # pg8000
            engine = create_engine(self.url, echo=self.echo, future=self.future)
        else:
            pass


class CreateSQLiteDb(CreateDb):  # Alles in een file of beter te splitsen?
    """Feed Forward of generic SQL functions to SQLite"""
    # https://docs.sqlalchemy.org/en/14/core/engines.html#sqlite

    def __init__(self):
        super(CreateSQLiteDb, self).__init__()
        print("SQLite")

    def start_engine(self, url: str, encoding: str, echo: bool, future: bool):
        """
        Start SQLite engine.
        Note: SQLite engine is not default developement database. 

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
        engine = create_engine(self.url, echo=self.echo, future=self.future)  # start from memory


def start_the_engine():
    """Start your engine"""
    pass


def create_default_tables():
    pass


if __name__ == "__main__":
    import sys

    # at least two variables required: filename; url_to_database.
    # if more variables need to be added it is always: filenname; database_type; url_to_database; **vargs
    if len(sys.argv) == 1:
        raise ValueError("Not enough parameters added")
    elif len(sys.argv) == 2:
        # default SQL engine choosen: PostgreSQL
        url_list = sys.argv[1].split(',')
        if len(url_list) == 1:
            print("Complete url received.")
            url = url_list
        elif len(url_list) == 6:
            print("Url shall be created")
            psql_engine = CreatePostgreSQLDb()
            url = psql_engine.make_url(url_list[0], url_list[1], url_list[2], url_list[3], url_list[4], int(url_list[5]))
        else:
            raise ValueError("{} is not the right amount of values for the url. [1 or 6]\n".format(len(url_list)))
        pass
    elif len(sys.argv) == 3:
        # Choose own SQL Engine
        pass
    else:
        # pass all parameters trough
        try:
            pass
        except:
            pass

    create_default_tables()
