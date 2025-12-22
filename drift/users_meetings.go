package drift

import (
	"context"
	"fmt"
	"net/http"
)

// MeetingsQuery is the query parameters for getting booked meetings
type MeetingsQuery struct {
	MinStartTime int64 `json:"min_start_time"` // Required, epoch milliseconds
	MaxStartTime int64 `json:"max_start_time"` // Required, epoch milliseconds
	Limit        int   `json:"limit"`          // Optional, 0-1000, default 100
}

// BuildURL will build a url for the meetings query
func (q *MeetingsQuery) BuildURL() (queryURL string, err error) {
	if q.MinStartTime == 0 {
		return "", ErrMissingMinStartTime
	}
	if q.MaxStartTime == 0 {
		return "", ErrMissingMaxStartTime
	}

	queryURL = fmt.Sprintf("%s/users/meetings/org?min_start_time=%d&max_start_time=%d",
		apiEndpoint, q.MinStartTime, q.MaxStartTime)

	if q.Limit > 0 {
		queryURL = fmt.Sprintf("%s&limit=%d", queryURL, q.Limit)
	}

	return queryURL, nil
}

// GetBookedMeetings will get booked meetings for the organization.
// This endpoint only returns meetings booked on dates up to 30 days in the past.
// specs: https://devdocs.drift.com/docs/get-booked-meetings
func (c *Client) GetBookedMeetings(ctx context.Context, query *MeetingsQuery) (meetings *Meetings, err error) {
	var response *RequestResponse
	if response, err = c.GetBookedMeetingsRaw(ctx, query); err != nil {
		return nil, err
	}

	meetings = new(Meetings)
	if err = response.UnmarshalTo(&meetings); err != nil {
		return nil, err
	}

	return meetings, nil
}

// GetBookedMeetingsRaw will fire the HTTP request to retrieve the raw meetings data
// specs: https://devdocs.drift.com/docs/get-booked-meetings
func (c *Client) GetBookedMeetingsRaw(ctx context.Context, query *MeetingsQuery) (*RequestResponse, error) {
	queryURL, err := query.BuildURL()
	if err != nil {
		return nil, err
	}

	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}
