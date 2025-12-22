package drift

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Static test errors to satisfy err113 linter
var (
	errTestRandom = errors.New("some error")
	errGetBody    = errors.New("GetBody error")
)

// mockTransport is a test transport that can simulate various responses
type mockTransport struct {
	responses []*http.Response
	errors    []error
	callCount atomic.Int32
}

func (m *mockTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	idx := int(m.callCount.Add(1)) - 1

	if idx < len(m.errors) && m.errors[idx] != nil {
		return nil, m.errors[idx]
	}

	if idx < len(m.responses) {
		return m.responses[idx], nil
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
	}, nil
}

// mockBackoff is a test backoff that returns zero delay
type mockBackoff struct {
	delays []time.Duration
	calls  atomic.Int32
}

func (m *mockBackoff) Next(_ int) time.Duration {
	idx := int(m.calls.Add(1)) - 1
	if idx < len(m.delays) {
		return m.delays[idx]
	}
	return 0
}

func TestNewResilientClient(t *testing.T) {
	t.Parallel()

	t.Run("creates client with defaults", func(t *testing.T) {
		client := NewResilientClient(&http.Client{})
		assert.NotNil(t, client)
		assert.Equal(t, 0, client.retryCount)
		assert.Nil(t, client.backoff)
	})

	t.Run("applies options", func(t *testing.T) {
		backoff := NewExponentialBackoff(1*time.Millisecond, 10*time.Millisecond, 2.0, 0)
		client := NewResilientClient(
			&http.Client{},
			WithBackoff(backoff),
			WithRetryCount(3),
		)
		assert.Equal(t, 3, client.retryCount)
		assert.Equal(t, backoff, client.backoff)
	})
}

func TestResilientClientDoNoRetry(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(`ok`))},
		},
	}

	client := NewResilientClient(&http.Client{Transport: transport})

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(1), transport.callCount.Load())
}

func TestResilientClientDoRetriesOn5xx(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewBufferString(`error`))},
			{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewBufferString(`error`))},
			{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(`ok`))},
		},
	}

	backoff := &mockBackoff{delays: []time.Duration{0, 0}}
	client := NewResilientClient(
		&http.Client{Transport: transport},
		WithBackoff(backoff),
		WithRetryCount(2),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(3), transport.callCount.Load())
}

func TestResilientClientDoRetriesOn429(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusTooManyRequests, Body: io.NopCloser(bytes.NewBufferString(`rate limited`))},
			{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(`ok`))},
		},
	}

	backoff := &mockBackoff{delays: []time.Duration{0}}
	client := NewResilientClient(
		&http.Client{Transport: transport},
		WithBackoff(backoff),
		WithRetryCount(1),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(2), transport.callCount.Load())
}

func TestResilientClientDoRetriesOn408(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusRequestTimeout, Body: io.NopCloser(bytes.NewBufferString(`timeout`))},
			{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(`ok`))},
		},
	}

	backoff := &mockBackoff{delays: []time.Duration{0}}
	client := NewResilientClient(
		&http.Client{Transport: transport},
		WithBackoff(backoff),
		WithRetryCount(1),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(2), transport.callCount.Load())
}

func TestResilientClientDoNoRetryOn4xx(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusBadRequest, Body: io.NopCloser(bytes.NewBufferString(`bad`))},
		},
	}

	backoff := &mockBackoff{delays: []time.Duration{0, 0}}
	client := NewResilientClient(
		&http.Client{Transport: transport},
		WithBackoff(backoff),
		WithRetryCount(2),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, int32(1), transport.callCount.Load()) // No retries
}

func TestResilientClientDoRespectsContextCancellation(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewBufferString(`error`))},
		},
	}

	// Use a longer delay so we can cancel during it
	backoff := &mockBackoff{delays: []time.Duration{100 * time.Millisecond}}
	client := NewResilientClient(
		&http.Client{Transport: transport},
		WithBackoff(backoff),
		WithRetryCount(5),
	)

	ctx, cancel := context.WithCancel(context.Background())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	// Cancel after first request
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	_, err = client.Do(req) //nolint:bodyclose // Error case: no response body to close
	require.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, int32(1), transport.callCount.Load()) // Only initial attempt
}

func TestResilientClientDoMaxRetriesExhausted(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusServiceUnavailable, Body: io.NopCloser(bytes.NewBufferString(`error`))},
			{StatusCode: http.StatusServiceUnavailable, Body: io.NopCloser(bytes.NewBufferString(`error`))},
			{StatusCode: http.StatusServiceUnavailable, Body: io.NopCloser(bytes.NewBufferString(`error`))},
		},
	}

	backoff := &mockBackoff{delays: []time.Duration{0, 0}}
	client := NewResilientClient(
		&http.Client{Transport: transport},
		WithBackoff(backoff),
		WithRetryCount(2),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	assert.Equal(t, int32(3), transport.callCount.Load()) // 1 initial + 2 retries
}

func TestResilientClientDoRetriesOnNetworkError(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{
		errors: []error{
			syscall.ECONNREFUSED,
			nil, // success on second try
		},
		responses: []*http.Response{
			nil,
			{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(`ok`))},
		},
	}

	backoff := &mockBackoff{delays: []time.Duration{0}}
	client := NewResilientClient(
		&http.Client{Transport: transport},
		WithBackoff(backoff),
		WithRetryCount(1),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(2), transport.callCount.Load())
}

func TestResilientClientDoRetriesOnEOF(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{
		errors: []error{
			io.EOF,
			nil, // success on second try
		},
		responses: []*http.Response{
			nil,
			{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(`ok`))},
		},
	}

	backoff := &mockBackoff{delays: []time.Duration{0}}
	client := NewResilientClient(
		&http.Client{Transport: transport},
		WithBackoff(backoff),
		WithRetryCount(1),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(2), transport.callCount.Load())
}

func TestResilientClientDoNoRetryOnContextCanceled(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{
		errors: []error{
			context.Canceled,
		},
	}

	backoff := &mockBackoff{delays: []time.Duration{0}}
	client := NewResilientClient(
		&http.Client{Transport: transport},
		WithBackoff(backoff),
		WithRetryCount(2),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	_, err = client.Do(req) //nolint:bodyclose // Error case: no response body to close
	require.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, int32(1), transport.callCount.Load()) // No retries
}

func TestIsRetryableError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{"nil error", nil, false},
		{"context canceled", context.Canceled, false},
		{"context deadline exceeded", context.DeadlineExceeded, false},
		{"connection refused", syscall.ECONNREFUSED, true},
		{"connection reset", syscall.ECONNRESET, true},
		{"EOF", io.EOF, true},
		{"unexpected EOF", io.ErrUnexpectedEOF, true},
		{"random error", errTestRandom, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.retryable, isRetryableError(tt.err), "error: %v", tt.err)
		})
	}
}

func TestIsRetryableStatusCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		code      int
		retryable bool
	}{
		{http.StatusOK, false},                 // 200
		{http.StatusCreated, false},            // 201
		{http.StatusBadRequest, false},         // 400
		{http.StatusUnauthorized, false},       // 401
		{http.StatusForbidden, false},          // 403
		{http.StatusNotFound, false},           // 404
		{http.StatusRequestTimeout, true},      // 408
		{http.StatusConflict, false},           // 409
		{http.StatusTooManyRequests, true},     // 429
		{http.StatusInternalServerError, true}, // 500
		{http.StatusBadGateway, true},          // 502
		{http.StatusServiceUnavailable, true},  // 503
		{http.StatusGatewayTimeout, true},      // 504
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.code), func(t *testing.T) {
			assert.Equal(t, tt.retryable, isRetryableStatusCode(tt.code), "status %d", tt.code)
		})
	}
}

func TestResilientClientDoWithPostBody(t *testing.T) {
	t.Parallel()

	bodyCaptured := make([]string, 0, 2)
	transport := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewBufferString(`error`))},
			{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(`ok`))},
		},
	}

	// Custom transport to capture body
	captureTransport := &bodyCapturingTransport{
		inner:  transport,
		bodies: &bodyCaptured,
	}

	backoff := &mockBackoff{delays: []time.Duration{0}}
	client := NewResilientClient(
		&http.Client{Transport: captureTransport},
		WithBackoff(backoff),
		WithRetryCount(1),
	)

	body := bytes.NewBufferString(`{"test":"data"}`)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://example.com", body)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify body was sent on both requests
	require.Len(t, bodyCaptured, 2)
	assert.JSONEq(t, `{"test":"data"}`, bodyCaptured[0])
	assert.JSONEq(t, `{"test":"data"}`, bodyCaptured[1])
}

type bodyCapturingTransport struct {
	inner  http.RoundTripper
	bodies *[]string
}

func (t *bodyCapturingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		*t.bodies = append(*t.bodies, string(body))
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	return t.inner.RoundTrip(req)
}

func TestResilientClientImplementsHTTPInterface(t *testing.T) {
	t.Parallel()

	var _ httpInterface = (*ResilientClient)(nil)
}

// mockTimeoutError implements net.Error for testing timeout scenarios
type mockTimeoutError struct{}

func (e *mockTimeoutError) Error() string   { return "mock timeout" }
func (e *mockTimeoutError) Timeout() bool   { return true }
func (e *mockTimeoutError) Temporary() bool { return true }

func TestIsRetryableErrorTimeout(t *testing.T) {
	t.Parallel()

	var netErr net.Error = &mockTimeoutError{}
	assert.True(t, isRetryableError(netErr))
}

func TestIsRetryableErrorDNS(t *testing.T) {
	t.Parallel()

	t.Run("temporary DNS error is retryable", func(t *testing.T) {
		dnsErr := &net.DNSError{
			Err:         "lookup failed",
			Name:        "example.com",
			IsTemporary: true,
		}
		assert.True(t, isRetryableError(dnsErr))
	})

	t.Run("permanent DNS error is not retryable", func(t *testing.T) {
		dnsErr := &net.DNSError{
			Err:         "no such host",
			Name:        "example.com",
			IsTemporary: false,
		}
		assert.False(t, isRetryableError(dnsErr))
	})
}

func TestResilientClientDoGetBodyError(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewBufferString(`error`))},
		},
	}

	backoff := &mockBackoff{delays: []time.Duration{0}}
	client := NewResilientClient(
		&http.Client{Transport: transport},
		WithBackoff(backoff),
		WithRetryCount(1),
	)

	body := bytes.NewBufferString(`{"test":"data"}`)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://example.com", body)
	require.NoError(t, err)

	// Set GetBody to return an error on retry
	req.GetBody = func() (io.ReadCloser, error) {
		return nil, errGetBody
	}

	_, err = client.Do(req) //nolint:bodyclose // Error case: no response body to close
	require.Error(t, err)
	assert.Contains(t, err.Error(), "GetBody error")
}

func TestResilientClientDoContextAlreadyCanceled(t *testing.T) {
	t.Parallel()

	transport := &mockTransport{}
	client := NewResilientClient(
		&http.Client{Transport: transport},
		WithRetryCount(2),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel before request

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	_, err = client.Do(req) //nolint:bodyclose // Error case: no response body to close
	require.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, int32(0), transport.callCount.Load()) // No attempts made
}

func TestWaitForRetryNilResponse(t *testing.T) {
	t.Parallel()

	client := NewResilientClient(&http.Client{})

	// waitForRetry should handle nil response gracefully
	err := client.waitForRetry(context.Background(), nil, 0)
	assert.NoError(t, err)
}

func TestWaitForRetryNilBody(t *testing.T) {
	t.Parallel()

	client := NewResilientClient(&http.Client{})

	// waitForRetry should handle response with nil body
	resp := &http.Response{StatusCode: 500, Body: nil}
	err := client.waitForRetry(context.Background(), resp, 0)
	assert.NoError(t, err)
}

func TestWaitForRetryWithoutBackoff(t *testing.T) {
	t.Parallel()

	// Client without backoff configured
	client := NewResilientClient(&http.Client{})

	resp := &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewBufferString(`error`)),
	}
	err := client.waitForRetry(context.Background(), resp, 0)
	assert.NoError(t, err)
}

func BenchmarkResilientClientDo(b *testing.B) {
	transport := &mockTransport{
		responses: []*http.Response{
			{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString(`ok`))},
		},
	}

	client := NewResilientClient(&http.Client{Transport: transport})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transport.callCount.Store(0)
		resp, _ := client.Do(req)
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}
}
