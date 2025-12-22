package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockGetTranscript returns a multi-route mock for transcript operations
func mockGetTranscript() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/conversations/116119985/transcript", http.StatusOK,
			`{"data":"[2023-06-09 10:15:23] John Doe: Hello, how can I help you?\n[2023-06-09 10:15:45] Customer: I have a question about your product.\n[2023-06-09 10:16:02] John Doe: Sure, I'd be happy to help!"}`).
		addRoute(apiEndpoint+"/conversations/111111111/transcript", http.StatusBadRequest, "").
		addRoute(apiEndpoint+"/conversations/222222222/transcript", http.StatusUnauthorized, "").
		addRoute(apiEndpoint+"/conversations/333333333/transcript", http.StatusOK,
			`{"data":"[2023-06-09 10:15:23] John Doe: Hello`).
		addRoute(apiEndpoint+"/conversations/444444444/transcript", http.StatusNotFound, "")
}

// mockGetJSONTranscript returns a multi-route mock for JSON transcript operations
func mockGetJSONTranscript() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/conversations/116119985/json_transcript", http.StatusOK,
			`{"data":{"messages":[{"id":987654321,"body":"Hello, how can I help you?","type":"chat","author":{"id":243266,"type":"user","bot":false},"createdAt":1686304523000},{"id":987654322,"body":"I have a question about your product.","type":"chat","author":{"id":903182234,"type":"contact","bot":false},"createdAt":1686304545000},{"id":987654323,"body":"Sure, I'd be happy to help!","type":"chat","author":{"id":243266,"type":"user","bot":false},"createdAt":1686304562000}]}}`).
		addRoute(apiEndpoint+"/conversations/111111111/json_transcript", http.StatusBadRequest, "").
		addRoute(apiEndpoint+"/conversations/222222222/json_transcript", http.StatusUnauthorized, "").
		addRoute(apiEndpoint+"/conversations/333333333/json_transcript", http.StatusOK,
			`{"data":{"messages":[{"id":987654321"body":"Hello"}`).
		addRoute(apiEndpoint+"/conversations/444444444/json_transcript", http.StatusNotFound, "")
}

// TestClient_GetTranscript tests the method GetTranscript()
func TestClient_GetTranscript(t *testing.T) {
	t.Parallel()

	t.Run("get valid transcript", func(t *testing.T) {
		client := newTestClient(mockGetTranscript())

		transcript, err := client.GetTranscript(context.Background(), testConversationID)
		require.NoError(t, err)
		assert.NotEmpty(t, transcript)
		assert.Contains(t, transcript, "Hello, how can I help you?")
		assert.Contains(t, transcript, "I have a question about your product.")
	})

	t.Run("missing conversation id", func(t *testing.T) {
		client := newTestClient(mockGetTranscript())

		transcript, err := client.GetTranscript(context.Background(), 0)
		require.Error(t, err)
		assert.Equal(t, ErrMissingConversationID, err)
		assert.Empty(t, transcript)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(mockGetTranscript())

		transcript, err := client.GetTranscript(context.Background(), testConversationIDBadRequest)
		require.Error(t, err)
		assert.Empty(t, transcript)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(mockGetTranscript())

		transcript, err := client.GetTranscript(context.Background(), testConversationIDUnauthorized)
		require.Error(t, err)
		assert.Empty(t, transcript)
	})

	t.Run("not found response", func(t *testing.T) {
		client := newTestClient(mockGetTranscript())

		transcript, err := client.GetTranscript(context.Background(), testConversationIDNotFound)
		require.Error(t, err)
		assert.Empty(t, transcript)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockGetTranscript())

		transcript, err := client.GetTranscript(context.Background(), testConversationIDBadJSON)
		require.Error(t, err)
		assert.Empty(t, transcript)
	})
}

// TestClient_GetTranscriptRaw tests the method GetTranscriptRaw()
func TestClient_GetTranscriptRaw(t *testing.T) {
	t.Parallel()

	t.Run("missing conversation id", func(t *testing.T) {
		client := newTestClient(mockGetTranscript())

		response, err := client.GetTranscriptRaw(context.Background(), 0)
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingConversationID, err)
	})

	t.Run("get valid transcript", func(t *testing.T) {
		client := newTestClient(mockGetTranscript())

		response, err := client.GetTranscriptRaw(context.Background(), testConversationID)
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/conversations/116119985/transcript", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

// TestClient_GetJSONTranscript tests the method GetJSONTranscript()
func TestClient_GetJSONTranscript(t *testing.T) {
	t.Parallel()

	t.Run("get valid json transcript", func(t *testing.T) {
		client := newTestClient(mockGetJSONTranscript())

		transcript, err := client.GetJSONTranscript(context.Background(), testConversationID)
		require.NoError(t, err)
		assert.NotNil(t, transcript)
		assert.NotNil(t, transcript.Data)
		assert.Len(t, transcript.Data.Messages, 3)

		// Check first message
		assert.Equal(t, uint64(987654321), transcript.Data.Messages[0].ID)
		assert.Equal(t, "Hello, how can I help you?", transcript.Data.Messages[0].Body)
		assert.Equal(t, "chat", transcript.Data.Messages[0].Type)
		assert.Equal(t, uint64(243266), transcript.Data.Messages[0].Author.ID)
		assert.Equal(t, "user", transcript.Data.Messages[0].Author.Type)
		assert.False(t, transcript.Data.Messages[0].Author.Bot)

		// Check second message (from contact)
		assert.Equal(t, "contact", transcript.Data.Messages[1].Author.Type)
		assert.Equal(t, uint64(903182234), transcript.Data.Messages[1].Author.ID)
	})

	t.Run("missing conversation id", func(t *testing.T) {
		client := newTestClient(mockGetJSONTranscript())

		transcript, err := client.GetJSONTranscript(context.Background(), 0)
		require.Error(t, err)
		assert.Equal(t, ErrMissingConversationID, err)
		assert.Nil(t, transcript)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(mockGetJSONTranscript())

		transcript, err := client.GetJSONTranscript(context.Background(), testConversationIDBadRequest)
		require.Error(t, err)
		assert.Nil(t, transcript)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(mockGetJSONTranscript())

		transcript, err := client.GetJSONTranscript(context.Background(), testConversationIDUnauthorized)
		require.Error(t, err)
		assert.Nil(t, transcript)
	})

	t.Run("not found response", func(t *testing.T) {
		client := newTestClient(mockGetJSONTranscript())

		transcript, err := client.GetJSONTranscript(context.Background(), testConversationIDNotFound)
		require.Error(t, err)
		assert.Nil(t, transcript)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockGetJSONTranscript())

		transcript, err := client.GetJSONTranscript(context.Background(), testConversationIDBadJSON)
		require.Error(t, err)
		assert.Nil(t, transcript)
	})
}

// TestClient_GetJSONTranscriptRaw tests the method GetJSONTranscriptRaw()
func TestClient_GetJSONTranscriptRaw(t *testing.T) {
	t.Parallel()

	t.Run("missing conversation id", func(t *testing.T) {
		client := newTestClient(mockGetJSONTranscript())

		response, err := client.GetJSONTranscriptRaw(context.Background(), 0)
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingConversationID, err)
	})

	t.Run("get valid json transcript", func(t *testing.T) {
		client := newTestClient(mockGetJSONTranscript())

		response, err := client.GetJSONTranscriptRaw(context.Background(), testConversationID)
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/conversations/116119985/json_transcript", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

// BenchmarkClient_GetTranscript benchmarks the GetTranscript method
func BenchmarkClient_GetTranscript(b *testing.B) {
	client := newTestClient(mockGetTranscript())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetTranscript(context.Background(), testConversationID)
	}
}

// BenchmarkClient_GetTranscriptRaw benchmarks the GetTranscriptRaw method
func BenchmarkClient_GetTranscriptRaw(b *testing.B) {
	client := newTestClient(mockGetTranscript())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetTranscriptRaw(context.Background(), testConversationID)
	}
}

// BenchmarkClient_GetJSONTranscript benchmarks the GetJSONTranscript method
func BenchmarkClient_GetJSONTranscript(b *testing.B) {
	client := newTestClient(mockGetJSONTranscript())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetJSONTranscript(context.Background(), testConversationID)
	}
}

// BenchmarkClient_GetJSONTranscriptRaw benchmarks the GetJSONTranscriptRaw method
func BenchmarkClient_GetJSONTranscriptRaw(b *testing.B) {
	client := newTestClient(mockGetJSONTranscript())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetJSONTranscriptRaw(context.Background(), testConversationID)
	}
}
