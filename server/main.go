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
const version = "v0.5.0"

// listen_addr is the address and port the server binds to.
const listen_addr = ":8080"

// db_path is the SQLite database file location.
const db_path = "coolrmm.db"

func main() {
        fmt.Printf("CoolRMM Server %s starting on %s...\n", version, listen_addr)

        // Open (or create) the database — fail hard if it can't be opened.
        db, err := open_db(db_path)
        if err != nil {
                log.Fatalf("failed to open database: %v", err)
        }
        defer db.Close()

        fmt.Printf("Database ready: %s\n", db_path)

        // Register all routes with the DB wired in where needed.
        http.HandleFunc("/", make_dashboard_handler())
        http.HandleFunc("/checkin", make_checkin_handler(db))
        http.HandleFunc("/devices", make_devices_handler(db))

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
