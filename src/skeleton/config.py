"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
import configparser
import appdirs

# Reading and writing the configuration file.
# The location of the file is: [prefix location]/FreePDM/FreePDM.conf
#   [prefix location] depends on which OS you use...
#
#  The items are:
#    * The start of the project directory (or directories).
#      These directories are shown at the start of the main window.
#        "startupdir"
#

appname = 'FreePDM'
config_dir = appdirs.user_config_dir(appname)
config_name = os.path.join(config_dir, 'FreePDM.conf')

###
### Variables
###

startup_directory = ''
filter = ''
log_file = ''
log_level = ''
fast_loading_dir =  ''


def read():
    config = configparser.ConfigParser()
    config.read(config_name)

    # reading variables from section: 'DEFAULT'
    startup_directory = config['DEFAULT']['startup_directory']
    filter =            config['DEFAULT']['filter']
    log_file =          config['DEFAULT']['log_file']
    log_level =         config['DEFAULT']['log_levle']
    fast_loading_dir =  config['DEFAULT']['fast_loading_dir']

def write():
    config = configparser.ConfigParser()
    config['DEFAULT']['startup_directory'] = startup_directory
    config['DEFAULT']['filter'] = filter
    config['DEFAULT']['log_file'] = log_file
    config['DEFAULT']['log_levle'] = log_level
    config['DEFAULT']['fast_loading_dir'] = fast_loading_dir

    with open(config_name, 'w') as configfile:
        config.write(configfile)
