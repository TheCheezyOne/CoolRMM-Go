/*
	server/db_test.go — Tests for db.go.
	Uses an in-memory SQLite database so tests run fast and leave no files behind.
	Confirms the schema creates cleanly and check-ins can be inserted.
*/

package main

import (
	"testing"
	"time"
)

// Test_open_db confirms the database opens and the schema is created without error.
func Test_open_db(t *testing.T) {
	// Use :memory: for an in-memory SQLite DB — no file created, no cleanup needed.
	db, err := open_db(":memory:")
	if err != nil {
		t.Fatalf("open_db() returned error: %v", err)
	}
	defer db.Close()
}

// Test_insert_checkin confirms a check-in record can be written to the database.
func Test_insert_checkin(t *testing.T) {
	// Open a fresh in-memory DB for this test.
	db, err := open_db(":memory:")
	if err != nil {
		t.Fatalf("open_db() returned error: %v", err)
	}
	defer db.Close()

	// Insert a check-in and confirm no error is returned.
	err = insert_checkin(db, "test-host", "test-user", time.Now().UTC())
	if err != nil {
		t.Fatalf("insert_checkin() returned error: %v", err)
	}

	// Confirm exactly one row was written.
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM checkins").Scan(&count); err != nil {
		t.Fatalf("failed to query checkins count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 row in checkins, got %d", count)
	}
}

/*
	To run these tests:
	  go test ./server/...

	Expected output:
	  ok  	github.com/TheCheezyOne/CoolRMM-Go/server
*/
