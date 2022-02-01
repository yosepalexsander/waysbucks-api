package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBStore struct {
	DB *sqlx.DB
}

func Connect(db *DBStore) {
	var err error
	db.DB, err = sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	db.DB.SetMaxOpenConns(5)
	db.DB.SetMaxIdleConns(2)
}
