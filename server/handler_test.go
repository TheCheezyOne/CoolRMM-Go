/*
	server/handler_test.go — Tests for handler.go.
	Simulates HTTP requests to make_checkin_handler using an in-memory DB.
	No real network or disk I/O required.
*/

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// test_db is a helper that opens a fresh in-memory DB for handler tests.
func test_db(t *testing.T) *testDB {
	t.Helper()
	db, err := open_db(":memory:")
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}
	return &testDB{db: db, t: t}
}

// testDB wraps sql.DB for use in tests with auto-cleanup.
type testDB struct {
	db interface{ Close() error }
	t  *testing.T
}

// Test_make_checkin_handler_ok confirms a valid POST returns 200 OK and writes to DB.
func Test_make_checkin_handler_ok(t *testing.T) {
	db, err := open_db(":memory:")
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}
	defer db.Close()

	// Build a valid payload.
	payload := check_in_payload{
		Hostname:    "test-machine",
		LoggedUser:  "test-user",
		CheckedInAt: time.Now().UTC(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal test payload: %v", err)
	}

	// Simulate the POST request.
	req := httptest.NewRequest(http.MethodPost, "/checkin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	make_checkin_handler(db)(rr, req)

	// Expect 200 OK.
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	// Confirm the row was written to the DB.
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM checkins").Scan(&count); err != nil {
		t.Fatalf("failed to query checkins: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 row in checkins after check-in, got %d", count)
	}
}

// Test_make_checkin_handler_wrong_method confirms a GET returns 405.
func Test_make_checkin_handler_wrong_method(t *testing.T) {
	db, err := open_db(":memory:")
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/checkin", nil)
	rr := httptest.NewRecorder()
	make_checkin_handler(db)(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

// Test_make_checkin_handler_bad_json confirms garbled JSON returns 400.
func Test_make_checkin_handler_bad_json(t *testing.T) {
	db, err := open_db(":memory:")
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}
	defer db.Close()

	req := httptest.NewRequest(http.MethodPost, "/checkin", bytes.NewBufferString("not json"))
	rr := httptest.NewRecorder()
	make_checkin_handler(db)(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

/*
	To run these tests:
	  go test ./server/...

	Expected output:
	  ok  	github.com/TheCheezyOne/CoolRMM-Go/server
*/
