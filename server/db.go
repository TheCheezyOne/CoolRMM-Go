/*
        server/db.go — Database layer for the CoolRMM server.
        Opens (or creates) the SQLite database file, ensures the schema exists,
        and provides insert_checkin() to persist incoming agent check-ins.
*/

package main

import (
        "database/sql"
        "fmt"
        "strings"
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

        // Create the commands table for remote shell support.
        if err := create_commands_schema(db); err != nil {
                return nil, fmt.Errorf("failed to create commands schema: %w", err)
        }

        return db, nil
}

// create_schema creates the checkins table if it doesn't exist,
// then migrates any existing DB by adding new columns if they are missing.
func create_schema(db *sql.DB) error {
        // Create the base table — safe to run on an existing DB.
        create_query := `
        CREATE TABLE IF NOT EXISTS checkins (
                id             INTEGER PRIMARY KEY AUTOINCREMENT,
                hostname       TEXT NOT NULL,
                logged_user    TEXT NOT NULL,
                checked_in_at  DATETIME NOT NULL,
                cpu_percent    REAL NOT NULL DEFAULT 0,
                ram_percent    REAL NOT NULL DEFAULT 0,
                disk_percent   REAL NOT NULL DEFAULT 0,
                uptime_seconds INTEGER NOT NULL DEFAULT 0
        );`

        if _, err := db.Exec(create_query); err != nil {
                return err
        }

        // Migrate existing DBs — add any missing columns, ignore if already present.
        migrations := []string{
                `ALTER TABLE checkins ADD COLUMN cpu_percent    REAL    NOT NULL DEFAULT 0;`,
                `ALTER TABLE checkins ADD COLUMN ram_percent    REAL    NOT NULL DEFAULT 0;`,
                `ALTER TABLE checkins ADD COLUMN disk_percent   REAL    NOT NULL DEFAULT 0;`,
                `ALTER TABLE checkins ADD COLUMN uptime_seconds INTEGER NOT NULL DEFAULT 0;`,
        }
        for _, m := range migrations {
                _, err := db.Exec(m)
                if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
                        return fmt.Errorf("failed to migrate schema: %w", err)
                }
        }

        return nil
}

// insert_checkin writes a single check-in record to the database.
func insert_checkin(db *sql.DB, hostname, logged_user string, checked_in_at time.Time, cpu_percent, ram_percent, disk_percent float64, uptime_seconds uint64) error {
        query := `INSERT INTO checkins (hostname, logged_user, checked_in_at, cpu_percent, ram_percent, disk_percent, uptime_seconds) VALUES (?, ?, ?, ?, ?, ?, ?);`

        // Store time as RFC3339 string — avoids ambiguous formats on read-back.
        _, err := db.Exec(query, hostname, logged_user, checked_in_at.UTC().Format(time.RFC3339), cpu_percent, ram_percent, disk_percent, uptime_seconds)
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
