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

// Static test errors to satisfy err113 linter
var (
	errSimulatedRead = errors.New("simulated read error")
	errNetwork       = errors.New("network error")
	errPartial       = errors.New("partial error")
)

// mockHTTPRequest implements httpInterface for request.go tests
type mockHTTPRequest struct {
	statusCode int
	body       string
	doError    error
}

// Do is a mock http request
func (m *mockHTTPRequest) Do(_ *http.Request) (*http.Response, error) {
	if m.doError != nil {
		return nil, m.doError
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(m.body)),
	}, nil
}

// mockHTTPRequestWithDoErrorAndResponse returns both an error and a response
type mockHTTPRequestWithDoErrorAndResponse struct {
	statusCode int
	doError    error
}

// Do returns both error and response (for testing line 89-91)
func (m *mockHTTPRequestWithDoErrorAndResponse) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}, m.doError
}

// errorReader implements io.ReadCloser and always returns an error on Read
type errorReader struct{}

func (e *errorReader) Read(_ []byte) (n int, err error) {
	return 0, errSimulatedRead
}

func (e *errorReader) Close() error {
	return nil
}

// mockHTTPRequestBadBody returns a body that fails to read
type mockHTTPRequestBadBody struct {
	statusCode int
}

// Do returns a response with a body that fails to read
func (m *mockHTTPRequestBadBody) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       &errorReader{},
	}, nil
}

// mockHTTPRequestCapture captures the request for inspection
type mockHTTPRequestCapture struct {
	statusCode     int
	body           string
	capturedReq    *http.Request
	capturedMethod string
	capturedURL    string
	hasAuthHeader  bool
	contentType    string
}

// Do captures the request and returns a mock response
func (m *mockHTTPRequestCapture) Do(req *http.Request) (*http.Response, error) {
	m.capturedReq = req
	m.capturedMethod = req.Method
	m.capturedURL = req.URL.String()
	m.hasAuthHeader = req.Header.Get("Authorization") != ""
	m.contentType = req.Header.Get("Content-Type")
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(m.body)),
	}, nil
}

// TestHttpRequest tests the httpRequest function
func TestHttpRequest(t *testing.T) {
	t.Parallel()

	t.Run("returns ErrResourceNotFound on 404 status", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(&mockHTTPRequest{statusCode: http.StatusNotFound, body: ""})

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/999",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.Error(t, response.Error)
		require.ErrorIs(t, response.Error, ErrResourceNotFound)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("returns ErrUnauthorized on 401 status", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(&mockHTTPRequest{statusCode: http.StatusUnauthorized, body: ""})

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/123",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.Error(t, response.Error)
		require.ErrorIs(t, response.Error, ErrUnauthorized)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("returns ErrMalformedRequest on 400 status", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(&mockHTTPRequest{statusCode: http.StatusBadRequest, body: ""})

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/123",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.Error(t, response.Error)
		require.ErrorIs(t, response.Error, ErrMalformedRequest)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("returns ErrConflict on 409 status", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(&mockHTTPRequest{statusCode: http.StatusConflict, body: ""})

		payload := &httpPayload{
			Method:         http.MethodPost,
			URL:            apiEndpoint + "/contacts",
			Data:           []byte(`{"email":"test@example.com"}`),
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.Error(t, response.Error)
		require.ErrorIs(t, response.Error, ErrConflict)
		assert.Equal(t, http.StatusConflict, response.StatusCode)
	})

	t.Run("returns ErrUnexpectedStatus on other status codes", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(&mockHTTPRequest{statusCode: http.StatusTeapot, body: ""})

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/123",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.Error(t, response.Error)
		require.ErrorIs(t, response.Error, ErrUnexpectedStatus)
		assert.Equal(t, http.StatusTeapot, response.StatusCode)
		assert.Contains(t, response.Error.Error(), "418 does not match 200")
	})

	t.Run("handles client Do error", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(&mockHTTPRequest{doError: errNetwork})

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/123",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.Error(t, response.Error)
		assert.Equal(t, errNetwork, response.Error)
	})

	t.Run("handles client Do error with response", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(&mockHTTPRequestWithDoErrorAndResponse{
			statusCode: http.StatusInternalServerError,
			doError:    errPartial,
		})

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/123",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.Error(t, response.Error)
		assert.Equal(t, errPartial, response.Error)
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("handles response body read error", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(&mockHTTPRequestBadBody{statusCode: http.StatusOK})

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/123",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.Error(t, response.Error)
		assert.Contains(t, response.Error.Error(), "simulated read error")
	})

	t.Run("sets Content-Type for POST requests", func(t *testing.T) {
		t.Parallel()
		mock := &mockHTTPRequestCapture{statusCode: http.StatusOK, body: "{}"}
		client := newTestClient(mock)

		payload := &httpPayload{
			Method:         http.MethodPost,
			URL:            apiEndpoint + "/contacts",
			Data:           []byte(`{"email":"test@example.com"}`),
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.NoError(t, response.Error)
		assert.Equal(t, "application/json", mock.contentType)
		assert.Equal(t, http.MethodPost, mock.capturedMethod)
	})

	t.Run("sets Content-Type for PATCH requests", func(t *testing.T) {
		t.Parallel()
		mock := &mockHTTPRequestCapture{statusCode: http.StatusOK, body: "{}"}
		client := newTestClient(mock)

		payload := &httpPayload{
			Method:         http.MethodPatch,
			URL:            apiEndpoint + "/contacts/123",
			Data:           []byte(`{"name":"Updated Name"}`),
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.NoError(t, response.Error)
		assert.Equal(t, "application/json", mock.contentType)
		assert.Equal(t, http.MethodPatch, mock.capturedMethod)
	})

	t.Run("does not set Content-Type for GET requests", func(t *testing.T) {
		t.Parallel()
		mock := &mockHTTPRequestCapture{statusCode: http.StatusOK, body: "{}"}
		client := newTestClient(mock)

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/123",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.NoError(t, response.Error)
		assert.Empty(t, mock.contentType)
	})

	t.Run("omits Authorization header when no OAuth token", func(t *testing.T) {
		t.Parallel()
		mock := &mockHTTPRequestCapture{statusCode: http.StatusOK, body: "{}"}
		client := newTestClientNoToken(mock)

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/123",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.NoError(t, response.Error)
		assert.False(t, mock.hasAuthHeader)
	})

	t.Run("sets Authorization header when OAuth token present", func(t *testing.T) {
		t.Parallel()
		mock := &mockHTTPRequestCapture{statusCode: http.StatusOK, body: "{}"}
		client := newTestClient(mock)

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/123",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.NoError(t, response.Error)
		assert.True(t, mock.hasAuthHeader)
	})

	t.Run("stores PostData for POST requests", func(t *testing.T) {
		t.Parallel()
		mock := &mockHTTPRequestCapture{statusCode: http.StatusOK, body: "{}"}
		client := newTestClient(mock)

		postData := `{"email":"test@example.com"}`
		payload := &httpPayload{
			Method:         http.MethodPost,
			URL:            apiEndpoint + "/contacts",
			Data:           []byte(postData),
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.NoError(t, response.Error)
		assert.Equal(t, postData, response.PostData)
	})

	t.Run("stores PostData for PATCH requests", func(t *testing.T) {
		t.Parallel()
		mock := &mockHTTPRequestCapture{statusCode: http.StatusOK, body: "{}"}
		client := newTestClient(mock)

		postData := `{"name":"Updated"}`
		payload := &httpPayload{
			Method:         http.MethodPatch,
			URL:            apiEndpoint + "/contacts/123",
			Data:           []byte(postData),
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.NoError(t, response.Error)
		assert.Equal(t, postData, response.PostData)
	})

	t.Run("stores method and URL in response", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(&mockHTTPRequest{statusCode: http.StatusOK, body: "{}"})

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/123",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.NoError(t, response.Error)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, apiEndpoint+"/contacts/123", response.URL)
	})

	t.Run("success case reads body contents", func(t *testing.T) {
		t.Parallel()
		expectedBody := `{"data":{"id":123}}`
		client := newTestClient(&mockHTTPRequest{statusCode: http.StatusOK, body: expectedBody})

		payload := &httpPayload{
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/contacts/123",
			ExpectedStatus: http.StatusOK,
		}

		response := httpRequest(context.Background(), client, payload)

		require.NoError(t, response.Error)
		assert.Equal(t, []byte(expectedBody), response.BodyContents)
	})
}

// BenchmarkHttpRequest benchmarks the httpRequest function
func BenchmarkHttpRequest(b *testing.B) {
	client := newTestClient(&mockHTTPRequest{statusCode: http.StatusOK, body: "{}"})
	payload := &httpPayload{
		Method:         http.MethodGet,
		URL:            apiEndpoint + "/contacts/123",
		ExpectedStatus: http.StatusOK,
	}

	for i := 0; i < b.N; i++ {
		_ = httpRequest(context.Background(), client, payload)
	}
}
