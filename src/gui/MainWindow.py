"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import os
from pathlib import Path
import sys

from PySide2 import QtCore, QtWidgets, QtGui
from PySide2.QtWidgets import QMainWindow, QTableWidget, QTableWidgetItem, QMessageBox
from PySide2.QtCore import Qt

from .EditItem import *
from .filter import *

sys.path.append(os.fspath(Path(__file__).resolve().parents[1] / 'skeleton'))

from directorymodel import DirectoryModel



class Ui_MainWindow(object):
    def setup_ui(self, MainWindow):
        self.root_directory = os.path.expanduser('~')
        if len(sys.argv) == 2:
            self.root_directory = sys.argv[1]
        print("self.root_directory = ", self.root_directory)
        self.current_directory = self.root_directory

        MainWindow.setObjectName("MainWindow")
        MainWindow.resize(817, 600)
        self.centralwidget = QtWidgets.QWidget(MainWindow)
        self.centralwidget.setObjectName("centralwidget")
        self.gridLayout = QtWidgets.QGridLayout(self.centralwidget)
        self.gridLayout.setObjectName("gridLayout")
        self.layoutButtonbar = QtWidgets.QHBoxLayout()
        self.layoutButtonbar.setObjectName("layoutButtonbar")
        self.buttonCheckIn = QtWidgets.QPushButton(self.centralwidget)
        self.buttonCheckIn.setObjectName("buttonCheckIn")
        self.layoutButtonbar.addWidget(self.buttonCheckIn)
        self.buttonCheckOut = QtWidgets.QPushButton(self.centralwidget)
        self.buttonCheckOut.setObjectName("buttonCheckOut")
        self.layoutButtonbar.addWidget(self.buttonCheckOut)
        self.buttonCheckInOut = QtWidgets.QPushButton(self.centralwidget)
        self.buttonCheckInOut.setObjectName("buttonCheckInOut")
        self.layoutButtonbar.addWidget(self.buttonCheckInOut)
        self.lineCI2P = QtWidgets.QFrame(self.centralwidget)
        self.lineCI2P.setFrameShape(QtWidgets.QFrame.VLine)
        self.lineCI2P.setFrameShadow(QtWidgets.QFrame.Sunken)
        self.lineCI2P.setObjectName("lineCI2P")
        self.layoutButtonbar.addWidget(self.lineCI2P)
        self.buttonPurge = QtWidgets.QPushButton(self.centralwidget)
        self.buttonPurge.setObjectName("buttonPurge")
        self.layoutButtonbar.addWidget(self.buttonPurge)
        self.buttonFilter = QtWidgets.QPushButton(self.centralwidget)
        self.buttonFilter.setObjectName("buttonFilter")
        self.layoutButtonbar.addWidget(self.buttonFilter)
        spacerItem = QtWidgets.QSpacerItem(40, 20, QtWidgets.QSizePolicy.Expanding, QtWidgets.QSizePolicy.Minimum)
        self.layoutButtonbar.addItem(spacerItem)
        self.buttonSearch = QtWidgets.QPushButton(self.centralwidget)
        self.buttonSearch.setObjectName("buttonSearch")
        self.layoutButtonbar.addWidget(self.buttonSearch)
        self.gridLayout.addLayout(self.layoutButtonbar, 0, 0, 1, 1)
        self.tableWorkspace = QtWidgets.QTableWidget(self.centralwidget)
        self.tableWorkspace.setObjectName("tableWorkspace")
        self.tableWorkspace.setColumnCount(7)
        self.tableWorkspace.setRowCount(0)
        item = QtWidgets.QTableWidgetItem()
        self.tableWorkspace.setHorizontalHeaderItem(0, item)
        item = QtWidgets.QTableWidgetItem()
        self.tableWorkspace.setHorizontalHeaderItem(1, item)
        item = QtWidgets.QTableWidgetItem()
        self.tableWorkspace.setHorizontalHeaderItem(2, item)
        item = QtWidgets.QTableWidgetItem()
        self.tableWorkspace.setHorizontalHeaderItem(3, item)
        item = QtWidgets.QTableWidgetItem()
        self.tableWorkspace.setHorizontalHeaderItem(4, item)
        item = QtWidgets.QTableWidgetItem()
        self.tableWorkspace.setHorizontalHeaderItem(5, item)
        item = QtWidgets.QTableWidgetItem()
        self.tableWorkspace.setHorizontalHeaderItem(6, item)
        self.gridLayout.addWidget(self.tableWorkspace, 1, 0, 1, 1)
        MainWindow.setCentralWidget(self.centralwidget)
        self.menubar = QtWidgets.QMenuBar(MainWindow)
        self.menubar.setGeometry(QtCore.QRect(0, 0, 817, 22))
        self.menubar.setObjectName("menubar")
        MainWindow.setMenuBar(self.menubar)
        self.statusbar = QtWidgets.QStatusBar(MainWindow)
        self.statusbar.setObjectName("statusbar")
        MainWindow.setStatusBar(self.statusbar)

        self.tableWorkspace.setColumnWidth(0, 5)
        self.tableWorkspace.setColumnWidth(1, 200)
        self.tableWorkspace.setColumnWidth(2, 200)
        self.tableWorkspace.setColumnWidth(3, 80)
        self.tableWorkspace.setColumnWidth(4, 80)
        self.tableWorkspace.setColumnWidth(5, 80)
        self.tableWorkspace.verticalHeader().setVisible(False)
        self.tableWorkspace.setSelectionBehavior(QTableWidget.SelectRows)

        self.retranslateUi(MainWindow)
        QtCore.QMetaObject.connectSlotsByName(MainWindow)

        self.buttonPurge.clicked.connect(self.purge)

        self.tableWorkspace.setSelectionBehavior(QTableWidget.SelectRows)
        # self.ui.tableWorkspace.selectionModel().selectionChanged.connect(self.on_selection_changed)
        self.tableWorkspace.doubleClicked.connect(self.file_double_clicked)
        # self.ui.buttonCheckOutButton('Check In', clicked=self.retrieve_check_button_values)
        self.buttonFilter.clicked.connect(self.set_filter)
        self.buttonPurge.clicked.connect(self.purge)

        self.load_data()


    def load_data(self):
        self.dirmodel = DirectoryModel(self.current_directory, self.root_directory != self.current_directory)
        row = 0
        self.tableWorkspace.setRowCount(self.dirmodel.size())
        # https://stackoverflow.com/questions/39511181/python-add-checkbox-to-every-row-in-qtablewidget
        for item in self.dirmodel.directoryList:
            checkb = QTableWidgetItem("")
            checkb.setFlags(Qt.ItemFlag.ItemIsUserCheckable | Qt.ItemFlag.ItemIsEnabled)
            checkb.setCheckState(Qt.CheckState.Unchecked)
            self.tableWorkspace.setItem(row, 0, checkb)

            file = QTableWidgetItem(item['filename'])
            file.setFlags(file.flags() ^ Qt.ItemIsEditable)
            if item['type'] == 'Directory':
                file.setForeground(QtGui.QColor('blue'))
            self.tableWorkspace.setItem(row, 1, file)

            description = QTableWidgetItem('')
            description.setFlags(description.flags() ^ Qt.ItemIsEditable)
            self.tableWorkspace.setItem(row, 2, description)

            state = QTableWidgetItem('')
            state.setFlags(state.flags() ^ Qt.ItemIsEditable)
            self.tableWorkspace.setItem(row, 3, state)

            version = QTableWidgetItem('')
            version.setFlags(version.flags() ^ Qt.ItemIsEditable)
            self.tableWorkspace.setItem(row, 4, version)

            itemtype = QTableWidgetItem(item['type'])  # type is default function in python
            itemtype.setFlags(itemtype.flags() ^ Qt.ItemIsEditable)
            self.tableWorkspace.setItem(row, 5, itemtype)

            size = QTableWidgetItem(item['size'])
            size.setTextAlignment(Qt.AlignRight)
            size.setFlags(size.flags() ^ Qt.ItemIsEditable)
            self.tableWorkspace.setItem(row, 6, size)
            row += 1

        if len(self.dirmodel.purge_list) == 0:
            self.buttonPurge.setEnabled(False)
        else:
            self.buttonPurge.setEnabled(True)


    def purge(self):
        msg = QMessageBox()
        msg.setWindowTitle("Purge")
        msg.setText("Purging means deleting old FreeCAD files.")
        msg.setInformativeText("Are you sure you want to delete {} files?".format(len(self.dirmodel.purge_list)))
        msg.setIcon(QMessageBox.Question)
        msg.setStandardButtons(QMessageBox.Cancel|QMessageBox.Ok)
        msg.setDefaultButton(QMessageBox.Cancel)
        msg.buttonClicked.connect(self.popup_purge)
        x = msg.exec_()

    def popup_purge(self, i):
        if i.text() == '&OK':
         self.dirmodel.purge()

    
    def retranslateUi(self, MainWindow):
        _translate = QtCore.QCoreApplication.translate
        MainWindow.setWindowTitle(_translate("MainWindow", "MainWindow"))
        self.buttonCheckIn.setText(_translate("MainWindow", "Check In"))
        self.buttonCheckOut.setText(_translate("MainWindow", "Check Out"))
        self.buttonCheckInOut.setText(_translate("MainWindow", "Check In/Out"))
        self.buttonPurge.setText(_translate("MainWindow", "Purge"))
        self.buttonFilter.setText(_translate("MainWindow", "Filter"))
        self.buttonSearch.setText(_translate("MainWindow", "Search"))
        item = self.tableWorkspace.horizontalHeaderItem(1)
        item.setText(_translate("MainWindow", "Name"))
        item = self.tableWorkspace.horizontalHeaderItem(2)
        item.setText(_translate("MainWindow", "Description"))
        item = self.tableWorkspace.horizontalHeaderItem(3)
        item.setText(_translate("MainWindow", "State"))
        item = self.tableWorkspace.horizontalHeaderItem(4)
        item.setText(_translate("MainWindow", "Version"))
        item = self.tableWorkspace.horizontalHeaderItem(5)
        item.setText(_translate("MainWindow", "Type"))
        item = self.tableWorkspace.horizontalHeaderItem(6)
        item.setText(_translate("MainWindow", "Size"))


    def set_filter(self):
        filter_dialog()
        self.load_data()

    # deal with doubleclick
    def file_double_clicked(self, event):
        row = event.row()
        item = self.tableWorkspace.item(row, 5).text()

        # Change directory
        if item == 'Directory':
            dir = self.tableWorkspace.item(row, 1).text()
            self.current_directory = os.path.abspath(os.path.join(self.current_directory, dir))
            self.load_data()

        # Edit FC Item
        if item == 'FCStd' or item == 'A2-Assy' or item == 'A3-Assy'or item == 'A4-Assy':
            part = self.tableWorkspace.item(row, 1).text()
            part = os.path.abspath(os.path.join(self.current_directory, part))
            edit_item_dialog(part)
            self.load_data()
 
  

    # deal with checkbox click on a field
    def retrieve_check_button_values(self):
        for row in range(self.tableWidget.rowCount()):
            if self.tableWorkspace.item(row, 0).checkState() == Qt.CheckState.Checked:
                print("selected row: ", row)

    # deal with the rows
    def on_selection_changed(self, selected, deselected):
        for ix in selected.indexes():
            print('Selected Cell Location Row: {0}, Column: {1}'.format(ix.row(), ix.column()))

        for ix in deselected.indexes():
            print('Deselected Cell Location Row: {0}, Column: {1}'.format(ix.row(), ix.column()))



def main():
    QtCore.QCoreApplication.setAttribute(QtCore.Qt.AA_ShareOpenGLContexts)
    app = QtWidgets.QApplication(sys.argv)
    ex = Ui_MainWindow()
    w = QMainWindow()
    ex.setup_ui(w)
    w.show()
    sys.exit(app.exec_())

if __name__ == '__main__':
    mainw = main()
