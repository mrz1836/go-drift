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

// mockHTTPCreateContactError for testing error scenarios
type mockHTTPCreateContactError struct {
	statusCode int
	body       string
}

// Do returns a configurable error response
func (m *mockHTTPCreateContactError) Do(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errMissingRequest
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(m.body)),
	}, nil
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

	t.Run("returns error when CreateContactRaw fails", func(t *testing.T) {
		client := newTestClient(&mockHTTPCreateContactError{
			statusCode: http.StatusBadRequest,
			body:       "",
		})

		contact, err := client.CreateContact(
			context.Background(),
			&ContactFields{&StandardAttributes{
				Email: testContactEmail,
				Name:  testContactName,
			}})

		require.Error(t, err)
		assert.Nil(t, contact)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(&mockHTTPCreateContactError{
			statusCode: http.StatusUnauthorized,
			body:       "",
		})

		contact, err := client.CreateContact(
			context.Background(),
			&ContactFields{&StandardAttributes{
				Email: testContactEmail,
			}})

		require.Error(t, err)
		assert.Nil(t, contact)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(&mockHTTPCreateContactError{
			statusCode: http.StatusOK,
			body:       `{"data":{"invalid json`,
		})

		contact, err := client.CreateContact(
			context.Background(),
			&ContactFields{&StandardAttributes{
				Email: testContactEmail,
			}})

		require.Error(t, err)
		assert.Nil(t, contact)
	})
}

// TestClient_CreateContactRaw tests the method CreateContactRaw()
func TestClient_CreateContactRaw(t *testing.T) {
	t.Parallel()

	t.Run("creates contact successfully", func(t *testing.T) {
		client := newTestClient(&mockHTTPCreateContact{})

		response, err := client.CreateContactRaw(
			context.Background(),
			&ContactFields{&StandardAttributes{
				Email: testContactEmail,
				Name:  testContactName,
			}})

		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodPost, response.Method)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(&mockHTTPCreateContactError{
			statusCode: http.StatusBadRequest,
			body:       "",
		})

		response, err := client.CreateContactRaw(
			context.Background(),
			&ContactFields{&StandardAttributes{
				Email: testContactEmail,
			}})

		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
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
