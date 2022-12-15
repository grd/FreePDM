#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import sys
import os
from pathlib import Path

from PySide2.QtWidgets import QDialog  # type: ignore
from PySide2.QtCore import QFile  # type: ignore
from PySide2.QtUiTools import QUiLoader  # type: ignore

sys.path.append(os.fspath(Path(__file__).resolve().parents[1] / 'authenticate'))
# import


class Authenticate(QDialog):
    def __init__(self):
        super(Authenticate, self).__init__()

        self.load_ui()
        # self.buttonBox.accepted.connect(self.accept)
        # self.buttonBox.rejected.connect(self.reject)

    def load_ui(self):
        loader = QUiLoader()
        path = os.fspath(Path(__file__).resolve().parents[1] / "gui/Authenticate.ui")
        print(path)
        ui_file = QFile(path)
        ui_file.open(QFile.ReadOnly)
        self.ui = loader.load(ui_file, self)
        self.ui.setGeometry(400, 300, 400, 300)
        self.ui.setWindowTitle("Login to Db")
        print("before show")
        self.ui.show()
        print("after show")


def authenticate_dialog():
    """Start Login / Authentication dialog"""
    print("dialog is running")
    auth = Authenticate()
    auth.show()


if __name__ == "__main__":
    # just for testing
    authenticate_dialog()
