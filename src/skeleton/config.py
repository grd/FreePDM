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
    def __init__(self):
        self.load_config()

    def load_config(self):
        config = configparser.ConfigParser()
        config['DEFAULT'] = {'StartupDirectory': '',
                            'Filter': '',
                            'LogFile': '',
                            'LogLevel': ''}
        with open('example.ini', 'w') as configfile:
            config.write(configfile)
