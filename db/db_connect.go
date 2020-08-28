package db

import (
	"database/sql"
	"log"
	"os"
)

type DB struct {
	Connection *sql.DB
}

func OpenDB() (*DB, error) {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Print(err)
		OpenDB()
	}

	return &DB{Connection: db}, nil
}
