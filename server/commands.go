/*
        server/commands.go — Database layer for remote command execution.
        Manages the commands table: inserting new commands, fetching pending ones,
        saving results, and retrieving by ID.
*/

package main

import (
        "database/sql"
        "fmt"
        "time"
)

// command_row represents a single row in the commands table.
type command_row struct {
        ID        int64
        Hostname  string
        Command   string
        Status    string
        Output    string
        ExitCode  int
        CreatedAt time.Time
}

// create_commands_schema creates the commands table if it doesn't exist.
func create_commands_schema(db *sql.DB) error {
        query := `
        CREATE TABLE IF NOT EXISTS commands (
                id           INTEGER PRIMARY KEY AUTOINCREMENT,
                hostname     TEXT NOT NULL,
                command      TEXT NOT NULL,
                status       TEXT NOT NULL DEFAULT 'pending',
                output       TEXT NOT NULL DEFAULT '',
                exit_code    INTEGER NOT NULL DEFAULT -1,
                created_at   DATETIME NOT NULL,
                completed_at DATETIME
        );`
        if _, err := db.Exec(query); err != nil {
                return fmt.Errorf("failed to create commands table: %w", err)
        }
        return nil
}

// insert_command adds a new pending command for the given hostname.
// Returns the new row's ID.
func insert_command(db *sql.DB, hostname, command string) (int64, error) {
        query := `INSERT INTO commands (hostname, command, status, created_at) VALUES (?, ?, 'pending', ?);`
        result, err := db.Exec(query, hostname, command, time.Now().UTC().Format(time.RFC3339))
        if err != nil {
                return 0, fmt.Errorf("failed to insert command: %w", err)
        }
        id, err := result.LastInsertId()
        if err != nil {
                return 0, fmt.Errorf("failed to get command id: %w", err)
        }
        return id, nil
}

// get_pending_command returns the oldest pending command for the given hostname.
// found=false means nothing is queued — not an error.
func get_pending_command(db *sql.DB, hostname string) (command_row, bool, error) {
        query := `SELECT id, hostname, command FROM commands WHERE hostname=? AND status='pending' ORDER BY created_at ASC LIMIT 1;`
        row := db.QueryRow(query, hostname)

        var c command_row
        if err := row.Scan(&c.ID, &c.Hostname, &c.Command); err != nil {
                if err == sql.ErrNoRows {
                        return command_row{}, false, nil
                }
                return command_row{}, false, fmt.Errorf("failed to scan command: %w", err)
        }
        return c, true, nil
}

// save_command_result marks a command complete and stores its output and exit code.
func save_command_result(db *sql.DB, id int64, output string, exit_code int) error {
        query := `UPDATE commands SET status='complete', output=?, exit_code=?, completed_at=? WHERE id=?;`
        _, err := db.Exec(query, output, exit_code, time.Now().UTC().Format(time.RFC3339), id)
        if err != nil {
                return fmt.Errorf("failed to save command result: %w", err)
        }
        return nil
}

// get_command_by_id retrieves a command row by its ID.
// found=false means the ID doesn't exist.
func get_command_by_id(db *sql.DB, id int64) (command_row, bool, error) {
        query := `SELECT id, hostname, command, status, output, exit_code FROM commands WHERE id=?;`
        row := db.QueryRow(query, id)

        var c command_row
        if err := row.Scan(&c.ID, &c.Hostname, &c.Command, &c.Status, &c.Output, &c.ExitCode); err != nil {
                if err == sql.ErrNoRows {
                        return command_row{}, false, nil
                }
                return command_row{}, false, fmt.Errorf("failed to scan command by id: %w", err)
        }
        return c, true, nil
}

/*
        This file is compiled as part of the server binary — it is not run directly.
        To run tests:
          go test ./server/...
*/
