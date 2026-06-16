/*
        agent/system_info.go — Collects basic system information from the local machine.
        Provides hostname, logged-in username, and CPU usage %.
        These are the data points the agent gathers before phoning home.
*/

package main

import (
        "fmt"
        "os"
        "os/user"
        "time"

        "github.com/shirou/gopsutil/v3/cpu"
)

// get_hostname returns the machine's hostname or an error if it can't be read.
func get_hostname() (string, error) {
        return os.Hostname()
}

// get_logged_user returns the logged-in username.
// Checks USERNAME (Windows), then USER (Linux/macOS), then falls back to os/user.Current().
func get_logged_user() (string, error) {
        // Windows sets USERNAME; Linux/macOS sets USER.
        for _, env_key := range []string{"USERNAME", "USER"} {
                if u := os.Getenv(env_key); u != "" {
                        return u, nil
                }
        }

        // Last resort: ask the OS directly via the standard library.
        current_user, err := user.Current()
        if err != nil {
                return "", fmt.Errorf("could not determine logged-in user: %v", err)
        }
        return current_user.Username, nil
}

// get_cpu_percent returns the overall CPU usage as a percentage (0–100).
// It samples over 500ms — the call blocks for that duration by design.
func get_cpu_percent() (float64, error) {
        // cpu.Percent(interval, percpu) — false means we want one combined value, not per-core.
        percents, err := cpu.Percent(500*time.Millisecond, false)
        if err != nil {
                return 0, fmt.Errorf("could not read cpu usage: %v", err)
        }

        // Percent() returns a slice — index 0 is the combined value when percpu is false.
        if len(percents) == 0 {
                return 0, fmt.Errorf("cpu.Percent() returned no data")
        }

        return percents[0], nil
}

/*
        No build or run instructions for this file — it is a package, not a main.
        It is compiled as part of the agent binary:
          go build ./agent
          GOOS=windows GOARCH=amd64 go build -o coolrmm-agent.exe ./agent

        To run tests for this file:
          go test ./agent/...
*/
