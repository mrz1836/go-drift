package drift

import (
	"context"
	"fmt"
	"net/http"
)

// AppUninstall triggers app uninstallation for a client.
// This only works for "public" integrations that clients have connected through OAuth.
// Upon success, the client's access token used to make the request will be invalidated.
// specs: https://devdocs.drift.com/docs/app-uninstall
func (c *Client) AppUninstall(ctx context.Context, clientID, clientSecret string) (response *StandardResponse, err error) {
	// Create and fire the request
	var reqResponse *RequestResponse
	if reqResponse, err = c.AppUninstallRaw(ctx, clientID, clientSecret); err != nil {
		return nil, err
	}

	// Parse the response
	err = reqResponse.UnmarshalTo(&response)
	return response, err
}

// AppUninstallRaw triggers app uninstallation and returns the raw response
// specs: https://devdocs.drift.com/docs/app-uninstall
func (c *Client) AppUninstallRaw(ctx context.Context, clientID, clientSecret string) (*RequestResponse, error) {
	// Validate required fields
	if err := requireString(clientID, ErrMissingClientID); err != nil {
		return nil, err
	}
	if err := requireString(clientSecret, ErrMissingClientSecret); err != nil {
		return nil, err
	}

	// Build the URL with query parameters
	queryURL := fmt.Sprintf("%s/app/uninstall?clientId=%s&clientSecret=%s", apiEndpoint, clientID, clientSecret)

	// Fire the request
	response := httpRequest(ctx, c, &httpPayload{
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodPost,
		URL:            queryURL,
	})

	return response, response.Error
}
