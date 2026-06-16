/*
	agent/sender.go — Sends the check-in payload to the CoolRMM server over HTTP.
	Builds the payload, encodes it as JSON, and POSTs it to the server's /checkin endpoint.
	Returns an error if anything in that chain fails — no silent failures.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// send_checkin builds a payload and POSTs it as JSON to the given server URL.
func send_checkin(server_url string) error {
	// Build the payload from live system info.
	payload, err := build_payload()
	if err != nil {
		return fmt.Errorf("failed to build payload: %w", err)
	}

	// Encode the payload to JSON.
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// POST the JSON to the server's /checkin endpoint.
	resp, err := http.Post(server_url+"/checkin", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to reach server: %w", err)
	}
	defer resp.Body.Close()

	// Anything other than 200 OK is a server-side problem — report it.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned unexpected status: %d", resp.StatusCode)
	}

	return nil
}

/*
	This file is compiled as part of the agent binary — it is not run directly.
	To build:
	  go build ./agent
	  GOOS=windows GOARCH=amd64 go build -o coolrmm-agent.exe ./agent

	To run tests:
	  go test ./agent/...
*/
