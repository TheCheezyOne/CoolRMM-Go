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
const version = "v1.0.0"

// check_in_interval is how often the agent phones home.
const check_in_interval = 60 * time.Second

// command_poll_interval is how often the agent checks for pending remote commands.
const command_poll_interval = 5 * time.Second

func main() {
        fmt.Printf("CoolRMM Agent %s starting...\n", version)

        // Load config from coolrmm.conf — fail hard if it's missing or malformed.
        cfg, err := load_config()
        if err != nil {
                log.Fatalf("failed to load config: %v", err)
        }

        fmt.Printf("Phoning home to %s every %s.\n", cfg.ServerURL, check_in_interval)

        // Start command poll loop in a goroutine — runs independently of check-in.
        go func() {
                poll_ticker := time.NewTicker(command_poll_interval)
                defer poll_ticker.Stop()
                for range poll_ticker.C {
                        if err := poll_commands(cfg.ServerURL, cfg.ApiKey); err != nil {
                                log.Printf("command poll failed: %v", err)
                        }
                }
        }()

        // Send the first check-in immediately — don't wait for the first tick.
        do_checkin(cfg.ServerURL, cfg.ApiKey)

        // Start the check-in ticker — blocks the main goroutine forever.
        ticker := time.NewTicker(check_in_interval)
        defer ticker.Stop()

        for range ticker.C {
                do_checkin(cfg.ServerURL, cfg.ApiKey)
        }
}

// do_checkin calls send_checkin and logs the result either way.
func do_checkin(server_url, api_key string) {
        if err := send_checkin(server_url, api_key); err != nil {
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
