package drift

import (
	"context"
	"net/http"
)

// DeleteAccount will fire the HTTP request to delete an existing account
// This performs a soft delete - the account will still be queryable by ID with deleted:true
// specs: https://devdocs.drift.com/docs/deleting-accounts
func (c *Client) DeleteAccount(ctx context.Context, accountID string) (response *StandardResponse, err error) {
	var reqResponse *RequestResponse
	if reqResponse, err = c.DeleteAccountRaw(ctx, accountID); err != nil {
		return nil, err
	}

	err = reqResponse.UnmarshalTo(&response)
	return response, err
}

// DeleteAccountRaw will delete an account and return the raw response
// specs: https://devdocs.drift.com/docs/deleting-accounts
func (c *Client) DeleteAccountRaw(ctx context.Context, accountID string) (*RequestResponse, error) {
	if err := requireString(accountID, ErrMissingAccountID); err != nil {
		return nil, err
	}

	response := httpRequest(ctx, c, &httpPayload{
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodDelete,
		URL:            apiEndpoint + "/accounts/" + accountID,
	})

	return response, response.Error
}
