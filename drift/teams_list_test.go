package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testTeamUserID         = uint64(228225)
	testTeamUserIDNotFound = uint64(999999)
)

// mockListTeams returns a mock for listing all teams
func mockListTeams() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":[{"id":1001,"orgId":12345,"workspaceId":"ws-abc-123","name":"Sales Team","updatedAt":1606273669631,"members":[228225,243266,300100],"owner":228225,"status":"ENABLED","main":true,"autoOffline":false,"teamCsatEnabled":true,"teamAvailabilityMode":"ALWAYS_ONLINE","responseTimerEnabled":true},{"id":1002,"orgId":12345,"workspaceId":"ws-abc-123","name":"Support Team","updatedAt":1614550516644,"members":[243266,300100],"owner":243266,"status":"ENABLED","main":false,"autoOffline":true,"teamCsatEnabled":false,"teamAvailabilityMode":"CUSTOM_HOURS","responseTimerEnabled":false}]}`),
	)
}

// mockListTeamsEmpty returns a mock for empty teams list
func mockListTeamsEmpty() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":[]}`),
	)
}

// mockListTeamsBadJSON returns a mock for bad JSON response
func mockListTeamsBadJSON() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":[{"id":1001"name":"Bad JSON"}]}`),
	)
}

// mockListTeamsByUser returns a multi-route mock for teams by user operations
func mockListTeamsByUser() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/teams/users/228225", http.StatusOK,
			`{"data":[{"id":1001,"orgId":12345,"workspaceId":"ws-abc-123","name":"Sales Team","updatedAt":1606273669631,"members":[228225,243266,300100],"owner":228225,"status":"ENABLED","main":true,"autoOffline":false,"teamCsatEnabled":true,"teamAvailabilityMode":"ALWAYS_ONLINE","responseTimerEnabled":true}]}`).
		addRoute(apiEndpoint+"/teams/users/999999", http.StatusNotFound, "")
}

// mockListTeamsByUserEmpty returns a mock for empty teams by user
func mockListTeamsByUserEmpty() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/teams/users/228225", http.StatusOK, `{"data":[]}`)
}

// mockListTeamsByUserBadJSON returns a mock for bad JSON teams by user
func mockListTeamsByUserBadJSON() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/teams/users/228225", http.StatusOK, `{"data":[{"id":1001"name":"Bad JSON"}]}`)
}

// TestClient_ListTeams tests the method ListTeams()
func TestClient_ListTeams(t *testing.T) {
	t.Parallel()

	t.Run("list all teams", func(t *testing.T) {
		client := newTestClient(mockListTeams())

		teams, err := client.ListTeams(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, teams)
		assert.Len(t, teams.Data, 2)

		// Check first team
		assert.Equal(t, uint64(1001), teams.Data[0].ID)
		assert.Equal(t, uint64(12345), teams.Data[0].OrgID)
		assert.Equal(t, "ws-abc-123", teams.Data[0].WorkspaceID)
		assert.Equal(t, "Sales Team", teams.Data[0].Name)
		assert.Equal(t, int64(1606273669631), teams.Data[0].UpdatedAt)
		assert.Len(t, teams.Data[0].Members, 3)
		assert.Equal(t, uint64(228225), teams.Data[0].Owner)
		assert.Equal(t, "ENABLED", teams.Data[0].Status)
		assert.True(t, teams.Data[0].Main)
		assert.False(t, teams.Data[0].AutoOffline)
		assert.True(t, teams.Data[0].TeamCsatEnabled)
		assert.Equal(t, "ALWAYS_ONLINE", teams.Data[0].TeamAvailabilityMode)
		assert.True(t, teams.Data[0].ResponseTimerEnabled)

		// Check second team
		assert.Equal(t, uint64(1002), teams.Data[1].ID)
		assert.Equal(t, "Support Team", teams.Data[1].Name)
		assert.Equal(t, "CUSTOM_HOURS", teams.Data[1].TeamAvailabilityMode)
		assert.True(t, teams.Data[1].AutoOffline)
		assert.False(t, teams.Data[1].Main)
	})

	t.Run("empty teams list", func(t *testing.T) {
		client := newTestClient(mockListTeamsEmpty())

		teams, err := client.ListTeams(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, teams)
		assert.Empty(t, teams.Data)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		teams, err := client.ListTeams(context.Background())
		require.Error(t, err)
		assert.Nil(t, teams)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockListTeamsBadJSON())

		teams, err := client.ListTeams(context.Background())
		require.Error(t, err)
		assert.Nil(t, teams)
	})
}

// TestClient_ListTeamsRaw tests the method ListTeamsRaw()
func TestClient_ListTeamsRaw(t *testing.T) {
	t.Parallel()

	t.Run("list all teams raw", func(t *testing.T) {
		client := newTestClient(mockListTeams())

		response, err := client.ListTeamsRaw(context.Background())
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/teams/org", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		response, err := client.ListTeamsRaw(context.Background())
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})
}

// TestClient_ListTeamsByUser tests the method ListTeamsByUser()
func TestClient_ListTeamsByUser(t *testing.T) {
	t.Parallel()

	t.Run("list teams for valid user", func(t *testing.T) {
		client := newTestClient(mockListTeamsByUser())

		teams, err := client.ListTeamsByUser(context.Background(), testTeamUserID)
		require.NoError(t, err)
		assert.NotNil(t, teams)
		assert.Len(t, teams.Data, 1)

		// Check team data
		assert.Equal(t, uint64(1001), teams.Data[0].ID)
		assert.Equal(t, "Sales Team", teams.Data[0].Name)
		assert.Equal(t, uint64(228225), teams.Data[0].Owner)
		assert.True(t, teams.Data[0].Main)
	})

	t.Run("missing user id", func(t *testing.T) {
		client := newTestClient(mockListTeamsByUser())

		teams, err := client.ListTeamsByUser(context.Background(), 0)
		require.Error(t, err)
		assert.Equal(t, ErrMissingUserID, err)
		assert.Nil(t, teams)
	})

	t.Run("user not found", func(t *testing.T) {
		client := newTestClient(mockListTeamsByUser())

		teams, err := client.ListTeamsByUser(context.Background(), testTeamUserIDNotFound)
		require.Error(t, err)
		assert.Nil(t, teams)
	})

	t.Run("empty teams for user", func(t *testing.T) {
		client := newTestClient(mockListTeamsByUserEmpty())

		teams, err := client.ListTeamsByUser(context.Background(), testTeamUserID)
		require.NoError(t, err)
		assert.NotNil(t, teams)
		assert.Empty(t, teams.Data)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockListTeamsByUserBadJSON())

		teams, err := client.ListTeamsByUser(context.Background(), testTeamUserID)
		require.Error(t, err)
		assert.Nil(t, teams)
	})
}

// TestClient_ListTeamsByUserRaw tests the method ListTeamsByUserRaw()
func TestClient_ListTeamsByUserRaw(t *testing.T) {
	t.Parallel()

	t.Run("list teams by user raw", func(t *testing.T) {
		client := newTestClient(mockListTeamsByUser())

		response, err := client.ListTeamsByUserRaw(context.Background(), testTeamUserID)
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/teams/users/228225", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("missing user id", func(t *testing.T) {
		client := newTestClient(mockListTeamsByUser())

		response, err := client.ListTeamsByUserRaw(context.Background(), 0)
		require.Error(t, err)
		assert.Equal(t, ErrMissingUserID, err)
		assert.Nil(t, response)
	})

	t.Run("user not found", func(t *testing.T) {
		client := newTestClient(mockListTeamsByUser())

		response, err := client.ListTeamsByUserRaw(context.Background(), testTeamUserIDNotFound)
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}

// BenchmarkClient_ListTeams benchmarks the ListTeams method
func BenchmarkClient_ListTeams(b *testing.B) {
	client := newTestClient(mockListTeams())
	for i := 0; i < b.N; i++ {
		_, _ = client.ListTeams(context.Background())
	}
}

// BenchmarkClient_ListTeamsRaw benchmarks the ListTeamsRaw method
func BenchmarkClient_ListTeamsRaw(b *testing.B) {
	client := newTestClient(mockListTeams())
	for i := 0; i < b.N; i++ {
		_, _ = client.ListTeamsRaw(context.Background())
	}
}

// BenchmarkClient_ListTeamsByUser benchmarks the ListTeamsByUser method
func BenchmarkClient_ListTeamsByUser(b *testing.B) {
	client := newTestClient(mockListTeamsByUser())
	for i := 0; i < b.N; i++ {
		_, _ = client.ListTeamsByUser(context.Background(), testTeamUserID)
	}
}

// BenchmarkClient_ListTeamsByUserRaw benchmarks the ListTeamsByUserRaw method
func BenchmarkClient_ListTeamsByUserRaw(b *testing.B) {
	client := newTestClient(mockListTeamsByUser())
	for i := 0; i < b.N; i++ {
		_, _ = client.ListTeamsByUserRaw(context.Background(), testTeamUserID)
	}
}
