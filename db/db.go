package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	Conn *sqlx.DB
}

func Connect(db *DB) {
	var err error
	db.Conn, err = sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
			log.Fatal(err)
	}
	log.Println("connect db")
	db.Conn.SetMaxOpenConns(5)
	db.Conn.SetMaxIdleConns(2)
}