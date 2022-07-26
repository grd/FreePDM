"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import configparser

# Reading and writing the configuration file.
# The location of the file is: (I don't know yet, and the name too).
# The file format is probably ini since that is the easiest.
#
#  The items are:
#    * The start of the project directory (or directories).
#      These directories are shown at the start of the main window.
#        "startupdir"
#

class generalOptions(object):
    def __init__():

config = configparser.ConfigParser()
config['DEFAULT'] = {'StartupDirectory': '',
                     'Filter': '',
                    'Compression': 'yes',
                    'CompressionLevel': '9'}
config['bitbucket.org'] = {}
config['bitbucket.org']['User'] = 'hg'
config['topsecret.server.com'] = {}
config['DEFAULT']['ForwardX11'] = 'yes'
with open('example.ini', 'w') as configfile:
   config.write(configfile)
