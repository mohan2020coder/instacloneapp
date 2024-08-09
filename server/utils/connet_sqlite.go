package utils

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectSQLiteDB() *sql.DB {
	db, err := sql.Open("sqlite3", os.Getenv("SQLITE_DB_PATH"))
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}

	// Optionally, set PRAGMA settings here for SQLite if needed

	log.Println("Connected to SQLite!")
	return db
}
