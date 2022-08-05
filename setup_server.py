#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

from setuptools import setup
import os

version_path = os.path.join(os.path.abspath(os.path.dirname(__file__)), 
                            "src", "version.py")
with open(version_path) as fp:
    exec(fp.read())

setup(name='freePDM-Server',
      version=str(__version__),
      packages=['src', ],
      maintainer=["grd", "Jee-Bee"],
      maintainer_email="",
      url="https://github.com/grd/FreePDM",
      description="Backend for FreePDM",
      install_requires=['SQLAlchemy', 'sqlalchemy-utils'],
      include_package_data=True)
