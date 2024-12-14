package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// errors is a map of error output value in ErrorHandler
var errors = map[string]string{
	"web":        "😮 Oops! Something went wrong",
	"restricted": "😣 Oops! this is a restricted path.\nplease use another path.",
}

// ErrorHandler outputs errors and safely exits the program
func ErrorHandler(errType string) {
	fmt.Println(errors[errType])
	os.Exit(0)
}

func Getdb(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create the users table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	createsessionquerry := `
	CREATE TABLE IF NOT EXISTS sessions (
		session_id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	if _, err := db.Exec(createsessionquerry); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return db, nil
}
