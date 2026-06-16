/*
        agent/checkin.go — Defines the check-in payload the agent will send to the server.
        Contains the data structure and the function that builds it from live system info.
        Nothing is sent over the network here — this is just the shape of the data.
*/

package main

import "time"

// check_in_payload holds everything the agent reports to the server on each check-in.
type check_in_payload struct {
        Hostname    string    `json:"hostname"`
        LoggedUser  string    `json:"logged_user"`
        CheckedInAt time.Time `json:"checked_in_at"`
        CpuPercent  float64   `json:"cpu_percent"`
}

// build_payload collects live system info and returns a populated check_in_payload.
func build_payload() (check_in_payload, error) {
        // Get the machine hostname.
        hostname, err := get_hostname()
        if err != nil {
                return check_in_payload{}, err
        }

        // Get the currently logged-in user.
        logged_user, err := get_logged_user()
        if err != nil {
                return check_in_payload{}, err
        }

        // Sample CPU usage — blocks for 500ms while gopsutil takes its reading.
        cpu_percent, err := get_cpu_percent()
        if err != nil {
                return check_in_payload{}, err
        }

        // Stamp the current UTC time so the server knows when this check-in was built.
        return check_in_payload{
                Hostname:    hostname,
                LoggedUser:  logged_user,
                CheckedInAt: time.Now().UTC(),
                CpuPercent:  cpu_percent,
        }, nil
}

/*
        This file is compiled as part of the agent binary — it is not run directly.
        To build:
          go build ./agent
          GOOS=windows GOARCH=amd64 go build -o coolrmm-agent.exe ./agent

        To run tests:
          go test ./agent/...
*/
