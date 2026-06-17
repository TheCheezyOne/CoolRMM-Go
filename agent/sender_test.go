/*
        agent/sender_test.go — Tests for sender.go.
        Spins up a local test HTTP server in memory — no real network required.
        Confirms send_checkin() succeeds on 200, and errors on bad status or bad URL.
*/

package main

import (
        "net/http"
        "net/http/httptest"
        "testing"
)

// Test_send_checkin_ok confirms a 200 response from the server means success.
func Test_send_checkin_ok(t *testing.T) {
        // Spin up a test server that always responds 200 OK.
        test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
        }))
        defer test_server.Close()

        // send_checkin should return nil on success.
        if err := send_checkin(test_server.URL, "test-token"); err != nil {
                t.Fatalf("send_checkin() returned unexpected error: %v", err)
        }
}

// Test_send_checkin_sends_auth confirms the Authorization header is sent with the right token.
func Test_send_checkin_sends_auth(t *testing.T) {
        var got_header string

        // Capture the Authorization header the agent sends.
        test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                got_header = r.Header.Get("Authorization")
                w.WriteHeader(http.StatusOK)
        }))
        defer test_server.Close()

        if err := send_checkin(test_server.URL, "my-secret"); err != nil {
                t.Fatalf("send_checkin() returned unexpected error: %v", err)
        }

        // Confirm the exact header value the server received.
        if got_header != "Bearer my-secret" {
                t.Errorf("expected Authorization header 'Bearer my-secret', got '%s'", got_header)
        }
}

// Test_send_checkin_bad_status confirms a non-200 response is treated as an error.
func Test_send_checkin_bad_status(t *testing.T) {
        // Spin up a test server that always responds 500 Internal Server Error.
        test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusInternalServerError)
        }))
        defer test_server.Close()

        // send_checkin should return an error when the server is unhappy.
        if err := send_checkin(test_server.URL, "test-token"); err == nil {
                t.Fatal("send_checkin() expected an error on 500 response, got nil")
        }
}

// Test_send_checkin_unreachable confirms an unreachable server returns an error.
func Test_send_checkin_unreachable(t *testing.T) {
        // Point at a port nothing is listening on.
        if err := send_checkin("http://127.0.0.1:19999", "test-token"); err == nil {
                t.Fatal("send_checkin() expected an error for unreachable server, got nil")
        }
}

/*
        To run these tests:
          go test ./agent/...

        Expected output:
          ok    github.com/TheCheezyOne/CoolRMM-Go/agent
*/
