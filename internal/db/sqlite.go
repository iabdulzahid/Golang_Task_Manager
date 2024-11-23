package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		priority TEXT,
		due_date TEXT,
		labels TEXT
	);`)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() *sql.DB {
	return db
}
