package drift

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// GetUser will get a single user by ID
// specs: https://devdocs.drift.com/docs/retrieving-user
func (c *Client) GetUser(ctx context.Context, userID uint64) (user *User, err error) {
	var response *RequestResponse
	if response, err = c.GetUserRaw(ctx, userID); err != nil {
		return nil, err
	}

	user = new(User)
	if err = response.UnmarshalTo(&user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserRaw will fire the HTTP request to retrieve the raw user data
// specs: https://devdocs.drift.com/docs/retrieving-user
func (c *Client) GetUserRaw(ctx context.Context, userID uint64) (*RequestResponse, error) {
	if err := requireID(userID, ErrMissingUserID); err != nil {
		return nil, err
	}

	queryURL := fmt.Sprintf("%s/users/%d", apiEndpoint, userID)
	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}

// GetUsers will get multiple users by their IDs (up to 20)
// specs: https://devdocs.drift.com/docs/retrieving-user
func (c *Client) GetUsers(ctx context.Context, userIDs []uint64) (users *Users, err error) {
	var response *RequestResponse
	if response, err = c.GetUsersRaw(ctx, userIDs); err != nil {
		return nil, err
	}

	// API returns a map structure for multiple users
	usersMap := new(UsersMap)
	if err = response.UnmarshalTo(&usersMap); err != nil {
		return nil, err
	}

	// Convert map to slice for consistent interface
	users = new(Users)
	users.Data = make([]*userData, 0, len(usersMap.Data))
	for _, user := range usersMap.Data {
		users.Data = append(users.Data, user)
	}

	return users, nil
}

// GetUsersRaw will fire the HTTP request to retrieve raw data for multiple users
// specs: https://devdocs.drift.com/docs/retrieving-user
func (c *Client) GetUsersRaw(ctx context.Context, userIDs []uint64) (*RequestResponse, error) {
	if len(userIDs) == 0 {
		return nil, ErrMissingUserID
	}

	if len(userIDs) > 20 {
		return nil, ErrTooManyUserIDs
	}

	// Build query string with multiple userId params
	params := make([]string, 0, len(userIDs))
	for _, id := range userIDs {
		params = append(params, "userId="+strconv.FormatUint(id, 10))
	}
	queryURL := fmt.Sprintf("%s/users?%s", apiEndpoint, strings.Join(params, "&"))

	response := httpRequest(
		ctx, c, &httpPayload{
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodGet,
			URL:            queryURL,
		},
	)

	return response, response.Error
}
