package drift

import (
	"context"
	"net/http"
)

// GetConversationalLandingPages retrieves all conversational landing pages for the organization
// specs: https://devdocs.drift.com/docs/retrieve-conversational-landing-pages
func (c *Client) GetConversationalLandingPages(ctx context.Context) (pages *ConversationalLandingPages, err error) {
	var response *RequestResponse
	if response, err = c.GetConversationalLandingPagesRaw(ctx); err != nil {
		return nil, err
	}

	// API returns an array directly, not wrapped in "data"
	var pageList []*ConversationalLandingPage
	if err = response.UnmarshalTo(&pageList); err != nil {
		return nil, err
	}

	pages = &ConversationalLandingPages{
		Data: pageList,
	}

	return pages, nil
}

// GetConversationalLandingPagesRaw will fire the HTTP request to retrieve the raw conversational landing pages data
// specs: https://devdocs.drift.com/docs/retrieve-conversational-landing-pages
func (c *Client) GetConversationalLandingPagesRaw(ctx context.Context) (*RequestResponse, error) {
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/playbooks/clp",
		},
	)

	return response, response.Error
}
