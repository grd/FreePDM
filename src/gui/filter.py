import sys
import os
from pathlib import Path

from PySide2.QtCore import *
from PySide2.QtGui import *
import PySide2.QtWidgets as qtw
from PySide2.QtWidgets import QDialog

sys.path.append(os.fspath(Path(__file__).resolve().parents[1] / 'skeleton'))
import config as conf

class FilterDialog(QDialog):
    def __init__(self):
        super(FilterDialog, self).__init__()

        self.resize(400, 300)
        self.buttonBox = qtw.QDialogButtonBox(self)
        self.buttonBox.setObjectName(u"buttonBox")
        self.buttonBox.setGeometry(QRect(30, 240, 341, 32))
        self.buttonBox.setOrientation(Qt.Horizontal)
        self.buttonBox.setStandardButtons(qtw.QDialogButtonBox.Cancel|qtw.QDialogButtonBox.Ok)

        self.show_fc_files_only = qtw.QCheckBox(self)
        self.show_fc_files_only.setObjectName(u"show_fc_files_only")
        self.show_fc_files_only.setGeometry(QRect(80, 20, 250, 17))

        self.show_versioned_fc_files = qtw.QCheckBox(self)
        self.show_versioned_fc_files.setObjectName(u"show_versioned_fc_files")
        self.show_versioned_fc_files.setGeometry(QRect(80, 50, 250, 17))

        self.show_all_files = qtw.QCheckBox(self)
        self.show_all_files.setObjectName(u"show_all_files")
        self.show_all_files.setGeometry(QRect(80, 80, 250, 17))

        self.retranslate_ui()
        self.buttonBox.accepted.connect(self.accept)
        self.buttonBox.rejected.connect(self.reject)

        QMetaObject.connectSlotsByName(self)

        self.retrieve_data()
        self.show()

    def retranslate_ui(self):
        self.setWindowTitle(QCoreApplication.translate("Filter Dialog", u"Filter Dialog", None))
        self.show_fc_files_only.setText(QCoreApplication.translate("Dialog", u"Show FreeCAD files only", None))
        self.show_versioned_fc_files.setText(QCoreApplication.translate("Dialog", u"Show versioned FreeCAD files", None))
        self.show_all_files.setText(QCoreApplication.translate("Dialog", u"Show all files", None))

    def retrieve_data(self):
        conf.read()
        filter = conf.filter & conf.show_fc_files_only
        if filter == conf.show_fc_files_only:
            self.show_fc_files_only.setCheckState(Qt.Checked)
        else:
            self.show_fc_files_only.setCheckState(Qt.Unchecked)

        filter = conf.filter & conf.show_versioned_fc_files
        if filter == conf.show_versioned_fc_files:
            self.show_versioned_fc_files.setCheckState(Qt.Checked)
        else:
            self.show_versioned_fc_files.setCheckState(Qt.Unchecked)

        filter = conf.filter & conf.show_all_files
        if filter == conf.show_all_files:
            self.show_all_files.setCheckState(Qt.Checked)
        else:
            self.show_all_files.setCheckState(Qt.Unchecked)


    
    def store_data(self):
        conf.filter = 0
        if self.show_fc_files_only.isChecked() == True:
            conf.filter = conf.filter | conf.show_fc_files_only
        if self.show_versioned_fc_files.isChecked() == True:
            conf.filter = conf.filter | conf.show_versioned_fc_files            
        if self.show_all_files.isChecked():
            conf.filter = conf.filter | conf.show_all_files
        conf.write()


def filter_dialog():
    filter = FilterDialog()
    if filter.exec_() == 1: # Ok button pushed
        filter.store_data()
