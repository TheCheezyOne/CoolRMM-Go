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
const version = "v0.7.0"

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

        // Dashboard and device list are read-only — no auth required yet.
        http.HandleFunc("/", make_dashboard_handler())
        http.HandleFunc("/devices", make_devices_handler(db))

        // Agent-facing endpoints are protected by the shared api_key.
        http.HandleFunc("/checkin", require_auth(cfg.ApiKey, make_checkin_handler(db)))

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
