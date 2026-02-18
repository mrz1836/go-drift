package drift

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"syscall"
	"time"
)

// ResilientClient wraps an http.Client with retry capabilities.
// It implements httpInterface for seamless integration.
type ResilientClient struct {
	client     *http.Client
	backoff    Backoff
	retryCount int
}

// ResilientClientOption configures a ResilientClient.
type ResilientClientOption func(*ResilientClient)

// WithBackoff sets the backoff strategy for retries.
func WithBackoff(b Backoff) ResilientClientOption {
	return func(rc *ResilientClient) {
		rc.backoff = b
	}
}

// WithRetryCount sets the maximum number of retry attempts.
func WithRetryCount(count int) ResilientClientOption {
	return func(rc *ResilientClient) {
		rc.retryCount = count
	}
}

// NewResilientClient creates a new resilient HTTP client.
func NewResilientClient(client *http.Client, opts ...ResilientClientOption) *ResilientClient {
	rc := &ResilientClient{
		client:     client,
		retryCount: 0, // No retries by default
	}

	for _, opt := range opts {
		opt(rc)
	}

	return rc
}

// Do executes the HTTP request with retry logic.
// It respects context cancellation during retry waits.
func (rc *ResilientClient) Do(req *http.Request) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response

	maxAttempts := 1 + rc.retryCount

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := req.Context().Err(); err != nil {
			return nil, err
		}

		reqToSend, err := rc.prepareRequest(req, attempt)
		if err != nil {
			return nil, err
		}

		resp, err := rc.client.Do(reqToSend) //nolint:gosec // G704: request originates from internal API calls, not user-controlled input

		if !rc.shouldRetry(err, resp) {
			return resp, err
		}

		lastErr = err
		lastResp = resp

		if attempt >= maxAttempts-1 {
			break
		}

		if err := rc.waitForRetry(req.Context(), resp, attempt); err != nil {
			return nil, err
		}
	}

	return lastResp, lastErr
}

// prepareRequest clones the request for retry if needed.
func (rc *ResilientClient) prepareRequest(req *http.Request, attempt int) (*http.Request, error) {
	if attempt > 0 && req.GetBody != nil {
		return cloneRequest(req)
	}
	return req, nil
}

// waitForRetry closes the response body and waits for the backoff delay.
func (rc *ResilientClient) waitForRetry(ctx context.Context, resp *http.Response, attempt int) error {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}

	if rc.backoff != nil {
		delay := rc.backoff.Next(attempt)
		return rc.sleep(ctx, delay)
	}

	return nil
}

// shouldRetry determines if the request should be retried based on
// the error and response.
func (rc *ResilientClient) shouldRetry(err error, resp *http.Response) bool {
	// No retries configured
	if rc.retryCount <= 0 {
		return false
	}

	// Check for retryable errors
	if err != nil {
		return isRetryableError(err)
	}

	// Check for retryable status codes
	if resp != nil {
		return isRetryableStatusCode(resp.StatusCode)
	}

	return false
}

// sleep waits for the specified duration, respecting context cancellation.
func (rc *ResilientClient) sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// cloneRequest creates a copy of the request with a fresh body.
func cloneRequest(req *http.Request) (*http.Request, error) {
	clone := req.Clone(req.Context())

	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		clone.Body = body
	}

	return clone, nil
}

// isRetryableError checks if an error is transient and worth retrying.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Context cancellation is not retryable
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// Network errors (connection refused, reset, timeout)
	var netErr net.Error
	if errors.As(err, &netErr) {
		// Timeout errors are retryable
		if netErr.Timeout() {
			return true
		}
	}

	// Connection refused
	if errors.Is(err, syscall.ECONNREFUSED) {
		return true
	}

	// Connection reset
	if errors.Is(err, syscall.ECONNRESET) {
		return true
	}

	// EOF during read (connection closed unexpectedly)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}

	// DNS errors - only temporary ones are retryable
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return dnsErr.Temporary()
	}

	return false
}

// isRetryableStatusCode determines if an HTTP status code warrants a retry.
// Retries on 5xx (server errors), 429 (rate limiting), and 408 (request timeout).
// Does NOT retry on other 4xx (client errors) as they indicate permanent failures.
func isRetryableStatusCode(code int) bool {
	switch code {
	case http.StatusRequestTimeout: // 408
		return true
	case http.StatusTooManyRequests: // 429
		return true
	default:
		// 5xx server errors
		return code >= 500 && code <= 599
	}
}
