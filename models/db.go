package models

import (
    "database/sql"
     _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB(connectionString string) {
    var err error
    db, err = sql.Open("mysql", connectionString)
    if err != nil {
        panic(err.Error())  //TODO fix error handling
    }

    if err = db.Ping(); err != nil {
        panic(err.Error())  //TODO fix error handling
    }
}