package drift

import (
	"context"
	"net/http"
	"strconv"
)

// DeleteContact will fire the HTTP request to delete an existing contact.
// This only removes a contact from indexing in your Drift account's Contacts view.
// For full GDPR-compliant deletion, use the GDPR deletion endpoint.
// specs: https://devdocs.drift.com/docs/removing-a-contact
func (c *Client) DeleteContact(ctx context.Context, contactID uint64) (response *StandardResponse, err error) {
	// Create and fire the request
	var reqResponse *RequestResponse
	if reqResponse, err = c.DeleteContactRaw(ctx, contactID); err != nil {
		return nil, err
	}

	// Parse the response
	err = reqResponse.UnmarshalTo(&response)
	return response, err
}

// DeleteContactRaw will delete a contact and return the raw response
// specs: https://devdocs.drift.com/docs/removing-a-contact
func (c *Client) DeleteContactRaw(ctx context.Context, contactID uint64) (*RequestResponse, error) {
	// Fire the request
	response := httpRequest(ctx, c, &httpPayload{
		ExpectedStatus: http.StatusAccepted,
		Method:         http.MethodDelete,
		URL:            apiEndpoint + "/contacts/" + strconv.FormatUint(contactID, 10),
	})

	return response, response.Error
}
