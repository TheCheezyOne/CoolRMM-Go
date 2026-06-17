/*
        agent/commands.go — Polls the server for pending remote commands and executes them.
        Hits GET /commands/pending every 5 seconds.
        Runs any returned command via the system shell, then POSTs output back to the server.
*/

package main

import (
        "bytes"
        "encoding/json"
        "fmt"
        "net/http"
        "os/exec"
        "runtime"
)

// pending_command_response is what the server returns when a command is queued.
type pending_command_response struct {
        ID      int64  `json:"id"`
        Command string `json:"command"`
}

// command_result_payload is what the agent sends back after executing a command.
type command_result_payload struct {
        ID       int64  `json:"id"`
        Output   string `json:"output"`
        ExitCode int    `json:"exit_code"`
}

// poll_commands checks for a pending command and runs it if one exists.
// Returns nil if nothing is queued — not an error condition.
func poll_commands(server_url, api_key string) error {
        // Get this machine's hostname to identify which commands belong to us.
        hostname, err := get_hostname()
        if err != nil {
                return fmt.Errorf("could not get hostname: %w", err)
        }

        // Ask the server for the next pending command for this host.
        req, err := http.NewRequest("GET", server_url+"/commands/pending?hostname="+hostname, nil)
        if err != nil {
                return fmt.Errorf("failed to build pending request: %w", err)
        }
        req.Header.Set("Authorization", "Bearer "+api_key)

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
                return fmt.Errorf("failed to reach server: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return fmt.Errorf("server returned %d on pending poll", resp.StatusCode)
        }

        // Decode the response — id == 0 means nothing pending.
        var pending pending_command_response
        if err := json.NewDecoder(resp.Body).Decode(&pending); err != nil {
                return fmt.Errorf("failed to decode pending response: %w", err)
        }
        if pending.ID == 0 {
                return nil
        }

        // Run the command and capture combined stdout + stderr.
        output, exit_code := run_command(pending.Command)

        // Post the result back so the dashboard can display it.
        return post_command_result(server_url, api_key, pending.ID, output, exit_code)
}

// run_command executes a shell command and returns combined output and exit code.
// Uses cmd /C on Windows and sh -c elsewhere.
func run_command(command string) (string, int) {
        var cmd *exec.Cmd
        if runtime.GOOS == "windows" {
                cmd = exec.Command("cmd", "/C", command)
        } else {
                cmd = exec.Command("sh", "-c", command)
        }

        out, err := cmd.CombinedOutput()
        output := string(out)

        exit_code := 0
        if err != nil {
                if exit_err, ok := err.(*exec.ExitError); ok {
                        exit_code = exit_err.ExitCode()
                } else {
                        // Non-exit error (e.g. command not found on PATH).
                        exit_code = -1
                        output += "\n" + err.Error()
                }
        }

        return output, exit_code
}

// post_command_result sends the execution output back to the server.
func post_command_result(server_url, api_key string, id int64, output string, exit_code int) error {
        result := command_result_payload{
                ID:       id,
                Output:   output,
                ExitCode: exit_code,
        }

        body, err := json.Marshal(result)
        if err != nil {
                return fmt.Errorf("failed to marshal result: %w", err)
        }

        req, err := http.NewRequest("POST", server_url+"/command_result", bytes.NewReader(body))
        if err != nil {
                return fmt.Errorf("failed to build result request: %w", err)
        }
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+api_key)

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
                return fmt.Errorf("failed to post result: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return fmt.Errorf("server returned %d on result post", resp.StatusCode)
        }

        return nil
}

/*
        This file is compiled as part of the agent binary — it is not run directly.
        To build:
          GOOS=windows GOARCH=amd64 go build -o coolrmm-agent.exe ./agent

        To run tests:
          go test ./agent/...
*/
