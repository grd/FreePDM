#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from sql import SQLUser, SQLRole, SQLProject, SQLItem, SQLModel, SQLDocument, SQLModelMaterial, SQLHistory, SQLPurchase, SQLManufacturer, SQLVendor

from sqlalchemy.orm import declarative_base
from sqlalchemy import create_engine
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

    def start_engine(self, encoding: str, echo: bool, future: bool, dialect: Optional[str]):
        self.encoding = encoding
        self.echo = echo
        self.future = future
        self.dialect = dialect
        # # https://docs.sqlalchemy.org/en/14/core/engines.html#database-urls
        if (self.dialect == "default") or (self.dialect is None):
            # default
            engine = create_engine('mysql://scott:tiger@localhost/foo')
            pass
        elif (self.dialect == "mysqlclient") or (self.dialect == "mysqldb"):
            # mysqlclient (a maintained fork of MySQL-Python)
            engine = create_engine('mysql+mysqldb://scott:tiger@localhost/foo')
            pass
        elif (self.dialect == "PyMySQL") or (self.dialect == "pymysql"):
            # PyMySQL
            engine = create_engine('mysql+pymysql://scott:tiger@localhost/foo')
        else:
            pass
        engine = create_engine("sqlite+pysqlite:///:memory:", echo=True, future=True)  # start from memory


class CreatePostgreSQLDb(CreateDb):
    """Feed Forward of generic SQL functions to PostgeSQL"""
    # https://docs.sqlalchemy.org/en/14/core/engines.html#postgresql

    def __init__(self):
        print("PostgreSQL")
        super(CreatePosgrSQLDb, self).__init__()

    def start_engine(self, encoding: str, echo: bool, future: bool, dialect: Optional[str]):
        self.encoding = encoding
        self.echo = echo
        self.future = future
        self.dialect = dialect
        # https://docs.sqlalchemy.org/en/14/core/engines.html#database-urls
        if (self.dialect == "default") or (self.dialect is None):
            # default
            engine = create_engine('postgresql://scott:tiger@localhost/mydatabase')
            pass
        elif self.dialect == "psycopg2":
            # psycopg2
            engine = create_engine('postgresql+psycopg2://scott:tiger@localhost/mydatabase')
            pass
        elif self.dialect == "pg8000":
            # pg8000
            engine = create_engine('postgresql+pg8000://scott:tiger@localhost/mydatabase')
        else:
            pass


class CreateSQLiteDb(CreateDb):  # Alles in een file of beter te splitsen?
    """Feed Forward of generic SQL functions to SQLite"""
    # https://docs.sqlalchemy.org/en/14/core/engines.html#sqlite

    def __init__(self):
        super(CreateSQLiteDb, self).__init__()
        print("SQLite")

    def start_engine(self, encoding: str, echo: bool, future: bool):
        self.encoding = encoding
        self.echo = echo
        self.future = future
        # https://docs.sqlalchemy.org/en/14/core/engines.html#database-urls
        engine = create_engine("sqlite+pysqlite:///:memory:", echo=True, future=True)  # start from memory


def start_the_engine():
    pass


def create_default_tables():
    pass


if __name__ == "__main__":
    import sys

    create_default_tables()
