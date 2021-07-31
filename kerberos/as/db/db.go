package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "./clients.db")
	if err != nil {
		log.Fatalln("Could not open SQLite database: ", err)
	}
}
