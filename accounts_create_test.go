package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAccountID     = "123458_domain.com"
	testAccountName   = "Test Company"
	testAccountDomain = "domain.com"
	testAccountOwner  = uint64(21965)
)

// mockCreateAccount returns a mock for account creation operations
func mockCreateAccount() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":{"accountId":"`+testAccountID+`","ownerId":21965,"name":"`+testAccountName+`","domain":"`+testAccountDomain+`","deleted":false,"targeted":true,"createDateTime":1614563742010,"updateDateTime":1614563742010}}`),
	)
}

// TestClient_CreateAccount tests the method CreateAccount()
func TestClient_CreateAccount(t *testing.T) {
	t.Parallel()

	t.Run("create an account successfully", func(t *testing.T) {
		client := newTestClient(mockCreateAccount())

		account, err := client.CreateAccount(
			context.Background(),
			&AccountFields{
				OwnerID:  testAccountOwner,
				Name:     testAccountName,
				Domain:   testAccountDomain,
				Targeted: true,
			})
		require.NoError(t, err)
		assert.NotNil(t, account)

		assert.Equal(t, testAccountID, account.Data.AccountID)
		assert.Equal(t, testAccountOwner, account.Data.OwnerID)
		assert.Equal(t, testAccountName, account.Data.Name)
		assert.Equal(t, testAccountDomain, account.Data.Domain)
		assert.False(t, account.Data.Deleted)
		assert.True(t, account.Data.Targeted)
		assert.Equal(t, int64(1614563742010), account.Data.CreateDateTime)
	})

	t.Run("returns error on 400 bad request", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		account, err := client.CreateAccount(
			context.Background(),
			&AccountFields{
				OwnerID: testAccountOwner,
				Name:    testAccountName,
			})

		require.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		account, err := client.CreateAccount(
			context.Background(),
			&AccountFields{
				OwnerID: testAccountOwner,
				Name:    testAccountName,
			})

		require.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on 409 conflict", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusConflict))

		account, err := client.CreateAccount(
			context.Background(),
			&AccountFields{
				OwnerID: testAccountOwner,
				Domain:  testAccountDomain,
			})

		require.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrConflict)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"data":{"invalid json`))

		account, err := client.CreateAccount(
			context.Background(),
			&AccountFields{
				OwnerID: testAccountOwner,
				Name:    testAccountName,
			})

		require.Error(t, err)
		assert.Nil(t, account)
	})
}

// TestClient_CreateAccountRaw tests the method CreateAccountRaw()
func TestClient_CreateAccountRaw(t *testing.T) {
	t.Parallel()

	t.Run("creates account successfully", func(t *testing.T) {
		client := newTestClient(mockCreateAccount())

		response, err := client.CreateAccountRaw(
			context.Background(),
			&AccountFields{
				OwnerID:  testAccountOwner,
				Name:     testAccountName,
				Domain:   testAccountDomain,
				Targeted: true,
			})

		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodPost, response.Method)
		assert.Equal(t, apiEndpoint+"/accounts/create", response.URL)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		response, err := client.CreateAccountRaw(
			context.Background(),
			&AccountFields{
				OwnerID: testAccountOwner,
				Name:    testAccountName,
			})

		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})
}

// BenchmarkClient_CreateAccount benchmarks the CreateAccount method
func BenchmarkClient_CreateAccount(b *testing.B) {
	client := newTestClient(mockCreateAccount())
	fields := &AccountFields{
		OwnerID:  testAccountOwner,
		Name:     testAccountName,
		Domain:   testAccountDomain,
		Targeted: true,
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateAccount(context.Background(), fields)
	}
}
