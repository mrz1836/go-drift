package drift

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockHTTPTimelineEvents for mocking requests
type mockHTTPTimelineEvents struct{}

// Do is a mock http request
func (m *mockHTTPTimelineEvents) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, errMissingRequest
	}

	// Valid response
	if req.URL.String() == apiEndpoint+"/contacts/timeline" {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBufferString(`{"data":{"attributes":{},"event":"` + testEventName + `","createdAt":1614571424495,"contactId":` + testContactID + `}}`))
	}

	// Default is valid
	return resp, nil
}

// TestClient_CreateTimelineEvent tests the method CreateTimelineEvent()
func TestClient_CreateTimelineEvent(t *testing.T) {
	t.Parallel()

	t.Run("create a timeline event", func(t *testing.T) {
		// Create a client
		client := newTestClient(&mockHTTPTimelineEvents{})

		id, err := strconv.ParseUint(testContactID, 10, 64)
		require.NoError(t, err)

		// Create a req
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
}

// BenchmarkClient_CreateTimelineEvent benchmarks the CreateTimelineEvent method
func BenchmarkClient_CreateTimelineEvent(b *testing.B) {
	client := newTestClient(&mockHTTPCreateContact{})
	id, _ := strconv.ParseUint(testContactID, 10, 64)
	fields := &TimelineEvent{
		ContactID: id,
		Event:     testEventName,
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateTimelineEvent(context.Background(), fields)
	}
}
