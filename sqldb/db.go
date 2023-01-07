package sqldb

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

var l *log.Logger

func init() {
	l = log.New(os.Stdout, "db", log.LstdFlags)

	db, err := sql.Open("sqlite3", "file:whatsappbot.db?foreign_keys=on")
	if err != nil {
		l.Fatal(err)
	}

	DB = db
}
