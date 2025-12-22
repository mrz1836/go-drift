package drift

import (
	"context"
	"net/http"
)

// GetAccount will get an account by its ID
// specs: https://devdocs.drift.com/docs/retrieving-an-account
func (c *Client) GetAccount(ctx context.Context, accountID string) (account *Account, err error) {
	var response *RequestResponse
	if response, err = c.GetAccountRaw(ctx, accountID); err != nil {
		return nil, err
	}

	err = response.UnmarshalTo(&account)
	return account, err
}

// GetAccountRaw will fire the HTTP request to retrieve the raw account data
// specs: https://devdocs.drift.com/docs/retrieving-an-account
func (c *Client) GetAccountRaw(ctx context.Context, accountID string) (*RequestResponse, error) {
	if err := requireString(accountID, ErrMissingAccountID); err != nil {
		return nil, err
	}

	response := httpRequest(ctx, c, &httpPayload{
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodGet,
		URL:            apiEndpoint + "/accounts/" + accountID,
	})

	return response, response.Error
}
