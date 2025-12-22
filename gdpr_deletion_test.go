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
	testGDPRDeleteMsg = "Your delete is processing"
)

// mockHTTPGDPRDeletion for mocking GDPR deletion requests
type mockHTTPGDPRDeletion struct{}

// Do is a mock http request
func (m *mockHTTPGDPRDeletion) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	if req == nil {
		return resp, errMissingRequest
	}

	if req.URL.String() == apiEndpoint+"/gdpr/delete" && req.Method == http.MethodPost {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`{"data":{"message":"` + testGDPRDeleteMsg + `"}}`))
	}

	return resp, nil
}

// mockHTTPGDPRDeletionError for testing error scenarios
type mockHTTPGDPRDeletionError struct {
	statusCode int
	body       string
}

// Do returns a configurable error response
func (m *mockHTTPGDPRDeletionError) Do(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errMissingRequest
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(m.body)),
	}, nil
}

// TestClient_DeleteGDPR tests the method DeleteGDPR()
func TestClient_DeleteGDPR(t *testing.T) {
	t.Parallel()

	t.Run("delete GDPR data successfully", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletion{})

		response, err := client.DeleteGDPR(context.Background(), testGDPREmail)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.Data)
		assert.Equal(t, testGDPRDeleteMsg, response.Data.Message)
	})

	t.Run("returns error when email is empty", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletion{})

		response, err := client.DeleteGDPR(context.Background(), "")
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingEmail)
	})

	t.Run("returns error on 400 bad request", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletionError{
			statusCode: http.StatusBadRequest,
			body:       "",
		})

		response, err := client.DeleteGDPR(context.Background(), testGDPREmail)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletionError{
			statusCode: http.StatusUnauthorized,
			body:       "",
		})

		response, err := client.DeleteGDPR(context.Background(), testGDPREmail)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on 404 not found", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletionError{
			statusCode: http.StatusNotFound,
			body:       "",
		})

		response, err := client.DeleteGDPR(context.Background(), testGDPREmail)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrResourceNotFound)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletionError{
			statusCode: http.StatusOK,
			body:       `{"invalid json`,
		})

		response, err := client.DeleteGDPR(context.Background(), testGDPREmail)
		require.Error(t, err)
		assert.Nil(t, response)
	})
}

// TestClient_DeleteGDPRWithRequest tests the method DeleteGDPRWithRequest()
func TestClient_DeleteGDPRWithRequest(t *testing.T) {
	t.Parallel()

	t.Run("delete GDPR data with request struct", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletion{})

		response, err := client.DeleteGDPRWithRequest(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.Data)
		assert.Equal(t, testGDPRDeleteMsg, response.Data.Message)
	})

	t.Run("returns error when request is nil", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletion{})

		response, err := client.DeleteGDPRWithRequest(context.Background(), nil)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingEmail)
	})

	t.Run("returns error when email is empty in request", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletion{})

		response, err := client.DeleteGDPRWithRequest(context.Background(), &GDPRRequest{
			Email: "",
		})
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingEmail)
	})
}

// TestClient_DeleteGDPRRaw tests the method DeleteGDPRRaw()
func TestClient_DeleteGDPRRaw(t *testing.T) {
	t.Parallel()

	t.Run("deletes GDPR data successfully", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletion{})

		response, err := client.DeleteGDPRRaw(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodPost, response.Method)
	})

	t.Run("returns error when request is nil", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletion{})

		response, err := client.DeleteGDPRRaw(context.Background(), nil)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingEmail)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletionError{
			statusCode: http.StatusBadRequest,
			body:       "",
		})

		response, err := client.DeleteGDPRRaw(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("uses correct endpoint URL", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletion{})

		response, err := client.DeleteGDPRRaw(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.NoError(t, err)
		assert.Contains(t, response.URL, "/gdpr/delete")
	})

	t.Run("uses POST method", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletion{})

		response, err := client.DeleteGDPRRaw(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, response.Method)
	})

	t.Run("includes email in post data", func(t *testing.T) {
		client := newTestClient(&mockHTTPGDPRDeletion{})

		response, err := client.DeleteGDPRRaw(context.Background(), &GDPRRequest{
			Email: testGDPREmail,
		})
		require.NoError(t, err)
		assert.Contains(t, response.PostData, testGDPREmail)
	})
}

// BenchmarkClient_DeleteGDPR benchmarks the DeleteGDPR method
func BenchmarkClient_DeleteGDPR(b *testing.B) {
	client := newTestClient(&mockHTTPGDPRDeletion{})
	for i := 0; i < b.N; i++ {
		_, _ = client.DeleteGDPR(context.Background(), testGDPREmail)
	}
}

// BenchmarkClient_DeleteGDPRRaw benchmarks the DeleteGDPRRaw method
func BenchmarkClient_DeleteGDPRRaw(b *testing.B) {
	client := newTestClient(&mockHTTPGDPRDeletion{})
	request := &GDPRRequest{Email: testGDPREmail}
	for i := 0; i < b.N; i++ {
		_, _ = client.DeleteGDPRRaw(context.Background(), request)
	}
}
