package drift

import (
	"context"
	"encoding/json"
	"net/http"
)

// UnsubscribeResponse represents the response from unsubscribing emails
type UnsubscribeResponse struct {
	OK     bool   `json:"ok,omitempty"`
	Result string `json:"result,omitempty"`
}

// UnsubscribeEmails will unsubscribe a list of email addresses from Drift emails
// specs: https://devdocs.drift.com/docs/unsubscribe-contacts-from-emails
func (c *Client) UnsubscribeEmails(ctx context.Context, emails []string) (response *UnsubscribeResponse, err error) {
	// Create and fire the request
	var reqResponse *RequestResponse
	if reqResponse, err = c.UnsubscribeEmailsRaw(ctx, emails); err != nil {
		return nil, err
	}

	// Parse the response
	err = json.Unmarshal(reqResponse.BodyContents, &response)
	return response, err
}

// UnsubscribeEmailsRaw will unsubscribe emails and return the raw response
// specs: https://devdocs.drift.com/docs/unsubscribe-contacts-from-emails
func (c *Client) UnsubscribeEmailsRaw(ctx context.Context, emails []string) (*RequestResponse, error) {
	// Marshal the email list
	data, err := json.Marshal(emails)
	if err != nil {
		return nil, err
	}

	// Fire the request
	response := httpRequest(ctx, c, &httpPayload{
		Data:           data,
		ExpectedStatus: http.StatusOK,
		Method:         http.MethodPost,
		URL:            apiEndpoint + "/emails/unsubscribe",
	})

	return response, response.Error
}
