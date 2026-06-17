/*
        server/commands_test.go — Tests for commands.go.
        Uses an in-memory SQLite database for isolation.
        Covers: schema creation, insert, get pending, save result, get by ID.
*/

package main

import (
        "testing"
        "time"
)

// Test_create_commands_schema confirms the commands table can be created cleanly.
func Test_create_commands_schema(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        // create_commands_schema is called in open_db — just confirm it didn't error.
        if err := create_commands_schema(db); err != nil {
                t.Fatalf("create_commands_schema() returned error: %v", err)
        }
}

// Test_insert_command confirms a new command is persisted and returns a non-zero ID.
func Test_insert_command(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        id, err := insert_command(db, "test-host", "ipconfig")
        if err != nil {
                t.Fatalf("insert_command() returned error: %v", err)
        }
        if id == 0 {
                t.Error("insert_command() returned id 0 — expected a real row ID")
        }
}

// Test_get_pending_command confirms the oldest pending command is returned.
func Test_get_pending_command(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        // Insert two commands — oldest should come back first.
        insert_command(db, "host-a", "first command")
        time.Sleep(1 * time.Millisecond)
        insert_command(db, "host-a", "second command")

        cmd, found, err := get_pending_command(db, "host-a")
        if err != nil {
                t.Fatalf("get_pending_command() returned error: %v", err)
        }
        if !found {
                t.Fatal("get_pending_command() returned found=false, expected a command")
        }
        if cmd.Command != "first command" {
                t.Errorf("expected 'first command', got '%s'", cmd.Command)
        }
}

// Test_get_pending_command_empty confirms found=false when nothing is queued.
func Test_get_pending_command_empty(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        _, found, err := get_pending_command(db, "host-a")
        if err != nil {
                t.Fatalf("get_pending_command() returned error: %v", err)
        }
        if found {
                t.Error("get_pending_command() returned found=true on an empty table")
        }
}

// Test_save_command_result confirms a command is marked complete with output stored.
func Test_save_command_result(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        id, err := insert_command(db, "host-a", "whoami")
        if err != nil {
                t.Fatalf("insert_command() failed: %v", err)
        }

        if err := save_command_result(db, id, "DESKTOP\\user", 0); err != nil {
                t.Fatalf("save_command_result() returned error: %v", err)
        }

        // Confirm the row is now marked complete.
        cmd, found, err := get_command_by_id(db, id)
        if err != nil || !found {
                t.Fatalf("get_command_by_id() failed after save: %v", err)
        }
        if cmd.Status != "complete" {
                t.Errorf("expected status 'complete', got '%s'", cmd.Status)
        }
        if cmd.Output != "DESKTOP\\user" {
                t.Errorf("expected output 'DESKTOP\\user', got '%s'", cmd.Output)
        }
}

// Test_get_command_by_id confirms retrieval by ID works and returns not-found correctly.
func Test_get_command_by_id(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        // Non-existent ID should return found=false.
        _, found, err := get_command_by_id(db, 9999)
        if err != nil {
                t.Fatalf("get_command_by_id() returned error: %v", err)
        }
        if found {
                t.Error("get_command_by_id() returned found=true for a non-existent ID")
        }
}

/*
        To run these tests:
          go test ./server/...
*/
