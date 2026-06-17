/*
        agent/config.go — Loads agent configuration from a file on disk.
        Reads coolrmm.conf from the same directory as the running binary.
        Required keys: server_url, api_key.
        Fails loudly if the file is missing or either key is not set.
*/

package main

import (
        "bufio"
        "fmt"
        "os"
        "path/filepath"
        "strings"
)

// agent_config holds all runtime configuration for the agent.
type agent_config struct {
        ServerURL string
        ApiKey    string
}

// load_config locates coolrmm.conf next to the binary and parses it.
func load_config() (agent_config, error) {
        // Locate the directory the running binary lives in.
        exe_path, err := os.Executable()
        if err != nil {
                return agent_config{}, fmt.Errorf("could not determine executable path: %w", err)
        }
        conf_path := filepath.Join(filepath.Dir(exe_path), "coolrmm.conf")

        return parse_config_file(conf_path)
}

// parse_config_file reads and parses a config file at the given path.
// Exported as its own function so tests can point it at temp files directly.
func parse_config_file(conf_path string) (agent_config, error) {
        // Open the config file — fail clearly if it isn't there.
        file, err := os.Open(conf_path)
        if err != nil {
                return agent_config{}, fmt.Errorf("config file not found at %s — create it with: server_url=http://<ip>:8080", conf_path)
        }
        defer file.Close()

        // Parse each line looking for key=value pairs.
        cfg := agent_config{}
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                line := strings.TrimSpace(scanner.Text())

                // Skip blank lines and comments.
                if line == "" || strings.HasPrefix(line, "#") {
                        continue
                }

                // Split on the first '=' only.
                parts := strings.SplitN(line, "=", 2)
                if len(parts) != 2 {
                        continue
                }

                key := strings.TrimSpace(parts[0])
                val := strings.TrimSpace(parts[1])

                // Store recognized keys.
                switch key {
                case "server_url":
                        cfg.ServerURL = val
                case "api_key":
                        cfg.ApiKey = val
                }
        }

        if err := scanner.Err(); err != nil {
                return agent_config{}, fmt.Errorf("error reading config file: %w", err)
        }

        // Both keys are required — fail hard if either is missing.
        if cfg.ServerURL == "" {
                return agent_config{}, fmt.Errorf("server_url is missing from %s", conf_path)
        }
        if cfg.ApiKey == "" {
                return agent_config{}, fmt.Errorf("api_key is missing from %s", conf_path)
        }

        return cfg, nil
}

/*
        This file is compiled as part of the agent binary — it is not run directly.
        To build:
          go build ./agent
          GOOS=windows GOARCH=amd64 go build -o coolrmm-agent.exe ./agent

        Deployment (per machine):
          1. Copy coolrmm.conf.example → coolrmm.conf in the same folder as coolrmm-agent.exe
          2. Edit coolrmm.conf — set server_url to the real server IP and port
          3. Run coolrmm-agent.exe

        coolrmm.conf format (# for comments):
          # CoolRMM agent config
          server_url=http://192.168.1.100:8080
          api_key=your_secret_here

        To run tests:
          go test ./agent/...
*/
