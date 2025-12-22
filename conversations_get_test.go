package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testConversationID             = uint64(116119985)
	testConversationIDBadRequest   = uint64(111111111)
	testConversationIDUnauthorized = uint64(222222222)
	testConversationIDBadJSON      = uint64(333333333)
	testConversationIDNotFound     = uint64(444444444)
)

// mockGetConversation returns a multi-route mock for conversation GET operations
func mockGetConversation() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/conversations/116119985", http.StatusOK,
			`{"data":{"id":116119985,"contactId":903182234,"inboxId":116983,"status":"closed","participants":[243266,252465],"conversationTags":[{"color":"0960C5","name":"test_tag"},{"color":"7695A5","name":"second_tag"}],"relatedPlaybookId":63505,"createdAt":1528293067340,"updatedAt":1565798942334}}`).
		addRoute(apiEndpoint+"/conversations/111111111", http.StatusBadRequest, "").
		addRoute(apiEndpoint+"/conversations/222222222", http.StatusUnauthorized, "").
		addRoute(apiEndpoint+"/conversations/333333333", http.StatusOK, `{"data":{"id":333333333"status":"open"}`).
		addRoute(apiEndpoint+"/conversations/444444444", http.StatusNotFound, "")
}

// TestClient_GetConversation tests the method GetConversation()
func TestClient_GetConversation(t *testing.T) {
	t.Parallel()

	t.Run("get a valid conversation by id", func(t *testing.T) {
		client := newTestClient(mockGetConversation())

		conversation, err := client.GetConversation(context.Background(), testConversationID)
		require.NoError(t, err)
		assert.NotNil(t, conversation)
		assert.NotNil(t, conversation.Data)

		// Check returned values
		assert.Equal(t, testConversationID, conversation.Data.ID)
		assert.Equal(t, uint64(903182234), conversation.Data.ContactID)
		assert.Equal(t, 116983, conversation.Data.InboxID)
		assert.Equal(t, "closed", conversation.Data.Status)
		assert.Len(t, conversation.Data.Participants, 2)
		assert.Equal(t, uint64(243266), conversation.Data.Participants[0])
		assert.Equal(t, uint64(252465), conversation.Data.Participants[1])
		assert.Len(t, conversation.Data.ConversationTags, 2)
		assert.Equal(t, "0960C5", conversation.Data.ConversationTags[0].Color)
		assert.Equal(t, "test_tag", conversation.Data.ConversationTags[0].Name)
		assert.Equal(t, 63505, conversation.Data.RelatedPlaybookID)
		assert.Equal(t, int64(1528293067340), conversation.Data.CreatedAt)
		assert.Equal(t, int64(1565798942334), conversation.Data.UpdatedAt)
	})

	t.Run("missing conversation id", func(t *testing.T) {
		client := newTestClient(mockGetConversation())

		conversation, err := client.GetConversation(context.Background(), 0)
		require.Error(t, err)
		assert.Equal(t, ErrMissingConversationID, err)
		assert.Nil(t, conversation)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(mockGetConversation())

		conversation, err := client.GetConversation(context.Background(), testConversationIDBadRequest)
		require.Error(t, err)
		assert.Nil(t, conversation)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(mockGetConversation())

		conversation, err := client.GetConversation(context.Background(), testConversationIDUnauthorized)
		require.Error(t, err)
		assert.Nil(t, conversation)
	})

	t.Run("not found response", func(t *testing.T) {
		client := newTestClient(mockGetConversation())

		conversation, err := client.GetConversation(context.Background(), testConversationIDNotFound)
		require.Error(t, err)
		assert.Nil(t, conversation)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockGetConversation())

		conversation, err := client.GetConversation(context.Background(), testConversationIDBadJSON)
		require.Error(t, err)
		assert.Nil(t, conversation)
	})
}

// TestClient_GetConversationRaw tests the method GetConversationRaw()
func TestClient_GetConversationRaw(t *testing.T) {
	t.Parallel()

	t.Run("missing conversation id", func(t *testing.T) {
		client := newTestClient(mockGetConversation())

		response, err := client.GetConversationRaw(context.Background(), 0)
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingConversationID, err)
	})

	t.Run("get a valid conversation by id", func(t *testing.T) {
		client := newTestClient(mockGetConversation())

		response, err := client.GetConversationRaw(context.Background(), testConversationID)
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/conversations/116119985", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

// BenchmarkClient_GetConversation benchmarks the GetConversation method
func BenchmarkClient_GetConversation(b *testing.B) {
	client := newTestClient(mockGetConversation())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetConversation(context.Background(), testConversationID)
	}
}

// BenchmarkClient_GetConversationRaw benchmarks the GetConversationRaw method
func BenchmarkClient_GetConversationRaw(b *testing.B) {
	client := newTestClient(mockGetConversation())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetConversationRaw(context.Background(), testConversationID)
	}
}
