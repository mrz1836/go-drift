package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testMessageID = uint64(987654321)

// mockGetMessages returns a multi-route mock for message GET operations
func mockGetMessages() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/conversations/116119985/messages", http.StatusOK,
			`{"data":{"messages":[{"id":987654321,"orgId":12345,"conversationId":116119985,"body":"Hello, how can I help you?","type":"chat","author":{"id":243266,"type":"user","bot":false},"createdAt":1686304523000,"context":{"ip":"192.168.1.1","userAgent":"Mozilla/5.0"}},{"id":987654322,"orgId":12345,"conversationId":116119985,"body":"I have a question about your product.","type":"chat","author":{"id":903182234,"type":"contact","bot":false},"createdAt":1686304545000}]},"pagination":{"next":"abc123"}}`).
		addRoute(apiEndpoint+"/conversations/116119985/messages?next=abc123", http.StatusOK,
			`{"data":{"messages":[{"id":987654323,"orgId":12345,"conversationId":116119985,"body":"Sure, I can help with that!","type":"chat","author":{"id":243266,"type":"user","bot":false},"createdAt":1686304562000}]}}`).
		addRoute(apiEndpoint+"/conversations/111111111/messages", http.StatusBadRequest, "").
		addRoute(apiEndpoint+"/conversations/222222222/messages", http.StatusUnauthorized, "").
		addRoute(apiEndpoint+"/conversations/333333333/messages", http.StatusOK,
			`{"data":{"messages":[{"id":987654321"body":"Hello"}`).
		addRoute(apiEndpoint+"/conversations/444444444/messages", http.StatusNotFound, "")
}

// mockGetMessagesEmpty returns a mock for empty messages
func mockGetMessagesEmpty() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/conversations/116119985/messages", http.StatusOK,
			`{"data":{"messages":[]}}`)
}

// mockGetMessagesWithAttachments returns a mock for messages with attachments
func mockGetMessagesWithAttachments() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/conversations/116119985/messages", http.StatusOK,
			`{"data":{"messages":[{"id":987654321,"orgId":12345,"conversationId":116119985,"body":"Here is the file you requested.","type":"chat","author":{"id":243266,"type":"user","bot":false},"attachments":[{"id":581264,"fileName":"document.pdf","mimeType":"application/pdf","url":"https://driftapi.com/attachments/581264/data"}],"createdAt":1686304523000}]}}`)
}

// TestClient_GetMessages tests the method GetMessages()
func TestClient_GetMessages(t *testing.T) {
	t.Parallel()

	t.Run("get valid messages", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		messages, err := client.GetMessages(context.Background(), testConversationID, "")
		require.NoError(t, err)
		assert.NotNil(t, messages)
		assert.NotNil(t, messages.Data)
		assert.Len(t, messages.Data.Messages, 2)

		// Check first message
		msg := messages.Data.Messages[0]
		assert.Equal(t, testMessageID, msg.ID)
		assert.Equal(t, 12345, msg.OrgID)
		assert.Equal(t, testConversationID, msg.ConversationID)
		assert.Equal(t, "Hello, how can I help you?", msg.Body)
		assert.Equal(t, "chat", msg.Type)
		assert.NotNil(t, msg.Author)
		assert.Equal(t, uint64(243266), msg.Author.ID)
		assert.Equal(t, "user", msg.Author.Type)
		assert.False(t, msg.Author.Bot)
		assert.NotNil(t, msg.Context)
		assert.Equal(t, "192.168.1.1", msg.Context.IP)
		assert.Equal(t, "Mozilla/5.0", msg.Context.UserAgent)

		// Check second message (from contact)
		msg2 := messages.Data.Messages[1]
		assert.Equal(t, "contact", msg2.Author.Type)

		// Check pagination
		assert.NotNil(t, messages.Pagination)
		assert.Equal(t, "abc123", messages.Pagination.Next)
	})

	t.Run("get messages with pagination token", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		messages, err := client.GetMessages(context.Background(), testConversationID, "abc123")
		require.NoError(t, err)
		assert.NotNil(t, messages)
		assert.Len(t, messages.Data.Messages, 1)
		assert.Equal(t, uint64(987654323), messages.Data.Messages[0].ID)
	})

	t.Run("get messages with attachments", func(t *testing.T) {
		client := newTestClient(mockGetMessagesWithAttachments())

		messages, err := client.GetMessages(context.Background(), testConversationID, "")
		require.NoError(t, err)
		assert.NotNil(t, messages)
		assert.Len(t, messages.Data.Messages, 1)

		msg := messages.Data.Messages[0]
		assert.Len(t, msg.Attachments, 1)
		assert.Equal(t, uint64(581264), msg.Attachments[0].ID)
		assert.Equal(t, "document.pdf", msg.Attachments[0].FileName)
		assert.Equal(t, "application/pdf", msg.Attachments[0].MimeType)
		assert.Equal(t, "https://driftapi.com/attachments/581264/data", msg.Attachments[0].URL)
	})

	t.Run("get empty messages", func(t *testing.T) {
		client := newTestClient(mockGetMessagesEmpty())

		messages, err := client.GetMessages(context.Background(), testConversationID, "")
		require.NoError(t, err)
		assert.NotNil(t, messages)
		assert.Empty(t, messages.Data.Messages)
	})

	t.Run("missing conversation id", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		messages, err := client.GetMessages(context.Background(), 0, "")
		require.Error(t, err)
		assert.Equal(t, ErrMissingConversationID, err)
		assert.Nil(t, messages)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		messages, err := client.GetMessages(context.Background(), testConversationIDBadRequest, "")
		require.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		messages, err := client.GetMessages(context.Background(), testConversationIDUnauthorized, "")
		require.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("not found response", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		messages, err := client.GetMessages(context.Background(), testConversationIDNotFound, "")
		require.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		messages, err := client.GetMessages(context.Background(), testConversationIDBadJSON, "")
		require.Error(t, err)
		assert.Nil(t, messages)
	})
}

// TestClient_GetMessagesRaw tests the method GetMessagesRaw()
func TestClient_GetMessagesRaw(t *testing.T) {
	t.Parallel()

	t.Run("missing conversation id", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		response, err := client.GetMessagesRaw(context.Background(), 0, "")
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingConversationID, err)
	})

	t.Run("get valid messages", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		response, err := client.GetMessagesRaw(context.Background(), testConversationID, "")
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/conversations/116119985/messages", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("get messages with pagination token", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		response, err := client.GetMessagesRaw(context.Background(), testConversationID, "abc123")
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.Equal(t, apiEndpoint+"/conversations/116119985/messages?next=abc123", response.URL)
	})
}

// TestClient_GetMessagesNext tests the method GetMessagesNext()
func TestClient_GetMessagesNext(t *testing.T) {
	t.Parallel()

	t.Run("get next page of messages", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		// First get the initial page
		messages, err := client.GetMessages(context.Background(), testConversationID, "")
		require.NoError(t, err)
		require.NotNil(t, messages.Pagination)

		// Get the next page
		nextMessages, err := client.GetMessagesNext(context.Background(), testConversationID, messages)
		require.NoError(t, err)
		assert.NotNil(t, nextMessages)
		assert.Len(t, nextMessages.Data.Messages, 1)
	})

	t.Run("nil messages returns error", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		nextMessages, err := client.GetMessagesNext(context.Background(), testConversationID, nil)
		require.Error(t, err)
		assert.Equal(t, ErrNoNextPage, err)
		assert.Nil(t, nextMessages)
	})

	t.Run("nil pagination returns error", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		messages := &Messages{
			Data:       &MessagesListData{Messages: []*MessageData{}},
			Pagination: nil,
		}
		nextMessages, err := client.GetMessagesNext(context.Background(), testConversationID, messages)
		require.Error(t, err)
		assert.Equal(t, ErrNoNextPage, err)
		assert.Nil(t, nextMessages)
	})

	t.Run("empty next token returns error", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		messages := &Messages{
			Data:       &MessagesListData{Messages: []*MessageData{}},
			Pagination: &MessagesPagination{Next: ""},
		}
		nextMessages, err := client.GetMessagesNext(context.Background(), testConversationID, messages)
		require.Error(t, err)
		assert.Equal(t, ErrNoNextPage, err)
		assert.Nil(t, nextMessages)
	})
}

// TestClient_GetLatestMessage tests the method GetLatestMessage()
func TestClient_GetLatestMessage(t *testing.T) {
	t.Parallel()

	t.Run("get latest message", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		msg, err := client.GetLatestMessage(context.Background(), testConversationID)
		require.NoError(t, err)
		assert.NotNil(t, msg)
		// The second message has the higher createdAt timestamp
		assert.Equal(t, uint64(987654322), msg.ID)
	})

	t.Run("empty messages returns error", func(t *testing.T) {
		client := newTestClient(mockGetMessagesEmpty())

		msg, err := client.GetLatestMessage(context.Background(), testConversationID)
		require.Error(t, err)
		assert.Equal(t, ErrNoMessages, err)
		assert.Nil(t, msg)
	})
}

// TestClient_GetFirstMessage tests the method GetFirstMessage()
func TestClient_GetFirstMessage(t *testing.T) {
	t.Parallel()

	t.Run("get first message", func(t *testing.T) {
		client := newTestClient(mockGetMessages())

		msg, err := client.GetFirstMessage(context.Background(), testConversationID)
		require.NoError(t, err)
		assert.NotNil(t, msg)
		// The first message has the lower createdAt timestamp
		assert.Equal(t, uint64(987654321), msg.ID)
	})

	t.Run("empty messages returns error", func(t *testing.T) {
		client := newTestClient(mockGetMessagesEmpty())

		msg, err := client.GetFirstMessage(context.Background(), testConversationID)
		require.Error(t, err)
		assert.Equal(t, ErrNoMessages, err)
		assert.Nil(t, msg)
	})
}

// BenchmarkClient_GetMessages benchmarks the GetMessages method
func BenchmarkClient_GetMessages(b *testing.B) {
	client := newTestClient(mockGetMessages())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetMessages(context.Background(), testConversationID, "")
	}
}

// BenchmarkClient_GetMessagesRaw benchmarks the GetMessagesRaw method
func BenchmarkClient_GetMessagesRaw(b *testing.B) {
	client := newTestClient(mockGetMessages())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetMessagesRaw(context.Background(), testConversationID, "")
	}
}
