"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
from pathlib import Path
import sys

from PySide2 import QtCore, QtWidgets, QtGui
from PySide2.QtWidgets import QApplication, QMainWindow, QTableWidget, QTableWidgetItem
from PySide2.QtCore import QFile, Qt
from PySide2.QtUiTools import QUiLoader

sys.path.append(os.fspath(Path(__file__).resolve().parents[1] / 'skeleton'))

from directorymodel import DirectoryModel


class MainWindow(QMainWindow):
    def __init__(self):
        super(MainWindow, self).__init__()

        self.root_directory = os.path.expanduser('~')
        if len(sys.argv) == 2:
            self.root_directory = sys.argv[1]
        print("self.root_directory = ", self.root_directory)
        self.current_directory = self.root_directory

        self.load_ui()
        self.load_data()

    def load_ui(self):
        loader = QUiLoader()
        path = os.fspath(Path(__file__).resolve().parents[1] / "gui/MainWindow.ui")
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
        self.ui.tableWorkspace.setColumnWidth(0, 10)
        self.ui.tableWorkspace.setColumnWidth(1, 150)
        self.ui.tableWorkspace.setColumnWidth(2, 200)
        self.ui.tableWorkspace.setColumnWidth(3, 80)
        self.ui.tableWorkspace.setColumnWidth(4, 80)
        self.ui.tableWorkspace.setColumnWidth(5, 80)
        self.ui.tableWorkspace.verticalHeader().setVisible(False)
        self.ui.tableWorkspace.setColumnWidth(0, 30)
        self.ui.tableWorkspace.setSelectionBehavior(QTableWidget.SelectRows)
        # self.ui.tableWorkspace.selectionModel().selectionChanged.connect(self.on_selectionChanged)
        self.ui.tableWorkspace.doubleClicked.connect(self.fileDoubleClicked)
        # self.ui.buttonCheckOutButton('Check In', clicked=self.retrieveCheckButtonValues)

    # deal with doubleclick
    def fileDoubleClicked(self, event):
        row = event.row()
        item = self.ui.tableWorkspace.item(row, 5).text()
        print('item = ', item)

        # Change directory
        if item == 'Directory':
            dir = self.ui.tableWorkspace.item(row, 1).text()
            self.current_directory = os.path.abspath(os.path.join(self.current_directory, dir))
            self.load_data()

    # deal with checkbox click on a field
    def retrieveCheckButtonValues(self):
        for row in range(self.ui.tableWidget.rowCount()):
            if self.ui.tableWorkspace.item(row, 0).checkState() == Qt.CheckState.Checked:
                print("selected row: ", row)

    # deal with the rows
    def on_selectionChanged(self, selected, deselected):
        for ix in selected.indexes():
            print('Selected Cell Location Row: {0}, Column: {1}'.format(ix.row(), ix.column()))

        for ix in deselected.indexes():
            print('Deselected Cell Location Row: {0}, Column: {1}'.format(ix.row(), ix.column()))

    def load_data(self):
        dm = DirectoryModel(self.current_directory, self.root_directory != self.current_directory)
        row = 0
        self.ui.tableWorkspace.setRowCount(dm.size())
        # https://stackoverflow.com/questions/39511181/python-add-checkbox-to-every-row-in-qtablewidget
        for item in dm.directoryList:
            cb = QTableWidgetItem("")
            cb.setFlags(Qt.ItemFlag.ItemIsUserCheckable | Qt.ItemFlag.ItemIsEnabled)
            cb.setCheckState(Qt.CheckState.Unchecked)
            self.ui.tableWorkspace.setItem(row, 0, cb)

            file = QTableWidgetItem(item['filename'])
            file.setFlags(file.flags() ^ Qt.ItemIsEditable)
            if item['type'] == 'Directory':
                file.setForeground(QtGui.QColor('blue')) 
            self.ui.tableWorkspace.setItem(row, 1, file)
            
            type = QTableWidgetItem(item['type'])
            type.setFlags(type.flags() ^ Qt.ItemIsEditable)
            self.ui.tableWorkspace.setItem(row, 5, type)

            size = QTableWidgetItem(item['size'])
            size.setFlags(size.flags() ^ Qt.ItemIsEditable)
            self.ui.tableWorkspace.setItem(row, 6, size)
            row += 1


def main():
    QtCore.QCoreApplication.setAttribute(QtCore.Qt.AA_ShareOpenGLContexts)
    app = QApplication(sys.argv)
    widget = MainWindow()
    sys.exit(app.exec_())


if __name__ == '__main__':
    mainw = main()
