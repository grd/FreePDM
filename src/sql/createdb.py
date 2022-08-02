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
            engine = create_engine('mysql://scott:tiger@localhost/foo', echo=self.echo, future=self.future)
            pass
        elif (self.dialect == "mysqlclient") or (self.dialect == "mysqldb"):
            # mysqlclient (a maintained fork of MySQL-Python)
            engine = create_engine('mysql+mysqldb://scott:tiger@localhost/foo', echo=self.echo, future=self.future)
            pass
        elif (self.dialect == "PyMySQL") or (self.dialect == "pymysql"):
            # PyMySQL
            engine = create_engine('mysql+pymysql://scott:tiger@localhost/foo', echo=self.echo, future=self.future)
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
            engine = create_engine('postgresql://scott:tiger@localhost/mydatabase', echo=self.echo, future=self.future)
            pass
        elif self.dialect == "psycopg2":
            # psycopg2
            engine = create_engine('postgresql+psycopg2://scott:tiger@localhost/mydatabase', echo=self.echo, future=self.future)
            pass
        elif self.dialect == "pg8000":
            # pg8000
            engine = create_engine('postgresql+pg8000://scott:tiger@localhost/mydatabase', echo=self.echo, future=self.future)
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
        engine = create_engine("sqlite+pysqlite:///:memory:", echo=self.echo, future=self.future)  # start from memory


def start_the_engine():
    """Start your engine"""
    pass


def create_default_tables():
    pass


if __name__ == "__main__":
    import sys

    create_default_tables()
