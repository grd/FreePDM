#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2023 by the FreePDM team
    :license:   MIT License.
"""

import os
from os import path
from typing import List
from datetime import datetime


"""This is a helper module of arbitrary functions."""

def today():
    """today returns the date of today in the format YYYY-MM-DD."""
    return datetime.today().isoformat()[0:10]


def get_key_value_list(fname: str) -> list[str, str]:
    """Returns a key-value list that is stored inside a text file."""

    if not path.isfile(fname):
        raise FileNotFoundError("File " + fname + " not found.")

    if path.getsize(fname) == 0:
        return []

    ret_list = []
    items = []
    with open(fname, "r") as file:
        items = file.readlines()

    for item in items:
        key, value = item.split("=")
        ret_list.append([key, value.strip()])

    return ret_list


def store_key_value_list(fname: str, key_value: list[str, str]):
    """Writes the key-value list to the text."""

    with open(fname, "w") as file:
        for item in key_value:
            file.write(item[0] + "=" + item[1] + "\n")


# def append_key_value_list(fname: str, key_value: list(str, str)):
#     """Appends the key-value list to the text file."""

#     if not path.isfile(fname):
#         raise FileNotFoundError("File " + fname + " not found.")

#     with open(fname, "a") as file:
#         for item in key_value:
#             file.write(item[0] + "=" + item[1] + "\n")
