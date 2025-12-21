package drift

import (
	"strings"
	"testing"
)

func FuzzContactQueryBuildURL(f *testing.F) {
	// Seed corpus with various string inputs
	f.Add("test@email.com", "", "", 0)
	f.Add("", "123", "", 0)
	f.Add("", "", "ext-id", 0)
	f.Add("user@domain.com", "", "", 10)
	f.Add("", "456789", "", 5)
	f.Add("", "", "external-123", 1)
	f.Add("", "", "", 0) // All empty - should error
	f.Add("special+chars@email.com", "", "", 0)
	f.Add("", "123 with space", "", 0)
	f.Add("", "", "id with\nnewline", 0)
	f.Add("email@with?query=param", "", "", 0)
	f.Add("", "id&with=special&chars", "", 0)

	f.Fuzz(func(t *testing.T, email, id, externalID string, limit int) {
		q := &ContactQuery{
			Email:      email,
			ID:         id,
			ExternalID: externalID,
			Limit:      limit,
		}

		// Should never panic
		url, err := q.BuildURL()

		// If all identifiers are empty, must return error
		allEmpty := len(email) == 0 && len(id) == 0 && len(externalID) == 0
		if allEmpty {
			if err == nil {
				t.Error("BuildURL() should return error when all identifiers are empty")
			}
			return
		}

		// If any identifier is present, must succeed
		if err != nil {
			t.Errorf("BuildURL() unexpected error with non-empty identifier: %v", err)
			return
		}

		// URL must not be empty when successful
		if url == "" {
			t.Error("BuildURL() returned empty URL on success")
		}

		// URL must start with the expected base
		if !strings.HasPrefix(url, apiEndpoint) {
			t.Errorf("BuildURL() URL doesn't start with apiEndpoint: %s", url)
		}

		// Verify priority: ID > Email > ExternalID
		verifyURLPriority(t, url, id, email, externalID)
	})
}

// verifyURLPriority checks that the URL correctly follows the ID > Email > ExternalID priority.
func verifyURLPriority(t *testing.T, url, id, email, externalID string) {
	t.Helper()

	switch {
	case len(id) > 0:
		expectedPath := "/contacts/" + id
		if !strings.Contains(url, expectedPath) {
			t.Errorf("BuildURL() with ID should contain %s, got: %s", expectedPath, url)
		}
	case len(email) > 0:
		if !strings.Contains(url, "email=") {
			t.Errorf("BuildURL() with Email should contain email=, got: %s", url)
		}
	case len(externalID) > 0:
		if !strings.Contains(url, "idType=external") {
			t.Errorf("BuildURL() with ExternalID should contain idType=external, got: %s", url)
		}
	}
}
