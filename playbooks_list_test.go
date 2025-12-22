package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPlaybookID            = uint64(12345)
	testPlaybookName          = "Welcome Campaign"
	testPlaybookOrgID         = uint64(999)
	testPlaybookReportType    = "BOOK_MORE_MEETINGS"
	testPlaybookInteractionID = uint64(5678)
)

// mockGetPlaybooks returns a mock for playbooks list operations
func mockGetPlaybooks() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`[{"id":12345,"name":"Welcome Campaign","orgId":999,"meta":{},"createdAt":1615000000000,"updatedAt":1620000000000,"createdAuthorId":1,"updatedAuthorId":2,"interactionId":5678,"reportType":"BOOK_MORE_MEETINGS","goals":[{"id":"goal_1","message":"Schedule a meeting"}]}]`),
	)
}

// mockGetPlaybooksEmpty returns a mock for empty playbooks list
func mockGetPlaybooksEmpty() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`[]`),
	)
}

// mockGetPlaybooksBadJSON returns a mock for bad JSON response
func mockGetPlaybooksBadJSON() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`[{"id":12345"name":"Bad JSON"}]`),
	)
}

// TestClient_GetPlaybooks tests the method GetPlaybooks()
func TestClient_GetPlaybooks(t *testing.T) {
	t.Parallel()

	t.Run("get valid playbooks", func(t *testing.T) {
		client := newTestClient(mockGetPlaybooks())

		playbooks, err := client.GetPlaybooks(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, playbooks)
		assert.Len(t, playbooks.Data, 1)

		// Check returned values
		assert.Equal(t, testPlaybookID, playbooks.Data[0].ID)
		assert.Equal(t, testPlaybookName, playbooks.Data[0].Name)
		assert.Equal(t, testPlaybookOrgID, playbooks.Data[0].OrgID)
		assert.Equal(t, testPlaybookReportType, playbooks.Data[0].ReportType)
		assert.Equal(t, testPlaybookInteractionID, playbooks.Data[0].InteractionID)
		assert.Equal(t, int64(1615000000000), playbooks.Data[0].CreatedAt)
		assert.Equal(t, int64(1620000000000), playbooks.Data[0].UpdatedAt)
		assert.Equal(t, uint64(1), playbooks.Data[0].CreatedAuthorID)
		assert.Equal(t, uint64(2), playbooks.Data[0].UpdatedAuthorID)

		// Check goals
		require.Len(t, playbooks.Data[0].Goals, 1)
		assert.Equal(t, "goal_1", playbooks.Data[0].Goals[0].ID)
		assert.Equal(t, "Schedule a meeting", playbooks.Data[0].Goals[0].Message)
	})

	t.Run("empty playbooks list", func(t *testing.T) {
		client := newTestClient(mockGetPlaybooksEmpty())

		playbooks, err := client.GetPlaybooks(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, playbooks)
		assert.Empty(t, playbooks.Data)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		playbooks, err := client.GetPlaybooks(context.Background())
		require.Error(t, err)
		assert.Nil(t, playbooks)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		playbooks, err := client.GetPlaybooks(context.Background())
		require.Error(t, err)
		assert.Nil(t, playbooks)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockGetPlaybooksBadJSON())

		playbooks, err := client.GetPlaybooks(context.Background())
		require.Error(t, err)
		assert.Nil(t, playbooks)
	})
}

// TestClient_GetPlaybooksRaw tests the method GetPlaybooksRaw()
func TestClient_GetPlaybooksRaw(t *testing.T) {
	t.Parallel()

	t.Run("get valid playbooks raw", func(t *testing.T) {
		client := newTestClient(mockGetPlaybooks())

		response, err := client.GetPlaybooksRaw(context.Background())
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/playbooks/list", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("bad request response raw", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		response, err := client.GetPlaybooksRaw(context.Background())
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}

// BenchmarkClient_GetPlaybooks benchmarks the GetPlaybooks method
func BenchmarkClient_GetPlaybooks(b *testing.B) {
	client := newTestClient(mockGetPlaybooks())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetPlaybooks(context.Background())
	}
}

// BenchmarkClient_GetPlaybooksRaw benchmarks the GetPlaybooksRaw method
func BenchmarkClient_GetPlaybooksRaw(b *testing.B) {
	client := newTestClient(mockGetPlaybooks())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetPlaybooksRaw(context.Background())
	}
}
