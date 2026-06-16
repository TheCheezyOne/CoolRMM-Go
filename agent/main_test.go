/*
	agent/main_test.go — Tests for the CoolRMM agent entry point.
	Confirms that the version constant is defined and non-empty.
*/

package main

import "testing"

// test_version_defined confirms the version constant is set and not blank.
func Test_version_defined(t *testing.T) {
	// Fail loudly if someone blanks out the version string by mistake.
	if version == "" {
		t.Fatal("version constant must not be empty")
	}
}

/*
	To run these tests:
	  go test ./agent/...

	Expected output:
	  ok  	github.com/TheCheezyOne/CoolRMM-Go/agent
*/
