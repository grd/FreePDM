"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
from pathlib import Path
import sys

from PySide2.QtWidgets import QApplication, QMainWindow, QTableWidget, QTableWidgetItem
from PySide2 import QtCore, QtWidgets
from PySide2.QtCore import QFile, Qt
from PySide2.QtUiTools import QUiLoader

sys.path.append(os.fspath(Path(__file__).resolve().parents[1] / 'Skeleton'))

from directorymodel import DirectoryModel


class MainWindow(QMainWindow):
    def __init__(self):
        self.dir = os.path.expanduser('~')
        if len(sys.argv) == 2:
            self.dir = sys.argv[1]
        print("self.dir =", self.dir)

        super(MainWindow, self).__init__()
        self.load_ui()
        self.load_data()
    
    def load_ui(self):
        loader = QUiLoader()
        path = os.fspath(Path(__file__).resolve().parents[1] / "GUI/MainWindow.ui")
        print(path)
        ui_file = QFile(path)
        ui_file.open(QFile.ReadOnly)
        self.ui = loader.load(ui_file, self)
        self.ui.setGeometry(50, 40, 800, 600)
        # Some change below based on https://pythonprogramming.net/basic-gui-pyqt-tutorial/
        self.ui.setWindowTitle("FreePDM")  # Done in ui file
        # self.ui.setWindowIcon(QtGui.QIcon(os.fspath(Path(__file__).resolve().parents[1] / "ui/logos/O_logo-32x32.png")))  # Probably done in ui file OSX don't show icon
        self.ui.show()
        ui_file.close()
        self.ui.tableWidget.verticalHeader().setVisible(False)
        self.ui.tableWidget.setSelectionBehavior(QTableWidget.SelectRows)
        self.ui.tableWidget.selectionModel().selectionChanged.connect(self.on_selectionChanged)
        self.ui.CheckInButton('Check In', clicked=self.retrieveCheckButtonValues)

    def retrieveCheckButtonValues(self):
        for row in range(self.ui.tableWidget.rowCount()):
            if self.ui.tableWidget.item(row, 0).checkState == Qt.CheckState.Checked:
                print("selected row: ", row)

    def on_selectionChanged(self, selected, deselected):
        for ix in selected.indexes():
            print('Selected Cell Location Row: {0}, Column: {1}'.format(ix.row(), ix.column()))

        for ix in deselected.indexes():
            print('Deselected Cell Location Row: {0}, Column: {1}'.format(ix.row(), ix.column()))


    def load_data(self):
        dm = DirectoryModel(self.dir)
        row = 0
        self.ui.tableWidget.setRowCount(dm.size())
        for item in dm.dirList:
            cb = QTableWidgetItem("")
            cb.setFlags(Qt.ItemFlag.ItemIsUserCheckable | Qt.ItemFlag.ItemIsEnabled)
            cb.setCheckState(Qt.CheckState.Unchecked)
            self.ui.tableWidget.setItem(row, 0, cb)

            self.ui.tableWidget.setItem(row, 1, QtWidgets.QTableWidgetItem(item["dirOrFile"]))
            self.ui.tableWidget.setItem(row, 2, QtWidgets.QTableWidgetItem(item["filename"]))
            self.ui.tableWidget.setItem(row, 3, QtWidgets.QTableWidgetItem(item["size"]))
            row=row+1


def main():
    QtCore.QCoreApplication.setAttribute(QtCore.Qt.AA_ShareOpenGLContexts)
    app = QApplication(sys.argv)
    widget = MainWindow()
    sys.exit(app.exec_())


if __name__ == '__main__':
    mainw = main()
