#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import sys
import os
from pathlib import Path

from PySide2 import QtCore  # type: ignore
from PySide2.QtWidgets import QDialog  # type: ignore
from PySide2.QtCore import QFile  # type: ignore
from PySide2.QtUiTools import QUiLoader  # type: ignore

# sys.path.append(os.fspath(Path(__file__).resolve().parents[1] / 'authenticate'))
# import


class Authenticate(QDialog):
    def __init__(self):
        super(Authenticate, self).__init__()

        loader = QUiLoader()
        path = os.fspath(Path(__file__).resolve().parents[1] / "gui/authenticate.ui")
        print(path)
        ui_file = QFile(path)
        ui_file.open(QFile.ReadOnly)
        self.ui = loader.load(ui_file, self)
        self.ui.setGeometry(400, 300, 400, 300)
        self.ui.setWindowTitle("Database Authentication")
        # self.ui.setModal(False)
        # self.ui.setWindowModality(QtCore.Qt.WindowModality(0))
        self.ui.buttonApply.setEnabled(False)

        # self.buttonBox.accepted.connect(self.accept)
        # self.buttonBox.rejected.connect(self.reject)

        # button clicked
        self.ui.buttonCancel.clicked.connect(self.cancel_auth)
        self.ui.buttonConnect.clicked.connect(self.connect_auth)
        self.ui.buttonApply.clicked.connect(self.apply_auth)

        # somehow show don't work out
        # self.ui.show()
        self.ui.exec()

    def load_ui(self):
        """Load default authentication dialog"""
        pass

    def cancel_auth(self):
        """Cancel authentication"""
        # Bring back in previous state
        print("in cancel_auth")
        self.ui.close()

    def connect_auth(self):
        """Connect authentication"""
        # TODO: Connect to server
        # If connect can be applied than
        # - set properety
        # - activater Apply
        # Else
        # - Show labReturnMessage - Connection failed
        pass

    def apply_auth(self):
        """Apply authentication"""
        # TODO: Apply connection == save connection state
        pass


def authenticate_dialog():
    """Start Login / Authentication dialog"""
    print("dialog is running")
    auth = Authenticate()


if __name__ == "__main__":
    # just for testing
    authenticate_dialog()
