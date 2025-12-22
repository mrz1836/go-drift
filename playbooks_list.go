package drift

import (
	"context"
	"net/http"
)

// GetPlaybooks retrieves all enabled and active playbooks for the organization.
// Playbook configuration is cached for 10 minutes.
// specs: https://devdocs.drift.com/docs/get-playbooks
func (c *Client) GetPlaybooks(ctx context.Context) (playbooks *Playbooks, err error) {
	var response *RequestResponse
	if response, err = c.GetPlaybooksRaw(ctx); err != nil {
		return nil, err
	}

	// API returns an array directly, not wrapped in "data"
	var playbookList []*playbookData
	if err = response.UnmarshalTo(&playbookList); err != nil {
		return nil, err
	}

	playbooks = &Playbooks{
		Data: playbookList,
	}

	return playbooks, nil
}

// GetPlaybooksRaw will fire the HTTP request to retrieve the raw playbooks data
// specs: https://devdocs.drift.com/docs/get-playbooks
func (c *Client) GetPlaybooksRaw(ctx context.Context) (*RequestResponse, error) {
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            apiEndpoint + "/playbooks/list",
		},
	)

	return response, response.Error
}
