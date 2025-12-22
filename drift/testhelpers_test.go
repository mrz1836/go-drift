package drift

import (
	"bytes"
	"io"
	"net/http"
)

// mockHTTP is a configurable mock HTTP client for testing
type mockHTTP struct {
	statusCode int
	body       string
}

// mockHTTPOption is a function that configures a mockHTTP
type mockHTTPOption func(*mockHTTP)

// newMockHTTP creates a new configurable mock HTTP client
func newMockHTTP(opts ...mockHTTPOption) *mockHTTP {
	m := &mockHTTP{
		statusCode: http.StatusOK,
		body:       "",
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withStatus sets the response status code
func withStatus(code int) mockHTTPOption {
	return func(m *mockHTTP) {
		m.statusCode = code
	}
}

// withBody sets the response body
func withBody(body string) mockHTTPOption {
	return func(m *mockHTTP) {
		m.body = body
	}
}

// Do implements the httpInterface for mockHTTP
func (m *mockHTTP) Do(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errMissingRequest
	}

	resp := &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(m.body)),
	}

	return resp, nil
}

// mockHTTPMulti is a mock that can handle multiple URL patterns
type mockHTTPMulti struct {
	routes map[string]*mockRoute
}

// mockRoute represents a single route configuration
type mockRoute struct {
	statusCode int
	body       string
}

// newMockHTTPMulti creates a mock that can handle multiple routes
func newMockHTTPMulti() *mockHTTPMulti {
	return &mockHTTPMulti{
		routes: make(map[string]*mockRoute),
	}
}

// Do implements the httpInterface for mockHTTPMulti
func (m *mockHTTPMulti) Do(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errMissingRequest
	}

	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}

	if route, ok := m.routes[req.URL.String()]; ok {
		resp.StatusCode = route.statusCode
		resp.Body = io.NopCloser(bytes.NewBufferString(route.body))
	}

	return resp, nil
}

// addRoute adds a route to the mock
func (m *mockHTTPMulti) addRoute(url string, statusCode int, body string) *mockHTTPMulti {
	m.routes[url] = &mockRoute{
		statusCode: statusCode,
		body:       body,
	}
	return m
}

// newMockError creates a mock that returns a specific error status code
func newMockError(statusCode int) *mockHTTP {
	return newMockHTTP(withStatus(statusCode))
}

// newMockSuccess creates a mock that returns success with the given body
func newMockSuccess(body string) *mockHTTP {
	return newMockHTTP(withStatus(http.StatusOK), withBody(body))
}
