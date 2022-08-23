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
filter = 0
log_file = ''
log_level = ''
fast_loading_dir =  ''

###
### filter flags
###
show_fc_files_only = 1
hide_versioned_fc_files = 2


def get_filter(filter_flag):
    global filter
    return filter & filter_flag == filter_flag


def read():
    global startup_directory
    global filter
    global log_file
    global log_level
    global fast_loading_dir

    config = configparser.ConfigParser()
    config.read(config_name)

    # reading variables from section: 'DEFAULT'
    startup_directory = config['DEFAULT']['startup_directory']
    filter =        int(config['DEFAULT']['filter'])
    log_file =          config['DEFAULT']['log_file']
    log_level =         config['DEFAULT']['log_levle']
    fast_loading_dir =  config['DEFAULT']['fast_loading_dir']


def write():
    global startup_directory
    global filter
    global log_file
    global log_level
    global fast_loading_dir

    config = configparser.ConfigParser()
    config['DEFAULT']['startup_directory'] = startup_directory
    config['DEFAULT']['filter'] = str(filter)
    config['DEFAULT']['log_file'] = log_file
    config['DEFAULT']['log_levle'] = log_level
    config['DEFAULT']['fast_loading_dir'] = fast_loading_dir

    with open(config_name, 'w') as configfile:
        config.write(configfile)


# create the new directory if it doesn't exist
if not os.path.exists(config_dir):
    os.makedirs(config_dir)

# create a new config file when it doesn't exist
if not os.path.isfile(config_name):
    write()
