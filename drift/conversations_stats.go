package drift

import (
	"context"
	"net/http"
)

// GetConversationStats will get the bulk conversation status counts
// specs: https://devdocs.drift.com/docs/bulk-conversation-statuses
func (c *Client) GetConversationStats(ctx context.Context) (stats *ConversationStats, err error) {
	var response *RequestResponse
	if response, err = c.GetConversationStatsRaw(ctx); err != nil {
		return nil, err
	}

	stats = new(ConversationStats)
	if err = response.UnmarshalTo(&stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// GetConversationStatsRaw will fire the HTTP request to retrieve the raw conversation stats
// specs: https://devdocs.drift.com/docs/bulk-conversation-statuses
func (c *Client) GetConversationStatsRaw(ctx context.Context) (*RequestResponse, error) {
	queryURL := apiEndpoint + "/conversations/stats"
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}

// GetOpenConversationCount returns the count of open conversations
func (c *Client) GetOpenConversationCount(ctx context.Context) (int, error) {
	stats, err := c.GetConversationStats(ctx)
	if err != nil {
		return 0, err
	}
	return stats.ConversationCount["OPEN"], nil
}

// GetClosedConversationCount returns the count of closed conversations
func (c *Client) GetClosedConversationCount(ctx context.Context) (int, error) {
	stats, err := c.GetConversationStats(ctx)
	if err != nil {
		return 0, err
	}
	return stats.ConversationCount["CLOSED"], nil
}

// GetPendingConversationCount returns the count of pending conversations
func (c *Client) GetPendingConversationCount(ctx context.Context) (int, error) {
	stats, err := c.GetConversationStats(ctx)
	if err != nil {
		return 0, err
	}
	return stats.ConversationCount["PENDING"], nil
}

// GetTotalConversationCount returns the total count of all conversations
func (c *Client) GetTotalConversationCount(ctx context.Context) (int, error) {
	stats, err := c.GetConversationStats(ctx)
	if err != nil {
		return 0, err
	}

	total := 0
	for _, count := range stats.ConversationCount {
		total += count
	}
	return total, nil
}
