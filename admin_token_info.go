package drift

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetTokenInfo retrieves information about an access token.
// specs: https://devdocs.drift.com/docs/get-token-information
func (c *Client) GetTokenInfo(ctx context.Context, accessToken string) (tokenInfo *TokenInfo, err error) {
	// Create and fire the request
	var reqResponse *RequestResponse
	if reqResponse, err = c.GetTokenInfoRaw(ctx, accessToken); err != nil {
		return nil, err
	}

	// Parse the response
	err = reqResponse.UnmarshalTo(&tokenInfo)
	return tokenInfo, err
}

// GetTokenInfoRaw retrieves token information and returns the raw response
// specs: https://devdocs.drift.com/docs/get-token-information
func (c *Client) GetTokenInfoRaw(ctx context.Context, accessToken string) (*RequestResponse, error) {
	// Validate required fields
	if err := requireString(accessToken, ErrMissingAccessToken); err != nil {
		return nil, err
	}

	// Build the request body
	requestBody := &TokenInfoRequest{
		AccessToken: accessToken,
	}

	data, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Fire the request
	response := httpRequest(ctx, c, &httpPayload{
		Data:           data,
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodPost,
		URL:            apiEndpoint + "/app/token_info",
	})

	return response, response.Error
}
