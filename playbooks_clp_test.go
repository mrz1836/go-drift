package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testCLPPlaybookID     = uint64(67890)
	testCLPPlaybookName   = "Product Demo"
	testCLPLandingPageURL = "https://example.drift.com/demo"
)

// mockGetCLP returns a mock for CLP list operations
func mockGetCLP() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`[{"playbookId":67890,"playbookName":"Product Demo","landingPageUrl":"https://example.drift.com/demo"},{"playbookId":67891,"playbookName":"Sales Inquiry","landingPageUrl":"https://example.drift.com/sales"}]`),
	)
}

// mockGetCLPEmpty returns a mock for empty CLP list
func mockGetCLPEmpty() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`[]`),
	)
}

// mockGetCLPBadJSON returns a mock for bad JSON response
func mockGetCLPBadJSON() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`[{"playbookId":67890"playbookName":"Bad JSON"}]`),
	)
}

// TestClient_GetConversationalLandingPages tests the method GetConversationalLandingPages()
func TestClient_GetConversationalLandingPages(t *testing.T) {
	t.Parallel()

	t.Run("get valid conversational landing pages", func(t *testing.T) {
		client := newTestClient(mockGetCLP())

		pages, err := client.GetConversationalLandingPages(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, pages)
		assert.Len(t, pages.Data, 2)

		// Check first returned value
		assert.Equal(t, testCLPPlaybookID, pages.Data[0].PlaybookID)
		assert.Equal(t, testCLPPlaybookName, pages.Data[0].PlaybookName)
		assert.Equal(t, testCLPLandingPageURL, pages.Data[0].LandingPageURL)

		// Check second returned value
		assert.Equal(t, uint64(67891), pages.Data[1].PlaybookID)
		assert.Equal(t, "Sales Inquiry", pages.Data[1].PlaybookName)
		assert.Equal(t, "https://example.drift.com/sales", pages.Data[1].LandingPageURL)
	})

	t.Run("empty conversational landing pages list", func(t *testing.T) {
		client := newTestClient(mockGetCLPEmpty())

		pages, err := client.GetConversationalLandingPages(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, pages)
		assert.Empty(t, pages.Data)
	})

	t.Run("bad request response", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		pages, err := client.GetConversationalLandingPages(context.Background())
		require.Error(t, err)
		assert.Nil(t, pages)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		pages, err := client.GetConversationalLandingPages(context.Background())
		require.Error(t, err)
		assert.Nil(t, pages)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockGetCLPBadJSON())

		pages, err := client.GetConversationalLandingPages(context.Background())
		require.Error(t, err)
		assert.Nil(t, pages)
	})
}

// TestClient_GetConversationalLandingPagesRaw tests the method GetConversationalLandingPagesRaw()
func TestClient_GetConversationalLandingPagesRaw(t *testing.T) {
	t.Parallel()

	t.Run("get valid conversational landing pages raw", func(t *testing.T) {
		client := newTestClient(mockGetCLP())

		response, err := client.GetConversationalLandingPagesRaw(context.Background())
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/playbooks/clp", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("bad request response raw", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		response, err := client.GetConversationalLandingPagesRaw(context.Background())
		require.Error(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}

// BenchmarkClient_GetConversationalLandingPages benchmarks the GetConversationalLandingPages method
func BenchmarkClient_GetConversationalLandingPages(b *testing.B) {
	client := newTestClient(mockGetCLP())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetConversationalLandingPages(context.Background())
	}
}

// BenchmarkClient_GetConversationalLandingPagesRaw benchmarks the GetConversationalLandingPagesRaw method
func BenchmarkClient_GetConversationalLandingPagesRaw(b *testing.B) {
	client := newTestClient(mockGetCLP())
	for i := 0; i < b.N; i++ {
		_, _ = client.GetConversationalLandingPagesRaw(context.Background())
	}
}
