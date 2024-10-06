// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/postgres"
)

// Generic SQL Database class
type Database struct {
	db *sql.DB

	drivername    string
	username      string
	password      string
	host          string
	port          int
	database_name string
}

type QueryResult struct {
	rowsAffected int64
	err          error
}

func NewDatabase(username, password, host string) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@%s", username, password, host)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Exec(query string, args ...interface{}) (*QueryResult, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return &QueryResult{}, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return &QueryResult{}, err
	}

	rowsAffected, _ := result.RowsAffected()
	return &QueryResult{rowsAffected: rowsAffected}, nil
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

/*
Create new url

Parameters
----------

drivername [str] :

	drivername. For example 'postgresql+psycopg2'

username [str] :

	username

password [str] :

	password.

host [str] :

	host address. For example localhost

port [int] :

	SQL port.

database_name [str] :

	name of the database

Returns
-------

url [URL object]
*/

// func (d *Database) MakeUrl(drivername, username, password, host string, port int, database_name string) (engine.URL, error) {
// 	d.drivername = drivername
// 	d.username = username
// 	d.password = password
// 	d.host = host
// 	d.port = port
// 	d.database_name = database_name

// 	new_url = engine.URL.create(self.drivername, self.username, self.password, self.host, self.port, self.database_name)
// 	return (new_url)

// }
