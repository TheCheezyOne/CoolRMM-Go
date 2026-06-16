/*
	server/devices_test.go — Tests for devices.go.
	Confirms status logic and DB query behavior using an in-memory SQLite database.
*/

package main

import (
	"testing"
	"time"
)

// Test_compute_status confirms each threshold returns the correct dot color.
func Test_compute_status(t *testing.T) {
	// Green: just checked in.
	if s := compute_status(time.Now().UTC()); s != "green" {
		t.Errorf("expected green for recent check-in, got %s", s)
	}

	// Yellow: 5 minutes ago.
	if s := compute_status(time.Now().UTC().Add(-5 * time.Minute)); s != "yellow" {
		t.Errorf("expected yellow for 5m ago, got %s", s)
	}

	// Red: 20 minutes ago.
	if s := compute_status(time.Now().UTC().Add(-20 * time.Minute)); s != "red" {
		t.Errorf("expected red for 20m ago, got %s", s)
	}
}

// Test_get_latest_devices_empty confirms an empty DB returns an empty slice without error.
func Test_get_latest_devices_empty(t *testing.T) {
	db, err := open_db(":memory:")
	if err != nil {
		t.Fatalf("open_db() failed: %v", err)
	}
	defer db.Close()

	devices, err := get_latest_devices(db)
	if err != nil {
		t.Fatalf("get_latest_devices() returned error: %v", err)
	}

	if len(devices) != 0 {
		t.Errorf("expected 0 devices from empty DB, got %d", len(devices))
	}
}

// Test_get_latest_devices_dedup confirms only the most recent check-in per hostname is returned.
func Test_get_latest_devices_dedup(t *testing.T) {
	db, err := open_db(":memory:")
	if err != nil {
		t.Fatalf("open_db() failed: %v", err)
	}
	defer db.Close()

	// Insert two check-ins for the same host — only the latest should come back.
	old_time := time.Now().UTC().Add(-5 * time.Minute)
	new_time := time.Now().UTC()

	if err := insert_checkin(db, "host-a", "user-a", old_time, 10.0); err != nil {
		t.Fatalf("insert_checkin() failed: %v", err)
	}
	if err := insert_checkin(db, "host-a", "user-a", new_time, 20.0); err != nil {
		t.Fatalf("insert_checkin() failed: %v", err)
	}

	devices, err := get_latest_devices(db)
	if err != nil {
		t.Fatalf("get_latest_devices() returned error: %v", err)
	}

	// Should be exactly one device.
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}

	// CPU should be from the newer check-in.
	if devices[0].CpuPercent != 20.0 {
		t.Errorf("expected cpu_percent 20.0 from latest check-in, got %.1f", devices[0].CpuPercent)
	}
}

/*
	To run these tests:
	  go test ./server/...
*/
