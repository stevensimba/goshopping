package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" //init db;
	"github.com/joho/godotenv"
)

func DbConn() (db *sql.DB, err error) {

	err2 := godotenv.Load()

	if err2 != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("name")
	userpwd := os.Getenv("password")
	host := os.Getenv("dbhost")
	port := os.Getenv("dbport")
	dbname := os.Getenv("dbname")

	path := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, userpwd, host, port, dbname)

	db, err = sql.Open("mysql", path)

	return db, err

}
