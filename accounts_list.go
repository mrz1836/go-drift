package drift

import (
	"context"
	"fmt"
	"net/http"
)

// BuildURL builds the URL for listing accounts with pagination parameters
func (q *AccountListQuery) BuildURL() string {
	baseURL := apiEndpoint + "/accounts"

	if q == nil {
		return baseURL
	}

	params := ""
	if q.Index > 0 {
		params = fmt.Sprintf("index=%d", q.Index)
	}
	if q.Size > 0 {
		if len(params) > 0 {
			params += "&"
		}
		params += fmt.Sprintf("size=%d", q.Size)
	}

	if len(params) > 0 {
		return baseURL + "?" + params
	}
	return baseURL
}

// ListAccounts will get a paginated list of accounts
// specs: https://devdocs.drift.com/docs/listing-accounts
func (c *Client) ListAccounts(ctx context.Context, query *AccountListQuery) (accounts *Accounts, err error) {
	var response *RequestResponse
	if response, err = c.ListAccountsRaw(ctx, query); err != nil {
		return nil, err
	}

	err = response.UnmarshalTo(&accounts)
	return accounts, err
}

// ListAccountsRaw will fire the HTTP request to retrieve the raw accounts list data
// specs: https://devdocs.drift.com/docs/listing-accounts
func (c *Client) ListAccountsRaw(ctx context.Context, query *AccountListQuery) (*RequestResponse, error) {
	queryURL := query.BuildURL()

	response := httpRequest(ctx, c, &httpPayload{
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodGet,
		URL:            queryURL,
	})

	return response, response.Error
}

// ListAccountsNext will get the next page of accounts using the Next URL from a previous response
// specs: https://devdocs.drift.com/docs/listing-accounts
func (c *Client) ListAccountsNext(ctx context.Context, accounts *Accounts) (*Accounts, error) {
	if accounts == nil || accounts.Data == nil || len(accounts.Data.Next) == 0 {
		return nil, ErrNoNextPage
	}

	// The Next field contains a relative URL like "/accounts?index=XXX&size=XXX"
	queryURL := apiEndpoint + accounts.Data.Next

	response := httpRequest(ctx, c, &httpPayload{
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodGet,
		URL:            queryURL,
	})

	if response.Error != nil {
		return nil, response.Error
	}

	var nextAccounts *Accounts
	err := response.UnmarshalTo(&nextAccounts)
	return nextAccounts, err
}
