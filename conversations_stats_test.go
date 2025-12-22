package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockGetConversationStats returns a mock for conversation stats operations
func mockGetConversationStats() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"conversationCount":{"CLOSED":282,"OPEN":125,"PENDING":43}}`),
	)
}

// mockGetConversationStatsEmpty returns a mock for empty conversation stats
func mockGetConversationStatsEmpty() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"conversationCount":{}}`),
	)
}

// TestClient_GetConversationStats tests the method GetConversationStats()
func TestClient_GetConversationStats(t *testing.T) {
	t.Parallel()

	t.Run("get valid conversation stats", func(t *testing.T) {
		client := newTestClient(mockGetConversationStats())

		stats, err := client.GetConversationStats(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, stats)

		// Check returned values
		assert.Equal(t, 282, stats.ConversationCount["CLOSED"])
		assert.Equal(t, 125, stats.ConversationCount["OPEN"])
		assert.Equal(t, 43, stats.ConversationCount["PENDING"])
	})

	t.Run("get empty conversation stats", func(t *testing.T) {
		client := newTestClient(mockGetConversationStatsEmpty())

		stats, err := client.GetConversationStats(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Empty(t, stats.ConversationCount)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"conversationCount":{"CLOSED":282"OPEN":125}}`))

		stats, err := client.GetConversationStats(context.Background())
		require.Error(t, err)
		assert.Nil(t, stats)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		stats, err := client.GetConversationStats(context.Background())
		require.Error(t, err)
		assert.Nil(t, stats)
	})
}

// TestClient_GetConversationStatsRaw tests the method GetConversationStatsRaw()
func TestClient_GetConversationStatsRaw(t *testing.T) {
	t.Parallel()

	t.Run("get valid response", func(t *testing.T) {
		client := newTestClient(mockGetConversationStats())

		response, err := client.GetConversationStatsRaw(context.Background())
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/conversations/stats", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

// TestClient_GetOpenConversationCount tests the convenience method GetOpenConversationCount()
func TestClient_GetOpenConversationCount(t *testing.T) {
	t.Parallel()

	t.Run("get open count", func(t *testing.T) {
		client := newTestClient(mockGetConversationStats())

		count, err := client.GetOpenConversationCount(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 125, count)
	})

	t.Run("get open count when empty", func(t *testing.T) {
		client := newTestClient(mockGetConversationStatsEmpty())

		count, err := client.GetOpenConversationCount(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

// TestClient_GetClosedConversationCount tests the convenience method GetClosedConversationCount()
func TestClient_GetClosedConversationCount(t *testing.T) {
	t.Parallel()

	t.Run("get closed count", func(t *testing.T) {
		client := newTestClient(mockGetConversationStats())

		count, err := client.GetClosedConversationCount(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 282, count)
	})
}

// TestClient_GetPendingConversationCount tests the convenience method GetPendingConversationCount()
func TestClient_GetPendingConversationCount(t *testing.T) {
	t.Parallel()

	t.Run("get pending count", func(t *testing.T) {
		client := newTestClient(mockGetConversationStats())

		count, err := client.GetPendingConversationCount(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 43, count)
	})
}

// TestClient_GetTotalConversationCount tests the convenience method GetTotalConversationCount()
func TestClient_GetTotalConversationCount(t *testing.T) {
	t.Parallel()

	t.Run("get total count", func(t *testing.T) {
		client := newTestClient(mockGetConversationStats())

		count, err := client.GetTotalConversationCount(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 450, count) // 282 + 125 + 43
	})

	t.Run("get total count when empty", func(t *testing.T) {
		client := newTestClient(mockGetConversationStatsEmpty())

		count, err := client.GetTotalConversationCount(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

// BenchmarkClient_GetConversationStats benchmarks the GetConversationStats method
func BenchmarkClient_GetConversationStats(b *testing.B) {
	client := newTestClient(mockGetConversationStats())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetConversationStats(context.Background())
	}
}

// BenchmarkClient_GetConversationStatsRaw benchmarks the GetConversationStatsRaw method
func BenchmarkClient_GetConversationStatsRaw(b *testing.B) {
	client := newTestClient(mockGetConversationStats())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetConversationStatsRaw(context.Background())
	}
}
