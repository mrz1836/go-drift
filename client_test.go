package drift

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

const (
	testContactEmail          = "johndoe@email.com"
	testContactID             = "123456789"
	testContactIDBadJSON      = "333333333"
	testContactIDBadRequest   = "111111111"
	testContactIDUnauthorized = "222222222"
	testContactName           = "John Doe"
	testContactPhone          = "15554443333"
	testDataOAuthToken        = "testKey1234567"
	testEventName             = "test-event-name-goes-here"
)

// newTestClient returns a client for mocking (using a custom HTTP interface)
func newTestClient(httpClient httpInterface) *Client {
	client := NewClient(testDataOAuthToken, nil, nil)
	client.httpClient = httpClient
	return client
}

// TestNewClient test new client
func TestNewClient(t *testing.T) {
	t.Parallel()

	client := NewClient(testDataOAuthToken, nil, nil)

	if len(client.Options.UserAgent) == 0 {
		t.Fatal("missing user agent")
	}
}

// TestNewClient_CustomHTTPClient test new client
func TestNewClient_CustomHTTPClient(t *testing.T) {
	t.Parallel()

	client := NewClient(testDataOAuthToken, nil, http.DefaultClient)

	if len(client.Options.UserAgent) == 0 {
		t.Fatal("user agent should be default even if custom HTTP")
	}
}

// ExampleNewClient example using NewClient()
func ExampleNewClient() {
	client := NewClient(testDataOAuthToken, nil, nil)
	fmt.Println(client.Options.UserAgent)
	// Output:go-drift: v0.0.2
}

// BenchmarkNewClient benchmarks the NewClient method
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewClient(testDataOAuthToken, nil, nil)
	}
}

// TestClientDefaultOptions tests setting ClientDefaultOptions()
func TestClientDefaultOptions(t *testing.T) {
	t.Parallel()

	options := DefaultClientOptions()

	if options.UserAgent != defaultUserAgent {
		t.Fatalf("expected value: %s got: %s", defaultUserAgent, options.UserAgent)
	}

	if options.BackOffExponentFactor != 2.0 {
		t.Fatalf("expected value: %f got: %f", 2.0, options.BackOffExponentFactor)
	}

	if options.BackOffInitialTimeout != 2*time.Millisecond {
		t.Fatalf("expected value: %v got: %v", 2*time.Millisecond, options.BackOffInitialTimeout)
	}

	if options.BackOffMaximumJitterInterval != 2*time.Millisecond {
		t.Fatalf("expected value: %v got: %v", 2*time.Millisecond, options.BackOffMaximumJitterInterval)
	}

	if options.BackOffMaxTimeout != 10*time.Millisecond {
		t.Fatalf("expected value: %v got: %v", 10*time.Millisecond, options.BackOffMaxTimeout)
	}

	if options.DialerKeepAlive != 20*time.Second {
		t.Fatalf("expected value: %v got: %v", 20*time.Second, options.DialerKeepAlive)
	}

	if options.DialerTimeout != 5*time.Second {
		t.Fatalf("expected value: %v got: %v", 5*time.Second, options.DialerTimeout)
	}

	if options.RequestRetryCount != 2 {
		t.Fatalf("expected value: %v got: %v", 2, options.RequestRetryCount)
	}

	if options.RequestTimeout != 10*time.Second {
		t.Fatalf("expected value: %v got: %v", 10*time.Second, options.RequestTimeout)
	}

	if options.TransportExpectContinueTimeout != 3*time.Second {
		t.Fatalf("expected value: %v got: %v", 3*time.Second, options.TransportExpectContinueTimeout)
	}

	if options.TransportIdleTimeout != 20*time.Second {
		t.Fatalf("expected value: %v got: %v", 20*time.Second, options.TransportIdleTimeout)
	}

	if options.TransportMaxIdleConnections != 10 {
		t.Fatalf("expected value: %v got: %v", 10, options.TransportMaxIdleConnections)
	}

	if options.TransportTLSHandshakeTimeout != 5*time.Second {
		t.Fatalf("expected value: %v got: %v", 5*time.Second, options.TransportTLSHandshakeTimeout)
	}
}

// TestClientDefaultOptions_NoRetry will set 0 retry counts
func TestClientDefaultOptions_NoRetry(t *testing.T) {
	options := DefaultClientOptions()
	options.RequestRetryCount = 0
	client := NewClient(testDataOAuthToken, options, nil)

	if client.Options.UserAgent != defaultUserAgent {
		t.Errorf("user agent mismatch")
	}
}
