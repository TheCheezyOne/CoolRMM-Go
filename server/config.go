/*
        server/config.go — Loads server configuration from coolrmm-server.conf.
        Reads the file from the same directory as the running binary.
        Required key: api_key.
        Fails loudly if the file is missing or the key is not set.
*/

package main

import (
        "bufio"
        "fmt"
        "os"
        "path/filepath"
        "strings"
)

// server_config holds all runtime configuration for the server.
type server_config struct {
        ApiKey string
}

// load_server_config locates coolrmm-server.conf next to the binary and parses it.
func load_server_config() (server_config, error) {
        exe_path, err := os.Executable()
        if err != nil {
                return server_config{}, fmt.Errorf("could not determine executable path: %w", err)
        }
        conf_path := filepath.Join(filepath.Dir(exe_path), "coolrmm-server.conf")
        return parse_server_config_file(conf_path)
}

// parse_server_config_file reads and parses a server config file at the given path.
// Exported as its own function so tests can point it at temp files directly.
func parse_server_config_file(conf_path string) (server_config, error) {
        // Open the config file — fail clearly if it isn't there.
        file, err := os.Open(conf_path)
        if err != nil {
                return server_config{}, fmt.Errorf("server config not found at %s — create it with: api_key=<your_secret>", conf_path)
        }
        defer file.Close()

        // Parse each line looking for key=value pairs.
        cfg := server_config{}
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

                if key == "api_key" {
                        cfg.ApiKey = val
                }
        }

        if err := scanner.Err(); err != nil {
                return server_config{}, fmt.Errorf("error reading server config: %w", err)
        }

        // api_key is required — fail hard if it wasn't found.
        if cfg.ApiKey == "" {
                return server_config{}, fmt.Errorf("api_key is missing from %s", conf_path)
        }

        return cfg, nil
}

/*
        This file is compiled as part of the server binary — it is not run directly.
        To build:
          GOOS=windows GOARCH=amd64 go build -o coolrmm-server.exe ./server

        coolrmm-server.conf format (# for comments):
          # CoolRMM Server Configuration
          api_key=your_secret_here

        The api_key must match the api_key in coolrmm.conf on every agent machine.

        To run tests:
          go test ./server/...
*/
