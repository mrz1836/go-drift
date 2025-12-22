package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testGDPREmail       = "user@example.com"
	testGDPRSentToEmail = "admin@company.com"
	testGDPRRetrieveMsg = "Your request is processing. When the data is gathered, it will be sent in an email to admin@company.com"
)

// mockGDPRRetrieval returns a mock for GDPR retrieval operations
func mockGDPRRetrieval() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":{"message":"`+testGDPRRetrieveMsg+`","sentToEmail":"`+testGDPRSentToEmail+`"}}`),
	)
}

// TestClient_RetrieveGDPR tests the method RetrieveGDPR()
func TestClient_RetrieveGDPR(t *testing.T) {
	t.Parallel()

	t.Run("retrieve GDPR data successfully", func(t *testing.T) {
		client := newTestClient(mockGDPRRetrieval())

		response, err := client.RetrieveGDPR(context.Background(), testGDPREmail)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.Data)
		assert.Equal(t, testGDPRRetrieveMsg, response.Data.Message)
		assert.Equal(t, testGDPRSentToEmail, response.Data.SentToEmail)
	})

	t.Run("returns error when email is empty", func(t *testing.T) {
		client := newTestClient(mockGDPRRetrieval())

		response, err := client.RetrieveGDPR(context.Background(), "")
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingEmail)
	})

	t.Run("returns error on 400 bad request", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		response, err := client.RetrieveGDPR(context.Background(), testGDPREmail)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		response, err := client.RetrieveGDPR(context.Background(), testGDPREmail)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on 404 not found", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusNotFound))

		response, err := client.RetrieveGDPR(context.Background(), testGDPREmail)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrResourceNotFound)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"data":{"invalid json`))

		response, err := client.RetrieveGDPR(context.Background(), testGDPREmail)
		require.Error(t, err)
		assert.Nil(t, response)
	})
}

// TestClient_RetrieveGDPRWithRequest tests the method RetrieveGDPRWithRequest()
func TestClient_RetrieveGDPRWithRequest(t *testing.T) {
	t.Parallel()

	t.Run("retrieve GDPR data with request struct", func(t *testing.T) {
		client := newTestClient(mockGDPRRetrieval())

		response, err := client.RetrieveGDPRWithRequest(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.Data)
		assert.Equal(t, testGDPRRetrieveMsg, response.Data.Message)
	})

	t.Run("returns error when request is nil", func(t *testing.T) {
		client := newTestClient(mockGDPRRetrieval())

		response, err := client.RetrieveGDPRWithRequest(context.Background(), nil)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingEmail)
	})

	t.Run("returns error when email is empty in request", func(t *testing.T) {
		client := newTestClient(mockGDPRRetrieval())

		response, err := client.RetrieveGDPRWithRequest(context.Background(), &GDPRRequest{
			Email: "",
		})
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingEmail)
	})
}

// TestClient_RetrieveGDPRRaw tests the method RetrieveGDPRRaw()
func TestClient_RetrieveGDPRRaw(t *testing.T) {
	t.Parallel()

	t.Run("retrieves GDPR data successfully", func(t *testing.T) {
		client := newTestClient(mockGDPRRetrieval())

		response, err := client.RetrieveGDPRRaw(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodPost, response.Method)
	})

	t.Run("returns error when request is nil", func(t *testing.T) {
		client := newTestClient(mockGDPRRetrieval())

		response, err := client.RetrieveGDPRRaw(context.Background(), nil)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingEmail)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		response, err := client.RetrieveGDPRRaw(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("uses correct endpoint URL", func(t *testing.T) {
		client := newTestClient(mockGDPRRetrieval())

		response, err := client.RetrieveGDPRRaw(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.NoError(t, err)
		assert.Contains(t, response.URL, "/gdpr/retrieve")
	})

	t.Run("uses POST method", func(t *testing.T) {
		client := newTestClient(mockGDPRRetrieval())

		response, err := client.RetrieveGDPRRaw(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, response.Method)
	})

	t.Run("includes email in post data", func(t *testing.T) {
		client := newTestClient(mockGDPRRetrieval())

		response, err := client.RetrieveGDPRRaw(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.NoError(t, err)
		assert.Contains(t, response.PostData, testGDPREmail)
	})
}

// BenchmarkClient_RetrieveGDPR benchmarks the RetrieveGDPR method
func BenchmarkClient_RetrieveGDPR(b *testing.B) {
	client := newTestClient(mockGDPRRetrieval())
	for i := 0; i < b.N; i++ {
		_, _ = client.RetrieveGDPR(context.Background(), testGDPREmail)
	}
}

// BenchmarkClient_RetrieveGDPRRaw benchmarks the RetrieveGDPRRaw method
func BenchmarkClient_RetrieveGDPRRaw(b *testing.B) {
	client := newTestClient(mockGDPRRetrieval())
	request := &GDPRRequest{Email: testGDPREmail}
	for i := 0; i < b.N; i++ {
		_, _ = client.RetrieveGDPRRaw(context.Background(), request)
	}
}
