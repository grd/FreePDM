#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from base import Base
from base import Session
from base import metadata_obj
from sqlalchemy import create_engine
from sqlalchemy.engine import URL
from sqlalchemy.engine import Engine
# from sqlalchemy import Table
from typing import Optional, Union, Literal


class GeneralDb():
    """Generic SQL DataBase class"""
    # https://docs.python.org/3/library/typing.html#typing.NamedTuple

    def __init__(self):
        print("Generic DataBase")

    def make_url(self, drivername: str, username: Optional[str], password:  Optional[str], host: Optional[str], port: Optional[int], database_name: Optional[str]) -> URL:
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

        url [URL object]
        """
        self.drivername = drivername
        self.username = username
        self.password = password
        self.host = host
        self.port = port
        self.database_name = database_name
        new_url = URL.create(self.drivername, self.username, self.password, self.host, self.port, self.database_name)
        return (new_url)

    def create_db(self):
        """Create new database"""
        # https://stackoverflow.com/questions/6506578/how-to-create-a-new-database-using-sqlalchemy
        raise NotImplementedError("Function create_db is not implemented yet")

    def create_table(self):
        """Create new Table"""
        raise NotImplementedError("Function create_table is not implemented yet")

    def delete_table(self):
        """Delete existing Table"""
        raise NotImplementedError("Function delete_table is not implemented yet")

    def get_tables(self):
        """Get Tables from existing database"""
        # https://stackoverflow.com/questions/44193823/get-existing-table-using-sqlalchemy-metadata
        raise NotImplementedError("Function get_tables is not implemented yet")


class MySQLDb(GeneralDb):  # Everything in a file or better to split it?
    """Feed Forward of generic SQL functions to MySQL"""
    # https://docs.sqlalchemy.org/en/14/core/engines.html#mysql

    def __init__(self):
        print("MySQL")
        super(MySQLDb, self).__init__()

    def start_engine(self, url: Union[str, URL], encoding: Optional[str], echo: Union[bool, Literal['debug'], None], future: Optional[bool], dialect: Optional[str]) -> Engine:
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

        if self.echo is None:
            self.echo = False

        if self.future is None:
            self.future = True

        # # https://docs.sqlalchemy.org/en/14/core/engines.html#database-urls
        if (self.dialect == "default") or (self.dialect is None):
            # Installing via `FreePDM-ServerInstaller.sh` installs default engine
            # default
            self.engine = create_engine(self.url, echo=self.echo, future=self.future)
            return (self.engine)
        elif (self.dialect == "mysqlclient") or (self.dialect == "mysqldb"):
            # mysqlclient (a maintained fork of MySQL-Python)
            self.engine = create_engine(self.url, echo=self.echo, future=self.future)
            return (self.engine)
        elif (self.dialect == "PyMySQL") or (self.dialect == "pymysql"):
            # PyMySQL
            self.engine = create_engine(self.url, echo=self.echo, future=self.future)
            return (self.engine)
        else:
            raise ValueError("{} not accepted value for dialect".format(self.dialect))


class PostgreSQLDb(GeneralDb):
    """Feed Forward of generic SQL functions to PostgeSQL"""
    # https://docs.sqlalchemy.org/en/14/core/engines.html#postgresql

    def __init__(self):
        print("PostgreSQL")
        super(PostgreSQLDb, self).__init__()

    def start_engine(self, url: Union[str, URL], encoding: Optional[str], echo: Union[bool, Literal['debug'], None], future: Optional[bool], dialect: Optional[str]) -> Engine:
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

        if self.echo is None:
            self.echo = False

        if self.encoding is None:
            self.encoding = "utf-8"

        if self.future is None:
            self.future = True

        # https://docs.sqlalchemy.org/en/14/core/engines.html#database-urls
        if (self.dialect == "default") or (self.dialect is None):
            # default
            self.engine = create_engine(self.url, echo=self.echo, encoding=self.encoding, future=self.future)
            return (self.engine)
        elif self.dialect == "psycopg2":
            # psycopg2
            self.engine = create_engine(self.url, echo=self.echo, encoding=self.encoding, future=self.future)
            return (self.engine)
        elif self.dialect == "pg8000":
            # pg8000
            self.engine = create_engine(self.url, echo=self.echo, encoding=self.encoding, future=self.future)
            return (self.engine)
        else:
            raise ValueError("{} not accepted value for dialect".format(self.dialect))


class SQLiteDb(GeneralDb):  # Everything in a file or better to split it?
    """Feed Forward of generic SQL functions to SQLite"""
    # https://docs.sqlalchemy.org/en/14/core/engines.html#sqlite

    def __init__(self):
        super(SQLiteDb, self).__init__()
        print("SQLite")

    def start_engine(self, url: Union[str, URL], encoding: Optional[str], echo: Union[bool, Literal['debug'], None], future: Optional[bool]) -> Engine:
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

        if self.echo is None:
            self.echo = False

        if self.encoding is None:
            self.encoding = "utf-8"

        if self.future is None:
            self.future = True

        # https://docs.sqlalchemy.org/en/14/core/engines.html#database-urls
        # exampleurl: "sqlite+pysqlite:///:memory:"
        self.engine = create_engine(self.url, echo=self.echo, encoding=self.encoding, future=self.future)  # start from memory
        return (self.engine)
