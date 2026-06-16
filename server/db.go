/*
	server/db.go — Database layer for the CoolRMM server.
	Opens (or creates) the SQLite database file, ensures the schema exists,
	and provides insert_checkin() to persist incoming agent check-ins.
*/

package main

import (
	"database/sql"
	"fmt"
	"time"

	// Import the pure-Go SQLite driver — registers itself with database/sql.
	_ "modernc.org/sqlite"
)

// open_db opens the SQLite database at the given file path.
// Creates the file and schema if they don't exist yet.
func open_db(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Confirm the connection is actually alive.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create the checkins table if it doesn't already exist.
	if err := create_schema(db); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return db, nil
}

// create_schema creates the checkins table if it doesn't exist.
func create_schema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS checkins (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		hostname      TEXT NOT NULL,
		logged_user   TEXT NOT NULL,
		checked_in_at DATETIME NOT NULL
	);`

	_, err := db.Exec(query)
	return err
}

// insert_checkin writes a single check-in record to the database.
func insert_checkin(db *sql.DB, hostname, logged_user string, checked_in_at time.Time) error {
	query := `INSERT INTO checkins (hostname, logged_user, checked_in_at) VALUES (?, ?, ?);`

	_, err := db.Exec(query, hostname, logged_user, checked_in_at.UTC())
	if err != nil {
		return fmt.Errorf("failed to insert check-in: %w", err)
	}

	return nil
}

/*
	This file is compiled as part of the server binary — it is not run directly.
	To build:
	  go build ./server
	  GOOS=windows GOARCH=amd64 go build -o coolrmm-server.exe ./server

	The database file is created automatically on first run.
	Default path: coolrmm.db (in the working directory of the server binary).

	To run tests:
	  go test ./server/...
*/
