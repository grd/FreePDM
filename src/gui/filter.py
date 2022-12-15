import sys
import os
from pathlib import Path

from PySide2.QtCore import *  # type: ignore
from PySide2.QtGui import *  # type: ignore
import PySide2.QtWidgets as qtw  # type: ignore
from PySide2.QtWidgets import QDialog  # type: ignore

sys.path.append(os.fspath(Path(__file__).resolve().parents[1] / 'skeleton'))
import config  # type: ignore


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

        self.hide_versioned_fc_files = qtw.QCheckBox(self)
        self.hide_versioned_fc_files.setObjectName(u"hide_versioned_fc_files")
        self.hide_versioned_fc_files.setGeometry(QRect(80, 50, 250, 17))

        self.retranslate_ui()
        self.buttonBox.accepted.connect(self.accept)
        self.buttonBox.rejected.connect(self.reject)

        QMetaObject.connectSlotsByName(self)

        self.retrieve_data()
        self.show()

    def retranslate_ui(self):
        self.setWindowTitle(QCoreApplication.translate("Filter Dialog", u"Filter Dialog", None))
        self.show_fc_files_only.setText(QCoreApplication.translate("Dialog", u"Show FreeCAD files only", None))
        self.hide_versioned_fc_files.setText(QCoreApplication.translate("Dialog", u"Hide versioned FreeCAD files", None))

    def retrieve_data(self):
        self.conf = config.conf()
        self.conf.read()

        state = Qt.Checked if self.conf.get_filter(config.show_fc_files_only) else Qt.Unchecked
        self.show_fc_files_only.setCheckState(state)

        state = Qt.Checked if self.conf.get_filter(config.hide_versioned_fc_files) else Qt.Unchecked
        self.hide_versioned_fc_files.setCheckState(state)

    def store_data(self):
        self.conf.filter = 0
        if self.show_fc_files_only.isChecked() == True:
            self.conf.set_filter(config.show_fc_files_only)
        if self.hide_versioned_fc_files.isChecked() == True:
            self.conf.set_filter(config.hide_versioned_fc_files)
        self.conf.write()


def filter_dialog():
    filter_dlg = FilterDialog()
    if filter_dlg.exec_() == 1: # Ok button pushed
        filter_dlg.store_data()
