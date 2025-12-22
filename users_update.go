package drift

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UpdateUser will update an existing user
// specs: https://devdocs.drift.com/docs/updating-users
func (c *Client) UpdateUser(ctx context.Context, userID uint64,
	fields *UserUpdateFields,
) (user *User, err error) {
	var response *RequestResponse
	if response, err = c.UpdateUserRaw(ctx, userID, fields); err != nil {
		return nil, err
	}

	user = new(User)
	if err = json.Unmarshal(response.BodyContents, &user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserRaw will update an existing user using a custom attribute struct
// specs: https://devdocs.drift.com/docs/updating-users
func (c *Client) UpdateUserRaw(ctx context.Context, userID uint64,
	fields interface{},
) (*RequestResponse, error) {
	if userID == 0 {
		return nil, ErrMissingUserID
	}

	// Marshal the fields to JSON
	data, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	queryURL := fmt.Sprintf("%s/users/update?userId=%d", apiEndpoint, userID)
	response := httpRequest(
		ctx, c, &httpPayload{
			Data:           data,
			ExpectedStatus: http.StatusOK,
			Method:         http.MethodPatch,
			URL:            queryURL,
		},
	)

	return response, response.Error
}
