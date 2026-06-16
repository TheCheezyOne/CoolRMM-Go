/*
        agent/main.go — Entry point for the CoolRMM agent.
        The agent runs on a remote Windows 11 device and reports
        status back to the CoolRMM server every 60 seconds.
        Server URL is loaded from coolrmm.conf at startup.
*/

package main

import (
        "fmt"
        "log"
        "time"
)

// version tracks the current release of the agent binary.
const version = "v0.7.0"

// check_in_interval is how often the agent phones home.
const check_in_interval = 60 * time.Second

func main() {
        fmt.Printf("CoolRMM Agent %s starting...\n", version)

        // Load config from coolrmm.conf — fail hard if it's missing or malformed.
        cfg, err := load_config()
        if err != nil {
                log.Fatalf("failed to load config: %v", err)
        }

        fmt.Printf("Phoning home to %s every %s.\n", cfg.ServerURL, check_in_interval)

        // Send the first check-in immediately — don't wait for the first tick.
        do_checkin(cfg.ServerURL)

        // Start the ticker and check in on every tick.
        ticker := time.NewTicker(check_in_interval)
        defer ticker.Stop()

        for range ticker.C {
                do_checkin(cfg.ServerURL)
        }
}

// do_checkin calls send_checkin and logs the result either way.
func do_checkin(server_url string) {
        if err := send_checkin(server_url); err != nil {
                // Log the failure but keep running — a missed check-in isn't fatal.
                log.Printf("check-in failed: %v", err)
                return
        }
        log.Printf("check-in sent successfully")
}

/*
        To build for Windows from Linux:
          GOOS=windows GOARCH=amd64 go build -o coolrmm-agent.exe ./agent

        To run locally on Linux (requires server running, and coolrmm.conf present):
          go run ./agent

        coolrmm.conf must sit next to the binary:
          server_url=http://192.168.1.100:8080

        To run tests:
          go test ./agent/...

        NOTE: main() loops forever by design. do_checkin() is the testable unit —
        the send_checkin() function it calls is fully covered in sender_test.go.
*/
