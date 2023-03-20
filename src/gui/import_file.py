"""
    :copyright: Copyright 2023 by the FreePDM team
    :license:   MIT License.
"""

import sys
import os
from pathlib import Path

from PySide2 import QtCore  # type: ignore
from PySide2.QtWidgets import QDialog, QFileDialog  # type: ignore
from PySide2.QtCore import QFile, QStringListModel      # type: ignore
from PySide2.QtUiTools import QUiLoader  # type: ignore

sys.path.append(os.fspath(Path(__file__).resolve().parents[1]))
from filesystem.filesystem import FileSystem 


class ImportFile(QDialog):
    def __init__(self, fs: FileSystem):
        super(ImportFile, self).__init__()

        self.fs = fs

        loader = QUiLoader()
        path = os.fspath(Path(__file__).resolve().parents[1] / "gui/import_file.ui")
        print(path)
        ui_file = QFile(path)
        ui_file.open(QFile.ReadOnly)
        self.ui = loader.load(ui_file, self)

        # button clicked
        self.ui.file_button.clicked.connect(self.file_button)
        self.ui.remove_file.clicked.connect(self.remove_file)
        self.ui.buttonBox.accepted.connect(self.apply_file)
        self.ui.buttonBox.rejected.connect(self.cancel_file)

        self.ui.exec()

    def remove_file(self):
        # TODO: Also ask whether to remove the backup file(s)
        print("pressed")

        
    def cancel_file(self):
        """Cancel file import"""
        print("cancel button pushed")
        self.ui.close()

    def file_button(self):
        """Select the file"""
        dlg = QFileDialog(self)
        dlg.setFileMode(QFileDialog.AnyFile)
        dlg.setNameFilters(["FreeCAD files (*.FCStd)", "All files (*.*)"])
        dlg.selectNameFilter("FreeCAD files (*.FCStd)")

        if dlg.exec_():
            filenames = dlg.selectedFiles()
            self.ui.file_edit.setText(filenames[0])

    def apply_file(self):
        """Apply file import"""
        print("Ok button pushed")
        # TODO: Apply connection == save connection state
        pass


def import_file_dialog(fs: FileSystem):
    """Start Import File dialog"""
    print("dialog is running")
    ImportFile(fs)
