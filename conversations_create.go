package drift

import (
	"context"
	"encoding/json"
	"net/http"
)

// CreateConversation will create a new conversation with a contact
// specs: https://devdocs.drift.com/docs/creating-a-conversation
func (c *Client) CreateConversation(ctx context.Context, request *NewConversationRequest) (conversation *Conversation, err error) {
	var response *RequestResponse
	if response, err = c.CreateConversationRaw(ctx, request); err != nil {
		return nil, err
	}

	conversation = new(Conversation)
	if err = response.UnmarshalTo(&conversation); err != nil {
		return nil, err
	}

	return conversation, nil
}

// CreateConversationRaw will fire the HTTP request to create a new conversation
// specs: https://devdocs.drift.com/docs/creating-a-conversation
func (c *Client) CreateConversationRaw(ctx context.Context, request *NewConversationRequest) (*RequestResponse, error) {
	// Validate required fields
	if request == nil {
		return nil, ErrMissingEmail
	}
	if err := requireString(request.Email, ErrMissingEmail); err != nil {
		return nil, err
	}
	if request.Message == nil {
		return nil, ErrMissingMessageBody
	}
	if err := requireString(request.Message.Body, ErrMissingMessageBody); err != nil {
		return nil, err
	}

	// Marshal the request to JSON
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	queryURL := apiEndpoint + "/conversations/new"
	response := httpRequest(
		ctx, c, &httpPayload{
			Data:           data,
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodPost,
			URL:            queryURL,
		},
	)

	return response, response.Error
}

// CreateConversationSimple is a convenience method to create a conversation with just email and body
func (c *Client) CreateConversationSimple(ctx context.Context, email, messageBody string) (*Conversation, error) {
	return c.CreateConversation(ctx, &NewConversationRequest{
		Email: email,
		Message: &NewConversationMessage{
			Body: messageBody,
		},
	})
}

// CreateConversationWithSource is a convenience method to create a conversation with an integration source
func (c *Client) CreateConversationWithSource(ctx context.Context, email, messageBody, integrationSource string) (*Conversation, error) {
	return c.CreateConversation(ctx, &NewConversationRequest{
		Email: email,
		Message: &NewConversationMessage{
			Body: messageBody,
			Attributes: map[string]interface{}{
				"integrationSource": integrationSource,
			},
		},
	})
}

// CreateConversationWithAssignee is a convenience method to create a conversation with auto-assignment
func (c *Client) CreateConversationWithAssignee(ctx context.Context, email, messageBody string, assigneeID uint64) (*Conversation, error) {
	return c.CreateConversation(ctx, &NewConversationRequest{
		Email: email,
		Message: &NewConversationMessage{
			Body: messageBody,
			Attributes: map[string]interface{}{
				"autoAssigneeId": assigneeID,
			},
		},
	})
}
