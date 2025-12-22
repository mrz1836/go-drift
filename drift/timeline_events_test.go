package drift

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTimelineEvents returns a mock for timeline event operations
func mockTimelineEvents() *mockHTTP {
	return newMockHTTP(
		withStatus(http.StatusOK),
		withBody(`{"data":{"attributes":{},"event":"`+testEventName+`","createdAt":1614571424495,"contactId":`+testContactID+`}}`),
	)
}

// TestClient_CreateTimelineEvent tests the method CreateTimelineEvent()
func TestClient_CreateTimelineEvent(t *testing.T) {
	t.Parallel()

	t.Run("create a timeline event", func(t *testing.T) {
		client := newTestClient(mockTimelineEvents())

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		var resp *TimelineResponse
		resp, err = client.CreateTimelineEvent(
			context.Background(), &TimelineEvent{
				ContactID: id,
				Event:     testEventName,
			})
		require.NoError(t, err)
		assert.NotNil(t, resp)

		// Got a contact
		assert.Equal(t, testEventName, resp.Data.Event)
		assert.Equal(t, uint64(1614571424495), resp.Data.CreatedAt)
		assert.Equal(t, id, resp.Data.ContactID)
	})

	t.Run("returns error on HTTP failure", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusBadRequest))

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		resp, err := client.CreateTimelineEvent(
			context.Background(), &TimelineEvent{
				ContactID: id,
				Event:     testEventName,
			})

		require.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, ErrMalformedRequest)
	})

	t.Run("returns error on 401 unauthorized", func(t *testing.T) {
		client := newTestClient(newMockError(http.StatusUnauthorized))

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		resp, err := client.CreateTimelineEvent(
			context.Background(), &TimelineEvent{
				ContactID: id,
				Event:     testEventName,
			})

		require.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("returns error on response unmarshal failure", func(t *testing.T) {
		client := newTestClient(newMockSuccess(`{"data":{"invalid json`))

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		resp, err := client.CreateTimelineEvent(
			context.Background(), &TimelineEvent{
				ContactID: id,
				Event:     testEventName,
			})

		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

// BenchmarkClient_CreateTimelineEvent benchmarks the CreateTimelineEvent method
func BenchmarkClient_CreateTimelineEvent(b *testing.B) {
	client := newTestClient(mockTimelineEvents())
	id, _ := strconv.ParseUint(testContactID, 10, 64)
	fields := &TimelineEvent{
		ContactID: id,
		Event:     testEventName,
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateTimelineEvent(context.Background(), fields)
	}
}
