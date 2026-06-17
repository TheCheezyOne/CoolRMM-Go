/*
        server/config_test.go — Tests for config.go.
        Writes temporary config files to disk to test parse_server_config_file() behavior.
        Covers: valid config, missing file, missing key, and comment/blank line handling.
*/

package main

import (
        "os"
        "path/filepath"
        "testing"
)

// write_temp_server_conf writes content to a temp file and returns its path.
func write_temp_server_conf(t *testing.T, content string) string {
        t.Helper()
        tmp, err := os.CreateTemp("", "coolrmm_server_test_*.conf")
        if err != nil {
                t.Fatalf("failed to create temp config file: %v", err)
        }
        if _, err := tmp.WriteString(content); err != nil {
                t.Fatalf("failed to write temp config: %v", err)
        }
        tmp.Close()
        return tmp.Name()
}

// Test_parse_server_config_valid confirms a well-formed config file loads correctly.
func Test_parse_server_config_valid(t *testing.T) {
        conf := write_temp_server_conf(t, "api_key=supersecret\n")
        defer os.Remove(conf)

        cfg, err := parse_server_config_file(conf)
        if err != nil {
                t.Fatalf("parse_server_config_file() returned error: %v", err)
        }
        if cfg.ApiKey != "supersecret" {
                t.Errorf("expected api_key 'supersecret', got '%s'", cfg.ApiKey)
        }
}

// Test_parse_server_config_missing_key confirms an error when api_key is absent.
func Test_parse_server_config_missing_key(t *testing.T) {
        conf := write_temp_server_conf(t, "# just a comment\n\n")
        defer os.Remove(conf)

        _, err := parse_server_config_file(conf)
        if err == nil {
                t.Fatal("parse_server_config_file() expected error for missing api_key, got nil")
        }
}

// Test_parse_server_config_comments confirms comments and blank lines are skipped.
func Test_parse_server_config_comments(t *testing.T) {
        conf := write_temp_server_conf(t, "# CoolRMM server config\n\napi_key=mykey\n")
        defer os.Remove(conf)

        cfg, err := parse_server_config_file(conf)
        if err != nil {
                t.Fatalf("parse_server_config_file() returned error: %v", err)
        }
        if cfg.ApiKey != "mykey" {
                t.Errorf("expected api_key 'mykey', got '%s'", cfg.ApiKey)
        }
}

// Test_parse_server_config_missing_file confirms a missing file returns an explicit error.
func Test_parse_server_config_missing_file(t *testing.T) {
        bogus := filepath.Join(os.TempDir(), "does_not_exist_coolrmm_server.conf")
        _, err := parse_server_config_file(bogus)
        if err == nil {
                t.Fatal("parse_server_config_file() expected error for missing file, got nil")
        }
}

/*
        To run these tests:
          go test ./server/...
*/
