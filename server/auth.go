/*
        server/auth.go — HTTP middleware that enforces the shared API key.
        All agent-facing endpoints are wrapped with require_auth() at startup.
        Requests missing or carrying the wrong Authorization header are rejected with 401.
        The dashboard (GET /) and device list (GET /devices) are read-only and not wrapped —
        they will be secured separately when the command feature is built.
*/

package main

import "net/http"

// require_auth wraps a handler and rejects requests without the correct Bearer token.
func require_auth(api_key string, next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                // Extract and compare the Authorization header — no hints on mismatch.
                if r.Header.Get("Authorization") != "Bearer "+api_key {
                        http.Error(w, "unauthorized", http.StatusUnauthorized)
                        return
                }
                next(w, r)
        }
}

/*
        This file is compiled as part of the server binary — it is not run directly.
        To run tests:
          go test ./server/...
*/
