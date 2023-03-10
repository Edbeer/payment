package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func NewPostgresDB() (*sql.DB, error) {
	connString := "host=paymentdb user=postgres password=postgres dbname=paymentdb sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}