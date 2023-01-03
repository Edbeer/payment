package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {}

type PostgresStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

