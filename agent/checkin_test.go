/*
	agent/checkin_test.go — Tests for checkin.go.
	Confirms that build_payload() returns a fully populated check_in_payload
	with no empty fields and a sane timestamp.
*/

package main

import (
	"testing"
	"time"
)

// Test_build_payload confirms all fields are populated and the timestamp is recent.
func Test_build_payload(t *testing.T) {
	before := time.Now().UTC()

	payload, err := build_payload()

	// A failure here means system info collection is broken — hard fail.
	if err != nil {
		t.Fatalf("build_payload() returned error: %v", err)
	}

	// Hostname must not be empty.
	if payload.Hostname == "" {
		t.Error("build_payload() returned empty Hostname")
	}

	// Logged user must not be empty.
	if payload.LoggedUser == "" {
		t.Error("build_payload() returned empty LoggedUser")
	}

	// Timestamp must be set and fall between before and now.
	after := time.Now().UTC()
	if payload.CheckedInAt.Before(before) || payload.CheckedInAt.After(after) {
		t.Errorf("build_payload() CheckedInAt %v is outside expected range [%v, %v]",
			payload.CheckedInAt, before, after)
	}
}

/*
	To run these tests:
	  go test ./agent/...

	Expected output:
	  ok  	github.com/TheCheezyOne/CoolRMM-Go/agent
*/
