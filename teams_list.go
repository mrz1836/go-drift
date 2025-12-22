package drift

import (
	"context"
	"fmt"
	"net/http"
)

// ListTeams will get the full list of teams in the organization
// specs: https://devdocs.drift.com/docs/listing-teams-org
func (c *Client) ListTeams(ctx context.Context) (teams *Teams, err error) {
	var response *RequestResponse
	if response, err = c.ListTeamsRaw(ctx); err != nil {
		return nil, err
	}

	teams = new(Teams)
	if err = response.UnmarshalTo(&teams); err != nil {
		return nil, err
	}

	return teams, nil
}

// ListTeamsRaw will fire the HTTP request to retrieve the raw teams list data
// specs: https://devdocs.drift.com/docs/listing-teams-org
func (c *Client) ListTeamsRaw(ctx context.Context) (*RequestResponse, error) {
	queryURL := apiEndpoint + "/teams/org"
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}

// ListTeamsByUser will get the list of teams for a specific user
// specs: https://devdocs.drift.com/docs/listing-teams-by-user
func (c *Client) ListTeamsByUser(ctx context.Context, userID uint64) (teams *Teams, err error) {
	var response *RequestResponse
	if response, err = c.ListTeamsByUserRaw(ctx, userID); err != nil {
		return nil, err
	}

	teams = new(Teams)
	if err = response.UnmarshalTo(&teams); err != nil {
		return nil, err
	}

	return teams, nil
}

// ListTeamsByUserRaw will fire the HTTP request to retrieve the raw teams data for a user
// specs: https://devdocs.drift.com/docs/listing-teams-by-user
func (c *Client) ListTeamsByUserRaw(ctx context.Context, userID uint64) (*RequestResponse, error) {
	if userID == 0 {
		return nil, ErrMissingUserID
	}

	queryURL := fmt.Sprintf("%s/teams/users/%d", apiEndpoint, userID)
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}
