package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockDeleteAccount returns a multi-route mock for account deletion operations
func mockDeleteAccount() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/accounts/"+testAccountID, http.StatusOK, `{"result":"OK","ok":true}`).
		addRoute(apiEndpoint+"/accounts/"+testAccountIDNotFound, http.StatusNotFound, `{"error":"not found"}`)
}

// TestClient_DeleteAccount tests the method DeleteAccount()
func TestClient_DeleteAccount(t *testing.T) {
	t.Parallel()

	t.Run("delete an account successfully", func(t *testing.T) {
		client := newTestClient(mockDeleteAccount())

		response, err := client.DeleteAccount(context.Background(), testAccountID)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.OK)
		assert.Equal(t, "OK", response.Result)
	})

	t.Run("returns error on missing account id", func(t *testing.T) {
		client := newTestClient(mockDeleteAccount())

		response, err := client.DeleteAccount(context.Background(), "")
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingAccountID)
	})

	t.Run("returns error on 404 not found", func(t *testing.T) {
		client := newTestClient(mockDeleteAccount())

		response, err := client.DeleteAccount(context.Background(), testAccountIDNotFound)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrResourceNotFound)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		response, err := client.DeleteAccount(context.Background(), testAccountID)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"invalid json`))

		response, err := client.DeleteAccount(context.Background(), testAccountID)
		require.Error(t, err)
		assert.Nil(t, response)
	})
}

// TestClient_DeleteAccountRaw tests the method DeleteAccountRaw()
func TestClient_DeleteAccountRaw(t *testing.T) {
	t.Parallel()

	t.Run("deletes account successfully", func(t *testing.T) {
		client := newTestClient(mockDeleteAccount())

		response, err := client.DeleteAccountRaw(context.Background(), testAccountID)
		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodDelete, response.Method)
		assert.Equal(t, apiEndpoint+"/accounts/"+testAccountID, response.URL)
	})

	t.Run("returns error on missing account id", func(t *testing.T) {
		client := newTestClient(mockDeleteAccount())

		response, err := client.DeleteAccountRaw(context.Background(), "")
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingAccountID)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		response, err := client.DeleteAccountRaw(context.Background(), testAccountID)
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

// BenchmarkClient_DeleteAccount benchmarks the DeleteAccount method
func BenchmarkClient_DeleteAccount(b *testing.B) {
	client := newTestClient(mockDeleteAccount())
	for i := 0; i < b.N; i++ {
		_, _ = client.DeleteAccount(context.Background(), testAccountID)
	}
}
