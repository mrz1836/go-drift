package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUserID             = 228225
	testUserIDBadRequest   = 111111
	testUserIDUnauthorized = 222222
	testUserIDBadJSON      = 333333
	testUserName           = "Test User"
	testUserEmail          = "testuser@example.com"
	testUserAlias          = "tuser"
	testUserPhone          = "555-123-4567"
)

// mockGetUser returns a multi-route mock for user GET operations
func mockGetUser() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/users/228225", http.StatusOK,
			`{"data":{"id":228225,"orgId":12345,"name":"Test User","alias":"tuser","email":"testuser@example.com","phone":"555-123-4567","locale":"en-US","availability":"AVAILABLE","role":"admin","timeZone":"America/New_York","avatarUrl":"https://example.com/avatar.png","verified":true,"bot":false,"createdAt":1606273669631,"updatedAt":1614550516644}}`).
		addRoute(apiEndpoint+"/users/111111", http.StatusBadRequest, "").
		addRoute(apiEndpoint+"/users/222222", http.StatusUnauthorized, "").
		addRoute(apiEndpoint+"/users/333333", http.StatusOK,
			`{"data":{"id":333333"name":"Bad JSON"}}`)
}

// mockGetUsers returns a multi-route mock for multiple user GET operations
func mockGetUsers() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/users?userId=228225&userId=243266", http.StatusOK,
			`{"data":{"228225":{"id":228225,"orgId":12345,"name":"Test User","email":"testuser@example.com","availability":"AVAILABLE","role":"admin","verified":true,"bot":false,"createdAt":1606273669631,"updatedAt":1614550516644},"243266":{"id":243266,"orgId":12345,"name":"Second User","email":"second@example.com","availability":"OFFLINE","role":"member","verified":true,"bot":false,"createdAt":1606273669631,"updatedAt":1614550516644}}}`).
		addRoute(apiEndpoint+"/users?userId=111111", http.StatusBadRequest, "").
		addRoute(apiEndpoint+"/users?userId=333333", http.StatusOK,
			`{"data":{"333333":{"id":333333"name":"Bad JSON"}}}`)
}

// TestClient_GetUser tests the method GetUser()
func TestClient_GetUser(t *testing.T) {
	t.Parallel()

	t.Run("get a valid user by id", func(t *testing.T) {
		client := newTestClient(mockGetUser())

		user, err := client.GetUser(context.Background(), testUserID)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotNil(t, user.Data)

		// Check returned values
		assert.Equal(t, uint64(228225), user.Data.ID)
		assert.Equal(t, uint64(12345), user.Data.OrgID)
		assert.Equal(t, testUserName, user.Data.Name)
		assert.Equal(t, testUserAlias, user.Data.Alias)
		assert.Equal(t, testUserEmail, user.Data.Email)
		assert.Equal(t, testUserPhone, user.Data.Phone)
		assert.Equal(t, "en-US", user.Data.Locale)
		assert.Equal(t, "AVAILABLE", user.Data.Availability)
		assert.Equal(t, "admin", user.Data.Role)
		assert.Equal(t, "America/New_York", user.Data.TimeZone)
		assert.Equal(t, "https://example.com/avatar.png", user.Data.AvatarURL)
		assert.True(t, user.Data.Verified)
		assert.False(t, user.Data.Bot)
		assert.Equal(t, int64(1606273669631), user.Data.CreatedAt)
		assert.Equal(t, int64(1614550516644), user.Data.UpdatedAt)
	})

	t.Run("missing user id", func(t *testing.T) {
		client := newTestClient(mockGetUser())

		user, err := client.GetUser(context.Background(), 0)
		require.Error(t, err)
		assert.Equal(t, ErrMissingUserID, err)
		assert.Nil(t, user)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(mockGetUser())

		user, err := client.GetUser(context.Background(), testUserIDBadRequest)
		require.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(mockGetUser())

		user, err := client.GetUser(context.Background(), testUserIDUnauthorized)
		require.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockGetUser())

		user, err := client.GetUser(context.Background(), testUserIDBadJSON)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

// TestClient_GetUserRaw tests the method GetUserRaw()
func TestClient_GetUserRaw(t *testing.T) {
	t.Parallel()

	t.Run("missing user id", func(t *testing.T) {
		client := newTestClient(mockGetUser())

		response, err := client.GetUserRaw(context.Background(), 0)
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingUserID, err)
	})

	t.Run("get a valid user by id", func(t *testing.T) {
		client := newTestClient(mockGetUser())

		response, err := client.GetUserRaw(context.Background(), testUserID)
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/users/228225", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

// TestClient_GetUsers tests the method GetUsers()
func TestClient_GetUsers(t *testing.T) {
	t.Parallel()

	t.Run("get multiple valid users", func(t *testing.T) {
		client := newTestClient(mockGetUsers())

		users, err := client.GetUsers(context.Background(), []uint64{228225, 243266})
		require.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users.Data, 2)
	})

	t.Run("missing user ids", func(t *testing.T) {
		client := newTestClient(mockGetUsers())

		users, err := client.GetUsers(context.Background(), []uint64{})
		require.Error(t, err)
		assert.Equal(t, ErrMissingUserID, err)
		assert.Nil(t, users)
	})

	t.Run("too many user ids", func(t *testing.T) {
		client := newTestClient(mockGetUsers())

		// Create 21 user IDs
		userIDs := make([]uint64, 21)
		for i := range userIDs {
			userIDs[i] = uint64(i) + 1
		}

		users, err := client.GetUsers(context.Background(), userIDs)
		require.Error(t, err)
		assert.Equal(t, ErrTooManyUserIDs, err)
		assert.Nil(t, users)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(mockGetUsers())

		users, err := client.GetUsers(context.Background(), []uint64{111111})
		require.Error(t, err)
		assert.Nil(t, users)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockGetUsers())

		users, err := client.GetUsers(context.Background(), []uint64{333333})
		require.Error(t, err)
		assert.Nil(t, users)
	})
}

// TestClient_GetUsersRaw tests the method GetUsersRaw()
func TestClient_GetUsersRaw(t *testing.T) {
	t.Parallel()

	t.Run("missing user ids", func(t *testing.T) {
		client := newTestClient(mockGetUsers())

		response, err := client.GetUsersRaw(context.Background(), []uint64{})
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingUserID, err)
	})

	t.Run("too many user ids", func(t *testing.T) {
		client := newTestClient(mockGetUsers())

		userIDs := make([]uint64, 21)
		for i := range userIDs {
			userIDs[i] = uint64(i) + 1
		}

		response, err := client.GetUsersRaw(context.Background(), userIDs)
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrTooManyUserIDs, err)
	})

	t.Run("get multiple valid users", func(t *testing.T) {
		client := newTestClient(mockGetUsers())

		response, err := client.GetUsersRaw(context.Background(), []uint64{228225, 243266})
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/users?userId=228225&userId=243266", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

// BenchmarkClient_GetUser benchmarks the GetUser method
func BenchmarkClient_GetUser(b *testing.B) {
	client := newTestClient(mockGetUser())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetUser(context.Background(), testUserID)
	}
}

// BenchmarkClient_GetUserRaw benchmarks the GetUserRaw method
func BenchmarkClient_GetUserRaw(b *testing.B) {
	client := newTestClient(mockGetUser())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetUserRaw(context.Background(), testUserID)
	}
}

// BenchmarkClient_GetUsers benchmarks the GetUsers method
func BenchmarkClient_GetUsers(b *testing.B) {
	client := newTestClient(mockGetUsers())
	userIDs := []uint64{228225, 243266}
	for i := 0; i < b.N; i++ {
		_, _ = client.GetUsers(context.Background(), userIDs)
	}
}

// BenchmarkClient_GetUsersRaw benchmarks the GetUsersRaw method
func BenchmarkClient_GetUsersRaw(b *testing.B) {
	client := newTestClient(mockGetUsers())
	userIDs := []uint64{228225, 243266}
	for i := 0; i < b.N; i++ {
		_, _ = client.GetUsersRaw(context.Background(), userIDs)
	}
}
