package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testClientID     = "test-client-id"
	testClientSecret = "test-client-secret"
)

// mockAppUninstall returns a mock for app uninstall requests
func mockAppUninstall() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"result":"OK","ok":true}`),
	)
}

// TestClient_AppUninstall tests the method AppUninstall()
func TestClient_AppUninstall(t *testing.T) {
	t.Parallel()

	t.Run("uninstall app successfully", func(t *testing.T) {
		client := newTestClient(mockAppUninstall())

		response, err := client.AppUninstall(context.Background(), testClientID, testClientSecret)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.OK)
		assert.Equal(t, "OK", response.Result)
	})

	t.Run("returns error when missing client id", func(t *testing.T) {
		client := newTestClient(mockAppUninstall())

		response, err := client.AppUninstall(context.Background(), "", testClientSecret)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingClientID)
	})

	t.Run("returns error when missing client secret", func(t *testing.T) {
		client := newTestClient(mockAppUninstall())

		response, err := client.AppUninstall(context.Background(), testClientID, "")
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingClientSecret)
	})

	t.Run("returns error when AppUninstallRaw fails", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		response, err := client.AppUninstall(context.Background(), testClientID, testClientSecret)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		response, err := client.AppUninstall(context.Background(), testClientID, testClientSecret)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"invalid json`))

		response, err := client.AppUninstall(context.Background(), testClientID, testClientSecret)
		require.Error(t, err)
		assert.Nil(t, response)
	})
}

// TestClient_AppUninstallRaw tests the method AppUninstallRaw()
func TestClient_AppUninstallRaw(t *testing.T) {
	t.Parallel()

	t.Run("uninstalls app successfully", func(t *testing.T) {
		client := newTestClient(mockAppUninstall())

		response, err := client.AppUninstallRaw(context.Background(), testClientID, testClientSecret)
		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodPost, response.Method)
	})

	t.Run("returns error when missing client id", func(t *testing.T) {
		client := newTestClient(mockAppUninstall())

		response, err := client.AppUninstallRaw(context.Background(), "", testClientSecret)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingClientID)
	})

	t.Run("returns error when missing client secret", func(t *testing.T) {
		client := newTestClient(mockAppUninstall())

		response, err := client.AppUninstallRaw(context.Background(), testClientID, "")
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingClientSecret)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		response, err := client.AppUninstallRaw(context.Background(), testClientID, testClientSecret)
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("uses correct endpoint URL with query params", func(t *testing.T) {
		client := newTestClient(mockAppUninstall())

		response, err := client.AppUninstallRaw(context.Background(), testClientID, testClientSecret)
		require.NoError(t, err)
		assert.Contains(t, response.URL, "/app/uninstall")
		assert.Contains(t, response.URL, "clientId="+testClientID)
		assert.Contains(t, response.URL, "clientSecret="+testClientSecret)
	})

	t.Run("uses POST method", func(t *testing.T) {
		client := newTestClient(mockAppUninstall())

		response, err := client.AppUninstallRaw(context.Background(), testClientID, testClientSecret)
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, response.Method)
	})
}

// BenchmarkClient_AppUninstall benchmarks the AppUninstall method
func BenchmarkClient_AppUninstall(b *testing.B) {
	client := newTestClient(mockAppUninstall())
	for i := 0; i < b.N; i++ {
		_, _ = client.AppUninstall(context.Background(), testClientID, testClientSecret)
	}
}

// BenchmarkClient_AppUninstallRaw benchmarks the AppUninstallRaw method
func BenchmarkClient_AppUninstallRaw(b *testing.B) {
	client := newTestClient(mockAppUninstall())
	for i := 0; i < b.N; i++ {
		_, _ = client.AppUninstallRaw(context.Background(), testClientID, testClientSecret)
	}
}
