package db

import (
	"database/sql"
	"log"
)

// db is the global database connection pool.
var DataBase *sql.DB

func Connect() {
	connStr := "user=elys password=elys dbname=indexer sslmode=disable"
	var err error
	DataBase, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
}

func Close() error {
	return DataBase.Close()
}

func Ping() error {
	return DataBase.Ping()
}
