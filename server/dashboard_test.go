/*
	server/dashboard_test.go — Tests for dashboard.go.
	Confirms the dashboard HTML and devices JSON endpoints respond correctly.
*/

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test_make_dashboard_handler_ok confirms GET / returns 200 with HTML content.
func Test_make_dashboard_handler_ok(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	make_dashboard_handler()(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	ct := rr.Header().Get("Content-Type")
	if ct != "text/html; charset=utf-8" {
		t.Errorf("expected HTML content-type, got %s", ct)
	}
}

// Test_make_dashboard_handler_not_found confirms unknown paths return 404.
func Test_make_dashboard_handler_not_found(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rr := httptest.NewRecorder()

	make_dashboard_handler()(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

// Test_make_devices_handler_ok confirms GET /devices returns valid JSON.
func Test_make_devices_handler_ok(t *testing.T) {
	db, err := open_db(":memory:")
	if err != nil {
		t.Fatalf("open_db() failed: %v", err)
	}
	defer db.Close()

	// Insert one device so the response isn't empty.
	if err := insert_checkin(db, "test-host", "test-user", time.Now().UTC(), 33.3); err != nil {
		t.Fatalf("insert_checkin() failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	rr := httptest.NewRecorder()

	make_devices_handler(db)(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	// Response must be valid JSON containing at least one device.
	var devices []device_status
	if err := json.Unmarshal(rr.Body.Bytes(), &devices); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	if len(devices) != 1 {
		t.Errorf("expected 1 device in response, got %d", len(devices))
	}
}

// Test_make_devices_handler_wrong_method confirms non-GET returns 405.
func Test_make_devices_handler_wrong_method(t *testing.T) {
	db, err := open_db(":memory:")
	if err != nil {
		t.Fatalf("open_db() failed: %v", err)
	}
	defer db.Close()

	req := httptest.NewRequest(http.MethodPost, "/devices", nil)
	rr := httptest.NewRecorder()

	make_devices_handler(db)(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

/*
	To run these tests:
	  go test ./server/...
*/
