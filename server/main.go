/*
        server/main.go — Entry point for the CoolRMM server.
        Opens the SQLite database, starts an HTTP server, and listens for agent check-ins.
*/

package main

import (
        "fmt"
        "log"
        "net/http"
)

// version tracks the current release of the server binary.
const version = "v0.8.0"

// listen_addr is the address and port the server binds to.
const listen_addr = ":8080"

// db_path is the SQLite database file location.
const db_path = "coolrmm.db"

func main() {
        fmt.Printf("CoolRMM Server %s starting on %s...\n", version, listen_addr)

        // Load the server config — fail hard if it's missing or incomplete.
        cfg, err := load_server_config()
        if err != nil {
                log.Fatalf("failed to load server config: %v", err)
        }
        fmt.Printf("Auth: api_key loaded.\n")

        // Open (or create) the database — fail hard if it can't be opened.
        db, err := open_db(db_path)
        if err != nil {
                log.Fatalf("failed to open database: %v", err)
        }
        defer db.Close()

        fmt.Printf("Database ready: %s\n", db_path)

        // Browser-facing endpoints — no Bearer auth (internal network only).
        http.HandleFunc("/", make_dashboard_handler())
        http.HandleFunc("/devices", make_devices_handler(db))
        http.HandleFunc("/shell", make_shell_handler())
        http.HandleFunc("/commands", make_submit_command_handler(db))
        http.HandleFunc("/command_output", make_command_output_handler(db))

        // Agent-facing endpoints — protected by the shared api_key.
        http.HandleFunc("/checkin", require_auth(cfg.ApiKey, make_checkin_handler(db)))
        http.HandleFunc("/commands/pending", require_auth(cfg.ApiKey, make_pending_command_handler(db)))
        http.HandleFunc("/command_result", require_auth(cfg.ApiKey, make_command_result_handler(db)))

        // Start the server — log.Fatal so any startup error prints and exits cleanly.
        log.Fatal(http.ListenAndServe(listen_addr, nil))
}

/*
        To run the server locally (Linux or Windows):
          go run ./server

        To build for Windows from Linux:
          GOOS=windows GOARCH=amd64 go build -o coolrmm-server.exe ./server

        To run tests:
          go test ./server/...

        The server listens on port 8080 by default.
        The database file (coolrmm.db) is created in the working directory on first run.
        Agents should POST JSON to http://<server-ip>:8080/checkin
*/
