/*
	server/dashboard.go — HTTP handlers for the web dashboard.
	GET / serves the dashboard HTML page.
	GET /devices returns the current device list as JSON for the page to consume.
*/

package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
)

// dashboard_html is the compiled-in dashboard page — served from the binary itself.
//
//go:embed dashboard.html
var dashboard_html []byte

// make_dashboard_handler serves the HTML dashboard at GET /.
func make_dashboard_handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Reject anything that isn't exactly GET /.
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(dashboard_html)
	}
}

// make_devices_handler returns the latest device status list as JSON at GET /devices.
func make_devices_handler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		devices, err := get_latest_devices(db)
		if err != nil {
			log.Printf("failed to get devices: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// Return an empty array rather than null when no devices have checked in yet.
		if devices == nil {
			devices = []device_status{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(devices); err != nil {
			log.Printf("failed to encode devices: %v", err)
		}
	}
}

/*
	This file is compiled as part of the server binary — it is not run directly.
	dashboard.html is embedded at compile time via go:embed.
	To build:
	  go build ./server
	  GOOS=windows GOARCH=amd64 go build -o coolrmm-server.exe ./server

	Routes added:
	  GET /         — dashboard HTML page
	  GET /devices  — device status JSON

	To run tests:
	  go test ./server/...
*/
