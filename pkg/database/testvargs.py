#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Tue Aug  2 16:01:12 2022

@author: jbtw
"""

if __name__ == "__main__":
    import sys
    
    print(sys.argv, len(sys.argv))
    new_list = sys.argv[2].split(',')
    print(new_list, len(new_list))

# test code: python testvargs.py 3 something,else,,  # sys.argv = length3; new_list =length4
