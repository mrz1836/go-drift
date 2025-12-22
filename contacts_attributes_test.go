package drift

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockHTTPListCustomAttributes for mocking requests
type mockHTTPListCustomAttributes struct{}

// Do is a mock http request
func (m *mockHTTPListCustomAttributes) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, errMissingRequest
	}

	// Valid response for list custom attributes
	if req.URL.String() == apiEndpoint+"/contacts/attributes" && req.Method == http.MethodGet {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`{
			"data": {
				"properties": [
					{"type": "STRING", "displayName": "Age", "name": "age"},
					{"type": "BOOLEAN", "displayName": "VIP Customer", "name": "vip_customer"},
					{"type": "NUMERIC", "displayName": "Score", "name": "score"}
				]
			}
		}`))
	}

	return resp, nil
}

// mockHTTPListCustomAttributesEmpty for testing empty response
type mockHTTPListCustomAttributesEmpty struct{}

// Do returns an empty properties array
func (m *mockHTTPListCustomAttributesEmpty) Do(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errMissingRequest
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"data": {"properties": []}}`)),
	}, nil
}

// mockHTTPListCustomAttributesError for testing error scenarios
type mockHTTPListCustomAttributesError struct {
	statusCode int
	body       string
}

// Do returns a configurable error response
func (m *mockHTTPListCustomAttributesError) Do(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errMissingRequest
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(m.body)),
	}, nil
}

// TestClient_ListCustomAttributes tests the method ListCustomAttributes()
func TestClient_ListCustomAttributes(t *testing.T) {
	t.Parallel()

	t.Run("list custom attributes successfully", func(t *testing.T) {
		client := newTestClient(&mockHTTPListCustomAttributes{})

		response, err := client.ListCustomAttributes(context.Background())
		require.NoError(t, err)
		require.NotNil(t, response)
		require.NotNil(t, response.Data)
		require.Len(t, response.Data.Properties, 3)

		// Verify first attribute
		assert.Equal(t, "STRING", response.Data.Properties[0].Type)
		assert.Equal(t, "Age", response.Data.Properties[0].DisplayName)
		assert.Equal(t, "age", response.Data.Properties[0].Name)

		// Verify second attribute
		assert.Equal(t, "BOOLEAN", response.Data.Properties[1].Type)
		assert.Equal(t, "VIP Customer", response.Data.Properties[1].DisplayName)
		assert.Equal(t, "vip_customer", response.Data.Properties[1].Name)

		// Verify third attribute
		assert.Equal(t, "NUMERIC", response.Data.Properties[2].Type)
		assert.Equal(t, "Score", response.Data.Properties[2].DisplayName)
		assert.Equal(t, "score", response.Data.Properties[2].Name)
	})

	t.Run("list custom attributes returns empty list", func(t *testing.T) {
		client := newTestClient(&mockHTTPListCustomAttributesEmpty{})

		response, err := client.ListCustomAttributes(context.Background())
		require.NoError(t, err)
		require.NotNil(t, response)
		require.NotNil(t, response.Data)
		assert.Empty(t, response.Data.Properties)
	})

	t.Run("returns error when ListCustomAttributesRaw fails", func(t *testing.T) {
		client := newTestClient(&mockHTTPListCustomAttributesError{
			statusCode: http.StatusBadRequest,
			body:       "",
		})

		response, err := client.ListCustomAttributes(context.Background())
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(&mockHTTPListCustomAttributesError{
			statusCode: http.StatusUnauthorized,
			body:       "",
		})

		response, err := client.ListCustomAttributes(context.Background())
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on 404 not found", func(t *testing.T) {
		client := newTestClient(&mockHTTPListCustomAttributesError{
			statusCode: http.StatusNotFound,
			body:       "",
		})

		response, err := client.ListCustomAttributes(context.Background())
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrResourceNotFound)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(&mockHTTPListCustomAttributesError{
			statusCode: http.StatusOK,
			body:       `{"invalid json`,
		})

		response, err := client.ListCustomAttributes(context.Background())
		require.Error(t, err)
		assert.Nil(t, response)
	})
}

// TestClient_ListCustomAttributesRaw tests the method ListCustomAttributesRaw()
func TestClient_ListCustomAttributesRaw(t *testing.T) {
	t.Parallel()

	t.Run("lists custom attributes successfully", func(t *testing.T) {
		client := newTestClient(&mockHTTPListCustomAttributes{})

		response, err := client.ListCustomAttributesRaw(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodGet, response.Method)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(&mockHTTPListCustomAttributesError{
			statusCode: http.StatusBadRequest,
			body:       "",
		})

		response, err := client.ListCustomAttributesRaw(context.Background())
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("uses correct endpoint URL", func(t *testing.T) {
		client := newTestClient(&mockHTTPListCustomAttributes{})

		response, err := client.ListCustomAttributesRaw(context.Background())
		require.NoError(t, err)
		assert.Equal(t, apiEndpoint+"/contacts/attributes", response.URL)
	})

	t.Run("uses GET method", func(t *testing.T) {
		client := newTestClient(&mockHTTPListCustomAttributes{})

		response, err := client.ListCustomAttributesRaw(context.Background())
		require.NoError(t, err)
		assert.Equal(t, http.MethodGet, response.Method)
	})
}

// BenchmarkClient_ListCustomAttributes benchmarks the ListCustomAttributes method
func BenchmarkClient_ListCustomAttributes(b *testing.B) {
	client := newTestClient(&mockHTTPListCustomAttributes{})
	for i := 0; i < b.N; i++ {
		_, _ = client.ListCustomAttributes(context.Background())
	}
}

// BenchmarkClient_ListCustomAttributesRaw benchmarks the ListCustomAttributesRaw method
func BenchmarkClient_ListCustomAttributesRaw(b *testing.B) {
	client := newTestClient(&mockHTTPListCustomAttributes{})
	for i := 0; i < b.N; i++ {
		_, _ = client.ListCustomAttributesRaw(context.Background())
	}
}
