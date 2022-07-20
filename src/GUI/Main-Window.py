"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import sys
from PyQt5 import QtWidgets
from PyQt5.QtWidgets import QMainWindow, QApplication, QMessageBox, QDialog
from PyQt5.uic import loadUi

sys.path.append('../Skeleton')       

import directorymodel

class MainWindow(QMainWindow):
    def __init__(self):
        super(MainWindow, self).__init__()
        loadUi("Main-Window.ui", self)
        self.loaddata()

    def loaddata(self):
        people=[{"name":"John","age":45,"address":"New York"}, {"name":"Mark", "age":40,"address":"LA"},
                {"name":"George","age":30,"address":"London"}]
        row=0
        self.tableWidget.setRowCount(len(people))
        for person in people:
            self.tableWidget.setItem(row, 0, QtWidgets.QTableWidgetItem(person["name"]))
            self.tableWidget.setItem(row, 1, QtWidgets.QTableWidgetItem(str(person["age"])))
            self.tableWidget.setItem(row, 2, QtWidgets.QTableWidgetItem(person["address"]))
            row=row+1


        
        
if __name__ == '__main__':
    app = QApplication(sys.argv)
    window = MainWindow()
    window.show()
    app.exec_()
