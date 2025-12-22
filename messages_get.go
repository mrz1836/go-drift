package drift

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// GetMessages will get messages for a conversation
// specs: https://devdocs.drift.com/docs/retrieve-a-conversations-messages
func (c *Client) GetMessages(ctx context.Context, conversationID uint64, next string) (messages *Messages, err error) {
	var response *RequestResponse
	if response, err = c.GetMessagesRaw(ctx, conversationID, next); err != nil {
		return nil, err
	}

	messages = new(Messages)
	if err = response.UnmarshalTo(&messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// GetMessagesRaw will fire the HTTP request to retrieve the raw messages data
// specs: https://devdocs.drift.com/docs/retrieve-a-conversations-messages
func (c *Client) GetMessagesRaw(ctx context.Context, conversationID uint64, next string) (*RequestResponse, error) {
	if err := requireID(conversationID, ErrMissingConversationID); err != nil {
		return nil, err
	}

	queryURL := fmt.Sprintf("%s/conversations/%d/messages", apiEndpoint, conversationID)
	if len(next) > 0 {
		queryURL += "?next=" + next
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

// GetMessagesNext will get the next page of messages using the pagination token
func (c *Client) GetMessagesNext(ctx context.Context, conversationID uint64, messages *Messages) (*Messages, error) {
	if messages == nil || messages.Pagination == nil || len(messages.Pagination.Next) == 0 {
		return nil, ErrNoNextPage
	}

	return c.GetMessages(ctx, conversationID, messages.Pagination.Next)
}

// GetAllMessages will get all messages for a conversation by following pagination
// Warning: This can be slow and memory-intensive for conversations with many messages
func (c *Client) GetAllMessages(ctx context.Context, conversationID uint64) (*Messages, error) {
	allMessages := new(Messages)
	allMessages.Data = &MessagesListData{
		Messages: make([]*MessageData, 0),
	}

	messages, err := c.GetMessages(ctx, conversationID, "")
	if err != nil {
		return nil, err
	}

	if messages.Data != nil {
		allMessages.Data.Messages = append(allMessages.Data.Messages, messages.Data.Messages...)
	}

	for messages.Pagination != nil && len(messages.Pagination.Next) > 0 {
		messages, err = c.GetMessagesNext(ctx, conversationID, messages)
		if err != nil {
			if errors.Is(err, ErrNoNextPage) {
				break
			}
			return nil, err
		}
		if messages.Data == nil {
			break
		}
		allMessages.Data.Messages = append(allMessages.Data.Messages, messages.Data.Messages...)
	}

	return allMessages, nil
}

// GetMessageCount returns the count of messages in a conversation
func (c *Client) GetMessageCount(ctx context.Context, conversationID uint64) (int, error) {
	allMessages, err := c.GetAllMessages(ctx, conversationID)
	if err != nil {
		return 0, err
	}
	if allMessages.Data == nil {
		return 0, nil
	}
	return len(allMessages.Data.Messages), nil
}

// GetLatestMessage returns the most recent message in a conversation
func (c *Client) GetLatestMessage(ctx context.Context, conversationID uint64) (*MessageData, error) {
	messages, err := c.GetMessages(ctx, conversationID, "")
	if err != nil {
		return nil, err
	}

	if messages.Data == nil || len(messages.Data.Messages) == 0 {
		return nil, ErrNoMessages
	}

	// Messages are typically returned in chronological order, so the last one is the most recent
	// But we'll find the one with the highest createdAt timestamp to be safe
	var latest *MessageData
	for _, msg := range messages.Data.Messages {
		if latest == nil || msg.CreatedAt > latest.CreatedAt {
			latest = msg
		}
	}

	return latest, nil
}

// GetFirstMessage returns the first message in a conversation
func (c *Client) GetFirstMessage(ctx context.Context, conversationID uint64) (*MessageData, error) {
	messages, err := c.GetMessages(ctx, conversationID, "")
	if err != nil {
		return nil, err
	}

	if messages.Data == nil || len(messages.Data.Messages) == 0 {
		return nil, ErrNoMessages
	}

	// Find the message with the lowest createdAt timestamp
	var first *MessageData
	for _, msg := range messages.Data.Messages {
		if first == nil || msg.CreatedAt < first.CreatedAt {
			first = msg
		}
	}

	return first, nil
}
