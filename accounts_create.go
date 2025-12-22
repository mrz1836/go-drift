package drift

import (
	"context"
	"encoding/json"
	"net/http"
)

// CreateAccount will fire the HTTP request to create a new account
// specs: https://devdocs.drift.com/docs/creating-an-account
func (c *Client) CreateAccount(ctx context.Context, fields *AccountFields) (account *Account, err error) {
	var response *RequestResponse
	if response, err = c.CreateAccountRaw(ctx, fields); err != nil {
		return nil, err
	}

	err = response.UnmarshalTo(&account)
	return account, err
}

// CreateAccountRaw will create an account and return the raw response
// specs: https://devdocs.drift.com/docs/creating-an-account
func (c *Client) CreateAccountRaw(ctx context.Context, fields *AccountFields) (*RequestResponse, error) {
	// Marshal the fields
	data, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	// Fire the request
	response := httpRequest(ctx, c, &httpPayload{
		Data:           data,
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodPost,
		URL:            apiEndpoint + "/accounts/create",
	})

	return response, response.Error
}
