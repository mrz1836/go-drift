package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAccountIDNotFound = "999999_notfound.com"
)

// mockGetAccount returns a multi-route mock for account GET operations
func mockGetAccount() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/accounts/"+testAccountID, http.StatusOK,
			`{"data":{"accountId":"`+testAccountID+`","ownerId":21965,"name":"`+testAccountName+`","domain":"`+testAccountDomain+`","deleted":false,"targeted":true,"createDateTime":1614563742010,"updateDateTime":1614563742010}}`).
		addRoute(apiEndpoint+"/accounts/"+testAccountIDNotFound, http.StatusNotFound,
			`{"error":"not found"}`)
}

// TestClient_GetAccount tests the method GetAccount()
func TestClient_GetAccount(t *testing.T) {
	t.Parallel()

	t.Run("get an account successfully", func(t *testing.T) {
		client := newTestClient(mockGetAccount())

		account, err := client.GetAccount(context.Background(), testAccountID)
		require.NoError(t, err)
		assert.NotNil(t, account)

		assert.Equal(t, testAccountID, account.Data.AccountID)
		assert.Equal(t, testAccountOwner, account.Data.OwnerID)
		assert.Equal(t, testAccountName, account.Data.Name)
		assert.Equal(t, testAccountDomain, account.Data.Domain)
		assert.False(t, account.Data.Deleted)
		assert.True(t, account.Data.Targeted)
	})

	t.Run("returns error on missing account id", func(t *testing.T) {
		client := newTestClient(mockGetAccount())

		account, err := client.GetAccount(context.Background(), "")
		require.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrMissingAccountID)
	})

	t.Run("returns error on 404 not found", func(t *testing.T) {
		client := newTestClient(mockGetAccount())

		account, err := client.GetAccount(context.Background(), testAccountIDNotFound)
		require.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrResourceNotFound)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		account, err := client.GetAccount(context.Background(), testAccountID)
		require.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"data":{"invalid json`))

		account, err := client.GetAccount(context.Background(), testAccountID)
		require.Error(t, err)
		assert.Nil(t, account)
	})
}

// TestClient_GetAccountRaw tests the method GetAccountRaw()
func TestClient_GetAccountRaw(t *testing.T) {
	t.Parallel()

	t.Run("gets account successfully", func(t *testing.T) {
		client := newTestClient(mockGetAccount())

		response, err := client.GetAccountRaw(context.Background(), testAccountID)
		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, apiEndpoint+"/accounts/"+testAccountID, response.URL)
	})

	t.Run("returns error on missing account id", func(t *testing.T) {
		client := newTestClient(mockGetAccount())

		response, err := client.GetAccountRaw(context.Background(), "")
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingAccountID)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		response, err := client.GetAccountRaw(context.Background(), testAccountID)
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

// BenchmarkClient_GetAccount benchmarks the GetAccount method
func BenchmarkClient_GetAccount(b *testing.B) {
	client := newTestClient(mockGetAccount())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetAccount(context.Background(), testAccountID)
	}
}
