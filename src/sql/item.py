#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

# from sqlalchemy import select
# from typing import Optional, Union


class Item():
    """Item related Class"""

    def __init__(self):
        print("Generic Item")

    def create_item(self):
        """Create new item"""
        raise NotImplementedError("Function create_item is not implemented yet")

    def remove_item(self):
        """Remove existing item"""
        # check if item is new (local == no state)
        # -> if True user can remove
        # -> if False Check if user == admin
        #   -> if True User can remove item
        #   -> if False warning message
        raise NotImplementedError("Function remove_item is not implemented yet")

    def update_item(self):
        """Update existing item"""
        raise NotImplementedError("Function update_item is not implemented yet")


# When inheritance not everything. Do i need a base class?
class Model():
    """Model related Class"""

    def __init__(self):
        print("Generic model")


    def create_model(self):
        """Create new model"""
        # create copy with iter: 0
        raise NotImplementedError("Function create_model is not implemented yet")
        
    def remove_model(self):
        """Remove existing model"""
        # check if model is new (local == no state)
        # -> if True user can remove
        # -> if False Check if user == admin
        #   -> if True User can remove item
        #   -> if False warning message
        raise NotImplementedError("Function remove_model is not implemented yet")

    def update_model(self):
        """Update existing model"""
        # create copy with iter: N
        raise NotImplementedError("Function update_model is not implemented yet")

    def get_version(self, save_iter):
        """Get model that is not latest version"""
        # in UI set only release versions or all
        # Optional two versions to compare
        # is FC able to reload when model is changed?
        raise NotImplementedError("Function get_version is not implemented yet")

# When inheritance not everything. Do i need a base class?
class Document():
    """Document related Class"""
    # How much difference is there between a document and a model?

    def __init__(self):
        print("Generic Document")
