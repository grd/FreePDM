#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import sql
import sqlalchemy as sqla
from sqlalchemy import create_engine
from typing import NewType, Optional


class CreateDb():
    """Create database"""

    def __init__(self):
        engine = create_engine("sqlite+pysqlite:///:memory:", echo=True, future=True)  # start from memory
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

def create_default_tables():
    pass

if __name__ == "__main__":
    create_default_tables()
