package drift

import (
	"context"
	"fmt"
	"net/http"
)

// GetConversation will get a single conversation by ID
// specs: https://devdocs.drift.com/docs/retrieve-a-conversation
func (c *Client) GetConversation(ctx context.Context, conversationID uint64) (conversation *Conversation, err error) {
	var response *RequestResponse
	if response, err = c.GetConversationRaw(ctx, conversationID); err != nil {
		return nil, err
	}

	conversation = new(Conversation)
	if err = response.UnmarshalTo(&conversation); err != nil {
		return nil, err
	}

	return conversation, nil
}

// GetConversationRaw will fire the HTTP request to retrieve the raw conversation data
// specs: https://devdocs.drift.com/docs/retrieve-a-conversation
func (c *Client) GetConversationRaw(ctx context.Context, conversationID uint64) (*RequestResponse, error) {
	if err := requireID(conversationID, ErrMissingConversationID); err != nil {
		return nil, err
	}

	queryURL := fmt.Sprintf("%s/conversations/%d", apiEndpoint, conversationID)
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}
