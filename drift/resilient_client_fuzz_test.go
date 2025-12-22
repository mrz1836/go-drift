package drift

import (
	"testing"
)

func FuzzIsRetryableStatusCode(f *testing.F) {
	// Seed corpus with important boundary values
	f.Add(0)
	f.Add(100)
	f.Add(200)
	f.Add(201)
	f.Add(299)
	f.Add(300)
	f.Add(399)
	f.Add(400)
	f.Add(407)
	f.Add(408) // Request Timeout - should retry
	f.Add(409)
	f.Add(428)
	f.Add(429) // Too Many Requests - should retry
	f.Add(430)
	f.Add(499)
	f.Add(500) // Server Error - should retry
	f.Add(501)
	f.Add(502)
	f.Add(503)
	f.Add(599) // Last 5xx - should retry
	f.Add(600)
	f.Add(-1)
	f.Add(-100)

	f.Fuzz(func(t *testing.T, code int) {
		// Should never panic
		result := isRetryableStatusCode(code)

		// Verify the result matches expected logic
		expected := code == 408 || code == 429 || (code >= 500 && code <= 599)

		if result != expected {
			t.Errorf("isRetryableStatusCode(%d) = %v, want %v", code, result, expected)
		}
	})
}
