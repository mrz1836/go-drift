package drift

import (
	"context"
	"encoding/json"
	"net/http"
)

// UpdateAccount will fire the HTTP request to update an existing account
// specs: https://devdocs.drift.com/docs/updating-accounts
func (c *Client) UpdateAccount(ctx context.Context, fields *AccountFields) (account *Account, err error) {
	var response *RequestResponse
	if response, err = c.UpdateAccountRaw(ctx, fields); err != nil {
		return nil, err
	}

	err = response.UnmarshalTo(&account)
	return account, err
}

// UpdateAccountRaw will update an account and return the raw response
// specs: https://devdocs.drift.com/docs/updating-accounts
func (c *Client) UpdateAccountRaw(ctx context.Context, fields *AccountFields) (*RequestResponse, error) {
	// Validate required fields
	if fields == nil {
		return nil, ErrMissingAccountID
	}
	if err := requireString(fields.AccountID, ErrMissingAccountID); err != nil {
		return nil, err
	}
	if err := requireID(fields.OwnerID, ErrMissingOwnerID); err != nil {
		return nil, err
	}

	// Marshal the fields
	data, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	// Fire the request
	response := httpRequest(ctx, c, &httpPayload{
		Data:           data,
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodPatch,
		URL:            apiEndpoint + "/accounts/update",
	})

	return response, response.Error
}
