package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAttachmentID             = uint64(581264)
	testAttachmentIDBadRequest   = uint64(111111)
	testAttachmentIDUnauthorized = uint64(222222)
	testAttachmentIDNotFound     = uint64(444444)
)

// mockGetAttachment returns a multi-route mock for attachment operations
func mockGetAttachment() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/attachments/581264/data", http.StatusOK, "%PDF-1.4 simulated pdf content here").
		addRoute(apiEndpoint+"/attachments/111111/data", http.StatusBadRequest, "").
		addRoute(apiEndpoint+"/attachments/222222/data", http.StatusUnauthorized, "").
		addRoute(apiEndpoint+"/attachments/444444/data", http.StatusNotFound, "")
}

// TestClient_GetAttachment tests the method GetAttachment()
func TestClient_GetAttachment(t *testing.T) {
	t.Parallel()

	t.Run("get valid attachment", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		data, err := client.GetAttachment(context.Background(), testAttachmentID)
		require.NoError(t, err)
		assert.NotNil(t, data)
		assert.NotEmpty(t, data.Data)
		assert.Contains(t, string(data.Data), "%PDF-1.4")
	})

	t.Run("missing attachment id", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		data, err := client.GetAttachment(context.Background(), 0)
		require.Error(t, err)
		assert.Equal(t, ErrMissingAttachmentID, err)
		assert.Nil(t, data)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		data, err := client.GetAttachment(context.Background(), testAttachmentIDBadRequest)
		require.Error(t, err)
		assert.Nil(t, data)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		data, err := client.GetAttachment(context.Background(), testAttachmentIDUnauthorized)
		require.Error(t, err)
		assert.Nil(t, data)
	})

	t.Run("not found response", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		data, err := client.GetAttachment(context.Background(), testAttachmentIDNotFound)
		require.Error(t, err)
		assert.Nil(t, data)
	})
}

// TestClient_GetAttachmentRaw tests the method GetAttachmentRaw()
func TestClient_GetAttachmentRaw(t *testing.T) {
	t.Parallel()

	t.Run("missing attachment id", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		response, err := client.GetAttachmentRaw(context.Background(), 0)
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingAttachmentID, err)
	})

	t.Run("get valid attachment", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		response, err := client.GetAttachmentRaw(context.Background(), testAttachmentID)
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/attachments/581264/data", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

// TestClient_GetAttachmentFromMessage tests the convenience method GetAttachmentFromMessage()
func TestClient_GetAttachmentFromMessage(t *testing.T) {
	t.Parallel()

	t.Run("get attachment from message attachment info", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		attachment := &MessageAttachment{
			ID:       testAttachmentID,
			FileName: "document.pdf",
			MimeType: "application/pdf",
			URL:      "https://driftapi.com/attachments/581264/data",
		}

		data, err := client.GetAttachmentFromMessage(context.Background(), attachment)
		require.NoError(t, err)
		assert.NotNil(t, data)
		assert.NotEmpty(t, data.Data)
		assert.Equal(t, "application/pdf", data.MimeType)
	})

	t.Run("nil attachment returns error", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		data, err := client.GetAttachmentFromMessage(context.Background(), nil)
		require.Error(t, err)
		assert.Equal(t, ErrMissingAttachmentID, err)
		assert.Nil(t, data)
	})
}

// TestClient_GetAllAttachmentsFromMessage tests the convenience method GetAllAttachmentsFromMessage()
func TestClient_GetAllAttachmentsFromMessage(t *testing.T) {
	t.Parallel()

	t.Run("get all attachments from message", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		message := &MessageData{
			ID:   123456789,
			Body: "Here are the files",
			Attachments: []*MessageAttachment{
				{
					ID:       testAttachmentID,
					FileName: "document.pdf",
					MimeType: "application/pdf",
					URL:      "https://driftapi.com/attachments/581264/data",
				},
			},
		}

		attachments, err := client.GetAllAttachmentsFromMessage(context.Background(), message)
		require.NoError(t, err)
		assert.Len(t, attachments, 1)
		assert.Equal(t, "application/pdf", attachments[0].MimeType)
	})

	t.Run("nil message returns nil", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		attachments, err := client.GetAllAttachmentsFromMessage(context.Background(), nil)
		require.NoError(t, err)
		assert.Nil(t, attachments)
	})

	t.Run("message with no attachments returns nil", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		message := &MessageData{
			ID:          123456789,
			Body:        "No attachments here",
			Attachments: nil,
		}

		attachments, err := client.GetAllAttachmentsFromMessage(context.Background(), message)
		require.NoError(t, err)
		assert.Nil(t, attachments)
	})

	t.Run("message with empty attachments returns nil", func(t *testing.T) {
		client := newTestClient(mockGetAttachment())

		message := &MessageData{
			ID:          123456789,
			Body:        "No attachments here",
			Attachments: []*MessageAttachment{},
		}

		attachments, err := client.GetAllAttachmentsFromMessage(context.Background(), message)
		require.NoError(t, err)
		assert.Nil(t, attachments)
	})
}

// BenchmarkClient_GetAttachment benchmarks the GetAttachment method
func BenchmarkClient_GetAttachment(b *testing.B) {
	client := newTestClient(mockGetAttachment())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetAttachment(context.Background(), testAttachmentID)
	}
}

// BenchmarkClient_GetAttachmentRaw benchmarks the GetAttachmentRaw method
func BenchmarkClient_GetAttachmentRaw(b *testing.B) {
	client := newTestClient(mockGetAttachment())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetAttachmentRaw(context.Background(), testAttachmentID)
	}
}
