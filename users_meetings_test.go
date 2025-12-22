package drift

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testMinStartTime = int64(1728946743006)
	testMaxStartTime = int64(1729551543006)
)

// mockGetMeetings returns a multi-route mock for meetings operations
func mockGetMeetings() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/users/meetings/org?min_start_time=1728946743006&max_start_time=1729551543006", http.StatusOK,
			`{"data":[{"conversationId":4019061071,"agentId":2322756,"orgId":12345,"status":"ACTIVE","meetingSource":"EMAIL_DROP","schedulerId":123456,"eventId":"event-123","slug":"meeting-slug","slotStart":1726491600000,"slotEnd":1726493400000,"updatedAt":1726491500000,"scheduledAt":1726490000000,"meetingType":"New Meeting","endUserTimeZone":"America/New_York","conferenceType":"ZOOM","isRecurring":false,"isPrivate":false},{"conversationId":4019061072,"agentId":2322757,"orgId":12345,"status":"CANCELED","meetingSource":"WIDGET","schedulerId":123457,"eventId":"event-456","slug":"meeting-slug-2","slotStart":1726578000000,"slotEnd":1726579800000,"updatedAt":1726577500000,"scheduledAt":1726576000000,"meetingType":"Follow Up","endUserTimeZone":"America/Los_Angeles","meetingNotes":"Test notes","bookedBy":2322756,"conferenceType":"GOOGLE_MEET","isRecurring":true,"isPrivate":true}]}`).
		addRoute(apiEndpoint+"/users/meetings/org?min_start_time=1728946743006&max_start_time=1729551543006&limit=50", http.StatusOK,
			`{"data":[{"conversationId":4019061071,"agentId":2322756,"orgId":12345,"status":"ACTIVE","meetingSource":"EMAIL_DROP","schedulerId":123456,"eventId":"event-123","slug":"meeting-slug","slotStart":1726491600000,"slotEnd":1726493400000,"updatedAt":1726491500000,"scheduledAt":1726490000000,"meetingType":"New Meeting","endUserTimeZone":"America/New_York","conferenceType":"ZOOM","isRecurring":false,"isPrivate":false}]}`)
}

// mockGetMeetingsEmpty returns a mock for empty meetings
func mockGetMeetingsEmpty() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/users/meetings/org?min_start_time=1728946743006&max_start_time=1729551543006", http.StatusOK,
			`{"data":[]}`)
}

// mockGetMeetingsUnauthorized returns a mock for unauthorized response
func mockGetMeetingsUnauthorized() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/users/meetings/org?min_start_time=1728946743006&max_start_time=1729551543006", http.StatusUnauthorized, "")
}

// mockGetMeetingsBadJSON returns a mock for bad JSON response
func mockGetMeetingsBadJSON() *mockHTTPMulti {
	return newMockHTTPMulti().
		addRoute(apiEndpoint+"/users/meetings/org?min_start_time=1728946743006&max_start_time=1729551543006", http.StatusOK,
			`{"data":[{"conversationId":4019061071"status":"Bad JSON"}]}`)
}

// TestMeetingsQuery_BuildURL tests the method BuildURL()
func TestMeetingsQuery_BuildURL(t *testing.T) {
	t.Parallel()

	t.Run("missing min_start_time", func(t *testing.T) {
		q := &MeetingsQuery{
			MaxStartTime: testMaxStartTime,
		}
		queryURL, err := q.BuildURL()
		require.Error(t, err)
		assert.Equal(t, ErrMissingMinStartTime, err)
		assert.Empty(t, queryURL)
	})

	t.Run("missing max_start_time", func(t *testing.T) {
		q := &MeetingsQuery{
			MinStartTime: testMinStartTime,
		}
		queryURL, err := q.BuildURL()
		require.Error(t, err)
		assert.Equal(t, ErrMissingMaxStartTime, err)
		assert.Empty(t, queryURL)
	})

	t.Run("valid query without limit", func(t *testing.T) {
		q := &MeetingsQuery{
			MinStartTime: testMinStartTime,
			MaxStartTime: testMaxStartTime,
		}
		queryURL, err := q.BuildURL()
		require.NoError(t, err)
		assert.Equal(t, apiEndpoint+"/users/meetings/org?min_start_time=1728946743006&max_start_time=1729551543006", queryURL)
	})

	t.Run("valid query with limit", func(t *testing.T) {
		q := &MeetingsQuery{
			MinStartTime: testMinStartTime,
			MaxStartTime: testMaxStartTime,
			Limit:        50,
		}
		queryURL, err := q.BuildURL()
		require.NoError(t, err)
		assert.Equal(t, apiEndpoint+"/users/meetings/org?min_start_time=1728946743006&max_start_time=1729551543006&limit=50", queryURL)
	})
}

// TestClient_GetBookedMeetings tests the method GetBookedMeetings()
func TestClient_GetBookedMeetings(t *testing.T) {
	t.Parallel()

	t.Run("get booked meetings", func(t *testing.T) {
		client := newTestClient(mockGetMeetings())

		query := &MeetingsQuery{
			MinStartTime: testMinStartTime,
			MaxStartTime: testMaxStartTime,
		}

		meetings, err := client.GetBookedMeetings(context.Background(), query)
		require.NoError(t, err)
		assert.NotNil(t, meetings)
		assert.Len(t, meetings.Data, 2)

		// Check first meeting
		assert.Equal(t, int64(4019061071), meetings.Data[0].ConversationID)
		assert.Equal(t, uint64(2322756), meetings.Data[0].AgentID)
		assert.Equal(t, uint64(12345), meetings.Data[0].OrgID)
		assert.Equal(t, "ACTIVE", meetings.Data[0].Status)
		assert.Equal(t, "EMAIL_DROP", meetings.Data[0].MeetingSource)
		assert.Equal(t, int64(123456), meetings.Data[0].SchedulerID)
		assert.Equal(t, "event-123", meetings.Data[0].EventID)
		assert.Equal(t, "meeting-slug", meetings.Data[0].Slug)
		assert.Equal(t, int64(1726491600000), meetings.Data[0].SlotStart)
		assert.Equal(t, int64(1726493400000), meetings.Data[0].SlotEnd)
		assert.Equal(t, "New Meeting", meetings.Data[0].MeetingType)
		assert.Equal(t, "America/New_York", meetings.Data[0].EndUserTimeZone)
		assert.Equal(t, "ZOOM", meetings.Data[0].ConferenceType)
		assert.False(t, meetings.Data[0].IsRecurring)
		assert.False(t, meetings.Data[0].IsPrivate)

		// Check second meeting
		assert.Equal(t, int64(4019061072), meetings.Data[1].ConversationID)
		assert.Equal(t, "CANCELED", meetings.Data[1].Status)
		assert.Equal(t, "Test notes", meetings.Data[1].MeetingNotes)
		assert.Equal(t, uint64(2322756), meetings.Data[1].BookedBy)
		assert.True(t, meetings.Data[1].IsRecurring)
		assert.True(t, meetings.Data[1].IsPrivate)
	})

	t.Run("get booked meetings with limit", func(t *testing.T) {
		client := newTestClient(mockGetMeetings())

		query := &MeetingsQuery{
			MinStartTime: testMinStartTime,
			MaxStartTime: testMaxStartTime,
			Limit:        50,
		}

		meetings, err := client.GetBookedMeetings(context.Background(), query)
		require.NoError(t, err)
		assert.NotNil(t, meetings)
	})

	t.Run("empty meetings", func(t *testing.T) {
		client := newTestClient(mockGetMeetingsEmpty())

		query := &MeetingsQuery{
			MinStartTime: testMinStartTime,
			MaxStartTime: testMaxStartTime,
		}

		meetings, err := client.GetBookedMeetings(context.Background(), query)
		require.NoError(t, err)
		assert.NotNil(t, meetings)
		assert.Empty(t, meetings.Data)
	})

	t.Run("missing min_start_time", func(t *testing.T) {
		client := newTestClient(mockGetMeetings())

		query := &MeetingsQuery{
			MaxStartTime: testMaxStartTime,
		}

		meetings, err := client.GetBookedMeetings(context.Background(), query)
		require.Error(t, err)
		assert.Equal(t, ErrMissingMinStartTime, err)
		assert.Nil(t, meetings)
	})

	t.Run("missing max_start_time", func(t *testing.T) {
		client := newTestClient(mockGetMeetings())

		query := &MeetingsQuery{
			MinStartTime: testMinStartTime,
		}

		meetings, err := client.GetBookedMeetings(context.Background(), query)
		require.Error(t, err)
		assert.Equal(t, ErrMissingMaxStartTime, err)
		assert.Nil(t, meetings)
	})

	t.Run("unauthorized response", func(t *testing.T) {
		client := newTestClient(mockGetMeetingsUnauthorized())

		query := &MeetingsQuery{
			MinStartTime: testMinStartTime,
			MaxStartTime: testMaxStartTime,
		}

		meetings, err := client.GetBookedMeetings(context.Background(), query)
		require.Error(t, err)
		assert.Nil(t, meetings)
	})

	t.Run("bad json response", func(t *testing.T) {
		client := newTestClient(mockGetMeetingsBadJSON())

		query := &MeetingsQuery{
			MinStartTime: testMinStartTime,
			MaxStartTime: testMaxStartTime,
		}

		meetings, err := client.GetBookedMeetings(context.Background(), query)
		require.Error(t, err)
		assert.Nil(t, meetings)
	})
}

// TestClient_GetBookedMeetingsRaw tests the method GetBookedMeetingsRaw()
func TestClient_GetBookedMeetingsRaw(t *testing.T) {
	t.Parallel()

	t.Run("missing min_start_time", func(t *testing.T) {
		client := newTestClient(mockGetMeetings())

		query := &MeetingsQuery{
			MaxStartTime: testMaxStartTime,
		}

		response, err := client.GetBookedMeetingsRaw(context.Background(), query)
		assert.Nil(t, response)
		require.Error(t, err)
		assert.Equal(t, ErrMissingMinStartTime, err)
	})

	t.Run("get booked meetings raw", func(t *testing.T) {
		client := newTestClient(mockGetMeetings())

		query := &MeetingsQuery{
			MinStartTime: testMinStartTime,
			MaxStartTime: testMaxStartTime,
		}

		response, err := client.GetBookedMeetingsRaw(context.Background(), query)
		assert.NotNil(t, response)
		require.NoError(t, err)
		assert.NoError(t, response.Error)

		assert.Equal(t, apiEndpoint+"/users/meetings/org?min_start_time=1728946743006&max_start_time=1729551543006", response.URL)
		assert.Equal(t, http.MethodGet, response.Method)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}

// BenchmarkClient_GetBookedMeetings benchmarks the GetBookedMeetings method
func BenchmarkClient_GetBookedMeetings(b *testing.B) {
	client := newTestClient(mockGetMeetings())
	query := &MeetingsQuery{
		MinStartTime: testMinStartTime,
		MaxStartTime: testMaxStartTime,
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.GetBookedMeetings(context.Background(), query)
	}
}

// BenchmarkClient_GetBookedMeetingsRaw benchmarks the GetBookedMeetingsRaw method
func BenchmarkClient_GetBookedMeetingsRaw(b *testing.B) {
	client := newTestClient(mockGetMeetings())
	query := &MeetingsQuery{
		MinStartTime: testMinStartTime,
		MaxStartTime: testMaxStartTime,
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.GetBookedMeetingsRaw(context.Background(), query)
	}
}
