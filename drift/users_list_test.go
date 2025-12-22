package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockListUsers returns a mock for user list operations
func mockListUsers() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":[{"id":228225,"orgId":12345,"name":"Test User","alias":"tuser","email":"testuser@example.com","phone":"555-123-4567","locale":"en-US","availability":"AVAILABLE","role":"admin","timeZone":"America/New_York","avatarUrl":"https://example.com/avatar.png","verified":true,"bot":false,"createdAt":1606273669631,"updatedAt":1614550516644},{"id":243266,"orgId":12345,"name":"Second User","email":"second@example.com","availability":"OFFLINE","role":"member","verified":true,"bot":false,"createdAt":1606273669631,"updatedAt":1614550516644}]}`),
	)
}

// mockListUsersEmpty returns a mock for empty user list
func mockListUsersEmpty() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":[]}`),
	)
}

// TestClient_ListUsers tests the method ListUsers()
func TestClient_ListUsers(t *testing.T) {
	t.Parallel()

	t.Run("list all users", func(t *testing.T) {
		client := newTestClient(mockListUsers())

		users, err := client.ListUsers(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users.Data, 2)

		// Check first user
		assert.Equal(t, uint64(228225), users.Data[0].ID)
		assert.Equal(t, "Test User", users.Data[0].Name)
		assert.Equal(t, "testuser@example.com", users.Data[0].Email)
		assert.Equal(t, "AVAILABLE", users.Data[0].Availability)
		assert.Equal(t, "admin", users.Data[0].Role)

		// Check second user
		assert.Equal(t, uint64(243266), users.Data[1].ID)
		assert.Equal(t, "Second User", users.Data[1].Name)
		assert.Equal(t, "OFFLINE", users.Data[1].Availability)
	})

	t.Run("empty user list", func(t *testing.T) {
		client := newTestClient(mockListUsersEmpty())

		users, err := client.ListUsers(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, users)
		assert.Empty(t, users.Data)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		users, err := client.ListUsers(context.Background())
		require.Error(t, err)
		assert.Nil(t, users)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"data":[{"id":228225"name":"Bad JSON"}]}`))

		users, err := client.ListUsers(context.Background())
		require.Error(t, err)
		assert.Nil(t, users)
	})
}

// TestClient_ListUsersRaw tests the method ListUsersRaw()
func TestClient_ListUsersRaw(t *testing.T) {
	t.Parallel()

	t.Run("list all users raw", func(t *testing.T) {
		client := newTestClient(mockListUsers())

		response, err := client.ListUsersRaw(context.Background())
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/users/list", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		response, err := client.ListUsersRaw(context.Background())
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})
}

// BenchmarkClient_ListUsers benchmarks the ListUsers method
func BenchmarkClient_ListUsers(b *testing.B) {
	client := newTestClient(mockListUsers())
	for i := 0; i < b.N; i++ {
		_, _ = client.ListUsers(context.Background())
	}
}

// BenchmarkClient_ListUsersRaw benchmarks the ListUsersRaw method
func BenchmarkClient_ListUsersRaw(b *testing.B) {
	client := newTestClient(mockListUsers())
	for i := 0; i < b.N; i++ {
		_, _ = client.ListUsersRaw(context.Background())
	}
}
