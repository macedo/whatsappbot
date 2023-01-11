package sqldb

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

var DSN string

var Provider string

var l *log.Logger

func init() {
	l = log.New(os.Stdout, "db", log.LstdFlags)

	if DSN = os.Getenv("DATABASE_URL"); DSN == "" {
		DSN = "postgres://postgres:postgres@localhost:5432/whatsappbot?sslmode=disable&timezone=UTC&connect_timeout=5"
	}

	Provider = "postgres"

	db, err := sql.Open(Provider, DSN)
	if err != nil {
		l.Fatal(err)
	}

	DB = db
}
