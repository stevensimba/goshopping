package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" //init db;
)

// Connect to a mysql database

func DbConn() (db *sql.DB, err error) {
	username := os.Getenv("name")
	userpwd := os.Getenv("password")
	host := os.Getenv("dbhost")
	port := os.Getenv("dbport")
	dbname := os.Getenv("dbname")

	path := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, userpwd, host, port, dbname)

	db, err = sql.Open("mysql", path)

	return db, err

}
