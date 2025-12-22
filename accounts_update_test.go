package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockUpdateAccount returns a mock for account update operations
func mockUpdateAccount() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":{"accountId":"`+testAccountID+`","ownerId":21965,"name":"Updated Company","domain":"`+testAccountDomain+`","deleted":false,"targeted":true,"createDateTime":1614563742010,"updateDateTime":1614563742020}}`),
	)
}

// TestClient_UpdateAccount tests the method UpdateAccount()
func TestClient_UpdateAccount(t *testing.T) {
	t.Parallel()

	t.Run("update an account successfully", func(t *testing.T) {
		client := newTestClient(mockUpdateAccount())

		account, err := client.UpdateAccount(
			context.Background(),
			&AccountFields{
				AccountID: testAccountID,
				OwnerID:   testAccountOwner,
				Name:      "Updated Company",
				Domain:    testAccountDomain,
				Targeted:  true,
			})
		require.NoError(t, err)
		assert.NotNil(t, account)

		assert.Equal(t, testAccountID, account.Data.AccountID)
		assert.Equal(t, testAccountOwner, account.Data.OwnerID)
		assert.Equal(t, "Updated Company", account.Data.Name)
		assert.Equal(t, testAccountDomain, account.Data.Domain)
		assert.True(t, account.Data.Targeted)
		assert.Equal(t, int64(1614563742020), account.Data.UpdateDateTime)
	})

	t.Run("returns error on missing account id", func(t *testing.T) {
		client := newTestClient(mockUpdateAccount())

		account, err := client.UpdateAccount(
			context.Background(),
			&AccountFields{
				OwnerID: testAccountOwner,
				Name:    "Updated Company",
			})

		require.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrMissingAccountID)
	})

	t.Run("returns error on nil fields", func(t *testing.T) {
		client := newTestClient(mockUpdateAccount())

		account, err := client.UpdateAccount(context.Background(), nil)

		require.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrMissingAccountID)
	})

	t.Run("returns error on missing owner id", func(t *testing.T) {
		client := newTestClient(mockUpdateAccount())

		account, err := client.UpdateAccount(
			context.Background(),
			&AccountFields{
				AccountID: testAccountID,
				Name:      "Updated Company",
			})

		require.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrMissingOwnerID)
	})

	t.Run("returns error on 400 bad request", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		account, err := client.UpdateAccount(
			context.Background(),
			&AccountFields{
				AccountID: testAccountID,
				OwnerID:   testAccountOwner,
				Name:      "Updated Company",
			})

		require.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		account, err := client.UpdateAccount(
			context.Background(),
			&AccountFields{
				AccountID: testAccountID,
				OwnerID:   testAccountOwner,
				Name:      "Updated Company",
			})

		require.Error(t, err)
		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"data":{"invalid json`))

		account, err := client.UpdateAccount(
			context.Background(),
			&AccountFields{
				AccountID: testAccountID,
				OwnerID:   testAccountOwner,
				Name:      "Updated Company",
			})

		require.Error(t, err)
		assert.Nil(t, account)
	})
}

// TestClient_UpdateAccountRaw tests the method UpdateAccountRaw()
func TestClient_UpdateAccountRaw(t *testing.T) {
	t.Parallel()

	t.Run("updates account successfully", func(t *testing.T) {
		client := newTestClient(mockUpdateAccount())

		response, err := client.UpdateAccountRaw(
			context.Background(),
			&AccountFields{
				AccountID: testAccountID,
				OwnerID:   testAccountOwner,
				Name:      "Updated Company",
				Domain:    testAccountDomain,
				Targeted:  true,
			})

		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodPatch, response.Method)
		assert.Equal(t, apiEndpoint+"/accounts/update", response.URL)
	})

	t.Run("returns error on missing account id", func(t *testing.T) {
		client := newTestClient(mockUpdateAccount())

		response, err := client.UpdateAccountRaw(
			context.Background(),
			&AccountFields{
				OwnerID: testAccountOwner,
				Name:    "Updated Company",
			})

		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingAccountID)
	})

	t.Run("returns error on missing owner id", func(t *testing.T) {
		client := newTestClient(mockUpdateAccount())

		response, err := client.UpdateAccountRaw(
			context.Background(),
			&AccountFields{
				AccountID: testAccountID,
				Name:      "Updated Company",
			})

		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingOwnerID)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		response, err := client.UpdateAccountRaw(
			context.Background(),
			&AccountFields{
				AccountID: testAccountID,
				OwnerID:   testAccountOwner,
				Name:      "Updated Company",
			})

		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})
}

// BenchmarkClient_UpdateAccount benchmarks the UpdateAccount method
func BenchmarkClient_UpdateAccount(b *testing.B) {
	client := newTestClient(mockUpdateAccount())
	fields := &AccountFields{
		AccountID: testAccountID,
		OwnerID:   testAccountOwner,
		Name:      "Updated Company",
		Domain:    testAccountDomain,
		Targeted:  true,
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.UpdateAccount(context.Background(), fields)
	}
}
