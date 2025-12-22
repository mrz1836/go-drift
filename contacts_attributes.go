package drift

import (
	"context"
	"encoding/json"
	"net/http"
)

// CustomAttribute represents a single custom contact attribute
type CustomAttribute struct {
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

// CustomAttributesData is the data wrapper for the response
type CustomAttributesData struct {
	Properties []*CustomAttribute `json:"properties"`
}

// CustomAttributesResponse represents the response from listing custom attributes
type CustomAttributesResponse struct {
	Data *CustomAttributesData `json:"data"`
}

// ListCustomAttributes retrieves all custom contact attributes for the organization.
// specs: https://devdocs.drift.com/docs/listing-custom-attributes
func (c *Client) ListCustomAttributes(ctx context.Context) (response *CustomAttributesResponse, err error) {
	// Create and fire the request
	var reqResponse *RequestResponse
	if reqResponse, err = c.ListCustomAttributesRaw(ctx); err != nil {
		return nil, err
	}

	// Parse the response
	err = json.Unmarshal(reqResponse.BodyContents, &response)
	return response, err
}

// ListCustomAttributesRaw retrieves raw response from custom attributes endpoint
// specs: https://devdocs.drift.com/docs/listing-custom-attributes
func (c *Client) ListCustomAttributesRaw(ctx context.Context) (*RequestResponse, error) {
	// Fire the request
	response := httpRequest(ctx, c, &httpPayload{
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodGet,
		URL:            apiEndpoint + "/contacts/attributes",
	})

	return response, response.Error
}
