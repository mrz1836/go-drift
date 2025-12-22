package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockListAccounts returns a mock for listing accounts with pagination
func mockListAccounts() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":{"accounts":[{"accountId":"`+testAccountID+`","ownerId":21965,"name":"`+testAccountName+`","domain":"`+testAccountDomain+`","deleted":false,"targeted":true,"createDateTime":1614563742010,"updateDateTime":1614563742010}],"total":1,"next":"/accounts?index=10&size=10"}}`),
	)
}

// mockListAccountsNoNext returns a mock for listing accounts without next page
func mockListAccountsNoNext() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":{"accounts":[{"accountId":"`+testAccountID+`","ownerId":21965,"name":"`+testAccountName+`","domain":"`+testAccountDomain+`"}],"total":1}}`),
	)
}

// TestClient_ListAccounts tests the method ListAccounts()
func TestClient_ListAccounts(t *testing.T) {
	t.Parallel()

	t.Run("list accounts successfully", func(t *testing.T) {
		client := newTestClient(mockListAccounts())

		accounts, err := client.ListAccounts(context.Background(), nil)
		require.NoError(t, err)
		assert.NotNil(t, accounts)
		assert.NotNil(t, accounts.Data)
		assert.Len(t, accounts.Data.Accounts, 1)
		assert.Equal(t, 1, accounts.Data.Total)
		assert.Equal(t, "/accounts?index=10&size=10", accounts.Data.Next)

		account := accounts.Data.Accounts[0]
		assert.Equal(t, testAccountID, account.AccountID)
		assert.Equal(t, testAccountOwner, account.OwnerID)
		assert.Equal(t, testAccountName, account.Name)
	})

	t.Run("list accounts with pagination parameters", func(t *testing.T) {
		client := newTestClient(mockListAccounts())

		accounts, err := client.ListAccounts(context.Background(), &AccountListQuery{
			Index: 0,
			Size:  25,
		})
		require.NoError(t, err)
		assert.NotNil(t, accounts)
		assert.NotNil(t, accounts.Data)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		accounts, err := client.ListAccounts(context.Background(), nil)
		require.Error(t, err)
		assert.Nil(t, accounts)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"data":{"invalid json`))

		accounts, err := client.ListAccounts(context.Background(), nil)
		require.Error(t, err)
		assert.Nil(t, accounts)
	})
}

// TestClient_ListAccountsRaw tests the method ListAccountsRaw()
func TestClient_ListAccountsRaw(t *testing.T) {
	t.Parallel()

	t.Run("lists accounts successfully", func(t *testing.T) {
		client := newTestClient(mockListAccounts())

		response, err := client.ListAccountsRaw(context.Background(), nil)
		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodGet, response.Method)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		response, err := client.ListAccountsRaw(context.Background(), nil)
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

// TestClient_ListAccountsNext tests the method ListAccountsNext()
func TestClient_ListAccountsNext(t *testing.T) {
	t.Parallel()

	t.Run("returns ErrNoNextPage when accounts is nil", func(t *testing.T) {
		client := newTestClient(mockListAccounts())

		nextAccounts, err := client.ListAccountsNext(context.Background(), nil)
		require.ErrorIs(t, err, ErrNoNextPage)
		assert.Nil(t, nextAccounts)
	})

	t.Run("returns ErrNoNextPage when data is nil", func(t *testing.T) {
		client := newTestClient(mockListAccounts())

		nextAccounts, err := client.ListAccountsNext(context.Background(), &Accounts{})
		require.ErrorIs(t, err, ErrNoNextPage)
		assert.Nil(t, nextAccounts)
	})

	t.Run("returns ErrNoNextPage when next is empty", func(t *testing.T) {
		client := newTestClient(mockListAccountsNoNext())

		accounts, err := client.ListAccounts(context.Background(), nil)
		require.NoError(t, err)

		nextAccounts, err := client.ListAccountsNext(context.Background(), accounts)
		require.ErrorIs(t, err, ErrNoNextPage)
		assert.Nil(t, nextAccounts)
	})

	t.Run("fetches next page successfully", func(t *testing.T) {
		client := newTestClient(mockListAccounts())

		accounts, err := client.ListAccounts(context.Background(), nil)
		require.NoError(t, err)
		require.NotNil(t, accounts)
		require.NotNil(t, accounts.Data)
		assert.NotEmpty(t, accounts.Data.Next)

		nextAccounts, err := client.ListAccountsNext(context.Background(), accounts)
		require.NoError(t, err)
		assert.NotNil(t, nextAccounts)
	})
}

// TestAccountListQuery_BuildURL tests the BuildURL method
func TestAccountListQuery_BuildURL(t *testing.T) {
	t.Parallel()

	t.Run("returns base URL when query is nil", func(t *testing.T) {
		var query *AccountListQuery
		url := query.BuildURL()
		assert.Equal(t, apiEndpoint+"/accounts", url)
	})

	t.Run("returns base URL when no params set", func(t *testing.T) {
		query := &AccountListQuery{}
		url := query.BuildURL()
		assert.Equal(t, apiEndpoint+"/accounts", url)
	})

	t.Run("adds index parameter", func(t *testing.T) {
		query := &AccountListQuery{Index: 10}
		url := query.BuildURL()
		assert.Equal(t, apiEndpoint+"/accounts?index=10", url)
	})

	t.Run("adds size parameter", func(t *testing.T) {
		query := &AccountListQuery{Size: 25}
		url := query.BuildURL()
		assert.Equal(t, apiEndpoint+"/accounts?size=25", url)
	})

	t.Run("adds both parameters", func(t *testing.T) {
		query := &AccountListQuery{Index: 10, Size: 25}
		url := query.BuildURL()
		assert.Equal(t, apiEndpoint+"/accounts?index=10&size=25", url)
	})
}

// BenchmarkClient_ListAccounts benchmarks the ListAccounts method
func BenchmarkClient_ListAccounts(b *testing.B) {
	client := newTestClient(mockListAccounts())
	for i := 0; i < b.N; i++ {
		_, _ = client.ListAccounts(context.Background(), nil)
	}
}
