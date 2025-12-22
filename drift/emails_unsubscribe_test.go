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

const (
	testUnsubscribeEmail  = "unsubscribe@test.com"
	testUnsubscribeEmail2 = "unsubscribe2@test.com"
)

// mockHTTPUnsubscribeEmails for mocking requests
type mockHTTPUnsubscribeEmails struct{}

// Do is a mock http request
func (m *mockHTTPUnsubscribeEmails) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, errMissingRequest
	}

	// Valid response for unsubscribe
	if req.URL.String() == apiEndpoint+"/emails/unsubscribe" && req.Method == http.MethodPost {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`{"result":"OK","ok":true}`))
	}

	// Default is bad request
	return resp, nil
}

// mockHTTPUnsubscribeEmailsError for testing error scenarios
type mockHTTPUnsubscribeEmailsError struct {
	statusCode int
	body       string
}

// Do returns a configurable error response
func (m *mockHTTPUnsubscribeEmailsError) Do(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errMissingRequest
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(m.body)),
	}, nil
}

// TestClient_UnsubscribeEmails tests the method UnsubscribeEmails()
func TestClient_UnsubscribeEmails(t *testing.T) {
	t.Parallel()

	t.Run("unsubscribe single email successfully", func(t *testing.T) {
		client := newTestClient(&mockHTTPUnsubscribeEmails{})

		response, err := client.UnsubscribeEmails(context.Background(), []string{testUnsubscribeEmail})
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.OK)
		assert.Equal(t, "OK", response.Result)
	})

	t.Run("unsubscribe multiple emails successfully", func(t *testing.T) {
		client := newTestClient(&mockHTTPUnsubscribeEmails{})

		emails := []string{testUnsubscribeEmail, testUnsubscribeEmail2}
		response, err := client.UnsubscribeEmails(context.Background(), emails)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.OK)
		assert.Equal(t, "OK", response.Result)
	})

	t.Run("unsubscribe empty list successfully", func(t *testing.T) {
		client := newTestClient(&mockHTTPUnsubscribeEmails{})

		response, err := client.UnsubscribeEmails(context.Background(), []string{})
		require.NoError(t, err)
		assert.NotNil(t, response)
	})

	t.Run("returns error when UnsubscribeEmailsRaw fails", func(t *testing.T) {
		client := newTestClient(&mockHTTPUnsubscribeEmailsError{
			statusCode: http.StatusBadRequest,
			body:       "",
		})

		response, err := client.UnsubscribeEmails(context.Background(), []string{testUnsubscribeEmail})
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(&mockHTTPUnsubscribeEmailsError{
			statusCode: http.StatusUnauthorized,
			body:       "",
		})

		response, err := client.UnsubscribeEmails(context.Background(), []string{testUnsubscribeEmail})
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(&mockHTTPUnsubscribeEmailsError{
			statusCode: http.StatusOK,
			body:       `{"invalid json`,
		})

		response, err := client.UnsubscribeEmails(context.Background(), []string{testUnsubscribeEmail})
		require.Error(t, err)
		assert.Nil(t, response)
	})
}

// TestClient_UnsubscribeEmailsRaw tests the method UnsubscribeEmailsRaw()
func TestClient_UnsubscribeEmailsRaw(t *testing.T) {
	t.Parallel()

	t.Run("unsubscribes emails successfully", func(t *testing.T) {
		client := newTestClient(&mockHTTPUnsubscribeEmails{})

		response, err := client.UnsubscribeEmailsRaw(context.Background(), []string{testUnsubscribeEmail})
		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodPost, response.Method)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(&mockHTTPUnsubscribeEmailsError{
			statusCode: http.StatusBadRequest,
			body:       "",
		})

		response, err := client.UnsubscribeEmailsRaw(context.Background(), []string{testUnsubscribeEmail})
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("uses correct endpoint URL", func(t *testing.T) {
		client := newTestClient(&mockHTTPUnsubscribeEmails{})

		response, err := client.UnsubscribeEmailsRaw(context.Background(), []string{testUnsubscribeEmail})
		require.NoError(t, err)
		assert.Contains(t, response.URL, "/emails/unsubscribe")
	})

	t.Run("uses POST method", func(t *testing.T) {
		client := newTestClient(&mockHTTPUnsubscribeEmails{})

		response, err := client.UnsubscribeEmailsRaw(context.Background(), []string{testUnsubscribeEmail})
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, response.Method)
	})

	t.Run("sends correct JSON payload", func(t *testing.T) {
		client := newTestClient(&mockHTTPUnsubscribeEmails{})

		emails := []string{testUnsubscribeEmail, testUnsubscribeEmail2}
		response, err := client.UnsubscribeEmailsRaw(context.Background(), emails)
		require.NoError(t, err)
		assert.Contains(t, response.PostData, testUnsubscribeEmail)
		assert.Contains(t, response.PostData, testUnsubscribeEmail2)
	})
}

// BenchmarkClient_UnsubscribeEmails benchmarks the UnsubscribeEmails method
func BenchmarkClient_UnsubscribeEmails(b *testing.B) {
	client := newTestClient(&mockHTTPUnsubscribeEmails{})
	emails := []string{testUnsubscribeEmail, testUnsubscribeEmail2}
	for i := 0; i < b.N; i++ {
		_, _ = client.UnsubscribeEmails(context.Background(), emails)
	}
}

// BenchmarkClient_UnsubscribeEmailsRaw benchmarks the UnsubscribeEmailsRaw method
func BenchmarkClient_UnsubscribeEmailsRaw(b *testing.B) {
	client := newTestClient(&mockHTTPUnsubscribeEmails{})
	emails := []string{testUnsubscribeEmail, testUnsubscribeEmail2}
	for i := 0; i < b.N; i++ {
		_, _ = client.UnsubscribeEmailsRaw(context.Background(), emails)
	}
}
