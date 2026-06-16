/*
        agent/system_info_test.go — Tests for system_info.go.
        Confirms that hostname, logged-in user, and CPU usage can be retrieved successfully.
        These tests run on Linux (dev) but the functions also target Windows 11 (prod).
*/

package main

import "testing"

// Test_get_hostname confirms the hostname is non-empty and retrievable.
func Test_get_hostname(t *testing.T) {
        hostname, err := get_hostname()

        // Any error here means the OS couldn't provide the hostname — that's a hard fail.
        if err != nil {
                t.Fatalf("get_hostname() returned error: %v", err)
        }

        // A blank hostname is also wrong — fail loudly.
        if hostname == "" {
                t.Fatal("get_hostname() returned an empty string")
        }
}

// Test_get_logged_user confirms a username is non-empty and retrievable.
func Test_get_logged_user(t *testing.T) {
        user, err := get_logged_user()

        // Any error means neither USERNAME nor USER was set — that's unexpected in any normal environment.
        if err != nil {
                t.Fatalf("get_logged_user() returned error: %v", err)
        }

        // A blank username is also wrong — fail loudly.
        if user == "" {
                t.Fatal("get_logged_user() returned an empty string")
        }
}

// Test_get_cpu_percent confirms CPU usage is readable and within a valid range.
func Test_get_cpu_percent(t *testing.T) {
        percent, err := get_cpu_percent()

        // Any error here means gopsutil couldn't sample the CPU — hard fail.
        if err != nil {
                t.Fatalf("get_cpu_percent() returned error: %v", err)
        }

        // A valid CPU reading is always between 0 and 100 inclusive.
        if percent < 0 || percent > 100 {
                t.Errorf("get_cpu_percent() returned out-of-range value: %v", percent)
        }
}

/*
        To run these tests:
          go test ./agent/...

        Expected output on Linux:
          ok    github.com/TheCheezyOne/CoolRMM-Go/agent

        NOTE: On Windows, USERNAME is set by the OS. On Linux, USER is used instead.
        Both paths are exercised by the same test — the environment determines which branch runs.

        NOTE: Test_get_cpu_percent blocks for ~500ms while gopsutil takes its sample.
        This is expected behavior — the sampling interval is intentional.
*/
