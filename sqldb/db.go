package sqldb

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

var Provider = "postgres"

var l *log.Logger

func init() {
	l = log.New(os.Stdout, "db", log.LstdFlags)

	db, err := sql.Open(Provider, "postgres://postgres:postgres@localhost:5432/whatsappbot?sslmode=disable&timezone=UTC&connect_timeout=5")
	if err != nil {
		l.Fatal(err)
	}

	DB = db
}
