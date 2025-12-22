package drift

import (
	"context"
	"fmt"
	"net/http"
)

// GetTranscript will get the formatted transcript of a conversation
// specs: https://devdocs.drift.com/docs/retrieving-a-conversations-transcript
func (c *Client) GetTranscript(ctx context.Context, conversationID uint64) (transcript string, err error) {
	var response *RequestResponse
	if response, err = c.GetTranscriptRaw(ctx, conversationID); err != nil {
		return "", err
	}

	// The API returns a data wrapper with the transcript string
	transcriptResponse := new(TranscriptResponse)
	if err = response.UnmarshalTo(&transcriptResponse); err != nil {
		return "", err
	}

	return transcriptResponse.Data, nil
}

// GetTranscriptRaw will fire the HTTP request to retrieve the raw transcript data
// specs: https://devdocs.drift.com/docs/retrieving-a-conversations-transcript
func (c *Client) GetTranscriptRaw(ctx context.Context, conversationID uint64) (*RequestResponse, error) {
	if conversationID == 0 {
		return nil, ErrMissingConversationID
	}

	queryURL := fmt.Sprintf("%s/conversations/%d/transcript", apiEndpoint, conversationID)
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}

// JSONTranscript represents the JSON transcript response
type JSONTranscript struct {
	Data *JSONTranscriptData `json:"data"`
}

// JSONTranscriptData contains the transcript messages
type JSONTranscriptData struct {
	Messages []*TranscriptMessage `json:"messages"`
}

// TranscriptMessage represents a message in the transcript
type TranscriptMessage struct {
	Author    *MessageAuthor `json:"author"`
	Body      string         `json:"body"`
	CreatedAt int64          `json:"createdAt"`
	ID        uint64         `json:"id"`
	Type      string         `json:"type"`
}

// GetJSONTranscript will get the JSON transcript of a conversation
// specs: https://devdocs.drift.com/docs/retrieving-a-conversations-transcript
func (c *Client) GetJSONTranscript(ctx context.Context, conversationID uint64) (transcript *JSONTranscript, err error) {
	var response *RequestResponse
	if response, err = c.GetJSONTranscriptRaw(ctx, conversationID); err != nil {
		return nil, err
	}

	transcript = new(JSONTranscript)
	if err = response.UnmarshalTo(&transcript); err != nil {
		return nil, err
	}

	return transcript, nil
}

// GetJSONTranscriptRaw will fire the HTTP request to retrieve the raw JSON transcript data
// specs: https://devdocs.drift.com/docs/retrieving-a-conversations-transcript
func (c *Client) GetJSONTranscriptRaw(ctx context.Context, conversationID uint64) (*RequestResponse, error) {
	if conversationID == 0 {
		return nil, ErrMissingConversationID
	}

	queryURL := fmt.Sprintf("%s/conversations/%d/json_transcript", apiEndpoint, conversationID)
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}
