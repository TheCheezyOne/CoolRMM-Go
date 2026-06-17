/*
        agent/commands_test.go — Tests for commands.go.
        Uses in-memory HTTP test servers to avoid real network calls.
        Covers: no pending command, pending command executed, result posted,
        run_command output capture.
*/

package main

import (
        "encoding/json"
        "net/http"
        "net/http/httptest"
        "testing"
)

// Test_poll_commands_nothing_pending confirms no error when the server returns id=0.
func Test_poll_commands_nothing_pending(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Return empty response — id=0 means nothing queued.
                json.NewEncoder(w).Encode(pending_command_response{ID: 0, Command: ""})
        }))
        defer server.Close()

        if err := poll_commands(server.URL, "test-key"); err != nil {
                t.Fatalf("poll_commands() returned error on empty queue: %v", err)
        }
}

// Test_poll_commands_runs_and_posts confirms a pending command is executed and result is posted.
func Test_poll_commands_runs_and_posts(t *testing.T) {
        var result_received command_result_payload

        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                switch r.URL.Path {
                case "/commands/pending":
                        // Return a simple echo command that works on both Windows and Linux.
                        json.NewEncoder(w).Encode(pending_command_response{ID: 42, Command: "echo hello"})
                case "/command_result":
                        // Capture the result the agent posts.
                        json.NewDecoder(r.Body).Decode(&result_received)
                        w.WriteHeader(http.StatusOK)
                }
        }))
        defer server.Close()

        if err := poll_commands(server.URL, "test-key"); err != nil {
                t.Fatalf("poll_commands() returned error: %v", err)
        }

        // Confirm the agent posted back the result for the right command.
        if result_received.ID != 42 {
                t.Errorf("expected result ID 42, got %d", result_received.ID)
        }
        if result_received.ExitCode != 0 {
                t.Errorf("expected exit code 0, got %d", result_received.ExitCode)
        }
}

// Test_run_command_captures_output confirms stdout is captured and returned.
func Test_run_command_captures_output(t *testing.T) {
        output, exit_code := run_command("echo hello")

        if exit_code != 0 {
                t.Errorf("expected exit code 0, got %d", exit_code)
        }
        if len(output) == 0 {
                t.Error("expected non-empty output from echo hello")
        }
}

// Test_run_command_bad_command confirms a non-zero exit code on a bad command.
func Test_run_command_bad_command(t *testing.T) {
        _, exit_code := run_command("commandthatdoesnotexist_xyz")

        if exit_code == 0 {
                t.Error("expected non-zero exit code for a bad command")
        }
}

// Test_post_command_result_sends_auth confirms the Bearer token is sent.
func Test_post_command_result_sends_auth(t *testing.T) {
        var got_auth string

        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                got_auth = r.Header.Get("Authorization")
                w.WriteHeader(http.StatusOK)
        }))
        defer server.Close()

        if err := post_command_result(server.URL, "my-key", 1, "output", 0); err != nil {
                t.Fatalf("post_command_result() returned error: %v", err)
        }

        if got_auth != "Bearer my-key" {
                t.Errorf("expected 'Bearer my-key', got '%s'", got_auth)
        }
}

/*
        To run these tests:
          go test ./agent/...
*/
