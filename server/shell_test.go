/*
        server/shell_test.go — Tests for shell.go HTTP handlers.
        Uses in-memory SQLite and httptest — no real network or file I/O.
        Covers: submit command, pending poll, command result, output poll.
*/

package main

import (
        "bytes"
        "encoding/json"
        "net/http"
        "net/http/httptest"
        "testing"
)

// Test_make_submit_command_handler_ok confirms a POST queues a command and returns an ID.
func Test_make_submit_command_handler_ok(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        payload, _ := json.Marshal(submit_command_payload{Hostname: "host-a", Command: "ipconfig"})
        req := httptest.NewRequest(http.MethodPost, "/commands", bytes.NewReader(payload))
        req.Header.Set("Content-Type", "application/json")
        rec := httptest.NewRecorder()

        make_submit_command_handler(db)(rec, req)

        if rec.Code != http.StatusOK {
                t.Fatalf("expected 200, got %d", rec.Code)
        }
        var resp submit_command_response
        json.NewDecoder(rec.Body).Decode(&resp)
        if resp.ID == 0 {
                t.Error("expected non-zero command ID in response")
        }
}

// Test_make_submit_command_handler_wrong_method confirms non-POST returns 405.
func Test_make_submit_command_handler_wrong_method(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        req := httptest.NewRequest(http.MethodGet, "/commands", nil)
        rec := httptest.NewRecorder()

        make_submit_command_handler(db)(rec, req)

        if rec.Code != http.StatusMethodNotAllowed {
                t.Errorf("expected 405, got %d", rec.Code)
        }
}

// Test_make_pending_command_handler_empty confirms id=0 returned when nothing is queued.
func Test_make_pending_command_handler_empty(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        req := httptest.NewRequest(http.MethodGet, "/commands/pending?hostname=host-a", nil)
        rec := httptest.NewRecorder()

        make_pending_command_handler(db)(rec, req)

        if rec.Code != http.StatusOK {
                t.Fatalf("expected 200, got %d", rec.Code)
        }
        var resp pending_command_response
        json.NewDecoder(rec.Body).Decode(&resp)
        if resp.ID != 0 {
                t.Errorf("expected id=0 for empty queue, got %d", resp.ID)
        }
}

// Test_make_pending_command_handler_found confirms a queued command is returned.
func Test_make_pending_command_handler_found(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        insert_command(db, "host-a", "whoami")

        req := httptest.NewRequest(http.MethodGet, "/commands/pending?hostname=host-a", nil)
        rec := httptest.NewRecorder()

        make_pending_command_handler(db)(rec, req)

        if rec.Code != http.StatusOK {
                t.Fatalf("expected 200, got %d", rec.Code)
        }
        var resp pending_command_response
        json.NewDecoder(rec.Body).Decode(&resp)
        if resp.ID == 0 {
                t.Error("expected a command ID, got 0")
        }
        if resp.Command != "whoami" {
                t.Errorf("expected command 'whoami', got '%s'", resp.Command)
        }
}

// Test_make_command_result_handler_ok confirms an agent result is saved correctly.
func Test_make_command_result_handler_ok(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        id, _ := insert_command(db, "host-a", "whoami")

        payload, _ := json.Marshal(command_result_payload{ID: id, Output: "DESKTOP\\user", ExitCode: 0})
        req := httptest.NewRequest(http.MethodPost, "/command_result", bytes.NewReader(payload))
        req.Header.Set("Content-Type", "application/json")
        rec := httptest.NewRecorder()

        make_command_result_handler(db)(rec, req)

        if rec.Code != http.StatusOK {
                t.Errorf("expected 200, got %d", rec.Code)
        }
}

// Test_make_command_output_handler_ok confirms polling returns the stored result.
func Test_make_command_output_handler_ok(t *testing.T) {
        db, err := open_db(":memory:")
        if err != nil {
                t.Fatalf("open_db() failed: %v", err)
        }
        defer db.Close()

        id, _ := insert_command(db, "host-a", "ipconfig")
        save_command_result(db, id, "192.168.1.1", 0)

        req := httptest.NewRequest(http.MethodGet, "/command_output?id="+string(rune('0'+id)), nil)
        rec := httptest.NewRecorder()

        make_command_output_handler(db)(rec, req)

        if rec.Code != http.StatusOK {
                t.Fatalf("expected 200, got %d", rec.Code)
        }
        var resp command_output_response
        json.NewDecoder(rec.Body).Decode(&resp)
        if resp.Status != "complete" {
                t.Errorf("expected status 'complete', got '%s'", resp.Status)
        }
}

/*
        To run these tests:
          go test ./server/...
*/
