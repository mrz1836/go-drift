package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testConversationEmail = "customer@example.com"
	testConversationBody  = "Hello, I need help with your product."
)

// mockCreateConversation returns a mock for conversation creation operations
func mockCreateConversation() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":{"id":464032472,"contactId":1261122150,"inboxId":116983,"status":"open","participants":[228224],"createdAt":1548700064840}}`),
	)
}

// TestClient_CreateConversation tests the method CreateConversation()
func TestClient_CreateConversation(t *testing.T) {
	t.Parallel()

	t.Run("create valid conversation", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		conversation, err := client.CreateConversation(context.Background(), &NewConversationRequest{
			Email: testConversationEmail,
			Message: &NewConversationMessage{
				Body: testConversationBody,
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, conversation)
		assert.NotNil(t, conversation.Data)

		// Check returned values
		assert.Equal(t, uint64(464032472), conversation.Data.ID)
		assert.Equal(t, uint64(1261122150), conversation.Data.ContactID)
		assert.Equal(t, 116983, conversation.Data.InboxID)
		assert.Equal(t, "open", conversation.Data.Status)
		assert.Len(t, conversation.Data.Participants, 1)
		assert.Equal(t, uint64(228224), conversation.Data.Participants[0])
	})

	t.Run("create conversation with integration source", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		conversation, err := client.CreateConversation(context.Background(), &NewConversationRequest{
			Email: testConversationEmail,
			Message: &NewConversationMessage{
				Body: testConversationBody,
				Attributes: map[string]interface{}{
					"integrationSource": "Support Portal",
				},
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, conversation)
	})

	t.Run("create conversation with auto assignee", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		conversation, err := client.CreateConversation(context.Background(), &NewConversationRequest{
			Email: testConversationEmail,
			Message: &NewConversationMessage{
				Body: testConversationBody,
				Attributes: map[string]interface{}{
					"autoAssigneeId": 12345,
				},
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, conversation)
	})

	t.Run("missing email", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		conversation, err := client.CreateConversation(context.Background(), &NewConversationRequest{
			Email: "",
			Message: &NewConversationMessage{
				Body: testConversationBody,
			},
		})
		require.Error(t, err)
		assert.Equal(t, ErrMissingEmail, err)
		assert.Nil(t, conversation)
	})

	t.Run("nil request", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		conversation, err := client.CreateConversation(context.Background(), nil)
		require.Error(t, err)
		assert.Equal(t, ErrMissingEmail, err)
		assert.Nil(t, conversation)
	})

	t.Run("missing message body", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		conversation, err := client.CreateConversation(context.Background(), &NewConversationRequest{
			Email: testConversationEmail,
			Message: &NewConversationMessage{
				Body: "",
			},
		})
		require.Error(t, err)
		assert.Equal(t, ErrMissingMessageBody, err)
		assert.Nil(t, conversation)
	})

	t.Run("nil message", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		conversation, err := client.CreateConversation(context.Background(), &NewConversationRequest{
			Email:   testConversationEmail,
			Message: nil,
		})
		require.Error(t, err)
		assert.Equal(t, ErrMissingMessageBody, err)
		assert.Nil(t, conversation)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		conversation, err := client.CreateConversation(context.Background(), &NewConversationRequest{
			Email: testConversationEmail,
			Message: &NewConversationMessage{
				Body: testConversationBody,
			},
		})
		require.Error(t, err)
		assert.Nil(t, conversation)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"data":{"id":464032472"contactId":1261122150}}`))

		conversation, err := client.CreateConversation(context.Background(), &NewConversationRequest{
			Email: testConversationEmail,
			Message: &NewConversationMessage{
				Body: testConversationBody,
			},
		})
		require.Error(t, err)
		assert.Nil(t, conversation)
	})
}

// TestClient_CreateConversationRaw tests the method CreateConversationRaw()
func TestClient_CreateConversationRaw(t *testing.T) {
	t.Parallel()

	t.Run("missing email", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		response, err := client.CreateConversationRaw(context.Background(), &NewConversationRequest{})
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingEmail, err)
	})

	t.Run("missing message body", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		response, err := client.CreateConversationRaw(context.Background(), &NewConversationRequest{
			Email: testConversationEmail,
		})
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingMessageBody, err)
	})

	t.Run("create valid conversation", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		response, err := client.CreateConversationRaw(context.Background(), &NewConversationRequest{
			Email: testConversationEmail,
			Message: &NewConversationMessage{
				Body: testConversationBody,
			},
		})
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/conversations/new", response.URL)
		assert.Equal(t, http.MethodPost, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

// TestClient_CreateConversationSimple tests the convenience method CreateConversationSimple()
func TestClient_CreateConversationSimple(t *testing.T) {
	t.Parallel()

	t.Run("create conversation with simple method", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		conversation, err := client.CreateConversationSimple(context.Background(), testConversationEmail, testConversationBody)
		require.NoError(t, err)
		assert.NotNil(t, conversation)
	})
}

// TestClient_CreateConversationWithSource tests the convenience method CreateConversationWithSource()
func TestClient_CreateConversationWithSource(t *testing.T) {
	t.Parallel()

	t.Run("create conversation with source", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		conversation, err := client.CreateConversationWithSource(context.Background(), testConversationEmail, testConversationBody, "Support Portal")
		require.NoError(t, err)
		assert.NotNil(t, conversation)
	})
}

// TestClient_CreateConversationWithAssignee tests the convenience method CreateConversationWithAssignee()
func TestClient_CreateConversationWithAssignee(t *testing.T) {
	t.Parallel()

	t.Run("create conversation with assignee", func(t *testing.T) {
		client := newTestClient(mockCreateConversation())

		conversation, err := client.CreateConversationWithAssignee(context.Background(), testConversationEmail, testConversationBody, 12345)
		require.NoError(t, err)
		assert.NotNil(t, conversation)
	})
}

// BenchmarkClient_CreateConversation benchmarks the CreateConversation method
func BenchmarkClient_CreateConversation(b *testing.B) {
	client := newTestClient(mockCreateConversation())
	request := &NewConversationRequest{
		Email: testConversationEmail,
		Message: &NewConversationMessage{
			Body: testConversationBody,
		},
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateConversation(context.Background(), request)
	}
}

// BenchmarkClient_CreateConversationRaw benchmarks the CreateConversationRaw method
func BenchmarkClient_CreateConversationRaw(b *testing.B) {
	client := newTestClient(mockCreateConversation())
	request := &NewConversationRequest{
		Email: testConversationEmail,
		Message: &NewConversationMessage{
			Body: testConversationBody,
		},
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateConversationRaw(context.Background(), request)
	}
}
