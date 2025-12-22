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

// mockHTTPUpdateUser for mocking requests
type mockHTTPUpdateUser struct{}

// Do is a mock http request
func (m *mockHTTPUpdateUser) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	if req == nil {
		return resp, errMissingRequest
	}

	// Update user endpoint
	if req.URL.String() == apiEndpoint+"/users/update?userId=228225" && req.Method == http.MethodPatch {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`{"data":{"id":228225,"orgId":12345,"name":"Updated User","alias":"updateduser","email":"updated@example.com","phone":"555-999-8888","locale":"en-US","availability":"OFFLINE","role":"admin","timeZone":"America/New_York","avatarUrl":"https://example.com/new-avatar.png","verified":true,"bot":false,"createdAt":1606273669631,"updatedAt":1614550516644}}`))
	} else if req.URL.String() == apiEndpoint+"/users/update?userId=111111" && req.Method == http.MethodPatch {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = io.NopCloser(nil)
	} else if req.URL.String() == apiEndpoint+"/users/update?userId=222222" && req.Method == http.MethodPatch {
		resp.StatusCode = http.StatusUnauthorized
		resp.Body = io.NopCloser(nil)
	} else if req.URL.String() == apiEndpoint+"/users/update?userId=333333" && req.Method == http.MethodPatch {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`{"data":{"id":333333"name":"Bad JSON"}}`))
	}

	return resp, nil
}

// TestClient_UpdateUser tests the method UpdateUser()
func TestClient_UpdateUser(t *testing.T) {
	t.Parallel()

	t.Run("update a valid user", func(t *testing.T) {
		client := newTestClient(&mockHTTPUpdateUser{})

		fields := &UserUpdateFields{
			Name:         "Updated User",
			Alias:        "updateduser",
			Email:        "updated@example.com",
			Phone:        "555-999-8888",
			Availability: "OFFLINE",
			AvatarURL:    "https://example.com/new-avatar.png",
		}

		user, err := client.UpdateUser(context.Background(), testUserID, fields)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotNil(t, user.Data)

		// Check returned values
		assert.Equal(t, uint64(228225), user.Data.ID)
		assert.Equal(t, "Updated User", user.Data.Name)
		assert.Equal(t, "updateduser", user.Data.Alias)
		assert.Equal(t, "updated@example.com", user.Data.Email)
		assert.Equal(t, "555-999-8888", user.Data.Phone)
		assert.Equal(t, "OFFLINE", user.Data.Availability)
		assert.Equal(t, "https://example.com/new-avatar.png", user.Data.AvatarURL)
	})

	t.Run("missing user id", func(t *testing.T) {
		client := newTestClient(&mockHTTPUpdateUser{})

		fields := &UserUpdateFields{
			Name: "Updated User",
		}

		user, err := client.UpdateUser(context.Background(), 0, fields)
		require.Error(t, err)
		assert.Equal(t, ErrMissingUserID, err)
		assert.Nil(t, user)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(&mockHTTPUpdateUser{})

		fields := &UserUpdateFields{
			Name: "Updated User",
		}

		user, err := client.UpdateUser(context.Background(), testUserIDBadRequest, fields)
		require.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(&mockHTTPUpdateUser{})

		fields := &UserUpdateFields{
			Name: "Updated User",
		}

		user, err := client.UpdateUser(context.Background(), testUserIDUnauthorized, fields)
		require.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(&mockHTTPUpdateUser{})

		fields := &UserUpdateFields{
			Name: "Updated User",
		}

		user, err := client.UpdateUser(context.Background(), testUserIDBadJSON, fields)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

// TestClient_UpdateUserRaw tests the method UpdateUserRaw()
func TestClient_UpdateUserRaw(t *testing.T) {
	t.Parallel()

	t.Run("missing user id", func(t *testing.T) {
		client := newTestClient(&mockHTTPUpdateUser{})

		fields := &UserUpdateFields{
			Name: "Updated User",
		}

		response, err := client.UpdateUserRaw(context.Background(), 0, fields)
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingUserID, err)
	})

	t.Run("update a valid user", func(t *testing.T) {
		client := newTestClient(&mockHTTPUpdateUser{})

		fields := &UserUpdateFields{
			Name:         "Updated User",
			Availability: "OFFLINE",
		}

		response, err := client.UpdateUserRaw(context.Background(), testUserID, fields)
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/users/update?userId=228225", response.URL)
		assert.Equal(t, http.MethodPatch, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

// BenchmarkClient_UpdateUser benchmarks the UpdateUser method
func BenchmarkClient_UpdateUser(b *testing.B) {
	client := newTestClient(&mockHTTPUpdateUser{})
	fields := &UserUpdateFields{
		Name: "Updated User",
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.UpdateUser(context.Background(), testUserID, fields)
	}
}

// BenchmarkClient_UpdateUserRaw benchmarks the UpdateUserRaw method
func BenchmarkClient_UpdateUserRaw(b *testing.B) {
	client := newTestClient(&mockHTTPUpdateUser{})
	fields := &UserUpdateFields{
		Name: "Updated User",
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.UpdateUserRaw(context.Background(), testUserID, fields)
	}
}
