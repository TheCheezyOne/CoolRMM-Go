/*
        server/handler.go — HTTP handler for the /checkin endpoint.
        Decodes the JSON payload sent by an agent, persists it to the DB,
        logs it to the console, and responds with 200 OK.
*/

package main

import (
        "database/sql"
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "time"
)

// check_in_payload mirrors the struct the agent sends on each check-in.
type check_in_payload struct {
        Hostname      string    `json:"hostname"`
        LoggedUser    string    `json:"logged_user"`
        CheckedInAt   time.Time `json:"checked_in_at"`
        CpuPercent    float64   `json:"cpu_percent"`
        RamPercent    float64   `json:"ram_percent"`
        DiskPercent   float64   `json:"disk_percent"`
        UptimeSeconds uint64    `json:"uptime_seconds"`
}

// make_checkin_handler returns an http.HandlerFunc with the DB wired in via closure.
func make_checkin_handler(db *sql.DB) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                // Only allow POST — reject anything else with 405 Method Not Allowed.
                if r.Method != http.MethodPost {
                        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
                        return
                }

                // Decode the incoming JSON body into our payload struct.
                var payload check_in_payload
                if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                        log.Printf("failed to decode check-in payload: %v", err)
                        http.Error(w, "bad request", http.StatusBadRequest)
                        return
                }

                // Persist the check-in to the database.
                if err := insert_checkin(db, payload.Hostname, payload.LoggedUser, payload.CheckedInAt, payload.CpuPercent, payload.RamPercent, payload.DiskPercent, payload.UptimeSeconds); err != nil {
                        log.Printf("failed to save check-in: %v", err)
                        http.Error(w, "internal server error", http.StatusInternalServerError)
                        return
                }

                // Log it to the console so we can watch check-ins arrive in real time.
                fmt.Printf("[check-in] host=%s user=%s cpu=%.1f%% ram=%.1f%% disk=%.1f%% uptime=%ds time=%s\n",
                        payload.Hostname,
                        payload.LoggedUser,
                        payload.CpuPercent,
                        payload.RamPercent,
                        payload.DiskPercent,
                        payload.UptimeSeconds,
                        payload.CheckedInAt.Format(time.RFC3339),
                )

                // Respond with 200 OK — the agent just needs to know we got it.
                w.WriteHeader(http.StatusOK)
        }
}

/*
        This file is compiled as part of the server binary — it is not run directly.
        To build:
          go build ./server
          GOOS=windows GOARCH=amd64 go build -o coolrmm-server.exe ./server

        To run tests:
          go test ./server/...
*/
