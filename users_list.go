package drift

import (
	"context"
	"encoding/json"
	"net/http"
)

// ListUsers will get the full list of users in the organization
// specs: https://devdocs.drift.com/docs/listing-users
func (c *Client) ListUsers(ctx context.Context) (users *Users, err error) {
	var response *RequestResponse
	if response, err = c.ListUsersRaw(ctx); err != nil {
		return nil, err
	}

	users = new(Users)
	if err = json.Unmarshal(response.BodyContents, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// ListUsersRaw will fire the HTTP request to retrieve the raw user list data
// specs: https://devdocs.drift.com/docs/listing-users
func (c *Client) ListUsersRaw(ctx context.Context) (*RequestResponse, error) {
	queryURL := apiEndpoint + "/users/list"
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}
