package drift

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Conversation status IDs for filtering
const (
	ConversationStatusOpen    = 1
	ConversationStatusClosed  = 2
	ConversationStatusPending = 3
)

// apiEndpointList is the alternate endpoint for listing conversations
const apiEndpointList = "https://api.drift.com"

// ConversationListQuery is the query parameters for listing conversations
type ConversationListQuery struct {
	Limit     int    // Max number of conversations to retrieve (max 100, default 25)
	PageToken string // Pagination token for next page
	StatusIDs []int  // Filter by status: 1=OPEN, 2=CLOSED, 3=PENDING
}

// BuildURL will build a URL for the list conversations query
func (q *ConversationListQuery) BuildURL() string {
	queryURL := apiEndpointList + "/conversations/list"

	params := make([]string, 0)

	// Add limit if specified
	if q.Limit > 0 {
		if q.Limit > 100 {
			q.Limit = 100
		}
		params = append(params, "limit="+strconv.Itoa(q.Limit))
	}

	// Add status filters
	for _, statusID := range q.StatusIDs {
		params = append(params, "statusId="+strconv.Itoa(statusID))
	}

	// Add page token if specified
	if len(q.PageToken) > 0 {
		params = append(params, "page_token="+q.PageToken)
	}

	if len(params) > 0 {
		queryURL += "?" + strings.Join(params, "&")
	}

	return queryURL
}

// ListConversations will list conversations with optional filters and pagination
// specs: https://devdocs.drift.com/docs/list-conversations
func (c *Client) ListConversations(ctx context.Context, query *ConversationListQuery) (conversations *Conversations, err error) {
	var response *RequestResponse
	if response, err = c.ListConversationsRaw(ctx, query); err != nil {
		return nil, err
	}

	conversations = new(Conversations)
	if err = response.UnmarshalTo(&conversations); err != nil {
		return nil, err
	}

	return conversations, nil
}

// ListConversationsRaw will fire the HTTP request to list conversations
// specs: https://devdocs.drift.com/docs/list-conversations
func (c *Client) ListConversationsRaw(ctx context.Context, query *ConversationListQuery) (*RequestResponse, error) {
	if query == nil {
		query = &ConversationListQuery{}
	}

	queryURL := query.BuildURL()
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}

// ListConversationsNext will get the next page of conversations using the pagination links
func (c *Client) ListConversationsNext(ctx context.Context, conversations *Conversations) (*Conversations, error) {
	if conversations == nil || conversations.Links == nil || len(conversations.Links.Next) == 0 {
		return nil, ErrNoNextPage
	}

	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            conversations.Links.Next,
		},
	)

	if response.Error != nil {
		return nil, response.Error
	}

	nextConversations := new(Conversations)
	if err := response.UnmarshalTo(&nextConversations); err != nil {
		return nil, err
	}

	return nextConversations, nil
}

// ListAllConversations will get all conversations by following pagination links
// Warning: This can be slow and memory-intensive for large datasets
func (c *Client) ListAllConversations(ctx context.Context, query *ConversationListQuery) (*Conversations, error) {
	allConversations := new(Conversations)
	allConversations.Data = make([]*conversationData, 0)

	conversations, err := c.ListConversations(ctx, query)
	if err != nil {
		return nil, err
	}

	allConversations.Data = append(allConversations.Data, conversations.Data...)

	for conversations.Links != nil && len(conversations.Links.Next) > 0 {
		conversations, err = c.ListConversationsNext(ctx, conversations)
		if err != nil {
			if errors.Is(err, ErrNoNextPage) {
				break
			}
			return nil, err
		}
		allConversations.Data = append(allConversations.Data, conversations.Data...)
	}

	return allConversations, nil
}

// ListConversationsByStatus is a convenience method to list conversations by a single status
func (c *Client) ListConversationsByStatus(ctx context.Context, statusID, limit int) (*Conversations, error) {
	query := &ConversationListQuery{
		Limit:     limit,
		StatusIDs: []int{statusID},
	}
	return c.ListConversations(ctx, query)
}

// ListOpenConversations is a convenience method to list open conversations
func (c *Client) ListOpenConversations(ctx context.Context, limit int) (*Conversations, error) {
	return c.ListConversationsByStatus(ctx, ConversationStatusOpen, limit)
}

// ListClosedConversations is a convenience method to list closed conversations
func (c *Client) ListClosedConversations(ctx context.Context, limit int) (*Conversations, error) {
	return c.ListConversationsByStatus(ctx, ConversationStatusClosed, limit)
}

// ListPendingConversations is a convenience method to list pending conversations
func (c *Client) ListPendingConversations(ctx context.Context, limit int) (*Conversations, error) {
	return c.ListConversationsByStatus(ctx, ConversationStatusPending, limit)
}

// ListConversationsByContactID lists all conversations for a specific contact
// This requires iterating through conversations as the API doesn't support direct filtering
func (c *Client) ListConversationsByContactID(ctx context.Context, contactID uint64, limit int) (*Conversations, error) {
	allConversations, err := c.ListAllConversations(ctx, &ConversationListQuery{Limit: limit})
	if err != nil {
		return nil, err
	}

	filtered := &Conversations{
		Data: make([]*conversationData, 0),
	}

	for _, conv := range allConversations.Data {
		if conv.ContactID == contactID {
			filtered.Data = append(filtered.Data, conv)
		}
	}

	return filtered, nil
}

// GetConversationCount returns the total number of conversations
// It is more efficient to use GetConversationStats for counts by status
func (c *Client) GetConversationCount(ctx context.Context) (int, error) {
	conversations, err := c.ListAllConversations(ctx, nil)
	if err != nil {
		return 0, err
	}
	return len(conversations.Data), nil
}

// internal helper to convert status code to string
func statusIDToString(statusID int) string {
	switch statusID {
	case ConversationStatusOpen:
		return "open"
	case ConversationStatusClosed:
		return "closed"
	case ConversationStatusPending:
		return "pending"
	default:
		return fmt.Sprintf("unknown(%d)", statusID)
	}
}
