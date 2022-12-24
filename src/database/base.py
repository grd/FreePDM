#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from sqlalchemy.orm import declarative_base
from sqlalchemy import MetaData
# from sqlalchemy import Table
from sqlalchemy.orm import sessionmaker


Base = declarative_base()

# Sessions / Connetions:
# - https://docs.sqlalchemy.org/en/14/orm/session_api.html#session-and-sessionmaker
# - https://www.fullstackpython.com/sqlalchemy-orm-session-examples.html
# - https://stackoverflow.com/questions/12223335/sqlalchemy-creating-vs-reusing-a-session
Session = sessionmaker()

# Metadata:
# - https://dataedo.com/kb/data-glossary/what-is-metadata
# - https://www.geeksforgeeks.org/difference-between-data-and-metadata/
# - https://www.geeksforgeeks.org/metadata-in-dbms-and-its-types/
metadata_obj = MetaData()
