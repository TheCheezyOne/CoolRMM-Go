/*
	agent/config_test.go — Tests for config.go.
	Writes temporary config files to disk to test load_config() behavior.
	Covers: valid config, missing file, missing key, and comment/blank line handling.
*/

package main

import (
	"os"
	"path/filepath"
	"testing"
)

// write_temp_conf writes content to a temp file and returns its path.
func write_temp_conf(t *testing.T, content string) string {
	t.Helper()
	tmp, err := os.CreateTemp("", "coolrmm_test_*.conf")
	if err != nil {
		t.Fatalf("failed to create temp config file: %v", err)
	}
	if _, err := tmp.WriteString(content); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	tmp.Close()
	return tmp.Name()
}

// load_config_from reads a config from an explicit path — used only in tests.
func load_config_from(conf_path string) (agent_config, error) {
	return parse_config_file(conf_path)
}

// Test_parse_config_valid confirms a well-formed config file loads correctly.
func Test_parse_config_valid(t *testing.T) {
	conf := write_temp_conf(t, "server_url=http://192.168.1.100:8080\n")
	defer os.Remove(conf)

	cfg, err := parse_config_file(conf)
	if err != nil {
		t.Fatalf("parse_config_file() returned error: %v", err)
	}
	if cfg.ServerURL != "http://192.168.1.100:8080" {
		t.Errorf("expected server_url 'http://192.168.1.100:8080', got '%s'", cfg.ServerURL)
	}
}

// Test_parse_config_missing_key confirms an error when server_url is absent.
func Test_parse_config_missing_key(t *testing.T) {
	conf := write_temp_conf(t, "# just a comment\n\n")
	defer os.Remove(conf)

	_, err := parse_config_file(conf)
	if err == nil {
		t.Fatal("parse_config_file() expected error for missing server_url, got nil")
	}
}

// Test_parse_config_comments confirms comment and blank lines are skipped cleanly.
func Test_parse_config_comments(t *testing.T) {
	conf := write_temp_conf(t, "# CoolRMM config\n\nserver_url=http://10.0.0.1:8080\n")
	defer os.Remove(conf)

	cfg, err := parse_config_file(conf)
	if err != nil {
		t.Fatalf("parse_config_file() returned error: %v", err)
	}
	if cfg.ServerURL != "http://10.0.0.1:8080" {
		t.Errorf("expected 'http://10.0.0.1:8080', got '%s'", cfg.ServerURL)
	}
}

// Test_parse_config_missing_file confirms a missing file returns an explicit error.
func Test_parse_config_missing_file(t *testing.T) {
	bogus := filepath.Join(os.TempDir(), "does_not_exist_coolrmm.conf")
	_, err := parse_config_file(bogus)
	if err == nil {
		t.Fatal("parse_config_file() expected error for missing file, got nil")
	}
}

/*
	To run these tests:
	  go test ./agent/...

	Expected output:
	  ok  	github.com/TheCheezyOne/CoolRMM-Go/agent
*/
