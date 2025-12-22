package drift

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Message types
const (
	MessageTypeChat        = "chat"
	MessageTypePrivateNote = "private_note"
)

// CreateMessage will create a new message in a conversation
// specs: https://devdocs.drift.com/docs/creating-a-message
func (c *Client) CreateMessage(ctx context.Context, conversationID uint64, request *CreateMessageRequest) (messages *Messages, err error) {
	var response *RequestResponse
	if response, err = c.CreateMessageRaw(ctx, conversationID, request); err != nil {
		return nil, err
	}

	messages = new(Messages)
	if err = response.UnmarshalTo(&messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// CreateMessageRaw will fire the HTTP request to create a new message
// specs: https://devdocs.drift.com/docs/creating-a-message
func (c *Client) CreateMessageRaw(ctx context.Context, conversationID uint64, request *CreateMessageRequest) (*RequestResponse, error) {
	if err := requireID(conversationID, ErrMissingConversationID); err != nil {
		return nil, err
	}
	if request == nil {
		return nil, ErrMissingMessageType
	}
	if err := requireString(request.Type, ErrMissingMessageType); err != nil {
		return nil, err
	}

	// Marshal the request to JSON
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	queryURL := fmt.Sprintf("%s/conversations/%d/messages", apiEndpoint, conversationID)
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

// SendChatMessage is a convenience method to send a simple chat message
func (c *Client) SendChatMessage(ctx context.Context, conversationID uint64, body string) (*Messages, error) {
	return c.CreateMessage(ctx, conversationID, &CreateMessageRequest{
		Type: MessageTypeChat,
		Body: body,
	})
}

// SendPrivateNote is a convenience method to send a private note (only visible to agents)
func (c *Client) SendPrivateNote(ctx context.Context, conversationID uint64, body string) (*Messages, error) {
	return c.CreateMessage(ctx, conversationID, &CreateMessageRequest{
		Type: MessageTypePrivateNote,
		Body: body,
	})
}

// SendChatMessageAsUser is a convenience method to send a chat message as a specific user
func (c *Client) SendChatMessageAsUser(ctx context.Context, conversationID uint64, body string, userID uint64) (*Messages, error) {
	return c.CreateMessage(ctx, conversationID, &CreateMessageRequest{
		Type:   MessageTypeChat,
		Body:   body,
		UserID: userID,
	})
}

// SendMessageWithButtons is a convenience method to send a message with interactive buttons
func (c *Client) SendMessageWithButtons(ctx context.Context, conversationID uint64, body string, buttons []*MessageButton) (*Messages, error) {
	return c.CreateMessage(ctx, conversationID, &CreateMessageRequest{
		Type:    MessageTypeChat,
		Body:    body,
		Buttons: buttons,
	})
}

// NewReplyButton creates a new reply button for use in messages
func NewReplyButton(label, value string) *MessageButton {
	return &MessageButton{
		Label: label,
		Value: value,
		Type:  "reply",
	}
}

// NewPrimaryButton creates a new primary-styled reply button
func NewPrimaryButton(label, value string) *MessageButton {
	return &MessageButton{
		Label: label,
		Value: value,
		Type:  "reply",
		Style: "primary",
	}
}

// NewDangerButton creates a new danger-styled reply button
func NewDangerButton(label, value string) *MessageButton {
	return &MessageButton{
		Label: label,
		Value: value,
		Type:  "reply",
		Style: "danger",
	}
}

// NewButtonWithReaction creates a button that shows a reaction message when clicked
func NewButtonWithReaction(label, value, reactionType, reactionMessage string) *MessageButton {
	return &MessageButton{
		Label: label,
		Value: value,
		Type:  "reply",
		Reaction: &ButtonReaction{
			Type:    reactionType,
			Message: reactionMessage,
		},
	}
}
