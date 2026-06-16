/*
        server/devices.go — Queries the database for the latest check-in per device
        and computes each device's status (green/yellow/red) based on how recently
        it phoned home.
*/

package main

import (
        "database/sql"
        "fmt"
        "time"
)

// sqlite_time_formats lists the datetime string formats SQLite may return.
// RFC3339 is the preferred format going forward; the others handle legacy rows.
var sqlite_time_formats = []string{
        "2006-01-02T15:04:05Z07:00",
        "2006-01-02T15:04:05Z",
        "2006-01-02 15:04:05+00:00",
        "2006-01-02 15:04:05",
        "2006-01-02 15:04:05.999999999 -0700 MST", // Go time.Time.String() format
}

// parse_sqlite_time tries each known SQLite datetime format until one succeeds.
func parse_sqlite_time(s string) (time.Time, error) {
        for _, layout := range sqlite_time_formats {
                if t, err := time.Parse(layout, s); err == nil {
                        return t, nil
                }
        }
        return time.Time{}, fmt.Errorf("could not parse datetime: %q", s)
}

// device_status holds the most recent check-in data for a single device,
// plus its computed dot color.
type device_status struct {
        Hostname    string    `json:"hostname"`
        LoggedUser  string    `json:"logged_user"`
        CpuPercent  float64   `json:"cpu_percent"`
        RamPercent  float64   `json:"ram_percent"`
        DiskPercent float64   `json:"disk_percent"`
        LastSeen    time.Time `json:"last_seen"`
        Status      string    `json:"status"` // "green", "yellow", or "red"
}

// Status thresholds — time since last check-in before the dot changes color.
const green_threshold = 2 * time.Minute
const yellow_threshold = 10 * time.Minute

// compute_status returns "green", "yellow", or "red" based on age of last check-in.
func compute_status(last_seen time.Time) string {
        age := time.Since(last_seen)
        if age <= green_threshold {
                return "green"
        }
        if age <= yellow_threshold {
                return "yellow"
        }
        return "red"
}

// get_latest_devices returns the most recent check-in for each unique hostname.
func get_latest_devices(db *sql.DB) ([]device_status, error) {
        query := `
        SELECT hostname, logged_user, cpu_percent, ram_percent, disk_percent, MAX(checked_in_at) AS last_seen
        FROM checkins
        GROUP BY hostname
        ORDER BY hostname;`

        rows, err := db.Query(query)
        if err != nil {
                return nil, fmt.Errorf("failed to query devices: %w", err)
        }
        defer rows.Close()

        var devices []device_status
        for rows.Next() {
                var d device_status
                var last_seen_str string
                if err := rows.Scan(&d.Hostname, &d.LoggedUser, &d.CpuPercent, &d.RamPercent, &d.DiskPercent, &last_seen_str); err != nil {
                        return nil, fmt.Errorf("failed to scan device row: %w", err)
                }
                // SQLite returns datetimes as strings — parse it manually.
                t, err := parse_sqlite_time(last_seen_str)
                if err != nil {
                        return nil, fmt.Errorf("failed to parse last_seen for %s: %w", d.Hostname, err)
                }
                d.LastSeen = t
                d.Status = compute_status(d.LastSeen)
                devices = append(devices, d)
        }

        if err := rows.Err(); err != nil {
                return nil, fmt.Errorf("error iterating device rows: %w", err)
        }

        return devices, nil
}

/*
        This file is compiled as part of the server binary — it is not run directly.
        To build:
          go build ./server
          GOOS=windows GOARCH=amd64 go build -o coolrmm-server.exe ./server

        Status thresholds:
          green  = checked in within the last 2 minutes
          yellow = 2–10 minutes ago
          red    = more than 10 minutes ago

        To run tests:
          go test ./server/...
*/
