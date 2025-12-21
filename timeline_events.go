package drift

import (
	"context"
	"encoding/json"
	"net/http"
)

// TimelineEvent is the timeline event object
type TimelineEvent struct {
	Attributes map[string]string `json:"attributes,omitempty"`
	ContactID  uint64            `json:"contactId"`
	CreatedAt  uint64            `json:"createdAt,omitempty"`
	Event      string            `json:"event"`
	ExternalID string            `json:"externalId,omitempty"`
}

// TimelineResponse is the response from creating a timeline event
type TimelineResponse struct {
	Data *TimelineEvent `json:"data"`
}

// CreateTimelineEvent will create a new timeline event
// specs: https://devdocs.drift.com/docs/posting-timeline-events
func (c *Client) CreateTimelineEvent(ctx context.Context,
	event *TimelineEvent,
) (response *TimelineResponse, err error) {
	// Marshall the attributes
	var data []byte
	if data, err = json.Marshal(event); err != nil {
		return response, err
	}

	// Create and fire the request
	var resp *RequestResponse
	if resp = httpRequest(
		ctx, c, &httpPayload{
			Data:           data,
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodPost,
			URL:            apiEndpoint + "/contacts/timeline",
		},
	); resp.Error != nil {
		err = resp.Error
		return response, err
	}

	// Parse the request
	err = json.Unmarshal(resp.BodyContents, &response)
	return response, err
}
