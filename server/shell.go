/*
        server/shell.go — HTTP handlers for the remote shell feature.

        Endpoints:
          GET  /shell              — serves the terminal UI page (browser, no auth)
          POST /commands           — queues a command for a device (browser, no auth)
          GET  /commands/pending   — agent polls for its next command (Bearer auth via main.go)
          POST /command_result     — agent posts execution output (Bearer auth via main.go)
          GET  /command_output     — browser polls for a command result (no auth)
*/

package main

import (
        "database/sql"
        "encoding/json"
        "log"
        "net/http"
        "strconv"
)

// submit_command_payload is what the shell page POSTs to queue a command.
type submit_command_payload struct {
        Hostname string `json:"hostname"`
        Command  string `json:"command"`
}

// submit_command_response carries the new command's ID back to the browser.
type submit_command_response struct {
        ID int64 `json:"id"`
}

// pending_command_response is what the agent receives when polling for work.
type pending_command_response struct {
        ID      int64  `json:"id"`
        Command string `json:"command"`
}

// command_result_payload is what the agent POSTs after running a command.
type command_result_payload struct {
        ID       int64  `json:"id"`
        Output   string `json:"output"`
        ExitCode int    `json:"exit_code"`
}

// command_output_response is what the browser receives when polling for output.
type command_output_response struct {
        ID       int64  `json:"id"`
        Status   string `json:"status"`
        Output   string `json:"output"`
        ExitCode int    `json:"exit_code"`
}

// make_shell_handler serves the terminal HTML page.
func make_shell_handler() http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodGet {
                        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
                        return
                }
                http.ServeFile(w, r, "shell.html")
        }
}

// make_submit_command_handler queues a command for a given device.
// Called by the browser — no Bearer auth (internal network only).
func make_submit_command_handler(db *sql.DB) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodPost {
                        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
                        return
                }

                var payload submit_command_payload
                if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                        http.Error(w, "bad request", http.StatusBadRequest)
                        return
                }

                if payload.Hostname == "" || payload.Command == "" {
                        http.Error(w, "hostname and command are required", http.StatusBadRequest)
                        return
                }

                id, err := insert_command(db, payload.Hostname, payload.Command)
                if err != nil {
                        log.Printf("failed to insert command: %v", err)
                        http.Error(w, "internal server error", http.StatusInternalServerError)
                        return
                }

                w.Header().Set("Content-Type", "application/json")
                json.NewEncoder(w).Encode(submit_command_response{ID: id})
        }
}

// make_pending_command_handler returns the oldest pending command for the agent's hostname.
// Returns id=0 when nothing is queued — not an error.
// Protected by Bearer auth in main.go.
func make_pending_command_handler(db *sql.DB) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodGet {
                        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
                        return
                }

                hostname := r.URL.Query().Get("hostname")
                if hostname == "" {
                        http.Error(w, "hostname is required", http.StatusBadRequest)
                        return
                }

                cmd, found, err := get_pending_command(db, hostname)
                if err != nil {
                        log.Printf("failed to get pending command: %v", err)
                        http.Error(w, "internal server error", http.StatusInternalServerError)
                        return
                }

                w.Header().Set("Content-Type", "application/json")
                if !found {
                        // Return id=0 to signal nothing pending — agent ignores it.
                        json.NewEncoder(w).Encode(pending_command_response{ID: 0, Command: ""})
                        return
                }

                json.NewEncoder(w).Encode(pending_command_response{ID: cmd.ID, Command: cmd.Command})
        }
}

// make_command_result_handler receives execution output from the agent.
// Protected by Bearer auth in main.go.
func make_command_result_handler(db *sql.DB) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodPost {
                        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
                        return
                }

                var payload command_result_payload
                if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                        http.Error(w, "bad request", http.StatusBadRequest)
                        return
                }

                if err := save_command_result(db, payload.ID, payload.Output, payload.ExitCode); err != nil {
                        log.Printf("failed to save command result: %v", err)
                        http.Error(w, "internal server error", http.StatusInternalServerError)
                        return
                }

                w.WriteHeader(http.StatusOK)
        }
}

// make_command_output_handler lets the browser poll for a command's result.
// Returns the current status and output for the given command ID.
func make_command_output_handler(db *sql.DB) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodGet {
                        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
                        return
                }

                id_str := r.URL.Query().Get("id")
                id, err := strconv.ParseInt(id_str, 10, 64)
                if err != nil || id == 0 {
                        http.Error(w, "valid id is required", http.StatusBadRequest)
                        return
                }

                cmd, found, err := get_command_by_id(db, id)
                if err != nil {
                        log.Printf("failed to get command by id: %v", err)
                        http.Error(w, "internal server error", http.StatusInternalServerError)
                        return
                }
                if !found {
                        http.Error(w, "not found", http.StatusNotFound)
                        return
                }

                w.Header().Set("Content-Type", "application/json")
                json.NewEncoder(w).Encode(command_output_response{
                        ID:       cmd.ID,
                        Status:   cmd.Status,
                        Output:   cmd.Output,
                        ExitCode: cmd.ExitCode,
                })
        }
}

/*
        This file is compiled as part of the server binary — it is not run directly.
        To run tests:
          go test ./server/...
*/
