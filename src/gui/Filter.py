"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
from pathlib import Path

from PySide2.QtWidgets import QDialog
from PySide2.QtUiTools import QUiLoader



class Filter(QDialog):
    def __init__(self):
        super(Filter, self).__init__()

        loader = QUiLoader()
        path = os.fspath(Path(__file__).resolve().parents[1] / "gui/Filter.ui")
        print(path)
        # Some change below based on https://pythonprogramming.net/basic-gui-pyqt-tutorial/
        # self.ui.setWindowIcon(QtGui.QIcon(os.fspath(Path(__file__).resolve().parents[1] / "ui/logos/O_logo-32x32.png")))  # Probably done in ui file OSX don't show icon

        self.ui.setWindowTitle("Filter")  # Done in ui file




