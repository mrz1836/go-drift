package drift

import (
	"context"
	"encoding/json"
	"net/http"
)

// DeleteGDPR triggers a GDPR data deletion request for the given email.
// WARNING: This permanently deletes all data (contacts, conversations, messages, etc.)
// for all contacts and users with this email. This action cannot be undone.
// specs: https://devdocs.drift.com/docs/gdpr-deletion
func (c *Client) DeleteGDPR(ctx context.Context, email string) (*GDPRDeletionResponse, error) {
	return c.DeleteGDPRWithRequest(ctx, &GDPRRequest{Email: email})
}

// DeleteGDPRWithRequest triggers a GDPR data deletion using a request struct.
// WARNING: This permanently deletes all data and cannot be undone.
// specs: https://devdocs.drift.com/docs/gdpr-deletion
func (c *Client) DeleteGDPRWithRequest(ctx context.Context, request *GDPRRequest) (response *GDPRDeletionResponse, err error) {
	var reqResponse *RequestResponse
	if reqResponse, err = c.DeleteGDPRRaw(ctx, request); err != nil {
		return nil, err
	}

	err = reqResponse.UnmarshalTo(&response)
	return response, err
}

// DeleteGDPRRaw triggers a GDPR data deletion and returns the raw response.
// WARNING: This permanently deletes all data and cannot be undone.
// specs: https://devdocs.drift.com/docs/gdpr-deletion
func (c *Client) DeleteGDPRRaw(ctx context.Context, request *GDPRRequest) (*RequestResponse, error) {
	if request == nil {
		return nil, ErrMissingEmail
	}
	if err := requireString(request.Email, ErrMissingEmail); err != nil {
		return nil, err
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response := httpRequest(ctx, c, &httpPayload{
		Data:           data,
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodPost,
		URL:            apiEndpoint + "/gdpr/delete",
	})

	return response, response.Error
}
