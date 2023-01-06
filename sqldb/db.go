package sqldb

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	var err error

	db, err := sql.Open("sqlite3", "whatsappbot.db?foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}
