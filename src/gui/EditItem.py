"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
from pathlib import Path
import sys

from PySide2 import QtGui
from PySide2.QtWidgets import QDialog
from PySide2.QtCore import QFile
from PySide2.QtUiTools import QUiLoader

sys.path.append(os.fspath(Path(__file__).resolve().parents[1] / 'skeleton'))

from itemdatamodel import ItemDataModel


class EditItem(QDialog):
    def __init__(self, file):
        super(EditItem, self).__init__()
        self.file = file
        print("File is ", self.file)
        self.idm = ItemDataModel(self.file)

        loader = QUiLoader()
        path = os.fspath(Path(__file__).resolve().parents[1] / "gui/EditItem.ui")
        print(path)
        ui_file = QFile(path)
        ui_file.open(QFile.ReadOnly)
        self.ui = loader.load(ui_file, self)
        self.ui.setWindowTitle("Edit Item")

        self.ui.weightEdit.setReadOnly(True)
        self.ui.unitEdit.setReadOnly(True)

        print('File ' + self.file + ' write has access ' + str(os.access(self.file, os.W_OK)))
        # if os.access(self.file, os.W_OK) == True:
        #     self.ui.buttonBox.Ok.setEnabled(False)

        if 'Unit' in self.idm.document_properties:
            self.ui.unitEdit.setText(self.idm.document_properties['Unit'])
        if "thumbnail" in self.idm.document_properties:
            pixmap = QtGui.QPixmap(self.idm.thumbnail)
            self.ui.lbl.setPixmap(pixmap.scaled(256, 256))

        self.ui.file_name.setText(file)
        self.ui.company_name.setText(self.idm.document_properties["Company"])
        self.ui.user_name.setText(self.idm.document_properties["CreatedBy"])
        self.ui.date_of_creation.setText(self.idm.document_properties["CreationDate"])


    def store_data(self):
        pass


def edit_item_dialog(item):
    edit = EditItem(item)
    if edit.ui.exec_() == 1: # Ok button pushed
        edit.store_data()
