package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCreateMessage returns a multi-route mock for message creation operations
func mockCreateMessage() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/conversations/116119985/messages", http.StatusOK,
			`{"data":{"messages":[{"id":123456789,"orgId":12345,"conversationId":116119985,"body":"Hello from the API!","type":"chat","author":{"id":228224,"type":"user","bot":false},"createdAt":1686304600000,"context":{"ip":"192.168.1.1","userAgent":"Go-Drift-Client"}}]}}`).
		addRoute(apiEndpoint+"/conversations/111111111/messages", http.StatusBadRequest, "").
		addRoute(apiEndpoint+"/conversations/222222222/messages", http.StatusUnauthorized, "").
		addRoute(apiEndpoint+"/conversations/333333333/messages", http.StatusOK,
			`{"data":{"messages":[{"id":123456789"body":"Bad JSON"}`).
		addRoute(apiEndpoint+"/conversations/444444444/messages", http.StatusNotFound, "")
}

// mockCreateMessageWithButtons returns a mock for messages with buttons
func mockCreateMessageWithButtons() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/conversations/116119985/messages", http.StatusOK,
			`{"data":{"messages":[{"id":123456789,"orgId":12345,"conversationId":116119985,"body":"Please choose an option:","type":"chat","author":{"id":228224,"type":"user","bot":true},"buttons":[{"label":"Yes","value":"yes","type":"reply","style":"primary"},{"label":"No","value":"no","type":"reply","style":"danger"}],"createdAt":1686304600000}]}}`)
}

// TestClient_CreateMessage tests the method CreateMessage()
func TestClient_CreateMessage(t *testing.T) {
	t.Parallel()

	t.Run("create valid chat message", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.CreateMessage(context.Background(), testConversationID, &CreateMessageRequest{
			Type: MessageTypeChat,
			Body: "Hello from the API!",
		})
		require.NoError(t, err)
		assert.NotNil(t, messages)
		assert.NotNil(t, messages.Data)
		assert.Len(t, messages.Data.Messages, 1)

		msg := messages.Data.Messages[0]
		assert.Equal(t, uint64(123456789), msg.ID)
		assert.Equal(t, testConversationID, msg.ConversationID)
		assert.Equal(t, "Hello from the API!", msg.Body)
		assert.Equal(t, MessageTypeChat, msg.Type)
		assert.NotNil(t, msg.Author)
		assert.Equal(t, "user", msg.Author.Type)
	})

	t.Run("create private note", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.CreateMessage(context.Background(), testConversationID, &CreateMessageRequest{
			Type: MessageTypePrivateNote,
			Body: "This is a private note",
		})
		require.NoError(t, err)
		assert.NotNil(t, messages)
	})

	t.Run("create message with user id", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.CreateMessage(context.Background(), testConversationID, &CreateMessageRequest{
			Type:   MessageTypeChat,
			Body:   "Hello from the API!",
			UserID: 228224,
		})
		require.NoError(t, err)
		assert.NotNil(t, messages)
	})

	t.Run("create message with buttons", func(t *testing.T) {
		client := newTestClient(mockCreateMessageWithButtons())

		messages, err := client.CreateMessage(context.Background(), testConversationID, &CreateMessageRequest{
			Type: MessageTypeChat,
			Body: "Please choose an option:",
			Buttons: []*MessageButton{
				NewPrimaryButton("Yes", "yes"),
				NewDangerButton("No", "no"),
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, messages)
		assert.Len(t, messages.Data.Messages, 1)

		msg := messages.Data.Messages[0]
		assert.Len(t, msg.Buttons, 2)
		assert.Equal(t, "Yes", msg.Buttons[0].Label)
		assert.Equal(t, "primary", msg.Buttons[0].Style)
		assert.Equal(t, "No", msg.Buttons[1].Label)
		assert.Equal(t, "danger", msg.Buttons[1].Style)
	})

	t.Run("missing conversation id", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.CreateMessage(context.Background(), 0, &CreateMessageRequest{
			Type: MessageTypeChat,
			Body: "Hello!",
		})
		require.Error(t, err)
		assert.Equal(t, ErrMissingConversationID, err)
		assert.Nil(t, messages)
	})

	t.Run("nil request", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.CreateMessage(context.Background(), testConversationID, nil)
		require.Error(t, err)
		assert.Equal(t, ErrMissingMessageType, err)
		assert.Nil(t, messages)
	})

	t.Run("missing message type", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.CreateMessage(context.Background(), testConversationID, &CreateMessageRequest{
			Body: "Hello!",
		})
		require.Error(t, err)
		assert.Equal(t, ErrMissingMessageType, err)
		assert.Nil(t, messages)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.CreateMessage(context.Background(), testConversationIDBadRequest, &CreateMessageRequest{
			Type: MessageTypeChat,
			Body: "Hello!",
		})
		require.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.CreateMessage(context.Background(), testConversationIDUnauthorized, &CreateMessageRequest{
			Type: MessageTypeChat,
			Body: "Hello!",
		})
		require.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("not found response", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.CreateMessage(context.Background(), testConversationIDNotFound, &CreateMessageRequest{
			Type: MessageTypeChat,
			Body: "Hello!",
		})
		require.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.CreateMessage(context.Background(), testConversationIDBadJSON, &CreateMessageRequest{
			Type: MessageTypeChat,
			Body: "Hello!",
		})
		require.Error(t, err)
		assert.Nil(t, messages)
	})
}

// TestClient_CreateMessageRaw tests the method CreateMessageRaw()
func TestClient_CreateMessageRaw(t *testing.T) {
	t.Parallel()

	t.Run("missing conversation id", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		response, err := client.CreateMessageRaw(context.Background(), 0, &CreateMessageRequest{
			Type: MessageTypeChat,
			Body: "Hello!",
		})
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingConversationID, err)
	})

	t.Run("missing message type", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		response, err := client.CreateMessageRaw(context.Background(), testConversationID, &CreateMessageRequest{
			Body: "Hello!",
		})
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingMessageType, err)
	})

	t.Run("create valid message", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		response, err := client.CreateMessageRaw(context.Background(), testConversationID, &CreateMessageRequest{
			Type: MessageTypeChat,
			Body: "Hello from the API!",
		})
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/conversations/116119985/messages", response.URL)
		assert.Equal(t, http.MethodPost, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

// TestClient_SendChatMessage tests the convenience method SendChatMessage()
func TestClient_SendChatMessage(t *testing.T) {
	t.Parallel()

	t.Run("send chat message", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.SendChatMessage(context.Background(), testConversationID, "Hello from the API!")
		require.NoError(t, err)
		assert.NotNil(t, messages)
	})
}

// TestClient_SendPrivateNote tests the convenience method SendPrivateNote()
func TestClient_SendPrivateNote(t *testing.T) {
	t.Parallel()

	t.Run("send private note", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.SendPrivateNote(context.Background(), testConversationID, "This is a private note")
		require.NoError(t, err)
		assert.NotNil(t, messages)
	})
}

// TestClient_SendChatMessageAsUser tests the convenience method SendChatMessageAsUser()
func TestClient_SendChatMessageAsUser(t *testing.T) {
	t.Parallel()

	t.Run("send chat message as user", func(t *testing.T) {
		client := newTestClient(mockCreateMessage())

		messages, err := client.SendChatMessageAsUser(context.Background(), testConversationID, "Hello!", 228224)
		require.NoError(t, err)
		assert.NotNil(t, messages)
	})
}

// TestClient_SendMessageWithButtons tests the convenience method SendMessageWithButtons()
func TestClient_SendMessageWithButtons(t *testing.T) {
	t.Parallel()

	t.Run("send message with buttons", func(t *testing.T) {
		client := newTestClient(mockCreateMessageWithButtons())

		buttons := []*MessageButton{
			NewPrimaryButton("Yes", "yes"),
			NewDangerButton("No", "no"),
		}

		messages, err := client.SendMessageWithButtons(context.Background(), testConversationID, "Please choose:", buttons)
		require.NoError(t, err)
		assert.NotNil(t, messages)
	})
}

// TestNewReplyButton tests the NewReplyButton helper
func TestNewReplyButton(t *testing.T) {
	t.Parallel()

	button := NewReplyButton("Click me", "clicked")
	assert.Equal(t, "Click me", button.Label)
	assert.Equal(t, "clicked", button.Value)
	assert.Equal(t, "reply", button.Type)
	assert.Empty(t, button.Style)
}

// TestNewPrimaryButton tests the NewPrimaryButton helper
func TestNewPrimaryButton(t *testing.T) {
	t.Parallel()

	button := NewPrimaryButton("Submit", "submit")
	assert.Equal(t, "Submit", button.Label)
	assert.Equal(t, "submit", button.Value)
	assert.Equal(t, "reply", button.Type)
	assert.Equal(t, "primary", button.Style)
}

// TestNewDangerButton tests the NewDangerButton helper
func TestNewDangerButton(t *testing.T) {
	t.Parallel()

	button := NewDangerButton("Delete", "delete")
	assert.Equal(t, "Delete", button.Label)
	assert.Equal(t, "delete", button.Value)
	assert.Equal(t, "reply", button.Type)
	assert.Equal(t, "danger", button.Style)
}

// TestNewButtonWithReaction tests the NewButtonWithReaction helper
func TestNewButtonWithReaction(t *testing.T) {
	t.Parallel()

	button := NewButtonWithReaction("Thanks", "thanks", "message", "Thank you for your feedback!")
	assert.Equal(t, "Thanks", button.Label)
	assert.Equal(t, "thanks", button.Value)
	assert.Equal(t, "reply", button.Type)
	assert.NotNil(t, button.Reaction)
	assert.Equal(t, "message", button.Reaction.Type)
	assert.Equal(t, "Thank you for your feedback!", button.Reaction.Message)
}

// BenchmarkClient_CreateMessage benchmarks the CreateMessage method
func BenchmarkClient_CreateMessage(b *testing.B) {
	client := newTestClient(mockCreateMessage())
	request := &CreateMessageRequest{
		Type: MessageTypeChat,
		Body: "Hello from the API!",
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateMessage(context.Background(), testConversationID, request)
	}
}

// BenchmarkClient_CreateMessageRaw benchmarks the CreateMessageRaw method
func BenchmarkClient_CreateMessageRaw(b *testing.B) {
	client := newTestClient(mockCreateMessage())
	request := &CreateMessageRequest{
		Type: MessageTypeChat,
		Body: "Hello from the API!",
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateMessageRaw(context.Background(), testConversationID, request)
	}
}
