package drift

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testContactIDNotFound = "999999999"
)

// mockHTTPDeleteContact for mocking requests
type mockHTTPDeleteContact struct{}

// Do is a mock http request
func (m *mockHTTPDeleteContact) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, errMissingRequest
	}

	// Valid response for delete
	if req.URL.String() == apiEndpoint+"/contacts/"+testContactID && req.Method == http.MethodDelete {
		resp.StatusCode = http.StatusAccepted
		resp.Body = io.NopCloser(bytes.NewBufferString(`{"result":"OK","ok":true}`))
	}

	// Not found response
	if req.URL.String() == apiEndpoint+"/contacts/"+testContactIDNotFound && req.Method == http.MethodDelete {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(bytes.NewBufferString(`{"error":"Contact not found"}`))
	}

	// Default is bad request
	return resp, nil
}

// mockHTTPDeleteContactError for testing error scenarios
type mockHTTPDeleteContactError struct {
	statusCode int
	body       string
}

// Do returns a configurable error response
func (m *mockHTTPDeleteContactError) Do(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errMissingRequest
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(m.body)),
	}, nil
}

// TestClient_DeleteContact tests the method DeleteContact()
func TestClient_DeleteContact(t *testing.T) {
	t.Parallel()

	t.Run("delete a contact successfully", func(t *testing.T) {
		client := newTestClient(&mockHTTPDeleteContact{})

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		response, err := client.DeleteContact(context.Background(), id)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.OK)
		assert.Equal(t, "OK", response.Result)
	})

	t.Run("returns error when DeleteContactRaw fails", func(t *testing.T) {
		client := newTestClient(&mockHTTPDeleteContactError{
			statusCode: http.StatusBadRequest,
			body:       "",
		})

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		response, err := client.DeleteContact(context.Background(), id)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(&mockHTTPDeleteContactError{
			statusCode: http.StatusUnauthorized,
			body:       "",
		})

		response, err := client.DeleteContact(context.Background(), 123456789)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on 404 not found", func(t *testing.T) {
		client := newTestClient(&mockHTTPDeleteContactError{
			statusCode: http.StatusNotFound,
			body:       "",
		})

		response, err := client.DeleteContact(context.Background(), 999999999)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrResourceNotFound)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(&mockHTTPDeleteContactError{
			statusCode: http.StatusAccepted,
			body:       `{"invalid json`,
		})

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		response, err := client.DeleteContact(context.Background(), id)
		require.Error(t, err)
		assert.Nil(t, response)
	})
}

// TestClient_DeleteContactRaw tests the method DeleteContactRaw()
func TestClient_DeleteContactRaw(t *testing.T) {
	t.Parallel()

	t.Run("deletes contact successfully", func(t *testing.T) {
		client := newTestClient(&mockHTTPDeleteContact{})

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		response, err := client.DeleteContactRaw(context.Background(), id)
		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusAccepted, response.StatusCode)
		assert.Equal(t, http.MethodDelete, response.Method)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(&mockHTTPDeleteContactError{
			statusCode: http.StatusBadRequest,
			body:       "",
		})

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		response, err := client.DeleteContactRaw(context.Background(), id)
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("uses correct endpoint URL with contact ID", func(t *testing.T) {
		client := newTestClient(&mockHTTPDeleteContact{})

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		response, err := client.DeleteContactRaw(context.Background(), id)
		require.NoError(t, err)
		assert.Contains(t, response.URL, "/contacts/"+testContactID)
	})

	t.Run("uses DELETE method", func(t *testing.T) {
		client := newTestClient(&mockHTTPDeleteContact{})

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		response, err := client.DeleteContactRaw(context.Background(), id)
		require.NoError(t, err)
		assert.Equal(t, http.MethodDelete, response.Method)
	})
}

// BenchmarkClient_DeleteContact benchmarks the DeleteContact method
func BenchmarkClient_DeleteContact(b *testing.B) {
	client := newTestClient(&mockHTTPDeleteContact{})
	id, _ := strconv.ParseUint(testContactID, 10, 64)
	for i := 0; i < b.N; i++ {
		_, _ = client.DeleteContact(context.Background(), id)
	}
}

// BenchmarkClient_DeleteContactRaw benchmarks the DeleteContactRaw method
func BenchmarkClient_DeleteContactRaw(b *testing.B) {
	client := newTestClient(&mockHTTPDeleteContact{})
	id, _ := strconv.ParseUint(testContactID, 10, 64)
	for i := 0; i < b.N; i++ {
		_, _ = client.DeleteContactRaw(context.Background(), id)
	}
}
