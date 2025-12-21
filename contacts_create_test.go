package drift

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errMissingRequest = errors.New("missing request")

// mockHTTPCreateContact for mocking requests
type mockHTTPCreateContact struct{}

// Do is a mock http request
func (m *mockHTTPCreateContact) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, errMissingRequest
	}

	// Valid response
	if req.URL.String() == apiEndpoint+"/contacts" {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`{"data":{"id":` + testContactID + `,"createdAt":1614563742010,"attributes":{"_END_USER_VERSION":3,"_end_user_version":3,"_calculated_version":3,"socialProfiles":{},"name":"` + testContactName + `","email":"` + testContactEmail + `","events":{},"tags":[],"start_date":1614563742010}}}`))
	}

	// Default is valid
	return resp, nil
}

// TestClient_CreateContact tests the method CreateContact()
func TestClient_CreateContact(t *testing.T) {
	t.Parallel()

	t.Run("create a standard contact", func(t *testing.T) {
		// Create a client
		client := newTestClient(&mockHTTPCreateContact{})

		// Create a req
		contact, err := client.CreateContact(
			context.Background(),
			&ContactFields{&StandardAttributes{
				Email: testContactEmail,
				Name:  testContactName,
				Phone: testContactPhone,
			}})
		require.NoError(t, err)
		assert.NotNil(t, contact)

		// Got a contact
		assert.Equal(t, uint64(123456789), contact.Data.ID)
		assert.Equal(t, int64(1614563742010), contact.Data.CreatedAt)
		assert.Equal(t, 3, contact.Data.Attributes.EndUserVersion)
		assert.Equal(t, 1614563742010, contact.Data.Attributes.StartDate)
	})
}

// BenchmarkClient_CreateContact benchmarks the CreateContact method
func BenchmarkClient_CreateContact(b *testing.B) {
	client := newTestClient(&mockHTTPCreateContact{})
	fields := &ContactFields{&StandardAttributes{
		Email: testContactEmail,
		Name:  testContactName,
		Phone: testContactPhone,
	}}
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateContact(context.Background(), fields)
	}
}
