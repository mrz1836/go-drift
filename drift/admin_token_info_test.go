package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAccessToken         = "test-access-token"
	testAuthenticatedUserID = "orgId:12345"
	testCredentialID        = "test-app-id"
	testTokenID             = "test-token-id"
	testTokenScope          = "read write"
)

// mockTokenInfo returns a mock for token info requests
func mockTokenInfo() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{
			"access_token": "`+testAccessToken+`",
			"authenticated_userid": "`+testAuthenticatedUserID+`",
			"credential_id": "`+testCredentialID+`",
			"token_type": "bearer",
			"expires_in": 7200000,
			"created_at": 1609459200000,
			"scope": "`+testTokenScope+`",
			"id": "`+testTokenID+`"
		}`),
	)
}

// TestClient_GetTokenInfo tests the method GetTokenInfo()
func TestClient_GetTokenInfo(t *testing.T) {
	t.Parallel()

	t.Run("get token info successfully", func(t *testing.T) {
		client := newTestClient(mockTokenInfo())

		response, err := client.GetTokenInfo(context.Background(), testAccessToken)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, testAccessToken, response.AccessToken)
		assert.Equal(t, testAuthenticatedUserID, response.AuthenticatedUserID)
		assert.Equal(t, testCredentialID, response.CredentialID)
		assert.Equal(t, "bearer", response.TokenType)
		assert.Equal(t, int64(7200000), response.ExpiresIn)
		assert.Equal(t, int64(1609459200000), response.CreatedAt)
		assert.Equal(t, testTokenScope, response.Scope)
		assert.Equal(t, testTokenID, response.ID)
	})

	t.Run("returns error when missing access token", func(t *testing.T) {
		client := newTestClient(mockTokenInfo())

		response, err := client.GetTokenInfo(context.Background(), "")
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingAccessToken)
	})

	t.Run("returns error when GetTokenInfoRaw fails", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		response, err := client.GetTokenInfo(context.Background(), testAccessToken)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		response, err := client.GetTokenInfo(context.Background(), testAccessToken)
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"invalid json`))

		response, err := client.GetTokenInfo(context.Background(), testAccessToken)
		require.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("handles token with no expiration", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{
			"access_token": "` + testAccessToken + `",
			"authenticated_userid": "` + testAuthenticatedUserID + `",
			"credential_id": "` + testCredentialID + `",
			"token_type": "bearer",
			"expires_in": 0,
			"created_at": 1609459200000,
			"scope": "` + testTokenScope + `",
			"id": "` + testTokenID + `"
		}`))

		response, err := client.GetTokenInfo(context.Background(), testAccessToken)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int64(0), response.ExpiresIn)
	})
}

// TestClient_GetTokenInfoRaw tests the method GetTokenInfoRaw()
func TestClient_GetTokenInfoRaw(t *testing.T) {
	t.Parallel()

	t.Run("gets token info successfully", func(t *testing.T) {
		client := newTestClient(mockTokenInfo())

		response, err := client.GetTokenInfoRaw(context.Background(), testAccessToken)
		require.NoError(t, err)
		assert.NotNil(t, response)
		require.NoError(t, response.Error)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, http.MethodPost, response.Method)
	})

	t.Run("returns error when missing access token", func(t *testing.T) {
		client := newTestClient(mockTokenInfo())

		response, err := client.GetTokenInfoRaw(context.Background(), "")
		require.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, ErrMissingAccessToken)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		response, err := client.GetTokenInfoRaw(context.Background(), testAccessToken)
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("uses correct endpoint URL", func(t *testing.T) {
		client := newTestClient(mockTokenInfo())

		response, err := client.GetTokenInfoRaw(context.Background(), testAccessToken)
		require.NoError(t, err)
		assert.Equal(t, apiEndpoint+"/app/token_info", response.URL)
	})

	t.Run("uses POST method", func(t *testing.T) {
		client := newTestClient(mockTokenInfo())

		response, err := client.GetTokenInfoRaw(context.Background(), testAccessToken)
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, response.Method)
	})

	t.Run("sends access token in request body", func(t *testing.T) {
		client := newTestClient(mockTokenInfo())

		response, err := client.GetTokenInfoRaw(context.Background(), testAccessToken)
		require.NoError(t, err)
		assert.Contains(t, response.PostData, testAccessToken)
		assert.Contains(t, response.PostData, "access_token")
	})
}

// BenchmarkClient_GetTokenInfo benchmarks the GetTokenInfo method
func BenchmarkClient_GetTokenInfo(b *testing.B) {
	client := newTestClient(mockTokenInfo())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetTokenInfo(context.Background(), testAccessToken)
	}
}

// BenchmarkClient_GetTokenInfoRaw benchmarks the GetTokenInfoRaw method
func BenchmarkClient_GetTokenInfoRaw(b *testing.B) {
	client := newTestClient(mockTokenInfo())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetTokenInfoRaw(context.Background(), testAccessToken)
	}
}
