package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockListConversations returns a multi-route mock for conversation list operations
func mockListConversations() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpointList+"/conversations/list", http.StatusOK,
			`{"data":[{"id":3782727146,"contactId":17035536800,"inboxId":62491,"status":"open","createdAt":1686303243241,"updatedAt":1686303381300},{"id":3782727147,"contactId":17035536801,"inboxId":62491,"status":"closed","createdAt":1686303243242,"updatedAt":1686303381301}],"links":{"self":"https://api.drift.com/conversations/list?page_token=abc123","next":"https://api.drift.com/conversations/list?page_token=def456"}}`).
		addRoute(apiEndpointList+"/conversations/list?limit=50", http.StatusOK,
			`{"data":[{"id":3782727146,"contactId":17035536800,"inboxId":62491,"status":"open","createdAt":1686303243241,"updatedAt":1686303381300}]}`).
		addRoute(apiEndpointList+"/conversations/list?limit=100", http.StatusOK,
			`{"data":[{"id":1,"contactId":1,"inboxId":1,"status":"open","createdAt":1,"updatedAt":1}]}`).
		addRoute(apiEndpointList+"/conversations/list?statusId=1", http.StatusOK,
			`{"data":[{"id":3782727146,"contactId":17035536800,"inboxId":62491,"status":"open","createdAt":1686303243241,"updatedAt":1686303381300}]}`).
		addRoute(apiEndpointList+"/conversations/list?statusId=2", http.StatusOK,
			`{"data":[{"id":3782727147,"contactId":17035536801,"inboxId":62491,"status":"closed","createdAt":1686303243242,"updatedAt":1686303381301}]}`).
		addRoute(apiEndpointList+"/conversations/list?statusId=3", http.StatusOK,
			`{"data":[{"id":3782727148,"contactId":17035536802,"inboxId":62491,"status":"pending","createdAt":1686303243243,"updatedAt":1686303381302}]}`).
		addRoute(apiEndpointList+"/conversations/list?statusId=1&statusId=2", http.StatusOK,
			`{"data":[{"id":3782727146,"contactId":17035536800,"inboxId":62491,"status":"open","createdAt":1686303243241,"updatedAt":1686303381300},{"id":3782727147,"contactId":17035536801,"inboxId":62491,"status":"closed","createdAt":1686303243242,"updatedAt":1686303381301}]}`).
		addRoute(apiEndpointList+"/conversations/list?page_token=abc123", http.StatusOK,
			`{"data":[{"id":3782727149,"contactId":17035536803,"inboxId":62491,"status":"open","createdAt":1686303243244,"updatedAt":1686303381303}]}`).
		addRoute("https://api.drift.com/conversations/list?page_token=def456", http.StatusOK,
			`{"data":[{"id":3782727150,"contactId":17035536804,"inboxId":62491,"status":"open","createdAt":1686303243245,"updatedAt":1686303381304}]}`).
		addRoute(apiEndpointList+"/conversations/list?limit=25&statusId=1", http.StatusOK,
			`{"data":[{"id":3782727146,"contactId":17035536800,"inboxId":62491,"status":"open","createdAt":1686303243241,"updatedAt":1686303381300}]}`).
		addRoute(apiEndpointList+"/conversations/list?limit=25&statusId=2", http.StatusOK,
			`{"data":[{"id":3782727147,"contactId":17035536801,"inboxId":62491,"status":"closed","createdAt":1686303243242,"updatedAt":1686303381301}]}`).
		addRoute(apiEndpointList+"/conversations/list?limit=25&statusId=3", http.StatusOK,
			`{"data":[{"id":3782727148,"contactId":17035536802,"inboxId":62491,"status":"pending","createdAt":1686303243243,"updatedAt":1686303381302}]}`)
}

// TestClient_ListConversations tests the method ListConversations()
func TestClient_ListConversations(t *testing.T) {
	t.Parallel()

	t.Run("list all conversations with no filters", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListConversations(context.Background(), nil)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 2)

		// Check first conversation
		assert.Equal(t, uint64(3782727146), conversations.Data[0].ID)
		assert.Equal(t, uint64(17035536800), conversations.Data[0].ContactID)
		assert.Equal(t, 62491, conversations.Data[0].InboxID)
		assert.Equal(t, "open", conversations.Data[0].Status)

		// Check pagination links
		assert.NotNil(t, conversations.Links)
		assert.Equal(t, "https://api.drift.com/conversations/list?page_token=abc123", conversations.Links.Self)
		assert.Equal(t, "https://api.drift.com/conversations/list?page_token=def456", conversations.Links.Next)
	})

	t.Run("list conversations with limit", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListConversations(context.Background(), &ConversationListQuery{
			Limit: 50,
		})
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
	})

	t.Run("limit capped at 100", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		query := &ConversationListQuery{Limit: 150}
		conversations, err := client.ListConversations(context.Background(), query)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Equal(t, 100, query.Limit)
	})

	t.Run("list open conversations", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListConversations(context.Background(), &ConversationListQuery{
			StatusIDs: []int{ConversationStatusOpen},
		})
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
		assert.Equal(t, "open", conversations.Data[0].Status)
	})

	t.Run("list closed conversations", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListConversations(context.Background(), &ConversationListQuery{
			StatusIDs: []int{ConversationStatusClosed},
		})
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
		assert.Equal(t, "closed", conversations.Data[0].Status)
	})

	t.Run("list with multiple status filters", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListConversations(context.Background(), &ConversationListQuery{
			StatusIDs: []int{ConversationStatusOpen, ConversationStatusClosed},
		})
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 2)
	})

	t.Run("list with page token", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListConversations(context.Background(), &ConversationListQuery{
			PageToken: "abc123",
		})
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
		assert.Equal(t, uint64(3782727149), conversations.Data[0].ID)
	})

	t.Run("bad json response", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list", http.StatusOK, `{"data":[{"id":}]}`)

		client := newTestClient(mock)

		conversations, err := client.ListConversations(context.Background(), nil)
		require.Error(t, err)
		assert.Nil(t, conversations)
	})

	t.Run("http error", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list", http.StatusInternalServerError, "")

		client := newTestClient(mock)

		conversations, err := client.ListConversations(context.Background(), nil)
		require.Error(t, err)
		assert.Nil(t, conversations)
	})
}

// TestClient_ListConversationsRaw tests the method ListConversationsRaw()
func TestClient_ListConversationsRaw(t *testing.T) {
	t.Parallel()

	t.Run("list with no filters", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		response, err := client.ListConversationsRaw(context.Background(), nil)
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpointList+"/conversations/list", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("list with limit", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		response, err := client.ListConversationsRaw(context.Background(), &ConversationListQuery{
			Limit: 50,
		})
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.Equal(t, apiEndpointList+"/conversations/list?limit=50", response.URL)
	})
}

// TestClient_ListConversationsNext tests the method ListConversationsNext()
func TestClient_ListConversationsNext(t *testing.T) {
	t.Parallel()

	t.Run("get next page of conversations", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		// First get the initial page
		conversations, err := client.ListConversations(context.Background(), nil)
		require.NoError(t, err)
		require.NotNil(t, conversations.Links)

		// Get the next page
		nextConversations, err := client.ListConversationsNext(context.Background(), conversations)
		require.NoError(t, err)
		assert.NotNil(t, nextConversations)
		assert.Len(t, nextConversations.Data, 1)
		assert.Equal(t, uint64(3782727150), nextConversations.Data[0].ID)
	})

	t.Run("nil conversations returns error", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		nextConversations, err := client.ListConversationsNext(context.Background(), nil)
		require.Error(t, err)
		assert.Equal(t, ErrNoNextPage, err)
		assert.Nil(t, nextConversations)
	})

	t.Run("nil links returns error", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations := &Conversations{
			Data:  []*conversationData{},
			Links: nil,
		}
		nextConversations, err := client.ListConversationsNext(context.Background(), conversations)
		require.Error(t, err)
		assert.Equal(t, ErrNoNextPage, err)
		assert.Nil(t, nextConversations)
	})

	t.Run("empty next link returns error", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations := &Conversations{
			Data:  []*conversationData{},
			Links: &PaginationLinks{Self: "test"},
		}
		nextConversations, err := client.ListConversationsNext(context.Background(), conversations)
		require.Error(t, err)
		assert.Equal(t, ErrNoNextPage, err)
		assert.Nil(t, nextConversations)
	})

	t.Run("http error on next page", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute("https://api.drift.com/conversations/list?page_token=error", http.StatusInternalServerError, "")

		client := newTestClient(mock)

		conversations := &Conversations{
			Data:  []*conversationData{},
			Links: &PaginationLinks{Next: "https://api.drift.com/conversations/list?page_token=error"},
		}
		nextConversations, err := client.ListConversationsNext(context.Background(), conversations)
		require.Error(t, err)
		assert.Nil(t, nextConversations)
	})

	t.Run("bad json on next page", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute("https://api.drift.com/conversations/list?page_token=badjson", http.StatusOK, `{"data":[{"id":}]}`)

		client := newTestClient(mock)

		conversations := &Conversations{
			Data:  []*conversationData{},
			Links: &PaginationLinks{Next: "https://api.drift.com/conversations/list?page_token=badjson"},
		}
		nextConversations, err := client.ListConversationsNext(context.Background(), conversations)
		require.Error(t, err)
		assert.Nil(t, nextConversations)
	})
}

// TestClient_ListOpenConversations tests the convenience method ListOpenConversations()
func TestClient_ListOpenConversations(t *testing.T) {
	t.Parallel()

	t.Run("list open conversations", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListOpenConversations(context.Background(), 25)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
		assert.Equal(t, "open", conversations.Data[0].Status)
		assert.Equal(t, uint64(3782727146), conversations.Data[0].ID)
	})

	t.Run("error handling", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list?limit=25&statusId=1", http.StatusUnauthorized, "")

		client := newTestClient(mock)

		conversations, err := client.ListOpenConversations(context.Background(), 25)
		require.Error(t, err)
		assert.Nil(t, conversations)
	})
}

// TestClient_ListClosedConversations tests the convenience method ListClosedConversations()
func TestClient_ListClosedConversations(t *testing.T) {
	t.Parallel()

	t.Run("list closed conversations", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListClosedConversations(context.Background(), 25)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
		assert.Equal(t, "closed", conversations.Data[0].Status)
		assert.Equal(t, uint64(3782727147), conversations.Data[0].ID)
	})

	t.Run("error handling", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list?limit=25&statusId=2", http.StatusUnauthorized, "")

		client := newTestClient(mock)

		conversations, err := client.ListClosedConversations(context.Background(), 25)
		require.Error(t, err)
		assert.Nil(t, conversations)
	})
}

// TestClient_ListPendingConversations tests the convenience method ListPendingConversations()
func TestClient_ListPendingConversations(t *testing.T) {
	t.Parallel()

	t.Run("list pending conversations", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListPendingConversations(context.Background(), 25)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
		assert.Equal(t, "pending", conversations.Data[0].Status)
		assert.Equal(t, uint64(3782727148), conversations.Data[0].ID)
	})

	t.Run("error handling", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list?limit=25&statusId=3", http.StatusUnauthorized, "")

		client := newTestClient(mock)

		conversations, err := client.ListPendingConversations(context.Background(), 25)
		require.Error(t, err)
		assert.Nil(t, conversations)
	})
}

// TestClient_ListAllConversations tests the method ListAllConversations()
func TestClient_ListAllConversations(t *testing.T) {
	t.Parallel()

	t.Run("single page no pagination", func(t *testing.T) {
		// Mock with no next link
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list", http.StatusOK,
				`{"data":[{"id":1,"contactId":100,"inboxId":1,"status":"open","createdAt":1,"updatedAt":1}]}`)

		client := newTestClient(mock)

		conversations, err := client.ListAllConversations(context.Background(), nil)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
		assert.Equal(t, uint64(1), conversations.Data[0].ID)
	})

	t.Run("multiple pages with pagination", func(t *testing.T) {
		// Use the existing mock which has pagination links
		client := newTestClient(mockListConversations())

		conversations, err := client.ListAllConversations(context.Background(), nil)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		// First page has 2 conversations, second page has 1
		assert.Len(t, conversations.Data, 3)
	})

	t.Run("empty response", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list", http.StatusOK,
				`{"data":[]}`)

		client := newTestClient(mock)

		conversations, err := client.ListAllConversations(context.Background(), nil)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Empty(t, conversations.Data)
	})

	t.Run("error on first request", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list", http.StatusInternalServerError, "")

		client := newTestClient(mock)

		conversations, err := client.ListAllConversations(context.Background(), nil)
		require.Error(t, err)
		assert.Nil(t, conversations)
	})

	t.Run("with query parameters", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list?limit=50&statusId=1", http.StatusOK,
				`{"data":[{"id":1,"contactId":100,"inboxId":1,"status":"open","createdAt":1,"updatedAt":1}]}`)

		client := newTestClient(mock)

		conversations, err := client.ListAllConversations(context.Background(), &ConversationListQuery{
			Limit:     50,
			StatusIDs: []int{ConversationStatusOpen},
		})
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
	})

	t.Run("error on subsequent page request", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list", http.StatusOK,
				`{"data":[{"id":1,"contactId":100,"inboxId":1,"status":"open","createdAt":1,"updatedAt":1}],"links":{"next":"https://api.drift.com/conversations/list?page_token=page2"}}`).
			addRoute("https://api.drift.com/conversations/list?page_token=page2", http.StatusInternalServerError, "")

		client := newTestClient(mock)

		conversations, err := client.ListAllConversations(context.Background(), nil)
		require.Error(t, err)
		assert.Nil(t, conversations)
	})
}

// TestClient_ListConversationsByStatus tests the method ListConversationsByStatus()
func TestClient_ListConversationsByStatus(t *testing.T) {
	t.Parallel()

	t.Run("filter by open status", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListConversationsByStatus(context.Background(), ConversationStatusOpen, 25)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
		assert.Equal(t, "open", conversations.Data[0].Status)
	})

	t.Run("filter by closed status", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListConversationsByStatus(context.Background(), ConversationStatusClosed, 25)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
		assert.Equal(t, "closed", conversations.Data[0].Status)
	})

	t.Run("filter by pending status", func(t *testing.T) {
		client := newTestClient(mockListConversations())

		conversations, err := client.ListConversationsByStatus(context.Background(), ConversationStatusPending, 25)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 1)
		assert.Equal(t, "pending", conversations.Data[0].Status)
	})

	t.Run("with custom limit", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list?limit=10&statusId=1", http.StatusOK,
				`{"data":[{"id":1,"contactId":100,"inboxId":1,"status":"open","createdAt":1,"updatedAt":1}]}`)

		client := newTestClient(mock)

		conversations, err := client.ListConversationsByStatus(context.Background(), ConversationStatusOpen, 10)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
	})

	t.Run("error handling", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list?limit=25&statusId=1", http.StatusInternalServerError, "")

		client := newTestClient(mock)

		conversations, err := client.ListConversationsByStatus(context.Background(), ConversationStatusOpen, 25)
		require.Error(t, err)
		assert.Nil(t, conversations)
	})
}

// TestClient_ListConversationsByContactID tests the method ListConversationsByContactID()
func TestClient_ListConversationsByContactID(t *testing.T) {
	t.Parallel()

	t.Run("find conversations matching contact id", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list?limit=50", http.StatusOK,
				`{"data":[{"id":1,"contactId":100,"inboxId":1,"status":"open","createdAt":1,"updatedAt":1},{"id":2,"contactId":200,"inboxId":1,"status":"closed","createdAt":2,"updatedAt":2},{"id":3,"contactId":100,"inboxId":1,"status":"pending","createdAt":3,"updatedAt":3}]}`)

		client := newTestClient(mock)

		conversations, err := client.ListConversationsByContactID(context.Background(), 100, 50)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Len(t, conversations.Data, 2)

		// Both conversations should have contactId 100
		for _, conv := range conversations.Data {
			assert.Equal(t, uint64(100), conv.ContactID)
		}
	})

	t.Run("no conversations match contact id", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list?limit=50", http.StatusOK,
				`{"data":[{"id":1,"contactId":100,"inboxId":1,"status":"open","createdAt":1,"updatedAt":1}]}`)

		client := newTestClient(mock)

		conversations, err := client.ListConversationsByContactID(context.Background(), 999, 50)
		require.NoError(t, err)
		assert.NotNil(t, conversations)
		assert.Empty(t, conversations.Data)
	})

	t.Run("error from list all conversations propagates", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list?limit=50", http.StatusInternalServerError, "")

		client := newTestClient(mock)

		conversations, err := client.ListConversationsByContactID(context.Background(), 100, 50)
		require.Error(t, err)
		assert.Nil(t, conversations)
	})
}

// TestClient_GetConversationCount tests the method GetConversationCount()
func TestClient_GetConversationCount(t *testing.T) {
	t.Parallel()

	t.Run("count with multiple conversations", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list", http.StatusOK,
				`{"data":[{"id":1,"contactId":100,"inboxId":1,"status":"open","createdAt":1,"updatedAt":1},{"id":2,"contactId":200,"inboxId":1,"status":"closed","createdAt":2,"updatedAt":2},{"id":3,"contactId":300,"inboxId":1,"status":"pending","createdAt":3,"updatedAt":3}]}`)

		client := newTestClient(mock)

		count, err := client.GetConversationCount(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("count with zero conversations", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list", http.StatusOK,
				`{"data":[]}`)

		client := newTestClient(mock)

		count, err := client.GetConversationCount(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("error from list all conversations propagates", func(t *testing.T) {
		mock := newMockHTTPMulti().
			addRoute(apiEndpointList+"/conversations/list", http.StatusInternalServerError, "")

		client := newTestClient(mock)

		count, err := client.GetConversationCount(context.Background())
		require.Error(t, err)
		assert.Equal(t, 0, count)
	})
}

// TestConversationListQuery_BuildURL tests the method BuildURL()
func TestConversationListQuery_BuildURL(t *testing.T) {
	t.Parallel()

	t.Run("empty query", func(t *testing.T) {
		q := &ConversationListQuery{}
		assert.Equal(t, apiEndpointList+"/conversations/list", q.BuildURL())
	})

	t.Run("with limit", func(t *testing.T) {
		q := &ConversationListQuery{Limit: 50}
		assert.Equal(t, apiEndpointList+"/conversations/list?limit=50", q.BuildURL())
	})

	t.Run("with status filter", func(t *testing.T) {
		q := &ConversationListQuery{StatusIDs: []int{ConversationStatusOpen}}
		assert.Equal(t, apiEndpointList+"/conversations/list?statusId=1", q.BuildURL())
	})

	t.Run("with multiple status filters", func(t *testing.T) {
		q := &ConversationListQuery{StatusIDs: []int{ConversationStatusOpen, ConversationStatusClosed}}
		assert.Equal(t, apiEndpointList+"/conversations/list?statusId=1&statusId=2", q.BuildURL())
	})

	t.Run("with page token", func(t *testing.T) {
		q := &ConversationListQuery{PageToken: "abc123"}
		assert.Equal(t, apiEndpointList+"/conversations/list?page_token=abc123", q.BuildURL())
	})

	t.Run("with all parameters", func(t *testing.T) {
		q := &ConversationListQuery{
			Limit:     25,
			StatusIDs: []int{ConversationStatusOpen},
			PageToken: "abc123",
		}
		assert.Equal(t, apiEndpointList+"/conversations/list?limit=25&statusId=1&page_token=abc123", q.BuildURL())
	})
}

// TestStatusIDToString tests the statusIDToString helper
func TestStatusIDToString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "open", statusIDToString(ConversationStatusOpen))
	assert.Equal(t, "closed", statusIDToString(ConversationStatusClosed))
	assert.Equal(t, "pending", statusIDToString(ConversationStatusPending))
	assert.Equal(t, "unknown(99)", statusIDToString(99))
}

// BenchmarkClient_ListConversations benchmarks the ListConversations method
func BenchmarkClient_ListConversations(b *testing.B) {
	client := newTestClient(mockListConversations())
	for i := 0; i < b.N; i++ {
		_, _ = client.ListConversations(context.Background(), nil)
	}
}

// BenchmarkClient_ListConversationsRaw benchmarks the ListConversationsRaw method
func BenchmarkClient_ListConversationsRaw(b *testing.B) {
	client := newTestClient(mockListConversations())
	for i := 0; i < b.N; i++ {
		_, _ = client.ListConversationsRaw(context.Background(), nil)
	}
}

// BenchmarkConversationListQuery_BuildURL benchmarks the BuildURL method
func BenchmarkConversationListQuery_BuildURL(b *testing.B) {
	q := &ConversationListQuery{
		Limit:     50,
		StatusIDs: []int{ConversationStatusOpen, ConversationStatusClosed},
		PageToken: "abc123",
	}
	for i := 0; i < b.N; i++ {
		_ = q.BuildURL()
	}
}
