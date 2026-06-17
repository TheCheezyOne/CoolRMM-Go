/*
        server/auth_test.go — Tests for auth.go.
        Confirms require_auth() passes valid requests and rejects invalid ones.
*/

package main

import (
        "net/http"
        "net/http/httptest"
        "testing"
)

// Test_require_auth_valid confirms a request with the correct token passes through.
func Test_require_auth_valid(t *testing.T) {
        handler := require_auth("correct-key", func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
        })

        req := httptest.NewRequest(http.MethodPost, "/checkin", nil)
        req.Header.Set("Authorization", "Bearer correct-key")
        rec := httptest.NewRecorder()

        handler(rec, req)

        if rec.Code != http.StatusOK {
                t.Errorf("expected 200 for valid token, got %d", rec.Code)
        }
}

// Test_require_auth_wrong_key confirms a wrong token gets a 401.
func Test_require_auth_wrong_key(t *testing.T) {
        handler := require_auth("correct-key", func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
        })

        req := httptest.NewRequest(http.MethodPost, "/checkin", nil)
        req.Header.Set("Authorization", "Bearer wrong-key")
        rec := httptest.NewRecorder()

        handler(rec, req)

        if rec.Code != http.StatusUnauthorized {
                t.Errorf("expected 401 for wrong token, got %d", rec.Code)
        }
}

// Test_require_auth_missing_header confirms a missing Authorization header gets a 401.
func Test_require_auth_missing_header(t *testing.T) {
        handler := require_auth("correct-key", func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
        })

        req := httptest.NewRequest(http.MethodPost, "/checkin", nil)
        rec := httptest.NewRecorder()

        handler(rec, req)

        if rec.Code != http.StatusUnauthorized {
                t.Errorf("expected 401 for missing header, got %d", rec.Code)
        }
}

/*
        To run these tests:
          go test ./server/...
*/
